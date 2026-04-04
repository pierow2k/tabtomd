// Package tabto provides utilities for converting tab-delimited text into
// Markdown and HTML tables. It handles parsing of tab-separated values,
// validation of consistent column counts across rows, and rendering into
// formatted table output.
//
// The package supports empty rows (which are padded with empty cells) and
// automatically trims whitespace from cell boundaries.
package tabto

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	// defaultOutputBufferSize is the initial capacity allocated for the
	// strings.Builder used when rendering table output. This value reduces
	// memory allocations for typical table sizes.
	defaultOutputBufferSize = 1024

	// defaultMaxColumnWidth is the maximum width applied to table columns
	// when rendering Markdown output. This prevents excessively wide columns
	// while allowing content to wrap naturally.
	defaultMaxColumnWidth = 50
)

// ErrInconsistentColumnCount is returned when rows in the input text have
// varying numbers of columns, indicating malformed tab-delimited data.
// The error message includes the row number and expected vs. actual column
// counts for debugging.
var ErrInconsistentColumnCount = errors.New("row has inconsistent column count")

// parseRow splits a tab-delimited line into individual cells and trims
// leading and trailing whitespace from each cell. If the input line is
// empty, parseRow returns nil.
func parseRow(line string) []string {
	if line == "" {
		return nil
	}

	cells := strings.Split(line, "\t")
	for i, cell := range cells {
		cells[i] = strings.TrimSpace(cell)
	}

	return cells
}

// trimEmptyLines removes leading and trailing lines that contain only
// whitespace from the input slice. It returns a new slice that may share
// storage with the original.
func trimEmptyLines(lines []string) []string {
	// Remove leading empty lines.
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}

	// Remove trailing empty lines.
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	return lines
}

// determineColumnCount determines the expected column count by examining
// the first non-empty row in the input lines. It returns 0 if no non-empty
// row is found.
func determineColumnCount(lines []string) int {
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}

		firstRow := parseRow(trimmedLine)

		return len(firstRow)
	}

	return 0
}

// validateAndParseRow validates that a line has the expected number of
// columns and parses it into a slice of cell strings. Lines containing
// only whitespace are treated as empty and padded with empty strings to
// match expectedCols.
//
// Parameters:
//   - line: the tab-delimited line to parse
//   - expectedCols: the required number of columns
//   - rowIndex: the 0-based row index, used for error reporting
//
// Returns the parsed row or an error wrapping ErrInconsistentColumnCount.
func validateAndParseRow(
	line string,
	expectedCols int,
	rowIndex int,
) ([]string, error) {
	// Treat lines containing only whitespace as empty rows.
	if strings.TrimSpace(line) == "" {
		return make([]string, expectedCols), nil
	}

	row := parseRow(line)

	if len(row) != expectedCols {
		return nil, fmt.Errorf(
			"%w: row %d (expected %d columns, got %d): %q",
			ErrInconsistentColumnCount,
			rowIndex+1,
			expectedCols,
			len(row),
			strings.TrimSpace(line),
		)
	}

	return row, nil
}

// parseTable parses tab-delimited text into a two-dimensional slice of
// strings representing table rows and columns. It trims leading and
// trailing empty lines from the input and removes whitespace from cell
// boundaries.
//
// Empty input returns an empty (non-nil) slice. All non-empty rows must
// have the same number of columns, or an error wrapping
// ErrInconsistentColumnCount is returned. Empty rows within the table
// are padded with empty strings to match the column count.
func parseTable(text string) ([][]string, error) {
	lines := strings.Split(text, "\n")
	lines = trimEmptyLines(lines)

	if len(lines) == 0 {
		return [][]string{}, nil
	}

	tableData := make([][]string, 0, len(lines))

	var expectedCols int

	for lineIndex, line := range lines {
		row := parseRow(line)

		// Handle empty rows: skip when determining column count,
		// otherwise pad with empty cells to maintain consistency.
		if len(row) == 0 || (len(row) == 1 && row[0] == "") {
			if expectedCols == 0 {
				continue
			}

			row = make([]string, expectedCols)
		}

		// Validate column consistency against the first non-empty row.
		if expectedCols == 0 {
			expectedCols = len(row)
		} else if len(row) != expectedCols {
			return nil, fmt.Errorf(
				"%w: row %d (expected %d columns, got %d)",
				ErrInconsistentColumnCount, lineIndex+1, expectedCols, len(row),
			)
		}

		tableData = append(tableData, row)
	}

	return tableData, nil
}

// renderTable parses tab-delimited text and renders it using the provided
// render callback. The first row is treated as the table header.
//
// The render callback receives a configured table.Writer and is responsible
// for calling the appropriate render method (e.g., RenderMarkdown or
// RenderHTML) and applying any format-specific configuration.
func renderTable(text string, render func(table.Writer)) (string, error) {
	tableData, err := parseTable(text)
	if err != nil {
		return "", fmt.Errorf("failed to parse table: %w", err)
	}

	if len(tableData) == 0 {
		return "", nil
	}

	// Pre-allocate buffer for output to reduce memory allocations.
	var out strings.Builder
	out.Grow(defaultOutputBufferSize)

	tableWriter := table.NewWriter()
	tableWriter.SetOutputMirror(&out)

	// Build the header row from the first parsed row.
	headerRow := make(table.Row, len(tableData[0]))
	for i, cell := range tableData[0] {
		headerRow[i] = cell
	}

	tableWriter.AppendHeader(headerRow)

	// Append remaining rows as table body.
	for _, rowData := range tableData[1:] {
		row := make(table.Row, len(rowData))
		for i, cell := range rowData {
			row[i] = cell
		}

		tableWriter.AppendRow(row)
	}

	render(tableWriter)

	return out.String(), nil
}

// Markdown converts tab-delimited text into a Markdown table. Column
// widths are capped at defaultMaxColumnWidth for uniform alignment. The
// resulting table uses properly spaced vertical bars for readability.
//
// If parsing fails due to inconsistent column counts, an error wrapping
// ErrInconsistentColumnCount is returned. The returned string does not
// include a trailing newline.
func Markdown(text string) (string, error) {
	result, err := renderTable(text, func(tableWriter table.Writer) {
		cols := tableWriter.Length()

		// Configure maximum width for each column to prevent
		// excessively wide tables.
		configs := make([]table.ColumnConfig, cols)
		for i := range cols {
			configs[i] = table.ColumnConfig{
				Number:   i + 1, // Column numbers are 1-indexed
				WidthMax: defaultMaxColumnWidth,
			}
		}

		tableWriter.SetColumnConfigs(configs)
		tableWriter.RenderMarkdown()
	})
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(result, "\n"), nil
}

// HTML converts tab-delimited text into an HTML table. If parsing fails
// due to inconsistent column counts, an error wrapping
// ErrInconsistentColumnCount is returned.
func HTML(text string) (string, error) {
	return renderTable(text, func(tw table.Writer) {
		tw.RenderHTML()
	})
}
