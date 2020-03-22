package cmd

import (
	"fmt"
	"os"

	"github.com/jespino/gitsub/config"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove subscription",
	Long:    `Remove one of your subscriptions`,
	Example: "remove my-subscription",
	Run:     removeCmdF,
	Args:    cobra.ExactArgs(1),
}

func removeCmdF(cmd *cobra.Command, args []string) {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Fatal error reading config: %s \n", err)
		os.Exit(1)
	}
	name := args[0]
	if _, ok := cfg.Subscriptions[name]; ok {
		delete(cfg.Subscriptions, name)
		err = config.Write(cfg)
		if err != nil {
			fmt.Printf("Fatal error writing config: %s \n", err)
			os.Exit(1)
		}
		return
	}

	if _, ok := cfg.CustomSubscriptions[name]; ok {
		delete(cfg.CustomSubscriptions, name)
		err = config.Write(cfg)
		if err != nil {
			fmt.Printf("Fatal error writing config: %s \n", err)
			os.Exit(1)
		}
		return
	}

	fmt.Printf("Subscription %s not found.\n", name)
	os.Exit(1)
}
