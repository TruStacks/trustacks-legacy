package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// current release version.
var currentVersion string

// versionCmd contains subcommands for managing factories.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show the cli version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(currentVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
