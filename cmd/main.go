package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// rootCmd is the cobra start command.
var rootCmd = &cobra.Command{
	Use:   "tsctl",
	Short: "Trustacks software delivery engine",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error executing the command: %s", err)
	}
}
