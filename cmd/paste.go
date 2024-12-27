// Package cmd handles conversion of tab-delimited text to
// a Markdown table.
package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/pierow2k/tabtomd/internal/convert"
	"github.com/pierow2k/tabtomd/internal/fileops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flags for optional arguments
var (
	prettyFlag     bool
	outputFilename string
	printFlag      bool
)

// pasteCmd represents the paste command. It configures the Cobra command
// to handle paste operations.
var pasteCmd = &cobra.Command{
	Use:   "paste",
	Short: "paste tab delimited data from the clipboard",
	Long:  "The paste command converts tab delimited data from the clipboard into a Markdown formatted table.",
	RunE:  func(_ *cobra.Command, _ []string) error { return pasteClipboard() },
}

func init() {
	pasteCmd.Flags().BoolVar(&printFlag, "print", false, "Print to screen")
	viper.BindPFlag("printFlag", pasteCmd.Flags().Lookup("print"))

	pasteCmd.Flags().StringVar(&outputFilename, "output", "", "Specify the output file to save the Markdown table")
	viper.BindPFlag("output", pasteCmd.Flags().Lookup("output"))

	pasteCmd.Flags().BoolVar(&prettyFlag, "pretty", false, "Use pretty Markdown table formatting")
	viper.BindPFlag("prettyFlag", pasteCmd.Flags().Lookup("pretty"))

	// Add the pasteCmd to the root command
	rootCmd.AddCommand(pasteCmd)
}

func pasteClipboard() error {
	text, err := clipboard.ReadAll() // Read text from clipboard
	if err != nil {
		return fmt.Errorf("failed to read from clipboard: %w", err)
	}

	var markdownTable string

	// Convert tab-delimited text to Markdown table
	if prettyFlag {
		markdownTable, _ = convert.ToMDTablePretty(text)
	} else {
		markdownTable, _ = convert.ToMDTable(text)
	}

	// Check if the --output flag was provided
	if outputFilename != "" {
		err = fileops.WriteMD(outputFilename, markdownTable)
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
		fmt.Printf("Markdown table successfully written to %s\n", outputFilename)
	}

	if viper.GetBool("printFlag") {
		// Print the Markdown table to stdout if printOutput is enabled
		fmt.Println(markdownTable)
	}

	return nil
}
