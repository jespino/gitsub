package catalog

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
	"path"
)

// Project represents each of the projects that you can subscribe from the projects catalog
type Project struct {
	Name         string
	Org          string
	Repo         string
	HelpWanted   []string
	FirstIssue   []string
	Languages    map[string]string
	Difficulties map[string]string
}

// Catalog defines the data structure that represents the projects catalog to
// subscribe to
type Catalog struct {
	Format  string
	Catalog map[string]Project
}

func Path() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(usr.HomeDir, ".gitsub.catalog"), nil
}

// Write write a catalog struct into the catalog file
func Write(c *Catalog) error {
	catalogPath, err := Path()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(catalogPath, data, 0600)
	if err != nil {
		return err
	}
	return nil
}

// Read get the catalog information from the default ~/.gitsub.catalog file
func Read() (*Catalog, error) {
	catalogPath, err := Path()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(catalogPath)
	if err != nil {
		return nil, err
	}
	catalog := Catalog{}
	err = json.Unmarshal(data, &catalog)
	if err != nil {
		return nil, err
	}
	return &catalog, nil
}

// GetEtag returns the etag of the current local catalog
func GetEtag() (string, error) {
	catalogPath, err := Path()
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadFile(catalogPath)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha1.Sum(data)), nil
}
