package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	queryParameter            string
	tagsParameter             []string
	languageParameter         string
	minContentLengthParameter int
	maxContentLengthParameter int
	fullShowParameter         bool
	listFormatParameter       formatParameterValue

	listCmd = &cobra.Command{
		Use:     "list [flags]",
		Aliases: []string{"l"},
		Short:   "List snippets filter/query options",
		Long: `List snippets.
Can filter by tags, language, content length min and max, and query for matching text in title and description.
Can show full data of the snippet. Can output in different formats.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("list called")
		},
	}
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&queryParameter, "query", "q", "", "Search for matching text in title and description")
	listCmd.Flags().StringArrayVar(&tagsParameter, "tag", []string{}, "Filter by tags, can be used multiple times. Multiple tags means that snippet must match ANY not ALL tags")
	listCmd.Flags().StringVar(&languageParameter, "language", "", "Filter by language")
	listCmd.Flags().IntVar(&minContentLengthParameter, "min-content-length", -1, "Filter by minimum content length (character count) (inclusive)")
	listCmd.Flags().IntVar(&maxContentLengthParameter, "max-content-length", -1, "Filter by maximum content length (character count) (inclusive)")
	listCmd.Flags().BoolVar(&fullShowParameter, "full", false, "Show full data of the snippet")
	listCmd.Flags().VarP(&listFormatParameter, "format", "f", "Output format (allowed values: 'json', 'yaml', 'default')")

	rootCmd.Flags().SortFlags = false
}
