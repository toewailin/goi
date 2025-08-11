package commands

import (
	"fmt"
	"goi/config"
	"goi/utils"

	"github.com/spf13/cobra"
)

// Version command to show the current version of GoI
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current version of the goi CLI",
	Long:  `The 'version' command displays the current version of the goi command-line interface tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintSuccess(fmt.Sprintf("goi CLI version %s", config.CLI_VERSION)) // Success message in green
	},
}