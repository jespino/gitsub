package cmd

import (
	"fmt"
	"os"

	"github.com/jespino/gitsub/catalog"
	"github.com/jespino/gitsub/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize user config",
	Long:  `Initialize the user config and catalog file`,
	Run:   initCmdF,
}

func initCmdF(cmd *cobra.Command, args []string) {
	cfg := config.Config{
		Main:                config.MainConfig{SleepSeconds: 600},
		Subscriptions:       map[string]config.Subscription{},
		CustomSubscriptions: map[string]config.CustomSubscription{},
	}
	if config.Exists() {
		if !askForConfirmation("Config file already exists, do you want to overwrite it?") {
			fmt.Printf("Aborting init command")
			os.Exit(1)
		}
	}
	if catalog.Exists() {
		if !askForConfirmation("Catalog file already exists, do you want to overwrite it?") {
			fmt.Printf("Aborting init command")
			os.Exit(1)
		}
	}
	err := config.Write(&cfg)
	if err != nil {
		fmt.Printf("Fatal error writing config: %s \n", err)
		os.Exit(1)
	}

	cat := catalog.Catalog{
		Format:  "1.0",
		Catalog: map[string]catalog.Project{},
	}
	err = catalog.Write(&cat)
	if err != nil {
		fmt.Printf("Fatal error writing catalog: %s \n", err)
		os.Exit(1)
	}
}
