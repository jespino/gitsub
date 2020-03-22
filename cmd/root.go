package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitsub",
	Short: "GitSub is a GitHub subscription tool",
	Long:  `GitHub subscription tool to be aware of new help wanted issues in multiple GitHub projects.`,
}

func init() {
	rootCmd.AddCommand(watchCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(addCustomCmd)
	rootCmd.AddCommand(catalogUpdateCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(initCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
