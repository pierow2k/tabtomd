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

// ErrAlignmentFailed is returned when rows in the input Markdown
// table have varying numbers of columns.
var ErrAlignmentFailed = errors.New("alignment failed")

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

// parseRow splits a single Markdown table row into its cells, with
// whitespace trimmed from each cell. It returns nil if the row is the
// header separator (which carries no data) or is not a valid table row
// (i.e., contains no pipe-delimited cells).
func parseRow(row string) []string {
	// strings.Split on a row containing at least one pipe yields at
	// least two parts: the empty segments before the first pipe and
	// after the last pipe. Fewer parts means the row has no pipes and
	// therefore no cells to parse.
	const minSplitParts = 2

	// Skip the header separator line, as it doesn't contain data.
	if isSeparatorRow(row) {
		return nil
	}

	// Split the row into columns based on the pipe delimiter. The first
	// and last elements are empty due to the leading and trailing
	// pipes, so we exclude them from the result.
	cells := strings.Split(row, "|")
	if len(cells) < minSplitParts {
		return nil
	}

	cells = cells[1 : len(cells)-1]

	// Trim whitespace from each cell.
	for i, cell := range cells {
		cells[i] = strings.TrimSpace(cell)
	}

	return cells
}

// updateColumnWidths widens the entries of columnWidths so that every
// cell in cells fits within its column. It uses rune counts to correctly
// handle multi-byte characters.
// Precondition: len(columnWidths) == len(cells), as guaranteed by the
// column-count check in parseMarkdownTable.
func updateColumnWidths(columnWidths []int, cells []string) {
	for i, cell := range cells {
		if width := utf8.RuneCountInString(cell); width > columnWidths[i] {
			columnWidths[i] = width
		}
	}
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
		cells := parseRow(row)
		if cells == nil {
			// Not a data row: either the header separator or a line that
			// doesn't look like a table row.
			continue
		}

		// The first valid data row (the header) defines the column layout;
		// every subsequent row must match it.
		if len(tableData) == 0 {
			columnWidths = make([]int, len(cells))
		} else if len(cells) != len(columnWidths) {
			return nil, nil, fmt.Errorf("%w: row %d has inconsistent column count",
				ErrAlignmentFailed, rowIndex+1)
		}

		updateColumnWidths(columnWidths, cells)
		tableData = append(tableData, cells)
	}

	return tableData, columnWidths, nil
}

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

// buildHeaderSeparator creates a Markdown header separator row using
// hyphens, sized according to the calculated column widths. It ensures
// that each separator has at least minSeparatorWidth hyphens to comply
// with the Markdown specification requiring a minimum of 3.
func buildHeaderSeparator(columnWidths []int) string {
	const minSeparatorWidth = 3

	separatorCells := make([]string, len(columnWidths))

	for index, width := range columnWidths {
		separatorWidth := max(width, minSeparatorWidth)
		separatorCells[index] = strings.Repeat("-", separatorWidth)
	}

	return "| " + strings.Join(separatorCells, " | ") + " |"
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
