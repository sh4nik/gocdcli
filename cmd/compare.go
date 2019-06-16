package cmd

import (
	"github.com/sh4nik/gopipe/client"
	"github.com/spf13/cobra"
)

// compareCmd represents the compare command
var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare two pipelines",
	Run: func(cmd *cobra.Command, args []string) {
		client.ComparePipes(args[0], args[1], false)
	},
}

func init() {
	rootCmd.AddCommand(compareCmd)
}
