// Package cmd implements the command-line interface for tabtomd, including
// the 'paste' command for converting clipboard data.
package cmd

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/pierow2k/tabtomd/internal/fileops"
	"github.com/pierow2k/tabtomd/internal/mdalign"
	"github.com/pierow2k/tabtomd/internal/tabto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flags for optional arguments.
var (
	outputFilename string
	prettyFlag     bool
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

//nolint:errcheck,gosec
func init() {
	pasteCmd.Flags().BoolVar(&prettyFlag, "pretty", false, "Align columns for better readability")
	viper.BindPFlag("prettyFlag", pasteCmd.Flags().Lookup("pretty"))

	pasteCmd.Flags().BoolVar(&printFlag, "print", false, "Print to screen")
	viper.BindPFlag("printFlag", pasteCmd.Flags().Lookup("print"))

	pasteCmd.Flags().StringVar(&outputFilename, "output", "", "Specify the output file to save the Markdown table")
	viper.BindPFlag("output", pasteCmd.Flags().Lookup("output"))

	// Add the pasteCmd to the root command
	rootCmd.AddCommand(pasteCmd)
}

// convertClipboardToMarkdown reads from the clipboard and converts the
// tab-delimited text to a Markdown table.
func convertClipboardToMarkdown() (string, error) {
	text, err := clipboard.ReadAll()
	if err != nil {
		return "", fmt.Errorf("failed to read from clipboard: %w", err)
	}

	markdownTable, err := tabto.Markdown(text)
	if err != nil {
		return "", fmt.Errorf("failed to convert clipboard contents to Markdown: %w", err)
	}

	return markdownTable, nil
}

// formatMarkdownTable applies pretty formatting to the table if requested.
func formatMarkdownTable(markdownTable string) (string, error) {
	if viper.GetBool("prettyFlag") {
		// Align the columns of the Markdown table for better readability.
		alignedRows, err := mdalign.Align(strings.Split(markdownTable, "\n"))
		if err != nil {
			return "", fmt.Errorf("failed to align Markdown table: %w", err)
		}

		return strings.Join(alignedRows, "\n"), nil
	}

	return markdownTable, nil
}

// handleOutput determines where to send the final Markdown table based on flags.
func handleOutput(markdownTable string) error {
	// Default behavior: copy to clipboard if no other output is specified.
	if outputFilename == "" && !viper.GetBool("printFlag") {
		if err := clipboard.WriteAll(markdownTable); err != nil {
			return fmt.Errorf("failed to write Markdown to clipboard: %w", err)
		}

		fmt.Println("Markdown table copied to clipboard.")

		return nil
	}

	// Output to a file if the --output flag is set.
	if outputFilename != "" {
		if err := fileops.WriteMD(outputFilename, markdownTable); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}

		fmt.Printf("Markdown table successfully written to %s\n", outputFilename)
	}

	// Print to stdout if the --print flag is set.
	if viper.GetBool("printFlag") {
		fmt.Println(markdownTable)
	}

	return nil
}

// pasteClipboard orchestrates the primary logic for the 'paste' command.
// It reads tab-delimited text from the system clipboard, converts it to a
// Markdown table, and then handles output based on user-provided flags.
//
// The output can be:
// 1. Written back to the clipboard (default behavior if no other output flags are set).
// 2. Printed to standard output (`--print` flag).
// 3. Saved to a file (`--output` flag).
func pasteClipboard() error {
	markdownTable, err := convertClipboardToMarkdown()
	if err != nil {
		return err
	}

	formattedTable, err := formatMarkdownTable(markdownTable)
	if err != nil {
		return err
	}

	return handleOutput(formattedTable)
}
