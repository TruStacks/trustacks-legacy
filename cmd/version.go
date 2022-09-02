package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd contains subcommands for managing factories.
var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cliVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
