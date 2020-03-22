package cmd

import (
	"fmt"
	"os"

	"github.com/jespino/gitsub/catalog"
	"github.com/jespino/gitsub/config"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new subscription",
	Long:  `Add a new subscription`,
	Run:   addCmdF,
}

func init() {
	addCmd.Flags().BoolP("first-issue", "f", false, "Labeled as good first issue")
	addCmd.Flags().StringArrayP("difficulty", "d", []string{}, "Difficulty level: easy, medium, hard (multiple allowed)")
	addCmd.Flags().StringArrayP("language", "l", []string{}, "Language: go, javascript, typescript, java, kotlin, rust... (multiple allowed)")
}

func addCmdF(cmd *cobra.Command, args []string) {
	name := args[0]
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Fatal error reading config: %s \n", err)
		os.Exit(1)
	}
	cat, err := catalog.Read()
	if err != nil {
		fmt.Printf("Fatal error reading catalog: %s \n", err)
		os.Exit(1)
	}

	if _, ok := cat.Catalog[name]; !ok {
		fmt.Printf("Project %s not found in the catalog.\n", name)
		os.Exit(1)
	}

	difficulties, _ := cmd.Flags().GetStringArray("difficulty")
	languages, _ := cmd.Flags().GetStringArray("language")
	firstIssue, _ := cmd.Flags().GetBool("first-issue")

	cfg.Subscriptions[name] = config.Subscription{
		Name:         name,
		Difficulties: difficulties,
		Languages:    languages,
		FirstIssue:   firstIssue,
	}
	err = config.Write(cfg)
	if err != nil {
		fmt.Printf("Fatal error writing config: %s \n", err)
		os.Exit(1)
	}
}
