// Package cmd handles conversion of tab-delimited text to
// a Markdown table.
package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/pierow2k/tabtomd/internal/fileops"
	"github.com/pierow2k/tabtomd/internal/tabto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flags for optional arguments
var (
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

	// Add the pasteCmd to the root command
	rootCmd.AddCommand(pasteCmd)
}

func pasteClipboard() error {
	text, err := clipboard.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read from clipboard: %w", err)
	}

	markdownTable, err := tabto.Markdown(text)
	if err != nil {
		return fmt.Errorf("failed to convert clipboard contents to Markdown: %w", err)
	}

	// If neither --print nor --output is set, default to copying the result
	// back onto the clipboard so you can paste the Markdown immediately.
	if outputFilename == "" && !viper.GetBool("printFlag") {
		if err := clipboard.WriteAll(markdownTable); err != nil {
			return fmt.Errorf("failed to write Markdown to clipboard: %w", err)
		}
		fmt.Println("Markdown table copied to clipboard.")
		return nil
	}

	if outputFilename != "" {
		if err := fileops.WriteMD(outputFilename, markdownTable); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
		fmt.Printf("Markdown table successfully written to %s\n", outputFilename)
	}

	if viper.GetBool("printFlag") {
		fmt.Println(markdownTable)
	}

	return nil
}
