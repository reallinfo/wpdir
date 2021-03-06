package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config contains useful App data.
type Config struct {
	Version  string
	Name     string
	Commit   string
	Date     string
	DataDir  string
	IndexDir string
	WD       string
	Ports    struct {
		HTTP  string
		HTTPS string
	}
}

// Setup creates and returns the Config struct.
func Setup() *Config {

	viper.SetDefault("version", "0.1.0")
	viper.SetDefault("name", "WP Directory")
	viper.SetDefault("commit", "")
	viper.SetDefault("date", "")
	viper.SetDefault("ports.http", "9077")
	viper.SetDefault("ports.https", "9078")
	viper.SetDefault("datadir", "F:\\wpdir\\data")
	viper.SetDefault("indexdir", "F:\\wpdir\\index")

	viper.AddConfigPath("/etc/wpdir/")
	viper.AddConfigPath(".")

	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Error reading config file: %s\n", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting Working Directory: %s\n", err)
	}

	config := &Config{
		Version:  viper.GetString("version"),
		Name:     viper.GetString("name"),
		Commit:   viper.GetString("commit"),
		Date:     viper.GetString("date"),
		DataDir:  viper.GetString("datadir"),
		IndexDir: viper.GetString("indexdir"),
		WD:       wd,
	}

	config.Ports.HTTP = viper.GetString("ports.http")
	config.Ports.HTTPS = viper.GetString("ports.https")

	return config

}
