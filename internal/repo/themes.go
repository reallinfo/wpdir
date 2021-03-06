package repo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/wpdirectory/wpdir/internal/config"
	"github.com/wpdirectory/wpdir/internal/db"
	"github.com/wpdirectory/wpdir/internal/index"
	"github.com/wpdirectory/wpdir/internal/svn"
	"github.com/wpdirectory/wpdir/internal/theme"
)

const (
	themeManagementUser = "theme-master"
)

// ThemeRepo ...
type ThemeRepo struct {
	Config *config.Config
	List   map[string]*theme.Theme
	sync.RWMutex
	Revision    int
	Updated     time.Time
	UpdateQueue chan string
}

// Len returns the number of Themes
func (tr *ThemeRepo) Len() int {
	return len(tr.List)
}

// Rev returns the current Revision
func (tr *ThemeRepo) Rev() int {
	return tr.Revision
}

// Exists checks if a Theme exists
func (tr *ThemeRepo) Exists(slug string) bool {
	tr.Lock()
	_, ok := tr.List[slug]
	tr.Unlock()
	return ok
}

// Get returns a Theme
func (tr *ThemeRepo) Get(slug string) Extension {
	tr.Lock()
	p := tr.List[slug]
	tr.Unlock()
	return p
}

// Add sets a new Theme
func (tr *ThemeRepo) Add(slug string) {
	tr.RLock()
	tr.List[slug] = &theme.Theme{
		Slug: slug,
	}
	tr.RUnlock()
	//tr.QueueUpdate(slug)
}

// Set ...
func (tr *ThemeRepo) Set(slug string, t *theme.Theme) {
	tr.RLock()
	tr.List[slug] = t
	tr.RUnlock()
}

// Remove deletes a current Theme
func (tr *ThemeRepo) Remove(slug string) {
	tr.RLock()
	delete(tr.List, slug)
	tr.RUnlock()
}

// UpdateIndex ...
func (tr *ThemeRepo) UpdateIndex(idx *index.Index) error {
	var slug string
	if slug = idx.Ref.Slug; slug == "" {
		// bad index, perhaps delete?
		return errors.New("Index contains empty slug")
	}
	tr.RLock()
	defer tr.RUnlock()

	if !tr.Exists(slug) {
		return errors.New("Index does not many an existing theme")
	}

	err := tr.List[slug].Searcher.SwapIndexes(idx)
	if err != nil {
		tr.List[slug].SetIndexed(false)
		return err
	}

	tr.List[slug].SetIndexed(true)

	return nil
}

// QueueUpdate adds a Theme to the update queue
func (tr *ThemeRepo) QueueUpdate(slug string) {
	tr.UpdateQueue <- slug
}

// UpdateWorker processes updates from the update queue
func (tr *ThemeRepo) UpdateWorker() {
	for {
		slug := <-tr.UpdateQueue
		err := tr.ProcessUpdate(slug)
		if err != nil {
			tr.UpdateQueue <- slug
		}
	}
}

// ProcessUpdate updates Theme data and indexes
func (tr *ThemeRepo) ProcessUpdate(slug string) error {
	t := tr.Get(slug).(*theme.Theme)
	err := t.LoadAPIData()
	if err != nil {
		return err
	}

	err = t.Update()
	if err != nil {
		t.SetIndexed(false)
		return err
	}

	t.SetIndexed(true)

	return nil
}

// UpdateList updates our list of themes.
func (tr *ThemeRepo) UpdateList() error {
	// Fetch list from SVN
	// https://themes.svn.wordpress.org/
	list, err := svn.GetList("themes", "")
	if err != nil {
		return err
	}

	for _, item := range list {
		if !utf8.Valid([]byte(item.Name)) {
			return errors.New("Theme slug is not valid utf8")
		}
		if !tr.Exists(item.Name) {
			tr.Add(item.Name)
		}
	}

	return nil
}

// LoadExisting ...
func (tr *ThemeRepo) LoadExisting() {

	tr.loadDBData()
	tr.loadIndexes()

}

// loadDBData loads all existing Theme data from the DB.
func (tr *ThemeRepo) loadDBData() {
	themes, err := db.GetAllFromBucket("themes")
	if err != nil {
		return
	}

	log.Printf("Found %d Theme(s) in DB.\n", len(themes))

	for slug, bytes := range themes {
		var t theme.Theme
		err := json.Unmarshal(bytes, &t)
		if err != nil {
			continue
		}

		tr.Set(slug, &t)
	}
}

// loadIndexes reads all existing Indexes and attempts to match them to a Theme.
func (tr *ThemeRepo) loadIndexes() {
	indexDir := filepath.Join(tr.Config.WD, "data", "index", "themes")

	dirs, err := ioutil.ReadDir(indexDir)
	if err != nil {
		log.Printf("Failed to read Theme index dir: %s\n", err)
	}

	log.Printf("Found %d existing Theme indexes.\n", len(dirs))

	for _, dir := range dirs {
		// If not Directory discard.
		if !dir.IsDir() {
			continue
		}

		path := filepath.Join(indexDir, dir.Name())

		// Read Index
		ref, err := index.Read(path)
		if err != nil {
			os.RemoveAll(path)
			continue
		}

		// Create Index
		idx, err := ref.Open()
		if err != nil {
			os.RemoveAll(path)
			continue
		}

		err = tr.UpdateIndex(idx)
		if err != nil {
			os.RemoveAll(path)
			continue
		}
	}
}

// Summary ...
func (tr *ThemeRepo) Summary() *RepoSummary {
	tr.Lock()
	defer tr.Unlock()

	rs := &RepoSummary{
		Revision: tr.Revision,
		Total:    len(tr.List),
		Queue:    len(tr.UpdateQueue),
	}

	for _, t := range tr.List {
		t.Lock()
		if t.Status == 1 {
			rs.Closed++
		}
		t.Unlock()
	}

	return rs
}
