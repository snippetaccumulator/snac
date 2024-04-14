package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	configFileParameter string
	urlParameter        string
	tokenParameter      string
	teamNameParameter   string
	passwordParameter   string
	rootCmd             = &cobra.Command{
		Use:   "snac",
		Short: "CLI to interact with a snac server",
		Long: `CLI to use snac tooling to manage and use snippets.
Can do all operations either interactively or with flags.`,
		Version: "0.0.1",
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFileParameter, "config", "", "Override path for config file to use (default is the UserConfigDir/snac/config.yaml")
	rootCmd.PersistentFlags().StringVar(&urlParameter, "url", "", "Override URL for database connection")
	rootCmd.PersistentFlags().StringVar(&tokenParameter, "token", "", "Override token for database connection")
	rootCmd.PersistentFlags().StringVar(&teamNameParameter, "team-name", "", "Override team name for connection")
	rootCmd.PersistentFlags().StringVar(&passwordParameter, "password", "", "Override password for connection")
	rootCmd.MarkFlagsRequiredTogether("url", "token")
	rootCmd.MarkFlagsRequiredTogether("team-name", "password")
}
