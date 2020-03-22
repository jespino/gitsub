package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

// CustomSubscription represents a user defined subscription to a project
type CustomSubscription struct {
	Name   string
	Org    string
	Repo   string
	Labels []string
}

// Subscription represents a subscription to a catalog entry
type Subscription struct {
	Name         string
	Difficulties []string
	Languages    []string
	FirstIssue   bool
}

// MainConfig represents the global settings parameters for the configuration
type MainConfig struct {
	SleepSeconds int
	GithubToken  string
}

// Config defines the data structure that represents the user config
type Config struct {
	Main                MainConfig
	Subscriptions       map[string]Subscription
	CustomSubscriptions map[string]CustomSubscription
}

func Path() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(usr.HomeDir, ".gitsub.conf"), nil
}

func Exists() bool {
	configPath, _ := Path()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func Write(cfg *Config) error {
	configPath, err := Path()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configPath, data, 0600)
	if err != nil {
		return err
	}
	return nil
}

// Read get the config information from the default ~/.gitsub.conf file
func Read() (*Config, error) {
	configPath, err := Path()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	config := Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
