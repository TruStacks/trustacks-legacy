package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/trustacks/trustacks/pkg"
)

// rootCmd is the cobra start command.
var rootCmd = &cobra.Command{
	Use:   "tsctl",
	Short: "Trustacks software delivery engine",
}

func main() {
	if err := os.Setenv("PATH", fmt.Sprintf("%s:%s", os.Getenv("PATH"), pkg.BinDir)); err != nil {
		fmt.Printf("error setting path: %s\n", err)
		os.Exit(1)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error executing the command: %s", err)
	}
}
