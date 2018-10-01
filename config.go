package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"

	"github.com/kardianos/osext"
)

// Debug debug mode.
var Debug bool

// Log log.println
func Log(args ...interface{}) {
	if Debug {
		log.Println(args...)
	}
}

// Config : config info
type Config struct {
	Source       string `json:"source"`
	Target       string `json:"target"`
	DiffType     string `json:"diff_type"`
	Include      string `json:"include"`
	Exclude      string `json:"exlude"`
	Output       string `json:"output"`
	IgnoreColumn string `json:"ignore_column"`
}

// IsValid : valid check config
func (c *Config) IsValid() bool {
	if c.DiffType == "schema" || c.DiffType == "data" {
		return c.Source != "" && c.Target != ""
	} else if c.DiffType == "doc" || c.DiffType == "sql" {
		return c.Source != ""
	}
	return false
}

// LoadConfig load config from file.
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		exePath, _ := osext.ExecutableFolder()
		path = exePath + "/conf.json"
	}

	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(dat, &config); err != nil {
		return nil, err
	}

	if config.IsValid() == false {
		return nil, errors.New("config file is invalid")
	}

	Log("load config file...", path, config)
	return &config, nil
}

// ParseArgs parse argument.
func ParseArgs() (*Config, error) {
	configPath := flag.String("conf", "", "config file path")

	var config Config
	flag.StringVar(&config.Source, "source", "", "source db connection string ex) [uid]:[pwd]@tcp([ip]:[port])/[dbname]")
	flag.StringVar(&config.Target, "target", "", "source db connection string ex) [uid]:[pwd]@tcp([ip]:[port])/[dbname]")
	flag.StringVar(&config.DiffType, "diff_type", "schema", "schema or data")
	flag.StringVar(&config.Include, "include", "", "include object name,name,....")
	flag.StringVar(&config.Exclude, "exclude", "", "exclude object name,name,name,...")
	flag.StringVar(&config.IgnoreColumn, "ignore_column", "", "ignore column split column,column,... ")
	flag.StringVar(&config.Output, "output", "", "result file")
	flag.BoolVar(&Debug, "debug", false, "debug mode")

	flag.Parse()

	if *configPath != "" || !config.IsValid() {
		return LoadConfig(*configPath)
	}

	Log("load config args", config)
	return &config, nil
}
