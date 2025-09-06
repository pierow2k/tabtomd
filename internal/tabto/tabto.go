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

// parseTable parses tab-delimited text into a two-dimensional slice of strings,
// representing the table's rows and columns. It ensures that all rows
// have the same number of columns.
//
// Returns the parsed table and an error if the data is inconsistent.
func parseTable(text string) ([][]string, error) {
	// Split the input text into lines to handle them individually.
	lines := strings.Split(text, "\n")

	// Find the start and end of the actual content, trimming empty or whitespace-only lines
	// from the beginning and end of the input.
	start := 0
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}
	end := len(lines)
	for end > start && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}
	lines = lines[start:end]

	// If there are no content lines after trimming, return empty results.
	if len(lines) == 0 {
		return nil, nil
	}

	// Pre-allocate a slice to hold the parsed table data.
	tableData := make([][]string, len(lines))
	var columnCount int

	// Iterate over each line to parse it into columns.
	for i, line := range lines {
		row := strings.Split(line, "\t")

		// For the first row (header), establish the expected number of columns.
		if i == 0 {
			columnCount = len(row)
		} else if len(row) != columnCount {
			// For subsequent rows, ensure the column count is consistent.
			return nil, fmt.Errorf("%w: row %d", ErrInconsistentColumnCount, i+1)
		}
		tableData[i] = row
	}

	// Return the parsed table and no error.
	return tableData, nil
}

// Markdown converts tab-delimited text into a Markdown table with
// adjusted column widths for uniform alignment. The resulting table includes
// properly spaced vertical bars for better readability.
//
// This function uses the 'go-pretty' library to generate the Markdown table.
//
// The function returns an error if the input has inconsistent column counts.
func Markdown(text string) (string, error) {
	tableData, err := parseTable(text)
	if err != nil {
		return "", err
	}

	// If there's no data after parsing (e.g., empty input), return an empty string.
	if len(tableData) == 0 {
		return "", nil
	}

	var out strings.Builder
	t := table.NewWriter()
	t.SetOutputMirror(&out)

	// The first row of the parsed data is treated as the table header.
	header := make(table.Row, 0, len(tableData[0]))
	for _, h := range tableData[0] {
		header = append(header, h)
	}
	t.AppendHeader(header)

	// The remaining rows are added to the table body.
	if len(tableData) > 1 {
		for _, rowData := range tableData[1:] {
			row := make(table.Row, 0, len(rowData))
			for _, cell := range rowData {
				row = append(row, cell)
			}
			t.AppendRow(row)
		}
	}

	t.RenderMarkdown()

	// The go-pretty library adds a trailing newline, which we trim to ensure
	// consistent output behavior.
	return strings.TrimSuffix(out.String(), "\n"), nil
}
