package cmd

import (
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove subscription",
	Long:  `Remove one of your subscriptions`,
	Run:   removeCmdF,
}

func removeCmdF(cmd *cobra.Command, args []string) {
	panic("Not implemented")
}
