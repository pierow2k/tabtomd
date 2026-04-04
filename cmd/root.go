// Package cmd provides the command-line interface for the tabtomd
// application. It defines the root command and its flags, handles
// input/output operations, and orchestrates the conversion of
// tab-delimited data to Markdown tables.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/pierow2k/tabtomd/internal/fileops"
	"github.com/pierow2k/tabtomd/internal/mdalign"
	"github.com/pierow2k/tabtomd/internal/tabto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// config holds the configuration options populated from command-line
// flags. It stores user preferences for input sources, output
// destinations, and formatting behavior.
type config struct {
	outputFilename string
	inputFileName  string
	pasteFlag      bool
	quietFlag      bool
	noPrettyFlag   bool
}

var (
	// cfg is the global configuration instance populated during command
	// initialization via flag bindings in init().
	cfg config

	// errFlagCombination is returned when mutually exclusive flags are
	// specified together (e.g., --output and --paste).
	errFlagCombination = errors.New("invalid flag combination")
)

// Build metadata variables. These values are overwritten by linker flags
// during the build process (see Makefile) to embed version information
// into the binary.
var (
	BuildDate     = "YYYY-MM-DDTHH:MM:SS-0000"
	CopyrightDate = "2026"
	License       = "Licensed under the MIT License <https://opensource.org/licenses/MIT>"
	Version       = "v0.0.0-dev"
)

// rootCmd is the base Cobra command for the tabtomd application.
// When invoked without subcommands, it performs the primary function of
// converting tab-delimited data to Markdown tables.
var rootCmd = &cobra.Command{
	Use:   "tabtomd",
	Short: "Convert tab delimited data to a Markdown table",
	Long: `tabtomd converts tab-delimited data to a Markdown table.

By default, it reads input from the system clipboard and prints the
resulting Markdown table to standard output. Alternative input sources
(files) and output destinations (files, clipboard) are available via flags.`,
	Example: `  # Convert clipboard content to aligned Markdown and print to stdout
  tabtomd

  # Read from a file and output aligned Markdown to stdout
  tabtomd -i data.tsv

  # Convert clipboard content and copy result back to clipboard
  tabtomd --paste

  # Read from file and write aligned Markdown to an output file
  tabtomd -i data.tsv -o table.md

  # Output compact (unaligned) Markdown table
  tabtomd --no-pretty`,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		if cfg.quietFlag {
			logrus.SetLevel(logrus.ErrorLevel)
		}
	},
	Version: fmt.Sprintf(
		"%s - built %s\nCopyright © %s Pierow2k\n%s",
		Version, BuildDate, CopyrightDate, License,
	),
	RunE: runFunction,
}

// Execute initializes the logger, configures default help and version
// flags, and runs the root command. If the command fails, Execute exits
// with status code 1.
func Execute() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})

	rootCmd.InitDefaultHelpFlag()
	rootCmd.Flags().Lookup("help").Usage = "Show help and usage information"
	rootCmd.InitDefaultVersionFlag()
	rootCmd.Flags().Lookup("version").Usage = "Show version, build details, and license"

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.Flags().BoolVarP(&cfg.pasteFlag, "paste", "p", false,
		"Paste results to the clipboard")
	rootCmd.Flags().BoolVarP(&cfg.noPrettyFlag, "no-pretty", "n", false,
		"Disable column alignment (output compact table)")
	rootCmd.Flags().BoolVarP(&cfg.quietFlag, "quiet", "q", false,
		"Suppress status messages")
	rootCmd.Flags().StringVarP(&cfg.outputFilename, "output", "o", "",
		"Write output to file")
	rootCmd.Flags().StringVarP(&cfg.inputFileName, "input", "i", "",
		"Read input from file")
}

// handleOutput processes and writes the generated Markdown table to the
// appropriate destination. It optionally aligns columns for readability,
// then outputs to one of: file (--output), clipboard (--paste), or stdout.
//
// The alignment step is skipped when --no-pretty is specified, producing
// a compact but unaligned table.
func handleOutput(markdownTable string) error {
	if !cfg.noPrettyFlag {
		alignedRows, err := mdalign.Align(strings.Split(markdownTable, "\n"))
		if err != nil {
			return fmt.Errorf("failed to align Markdown table: %w", err)
		}
		markdownTable = strings.Join(alignedRows, "\n")
	}

	switch {
	case cfg.outputFilename != "":
		if err := fileops.WriteMD(cfg.outputFilename, markdownTable); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
		logrus.Infof("Markdown table successfully written to %s\n", cfg.outputFilename)

	case cfg.pasteFlag:
		if err := clipboard.WriteAll(markdownTable); err != nil {
			return fmt.Errorf("failed to write Markdown to clipboard: %w", err)
		}
		logrus.Info("Markdown table copied to clipboard.")

	default:
		fmt.Println(markdownTable)
	}

	return nil
}

// runFunction is the main execution function for the root command.
// It reads tab-delimited data from either a file (--input) or the clipboard,
// converts it to a Markdown table, and delegates output handling to handleOutput.
//
// Returns an error if flag validation fails, input cannot be read, or
// conversion encounters an error.
func runFunction(_ *cobra.Command, _ []string) error {
	var (
		content string
		err     error
	)

	// Validate mutually exclusive output flags.
	if cfg.outputFilename != "" && cfg.pasteFlag {
		return fmt.Errorf("%w: cannot use --output and --paste together", errFlagCombination)
	}

	// Read input from file or clipboard.
	if cfg.inputFileName != "" {
		content, err = fileops.ReadTSV(cfg.inputFileName)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
	} else {
		content, err = clipboard.ReadAll()
		if err != nil {
			return fmt.Errorf("failed to read from clipboard: %w", err)
		}
	}

	// Convert tab-delimited content to Markdown.
	markdownTable, err := tabto.Markdown(content)
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	return handleOutput(markdownTable)
}
