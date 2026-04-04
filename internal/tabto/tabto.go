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

var (
	// ErrInconsistentColumnCount is returned when rows in the input text have
	// varying numbers of columns, indicating malformed tab-delimited data.
	ErrInconsistentColumnCount = errors.New("row has inconsistent column count")
	// ErrInternalConsistency is returned when an internal consistency error occurs
	// during table processing.
	ErrInternalConsistency = errors.New("internal consistency error")
)

// parseRow splits a tab-delimited line into cells and trims whitespace from each.
func parseRow(line string) []string {
	if line == "" {
		return []string{}
	}

	cells := strings.Split(line, "\t")

	row := make([]string, len(cells))
	for i, cell := range cells {
		row[i] = strings.TrimSpace(cell)
	}

	return row
}

// trimEmptyLines removes leading and trailing empty lines from lines.
func trimEmptyLines(lines []string) []string {
	start := 0
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}

	end := len(lines)
	for end > start && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}

	return lines[start:end]
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

// parseTableRows parses all rows after determining the expected column
// count. It validates each row and builds the complete table data
// structure.
func parseTableRows(lines []string, expectedCols int) ([][]string, error) {
	tableData := make([][]string, 0, len(lines))

	for index, line := range lines {
		row, err := validateAndParseRow(line, expectedCols, index)
		if err != nil {
			return nil, err
		}

		tableData = append(tableData, row)
	}

	return tableData, nil
}

// parseTable parses tab-delimited text into a two-dimensional slice
// of strings, representing the table's rows and columns. It trims
// leading/trailing empty lines and whitespace from cells. Empty input
// returns an empty table. It ensures that all non-empty rows have the same
// number of columns.
// Returns the parsed table or an error if the data is inconsistent.
func parseTable(text string) ([][]string, error) {
	lines := strings.Split(text, "\n")

	// Trim empty lines from beginning and end
	lines = trimEmptyLines(lines)

	if len(lines) == 0 {
		return [][]string{}, nil
	}

	// Determine expected column count from first valid row
	expectedCols := determineColumnCount(lines)

	// Parse and validate all rows
	tableData, err := parseTableRows(lines, expectedCols)
	if err != nil {
		return nil, err
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

	// Validate consistency (should be guaranteed by parseTable, but defensive)
	expectedCols := len(tableData[0])
	for i, row := range tableData {
		if len(row) != expectedCols {
			return "", fmt.Errorf(
				"%w: row %d has %d columns, expected %d",
				ErrInternalConsistency,
				i,
				len(row),
				expectedCols,
			)
		}
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
