package cmd

import (
	"github.com/sh4nik/gopipe/client"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all pipelines",
	Run: func(cmd *cobra.Command, args []string) {
		client.ListPipes()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
