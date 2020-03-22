package cmd

import (
	"fmt"
	"os"

	"github.com/jespino/gitsub/config"
	"github.com/spf13/cobra"
)

var addCustomCmd = &cobra.Command{
	Use:   "add-custom",
	Short: "Add new custom subscription",
	Long:  `Add a new custom subscription that is not present in the catalog`,
	Run:   addCustomCmdF,
	Args:  cobra.ExactArgs(1),
}

func init() {
	addCustomCmd.Flags().StringP("org", "o", "", "Github owner or organization")
	addCustomCmd.Flags().StringP("repo", "r", "", "Github repository")
	addCustomCmd.Flags().StringArrayP("label", "l", []string{}, "Labels to subscribe")
	addCustomCmd.MarkFlagRequired("org")
	addCustomCmd.MarkFlagRequired("repo")
	addCustomCmd.MarkFlagRequired("label")
}

func addCustomCmdF(cmd *cobra.Command, args []string) {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Fatal error reading config: %s \n", err)
		os.Exit(1)
	}
	name := args[0]
	if _, ok := cfg.CustomSubscriptions[name]; ok {
		fmt.Printf("the custom subscription %s already exists\n", name)
		os.Exit(1)
	}

	org, _ := cmd.Flags().GetString("org")
	repo, _ := cmd.Flags().GetString("org")
	labels, _ := cmd.Flags().GetStringArray("labels")
	cfg.CustomSubscriptions[name] = config.CustomSubscription{
		Name:   name,
		Org:    org,
		Repo:   repo,
		Labels: labels,
	}
	err = config.Write(cfg)
	if err != nil {
		fmt.Printf("Fatal error writing config: %s \n", err)
		os.Exit(1)
	}
}
