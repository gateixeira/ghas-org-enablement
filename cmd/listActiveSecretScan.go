/*
Package cmd provides a command-line interface for changing GHAS settings for a given organization.
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var listActiveSecretScanCmd = &cobra.Command{
	Use:   "list-secret-scan",
	Short: "List repositories that have secret scanning activated",
	Long: `If the enterprise slug is provided, this tool runs for all organizations in an enterprise.
	Provide only the organization slug if you want to do it against a single organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		enterprise, _ := cmd.Flags().GetString("enterprise")
		token, _ := cmd.Flags().GetString("token")
		organization, _ := cmd.Flags().GetString("organization")
		url, _ := cmd.Flags().GetString("url")

		log.Println("Listing repositories with secret scanning activated")

		err := ListSecretScanning(enterprise, organization, token, url, false)

		if err != nil {
			log.Println("Error deactivating GHAS features")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listActiveSecretScanCmd)
}
