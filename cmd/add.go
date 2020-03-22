package cmd

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new subscription",
	Long:  `Add a new subscription`,
	Run:   addCmdF,
}

func addCmdF(cmd *cobra.Command, args []string) {
	panic("Not implemented")
}
