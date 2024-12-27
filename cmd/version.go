package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the tabtomd version",
	Long:  "Print the tabtomd version number and build date",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tabtomd - Convert tab delimited data to Markdown tables")
		fmt.Println("Version:", viper.GetString("Version"), " - ", "Build Date:", viper.GetString("BuildDate"))
		fmt.Println("Copyright (c)", viper.GetString("CopyrightDate"), "Pierow2k")
		fmt.Println("Distributed under the MIT License")
	},
}
