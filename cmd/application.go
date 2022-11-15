package main

import (
	"fmt"
	"log"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	toolchain "github.com/trustacks/trustacks/pkg/provision"
)

// application cli command flags.
var (
	applicationName   string
	applicationConfig string
	applicationForce  bool
)

// applicationCmd contains subcommands for managing factories.
var applicationCmd = &cobra.Command{
	Use:   "application",
	Short: "manage applications",
}

// applicationCreateCmd creates a new application.
var applicationCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new application",
	Run: func(cmd *cobra.Command, args []string) {
		if err := toolchain.CreateApplication(applicationName, applicationForce, applicationConfig, git.PlainClone); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	applicationCmd.AddCommand(applicationCreateCmd)

	applicationCreateCmd.Flags().StringVar(&applicationName, "name", "", "application name")
	if err := applicationCreateCmd.MarkFlagRequired("name"); err != nil {
		log.Fatal(err)
	}
	applicationCreateCmd.Flags().StringVar(&applicationConfig, "config", "", "configuration file")
	if err := applicationCreateCmd.MarkFlagRequired("config"); err != nil {
		log.Fatal(err)
	}
	applicationCreateCmd.Flags().BoolVar(&applicationForce, "force", false, "force update (experimental: use at your own risk)")

	rootCmd.AddCommand(applicationCmd)
}
