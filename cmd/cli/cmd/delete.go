package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete <ID>",
	Aliases: []string{"d"},
	Args:    cobra.ExactArgs(1),
	Short:   "Deletes a snippet by ID",
	Long:    `Deletes a snippet by ID. Nothing more, nothing less. Why are you looking at this? This just deletes...`,
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		fmt.Printf("delete called with id: %s\n", id)
		fmt.Println("TODO: Implement delete command")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

}
