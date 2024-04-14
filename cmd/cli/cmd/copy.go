package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	silentParameter     bool
	outputFileParameter string

	copyCmd = &cobra.Command{
		Use:     "copy <ID>",
		Aliases: []string{"cp"},
		Args:    cobra.ExactArgs(1),
		Short:   "Copy snippet content to clipboard or file",
		Long: `Copies content of a snippet to the clipboard.
Can also write the content to a file.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("copy called")
		},
	}
)

func init() {
	rootCmd.AddCommand(copyCmd)

	copyCmd.Flags().BoolVarP(&silentParameter, "silent", "s", false, "Do not print anything to the console")
	copyCmd.Flags().StringVarP(&outputFileParameter, "output", "o", "", "Write content to a file instead of copying it to the clipboard. Can use '-' to write to stdout")

	copyCmd.Flags().SortFlags = false
}
