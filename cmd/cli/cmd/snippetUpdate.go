package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	updateNonInteractiveParameter bool
	updateTitleParameter          string
	updateDescriptionParameter    string
	updateLanguageParameter       string
	updateTagsToAddParameter      []string
	updateTagsToRemoveParameter   []string

	updateCmd = &cobra.Command{
		Use:     "update <ID>",
		Aliases: []string{"u"},
		Args:    cobra.ExactArgs(1),
		Short:   "Update a snippet by ID",
		Long: `Interactively update a snippet by ID.
Can be used non-interactively by providing all fields to override.
Can add and remove tags.`,
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			fmt.Printf("update called with id: %s\n", id)
			fmt.Println("TODO: Implement update command")
		},
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolVarP(&updateNonInteractiveParameter, "non-interactive", "n", false, "Non interactive mode if used, all other fields must be given as aruments")
	updateCmd.Flags().StringVar(&updateTitleParameter, "title", "", "New title of the snippet")
	updateCmd.Flags().StringVar(&updateDescriptionParameter, "description", "", "New description of the snippet")
	updateCmd.Flags().StringVar(&updateLanguageParameter, "language", "", "New language of the snippet")
	updateCmd.Flags().StringArrayVar(&updateTagsToAddParameter, "tag", []string{}, "Add a single tag to the snippet, can be used multiple times")
	updateCmd.Flags().StringArrayVar(&updateTagsToRemoveParameter, "untag", []string{}, "Remove a single tag from the snippet, can be used multiple times")

	createCmd.Flags().SortFlags = false
}
