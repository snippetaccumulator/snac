package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	createNonInteractiveParameter bool
	createTitleParameter          string
	createDescriptionParameter    string
	createTagsParameter           []string
	createLanguageParameter       string
	createContentParameter        string
	createContentFileParameter    string

	createCmd = &cobra.Command{
		Use:     "create [flags]",
		Aliases: []string{"c"},
		Args:    cobra.NoArgs,
		Short:   "Creates a new snippet",
		Long: `Interactively create a new snippet.
If non-interactive mode is used, all required fields must be given as arguments (title, at least one tag, some content).
If interactive mode is used, any of the field flags will be ignored.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if createNonInteractiveParameter {
				if createTitleParameter == "" || len(createTagsParameter) == 0 || (createContentParameter == "" && createContentFileParameter == "") {
					return fmt.Errorf("All required fields must be given in non-interactive mode")
				}
				fmt.Println("TODO: Implement non-interactive mode")
				return nil
			}
			fmt.Println("TODO: Implement interactive mode")
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVarP(&createNonInteractiveParameter, "non-interactive", "n", false, "Non interactive mode if used, all other fields must be given as aruments")
	createCmd.Flags().StringVar(&createTitleParameter, "title", "", "Title of the snippet")
	createCmd.Flags().StringVar(&createDescriptionParameter, "description", "", "Description of the snippet")
	createCmd.Flags().StringVar(&createLanguageParameter, "language", "", "Language of the snippet")
	createCmd.Flags().StringArrayVar(&createTagsParameter, "tag", []string{}, "Add a single tag to the snippet, can be used multiple times")

	createCmd.Flags().StringVar(&createContentParameter, "content", "", "Content of the snippet")
	createCmd.Flags().StringVar(&createContentFileParameter, "content-file", "", "File containing the content of the snippet")
	createCmd.MarkFlagsMutuallyExclusive("content", "content-file")

	createCmd.Flags().SortFlags = false
}
