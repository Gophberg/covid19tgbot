package bot

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Config ...
type Config struct {
	Token            string `yaml:"token"`
	DbDataSourceName string `yaml:"db_data_source_name"`
	URL              string `yaml:"url"`
}

// NewConfig ...
func NewConfig() *Config {
	var config Config
	f, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	cfg, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(cfg, &config)
	if err != nil {
		log.Fatal(err)
	}
	return &config
}
