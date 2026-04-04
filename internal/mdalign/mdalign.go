// Package mdalign provides utilities for aligning Markdown tables.
// It takes an existing, potentially unaligned, Markdown table and
// formats it for better readability by adjusting column widths.
package mdalign

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// ErrInconsistentColumnCount is returned when rows in the input Markdown
// table have varying numbers of columns.
var ErrInconsistentColumnCount = errors.New("row has inconsistent column count")

// formatRow formats a single table row with proper padding based on the
// provided column widths. Each cell is left-aligned and padded with spaces
// to match its corresponding column width.
func formatRow(cells []string, columnWidths []int) string {
	formattedCells := make([]string, len(cells))
	for i, cell := range cells {
		formattedCells[i] = fmt.Sprintf("%-*s", columnWidths[i], cell)
	}

	return "| " + strings.Join(formattedCells, " | ") + " |"
}

// Align formats a Markdown table for uniform column alignment. It accepts
// a slice of strings, where each string is a row of the table, and returns
// a new slice with each column padded for readability.
// It returns an error if the table has an inconsistent column count.
func Align(rows []string) ([]string, error) {
	table, columnWidths, err := parseMarkdownTable(rows)
	if err != nil {
		return nil, err
	}

	if len(table) == 0 {
		return rows, nil
	}

	result := make([]string, 0, len(table)+1)
	result = append(result, formatRow(table[0], columnWidths))

	result = append(result, buildHeaderSeparator(columnWidths))
	for _, row := range table[1:] {
		result = append(result, formatRow(row, columnWidths))
	}

	return result, nil
}

// isSeparatorRow checks if a given line from a Markdown table is the
// header separator (e.g., "|---|---|").
func isSeparatorRow(row string) bool {
	// Trim leading/trailing whitespace and pipe characters.
	trimmed := strings.Trim(row, " \t|")
	if trimmed == "" {
		return false
	}
	// If, after removing hyphens, colons, and whitespace, nothing remains,
	// the row is a separator. This correctly handles alignment indicators
	// like ":---", "---:", and ":---:".
	return strings.Trim(trimmed, "-:| \t") == ""
}

// parseMarkdownTable parses rows of a Markdown table into a 2D slice of
// cells, while calculating the maximum width of each column. It ignores
// header separator rows during parsing.
// It returns the parsed table data, column widths, and an error if the
// rows have inconsistent column counts.
func parseMarkdownTable(rows []string) ([][]string, []int, error) {
	var (
		columnWidths []int
		tableData    = make([][]string, 0, len(rows))
	)

	for rowIndex, row := range rows {
		// Skip the header separator line, as it doesn't contain data.
		if isSeparatorRow(row) {
			continue
		}

		// Split the row into cells based on the pipe delimiter.
		// The first and last elements will be empty due to the leading and
		// trailing pipes, so we exclude them from the result.
		cells := strings.Split(row, "|")
		if len(cells) > 1 {
			cells = cells[1 : len(cells)-1]
		} else {
			// Skip rows that don't appear to be valid table rows.
			continue
		}

		// Trim whitespace from each cell.
		for i, cell := range cells {
			cells[i] = strings.TrimSpace(cell)
		}

		// For the first valid data row (the header), initialize columnWidths.
		if len(tableData) == 0 {
			columnWidths = make([]int, len(cells))
		} else if len(cells) != len(columnWidths) {
			// Ensure subsequent rows have a consistent number of columns.
			return nil, nil, fmt.Errorf("%w: row %d",
				ErrInconsistentColumnCount, rowIndex+1)
		}

		// Update the maximum width for each column.
		for columnIndex, column := range cells {
			// Use RuneCountInString to correctly handle multi-byte characters.
			colWidth := utf8.RuneCountInString(column)
			if colWidth > columnWidths[columnIndex] {
				columnWidths[columnIndex] = colWidth
			}
		}

		tableData = append(tableData, cells)
	}

	return tableData, columnWidths, nil
}

// buildHeaderSeparator creates a Markdown header separator row using
// hyphens, sized according to the calculated column widths. It ensures
// that each separator has at least minSeparatorWidth hyphens to comply
// with the Markdown specification requiring a minimum of 3.
func buildHeaderSeparator(columnWidths []int) string {
	// minSeparatorWidth is the minimum width required for a Markdown
	// table separator to be compliant with the specification.
	const minSeparatorWidth = 3

	separatorCells := make([]string, len(columnWidths))

	for index, width := range columnWidths {
		separatorWidth := max(width, minSeparatorWidth)
		separatorCells[index] = strings.Repeat("-", separatorWidth)
	}

	return "| " + strings.Join(separatorCells, " | ") + " |"
}
