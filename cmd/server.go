package main

import (
	"github.com/spf13/cobra"
	"github.com/trustacks/trustacks/pkg/api/server"
)

var (
	serverHost string
	serverPort string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start the api server",
	Run: func(cmd *cobra.Command, args []string) {
		server.New(serverHost, serverPort)
	},
}

func init() {
	serverCmd.Flags().StringVar(&serverHost, "host", "", `server host (default "*")`)
	serverCmd.Flags().StringVar(&serverPort, "port", "8080", "server port")
	rootCmd.AddCommand(serverCmd)
}
