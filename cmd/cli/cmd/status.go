package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Args:  cobra.NoArgs,
	Short: "Checks the status to the backend with current credentials",
	Long: `Checks the status to the backend, using team name and password.
Will also show other relevant information like config file location`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status called")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
