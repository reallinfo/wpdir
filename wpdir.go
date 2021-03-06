package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wpdirectory/wpdir/internal/config"
	"github.com/wpdirectory/wpdir/internal/db"
	"github.com/wpdirectory/wpdir/internal/log"
	"github.com/wpdirectory/wpdir/internal/server"
)

func main() {

	// Set and Parse flags
	flagHelp := flag.Bool("help", false, "")
	flag.Parse()

	if *flagHelp {
		fmt.Println(helpText)
		os.Exit(1)
	}

	fmt.Println("Starting WPDirectory")

	// Setup Stats.
	//stats.Setup()

	// Create Logger
	l := log.New()

	// Create Config
	c := config.Setup()

	// Setup BoltDB
	db.Setup(c.WD)
	defer db.Close()

	// Setup server struct to hold all App data
	s := server.New(l, c)

	// Setup HTTP server.
	s.Setup()

}

const (
	helpText = `WPDirectory is a web service for lightning fast code searching of the WordPress Plugin & Theme Directories.

Usage:
  wpdir [flags]

Version: 0.5.0
	
Flags:
  --help      Help outputs help text and exits.
  
Config:
  WPDirectory requires a config file, located at /etc/wpdir/ or in the working directory, to successfully run. See the example-config.yml.`
)
