// root.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubectl-declan",
	Short: "A kubectl plugin to clean namespaces",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Execute()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
