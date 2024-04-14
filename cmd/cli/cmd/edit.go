package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	skipMetaParameter  bool
	fileParameter      string
	editorBinParameter string
	editCmd            = &cobra.Command{
		Use:     "edit <ID>",
		Aliases: []string{"e"},
		Args:    cobra.ExactArgs(1),
		Short:   "Edit snippet content in $EDITOR",
		Long: `Put snippet content into temporary file and open it in $EDITOR, or provided editor binary.
Afterwards update metadata as if using update command.
Options to skip metadata update or provide file directly.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("edit called")
		},
	}
)

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().BoolVarP(&skipMetaParameter, "skip-meta", "s", false, "Skip asking with any meta data should be changed")
	editCmd.Flags().StringVarP(&fileParameter, "file", "f", "", "File to use for content")

	editCmd.Flags().SortFlags = false
}
