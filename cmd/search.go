package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jespino/gitsub/catalog"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search a project in the catalog",
	Long:  `Search for a project in the projects catalog.`,
	Run:   searchCmdF,
	Args:  cobra.MinimumNArgs(1),
}

func printProject(id string, p *catalog.Project) {
	fmt.Printf("%s (%s):\n", id, p.Name)
	fmt.Printf("  - Url: https://github.com/%s/%s\n", p.Org, p.Repo)
	languages := []string{}
	for lang := range p.Languages {
		languages = append(languages, lang)
	}
	fmt.Printf("  - Languages: %s\n", strings.Join(languages, ", "))
}

func checkMatch(str string, words []string) bool {
	strLower := strings.ToLower(str)
	for _, word := range words {
		if strings.Contains(strLower, word) {
			return true
		}
	}
	return false
}

func searchCmdF(cmd *cobra.Command, args []string) {
	cat, err := catalog.Read()
	if err != nil {
		fmt.Printf("Fatal error reading catalog: %s \n", err)
		os.Exit(1)
	}

	words := []string{}
	for _, word := range args {
		words = append(words, strings.ToLower(word))
	}

	for id, p := range cat.Catalog {
		if checkMatch(id, words) || checkMatch(p.Name, words) || checkMatch(p.Org, words) || checkMatch(p.Repo, words) {
			printProject(id, &p)
		}
	}
}
