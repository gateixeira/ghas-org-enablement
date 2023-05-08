/*
Package cmd provides a command-line interface for changing GHAS settings for a given organization.
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var activateGhasFeaturesCmd = &cobra.Command{
	Use:   "activate-ghas-features",
	Short: "Activate GitHub Advanced Security features for all organizations in an enterprise",
	Long: `If the enterprise slug is provided, this tool activates GitHub Advanced Security features for all organizations in an enterprise.
	Provide only the organization slug if you want to enable it for a single organization.
	
	Advanced Security features will also be enabled for new repositories.`,
	Run: func(cmd *cobra.Command, args []string) {
		enterprise, _ := cmd.Flags().GetString("enterprise")
		token, _ := cmd.Flags().GetString("token")
		organization, _ := cmd.Flags().GetString("organization")
		url, _ := cmd.Flags().GetString("url")

		log.Println("Activating GHAS features")

		err := ManageGhasFeatures(enterprise, organization, token, url, true)

		if err != nil {
			log.Println("Error activating GHAS features")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(activateGhasFeaturesCmd)
}
