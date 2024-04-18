package cmd

import (
	"github.com/snippetaccumulator/snac/internal/backend/request"
	"github.com/snippetaccumulator/snac/internal/log"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Args:  cobra.NoArgs,
	Short: "Checks the status to the backend with current credentials",
	Long: `Checks the status to the backend, using team name and password.
Will also show other relevant information like config file location`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Checking status...")
		reqBuilder := request.NewRequestBuilder()
		req := reqBuilder.
			ForTeamByID(config.TeamName, config.Password, false).
			Check().Build()
		_ = req
		retData, retType, err := req.Execute(db)
		log.Err(true, err)
		err = request.TypeCheck(retData, retType)
		log.Err(true, err)

		log.Success("Connection check for team '%s' successful", config.TeamName)

		log.Info("Config file location: %s", configLoc)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
