package cmd

import (
	"fmt"
	"os"

	"github.com/jespino/gitsub/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List current subscriptions",
	Long:  `List the current existing subscriptions`,
	Run:   listCmdF,
}

func listCmdF(cmd *cobra.Command, args []string) {
	config, err := config.Read()
	if err != nil {
		fmt.Printf("Fatal error reading config: %s \n", err)
		os.Exit(1)
	}
	if len(config.Subscriptions) == 0 && len(config.CustomSubscriptions) == 0 {
		fmt.Println("No subscriptions found")
		os.Exit(0)
	}
	if len(config.Subscriptions) > 0 {
		fmt.Println("Catalog subscriptions:")
		for _, subscription := range config.Subscriptions {
			fmt.Printf("  - %s\n", subscription.Name)
		}
	}
	if len(config.CustomSubscriptions) > 0 {
		fmt.Println("Custom subscriptions:")
		for _, subscription := range config.CustomSubscriptions {
			fmt.Printf("  - %s\n", subscription.Name)
		}
	}
}
