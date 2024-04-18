package cmd

import (
	"os"
	"path/filepath"

	"github.com/snippetaccumulator/configloader"
	"github.com/snippetaccumulator/snac/internal/backend/database"
	"github.com/snippetaccumulator/snac/internal/cli"
	"github.com/snippetaccumulator/snac/internal/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	config    cli.Config
	configLoc string
	db        database.Database

	configFileParameter       string
	createConfigFileParameter bool
	urlParameter              string
	tokenParameter            string
	teamNameParameter         string
	passwordParameter         string
	verboseParameter          bool
	rootCmd                   = &cobra.Command{
		Use:   "snac",
		Short: "CLI to interact with a snac server",
		Long: `CLI to use snac tooling to manage and use snippets.
Can do all operations either interactively or with flags.`,
		Version: "0.0.1",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cfgPathDir, err := os.UserConfigDir()
			if err != nil {
				log.Err(true, err)
			}
			cfgDir := filepath.Join(cfgPathDir, "snac")
			configLoc = filepath.Join(cfgDir, "config.yaml")
			if configFileParameter != "" {
				configLoc = configFileParameter
			}

			//check if config file exists
			if _, err := os.Stat(configLoc); os.IsNotExist(err) {
				if createConfigFileParameter {
					dir := filepath.Dir(configLoc)
					err := os.MkdirAll(dir, os.ModePerm)
					if err != nil {
						log.Err(true, err)
					}
					file, err := os.Create(configLoc)
					if err != nil {
						log.Err(true, err)
					}

					//write empty config struct to config
					emptyConfig := cli.Config{}
					data, err := yaml.Marshal(emptyConfig)
					if err != nil {
						log.Error(true, "Error while writing new config file: %s", err)
					}
					_, err = file.Write(data)
					if err != nil {
						log.Error(true, "Error while writing new config file: %s", err)
					}

					log.Info("Created new empty config file at %s", configLoc)

					file.Close()
				} else {
					log.Info("Use --create-config flag to create a new config file")
					log.Error(true, "Config file does not exist at %s", configLoc)
				}
			}

			cfgLoader := configloader.NewConfigLoader("config.yaml",
				configloader.WithPath(cfgDir),
				configloader.WithDeserializer(&configloader.YAMLDeserializer{}),
			)

			if urlParameter != "" {
				cfgLoader.Override("Database.Url", urlParameter)
			}
			if tokenParameter != "" {
				cfgLoader.Override("Database.AuthToken", tokenParameter)
			}
			if teamNameParameter != "" {
				cfgLoader.Override("TeamName", teamNameParameter)
			}
			if passwordParameter != "" {
				cfgLoader.Override("Password", passwordParameter)
			}
			if verboseParameter {
				cfgLoader.Override("LogLevel", "DEBUG")
			}

			err = cfgLoader.Load(&config)
			if err != nil {
				log.Error(true, "Error while loading config file: %s", err)
			}

			log.SetLevel(log.FromString(config.LogLevel))

			log.Debug("Loaded config")

			db, err = database.NewDB(config.Database)
			if err != nil {
				log.Error(true, "Error while creating database connection: %s", err)
			}
		},
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
	rootCmd.PersistentFlags().BoolVar(&createConfigFileParameter, "create-config", false, "Create a new config file with the provided flags; has no effect if config file already exists")
	rootCmd.PersistentFlags().StringVar(&urlParameter, "url", "", "Override URL for database connection")
	rootCmd.PersistentFlags().StringVar(&tokenParameter, "token", "", "Override token for database connection")
	rootCmd.PersistentFlags().StringVar(&teamNameParameter, "team-name", "", "Override team name for connection")
	rootCmd.PersistentFlags().StringVar(&passwordParameter, "password", "", "Override password for connection")
	rootCmd.MarkFlagsRequiredTogether("url", "token")
	rootCmd.MarkFlagsRequiredTogether("team-name", "password")

	rootCmd.PersistentFlags().BoolVarP(&verboseParameter, "verbose", "V", false, "Enable verbose (debug) output")
}
