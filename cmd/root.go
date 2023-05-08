/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	_ "embed"
	"os"

	"github.com/spf13/cobra"
)

const VERSION = "0.1.2"

const (
	enterpriseFlagName   = "enterprise"
	tokenFlagName        = "token"
	organizationFlagName = "organization"
	urlFlagName          = "url"
)

//go:embed banner.txt
var banner []byte

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ghas-org-enablement",
	Short: "CLI Tool to change GHAS organization settings",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	rootCmd.Version = VERSION
	rootCmd.Println("\n\n" + string(banner) + "\n\n")

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().String(enterpriseFlagName, "", "The slug of the enterprise.")
	rootCmd.MarkFlagRequired(enterpriseFlagName)

	rootCmd.PersistentFlags().String(tokenFlagName, "", "The access token.")
	rootCmd.MarkFlagRequired(tokenFlagName)

	rootCmd.PersistentFlags().String(organizationFlagName, "", "[Optional] filter for a single organization")
	rootCmd.PersistentFlags().String(urlFlagName, "https://api.github.com", "[Optional] URL of the GitHub Enterprise instance. Defaults to https://api.github.com")
}
