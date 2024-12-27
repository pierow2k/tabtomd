// Package cmd handles the command line interface for tabtomd.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Build date is rewritten at build time by the make file.
	BuildDate = "YYYY-MM-DDTHH:MM:SSZ"
	// Version is rewritten at build time by the make file.
	Version = "X.X.X"
	// CopyrightDate is rewritten at build time by the make file.
	CopyrightDate = "YYYY"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "tabtomd",
	Short: "Convert tab delimited data to a Markdown table",
	Long:  "Convert tab delimited data to a Markdown table",
}

// Execute adds all child commands to the root command and sets flags
// appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	viper.Set("BuildDate", BuildDate)
	viper.Set("Version", Version)
	viper.Set("CopyrightDate", CopyrightDate)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // Read in environment variables that match.
}
