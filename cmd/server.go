package main

import (
	"github.com/spf13/cobra"
	"github.com/trustacks/trustacks/server"
)

// server cli command flags.
var (
	serverHost    string
	serverPort    string
	serverUseTLS  bool
	serverTLSCert string
	serverTLSKey  string
)

// applicationCmd contains subcommands for managing factories.
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start the api server",
	Run: func(cmd *cobra.Command, args []string) {
		if serverUseTLS && serverPort == "8080" {
			serverPort = "8443"
		}
		server.Start(serverHost, serverPort, serverUseTLS, serverTLSCert, serverTLSKey)
	},
}

func init() {
	serverCmd.Flags().StringVar(&serverHost, "host", "0.0.0.0", "server host")
	serverCmd.Flags().StringVar(&serverPort, "port", "8080", "server port")
	serverCmd.Flags().BoolVar(&serverUseTLS, "tls", false, "use secure tls")
	serverCmd.Flags().StringVar(&serverTLSCert, "tls-cert", "/tls.crt", "tls certificate path")
	serverCmd.Flags().StringVar(&serverTLSKey, "tls-key", "/tls.key", "tls key path")

	rootCmd.AddCommand(serverCmd)
}
