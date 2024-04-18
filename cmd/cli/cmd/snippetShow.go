package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	shortParameter       bool
	contentOnlyParameter bool
	cutOffParameter      int
	showFormatParameter  formatParameterValue

	showCmd = &cobra.Command{
		Use:     "show <ID>",
		Aliases: []string{"s", "get", "g"},
		Args:    cobra.ExactArgs(1),
		Short:   "",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("show called")
		},
	}
)

func init() {
	rootCmd.AddCommand(showCmd)

	showCmd.Flags().BoolVarP(&shortParameter, "short", "s", false, "Show only minimal information about the snippet")
	showCmd.Flags().BoolVar(&contentOnlyParameter, "content-only", false, "Show only the content of the snippet")
	showCmd.Flags().IntVar(&cutOffParameter, "cutoff", -1, "Cut off the content after a certain number of characters")
	showCmd.Flags().Var(&showFormatParameter, "format", "Output format (allowed values: 'json', 'yaml', 'default')")

	showCmd.Flags().SortFlags = false
}
