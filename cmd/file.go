package cmd

import (
	"fmt"

	"github.com/pierow2k/tabtomd/internal/convert"
	"github.com/pierow2k/tabtomd/internal/fileops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flags for optional arguments
// var (
// 	prettyFlag         bool
// 	outputFilename string
// 	print          bool
// )

// fileCmd represents the file command. It configures the Cobra command
// to handle file input operations.
var fileCmd = &cobra.Command{
	Use:   "file [filename]",
	Short: "convert a tab-delimited file to a Markdown table",
	Long:  "The file command reads a tab-delimited file and converts its content into a Markdown formatted table.",
	Args:  cobra.ExactArgs(1),
	RunE:  func(cmd *cobra.Command, args []string) error { return processFile(args[0]) },
}

func init() {
	fileCmd.Flags().BoolVar(&printFlag, "print", false, "Print to screen")
	viper.BindPFlag("printFlag", fileCmd.Flags().Lookup("print"))

	fileCmd.Flags().BoolVar(&prettyFlag, "pretty", false, "Use pretty Markdown table formatting")
	viper.BindPFlag("prettyFlag", fileCmd.Flags().Lookup("pretty"))

	fileCmd.Flags().StringVar(&outputFilename, "output", "", "Specify the output file to save the Markdown table")
	viper.BindPFlag("output", fileCmd.Flags().Lookup("output"))

	// Add the fileCmd to the root command
	rootCmd.AddCommand(fileCmd)
}

func processFile(filename string) error {
	// Read the file contents
	content, err := fileops.ReadTSV(filename)
	if err != nil {
		return err
	}

	var markdownTable string

	// Convert tab-delimited text to Markdown table based on the pretty flag
	if prettyFlag {
		markdownTable, _ = convert.ToMDTablePretty(content)
	} else {
		markdownTable, _ = convert.ToMDTable(content)
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
