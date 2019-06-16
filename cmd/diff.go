package cmd

import (
	"github.com/sh4nik/gopipe/client"
	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Diff two pipelines",
	Run: func(cmd *cobra.Command, args []string) {
		client.ComparePipes(args[0], args[1], true)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
