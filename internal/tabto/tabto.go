// Package tabto provides utilities for converting tab-delimited text
// into Markdown tables.
package tabto

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// ErrInconsistentColumnCount is returned when rows in the input text have
// varying numbers of columns, indicating malformed tab-delimited data.
var ErrInconsistentColumnCount = errors.New("row has inconsistent column count")

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

// parseTable parses tab-delimited text into a two-dimensional slice of strings,
// representing the table's rows and columns. It trims leading/trailing empty lines
// and whitespace from cells. Empty input returns an empty table. It ensures that
// all non-empty rows have the same number of columns.
//
// Returns the parsed table and an error if the data is inconsistent.
func parseTable(text string) ([][]string, error) {
	lines := strings.Split(text, "\n")

	// Find content boundaries
	start := 0
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}

	end := len(lines)
	for end > start && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}

	lines = lines[start:end]

	if len(lines) == 0 {
		return [][]string{}, nil
	}

	// Parse first non-empty row to determine column count
	var columnCount int

	for index, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}

		firstRow := parseRow(trimmedLine)
		if len(firstRow) == 0 {
			continue
		}

		columnCount = len(firstRow)
		lines[index] = trimmedLine // Update with trimmed version

		break
	}

	if columnCount == 0 {
		return [][]string{}, nil
	}

	// Parse all rows
	tableData := make([][]string, 0, len(lines))
	for index, line := range lines {
		row := parseRow(line)
		if len(row) == 0 {
			// Empty row - create with expected column count
			row = make([]string, columnCount)
		} else if len(row) != columnCount {
			return nil, fmt.Errorf("%w: row %d (expected %d columns, got %d): %q",
				ErrInconsistentColumnCount, index+1, columnCount, len(row), strings.TrimSpace(line))
		}

		tableData = append(tableData, row)
	}

	return tableData, nil
}

// Markdown converts tab-delimited text into a Markdown table with
// adjusted column widths for uniform alignment. The resulting table includes
// properly spaced vertical bars for better readability.
//
// This function uses the 'go-pretty' library to generate the Markdown table.
// The library handles pipe escaping automatically. No additional escaping
// is needed since cells from tab-delimited data won't contain newlines.
//
// Returns an empty string (not an error) for empty input. Errors occur only
// for parsing failures (inconsistent column counts).
func Markdown(text string) (string, error) {
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
			return "", fmt.Errorf("internal consistency error: row %d has %d columns, expected %d",
				i, len(row), expectedCols)
		}
	}

	var out strings.Builder
	out.Grow(1024) // Pre-allocate capacity for better performance

	t := table.NewWriter()
	t.SetOutputMirror(&out)

	// Add header - no escaping needed
	headerRow := make(table.Row, len(tableData[0]))
	for index, cell := range tableData[0] {
		headerRow[index] = cell // Already trimmed by parseRow
	}

	t.AppendHeader(headerRow)

	// Add body rows if they exist
	if len(tableData) > 1 {
		for _, rowData := range tableData[1:] {
			row := make(table.Row, len(rowData))
			for i, cell := range rowData {
				row[i] = cell // Already trimmed by parseRow
			}

			t.AppendRow(row)
		}
	}

	// Optional: Configure column widths for better readability
	columnConfigs := make([]table.ColumnConfig, expectedCols)
	for i := range expectedCols {
		columnConfigs[i] = table.ColumnConfig{
			Number:   i,
			WidthMax: 50, // Adjust based on your needs
		}
	}

	t.SetColumnConfigs(columnConfigs)

	t.RenderMarkdown()

	// Trim trailing newline for consistent output
	result := out.String()

	return strings.TrimSuffix(result, "\n"), nil
}
