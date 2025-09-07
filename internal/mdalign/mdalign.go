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

// minSeparatorWidth is the minimum number of hyphens required for a Markdown
// table separator to be compliant with the specification.
const minSeparatorWidth = 3

// Align formats a Markdown table for uniform column alignment. It accepts a
// slice of strings, where each string is a row of the table, and returns a
// new slice with each column padded for readability.
//
// The function returns an error if the table has an inconsistent column count.
func Align(rows []string) ([]string, error) {
	table, columnWidths, err := parseMarkdownTable(rows)
	if err != nil {
		return nil, err
	}

	// If the table was empty or only contained a separator, return as-is.
	if len(table) == 0 {
		return rows, nil
	}

	formattedRows := formatRows(table, columnWidths)
	headerSeparator := buildHeaderSeparator(columnWidths)
	// The original table must have at least a header and one data row
	// (or just a header) to be considered valid for inserting a separator.
	if len(formattedRows) > 0 {
		formattedRows = insertHeaderSeparator(formattedRows, headerSeparator)
	}

	return formattedRows, nil
}

// isSeparatorRow checks if a given line from a Markdown table is the
// header separator (e.g., "|---|---|").
func isSeparatorRow(row string) bool {
	// Trim leading/trailing whitespace and pipe characters.
	trimmed := strings.Trim(row, " \t|")
	if trimmed == "" {
		return false
	}
	// If, after removing hyphens, nothing is left, it's a separator.
	// This correctly handles variations like ":---", "---:", and ":---:".
	return strings.Trim(trimmed, "-:| \t") == ""
}

// parseMarkdownTable parses rows of a Markdown table into a 2D slice of
// cells, while calculating the maximum width of each column. It ignores the
// header separator row during parsing.
//
// Returns the parsed table, column widths, and an error if the data is inconsistent.
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
		// trailing pipes, so we trim them off.
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
			return nil, nil, fmt.Errorf("%w: row %d", ErrInconsistentColumnCount, rowIndex+1)
		}

		// Iterate over each column to calculate its maximum width.
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

// formatRows formats the rows of the Markdown table. Each column is padded
// with spaces to match the calculated maximum width for that column.
//
// Returns a slice of formatted Markdown rows.
func formatRows(table [][]string, columnWidths []int) []string {
	formattedRows := make([]string, len(table))

	for rowIndex, row := range table {
		formattedRow := make([]string, len(row))

		for columnIndex, column := range row {
			padding := columnWidths[columnIndex] - utf8.RuneCountInString(column)
			formattedRow[columnIndex] = column + strings.Repeat(" ", padding)
		}

		formattedRows[rowIndex] = "| " + strings.Join(formattedRow, " | ") + " |"
	}

	return formattedRows
}

// buildHeaderSeparator creates a Markdown header separator row using hyphens,
// sized according to the calculated column widths.
//
// Returns the formatted header separator string.
func buildHeaderSeparator(columnWidths []int) string {
	separatorCells := make([]string, len(columnWidths))

	for index, width := range columnWidths {
		// Ensure the separator width is at least 3 hyphens to be compliant
		// with the Markdown specification (at least 3 hyphens).
		separatorWidth := max(width, minSeparatorWidth)

		separatorCells[index] = strings.Repeat("-", separatorWidth)
	}

	return "| " + strings.Join(separatorCells, " | ") + " |"
}

// insertHeaderSeparator inserts the header separator row into the table
// after the header row (which is the first row).
//
// Returns the updated slice of table rows.
func insertHeaderSeparator(rows []string, separator string) []string {
	// If there's only a header row, append the separator.
	if len(rows) == 1 {
		return append(rows, separator)
	}
	// Otherwise, insert it after the header.
	return append([]string{rows[0], separator}, rows[1:]...)
}
