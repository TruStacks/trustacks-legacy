package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/trustacks/trustacks/pkg/workflows/worker"

	// import workflows
	_ "github.com/trustacks/trustacks/pkg/workflows/react"
)

var (
	workerApplication string
	workerKind        string
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "create a worker process",
	Run: func(cmd *cobra.Command, args []string) {
		if err := worker.New(workerApplication, workerKind); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	workerCmd.Flags().StringVar(&workerApplication, "application", "", "the name of the application")
	if err := workerCmd.MarkFlagRequired("application"); err != nil {
		log.Fatal(err)
	}
	workerCmd.Flags().StringVar(&workerKind, "kind", "", "the workflow type (ie. react)")
	if err := workerCmd.MarkFlagRequired("kind"); err != nil {
		log.Fatal(err)
	}
	rootCmd.AddCommand(workerCmd)
}
