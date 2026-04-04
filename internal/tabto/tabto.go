// Package tabto provides utilities for converting tab-delimited text
// into Markdown tables.
package tabto

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	defaultOutputBufferSize = 1024 // default initial buffer size for the output builder
	defaultMaxColumnWidth   = 50   // default maximum width for table columns
)

// ErrInconsistentColumnCount is returned when rows in the input text have
// varying numbers of columns, indicating malformed tab-delimited data.
var ErrInconsistentColumnCount = errors.New("row has inconsistent column count")

// parseRow splits a tab-delimited line into cells and trims whitespace
// from each.
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

// trimEmptyLines removes leading and trailing empty lines from lines.
func trimEmptyLines(lines []string) []string {
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}

	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	return lines
}

// determineColumnCount determines the expected column count from the first
// valid row.
// Returns 0 if no valid row is found.
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

// validateAndParseRow validates a single row against the expected column
// count and parses it. Empty rows are padded to match the expected column
// count.
func validateAndParseRow(
	line string,
	expectedCols int,
	rowIndex int,
) ([]string, error) {
	// Treat lines containing only whitespace as empty.
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

// parseTable parses tab-delimited text into a two-dimensional slice
// of strings, representing the table's rows and columns. It trims
// leading/trailing empty lines and whitespace from cells. Empty input
// returns an empty table. It ensures that all non-empty rows have the same
// number of columns.
// Returns the parsed table or an error if the data is inconsistent.
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

		// Skip empty rows for determining column count
		if len(row) == 0 || (len(row) == 1 && row[0] == "") {
			if expectedCols == 0 {
				continue
			}

			row = make([]string, expectedCols)
		}

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

// renderTable handles common table parsing and construction logic.
// The render callback performs format-specific rendering.
func renderTable(text string, render func(table.Writer)) (string, error) {
	tableData, err := parseTable(text)
	if err != nil {
		return "", fmt.Errorf("failed to parse table: %w", err)
	}

	if len(tableData) == 0 {
		return "", nil
	}

	var out strings.Builder
	out.Grow(defaultOutputBufferSize)

	tableWriter := table.NewWriter()
	tableWriter.SetOutputMirror(&out)

	headerRow := make(table.Row, len(tableData[0]))
	for i, cell := range tableData[0] {
		headerRow[i] = cell
	}

	tableWriter.AppendHeader(headerRow)

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

// Markdown converts tab-delimited text into a Markdown table with adjusted
// column widths for uniform alignment. The resulting table includes
// properly spaced vertical bars for better readability.
// Returns errors for parsing failures (inconsistent column counts).
func Markdown(text string) (string, error) {
	result, err := renderTable(text, func(tableWriter table.Writer) {
		cols := tableWriter.Length()

		configs := make([]table.ColumnConfig, cols)
		for i := range cols {
			configs[i] = table.ColumnConfig{
				Number:   i,
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

// HTML converts tab-delimited text into an HTML table.
// Returns errors for parsing failures (inconsistent column counts).
func HTML(text string) (string, error) {
	return renderTable(text, func(tw table.Writer) {
		tw.RenderHTML()
	})
}
