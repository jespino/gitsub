package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jespino/gitsub/catalog"
	"github.com/spf13/cobra"
)

const CATALOG_FILE_URL = "https://github.com/jespino/gitsub/catalog.json"
const CATALOG_PATH = "~/.gitsub-catalog"

var catalogUpdateCmd = &cobra.Command{
	Use:   "catalog-update",
	Short: "Update the projects catalog",
	Long:  `Update the list of available projects to contribute in the catalog`,
	Run:   catalogUpdateCmdF,
}

func catalogUpdateCmdF(cmd *cobra.Command, args []string) {
	catalogPath, err := catalog.Path()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	etag, err := catalog.GetEtag()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	resp, err := http.Get(fmt.Sprintf("%s?etag=%s", CATALOG_FILE_URL, etag))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(catalogPath, data, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Catalog updated")
}
