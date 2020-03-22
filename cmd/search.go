package cmd

import (
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search a project in the catalog",
	Long:  `Search for a project in the projects catalog.`,
	Run:   searchCmdF,
}

func searchCmdF(cmd *cobra.Command, args []string) {
	panic("Not implemented")
}
