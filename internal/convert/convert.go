// Package convert provides utilities for converting tab-delimited text
// into Markdown tables, with options for formatting column alignment.
package convert

import (
	"errors"
	"fmt"
	"strings"
)

// DefaultHeaderSeparatorWidth specifies the default number of hyphens used
// in the Markdown table header separator if column widths are not adjusted.
const DefaultHeaderSeparatorWidth = 5

// ErrInconsistentColumnCount is returned when rows in the input text have
// varying numbers of columns, indicating malformed tab-delimited data.
var ErrInconsistentColumnCount = errors.New("row has inconsistent column count")

// ToMDTablePretty converts tab-delimited text into a Markdown table with
// adjusted column widths for uniform alignment. The resulting table includes
// properly spaced vertical bars for better readability.
//
// The function returns an error if the input has inconsistent column counts.
func ToMDTablePretty(text string) (string, error) {
	table, columnWidths, err := parseAndAnalyze(text)
	if err != nil {
		return "", err
	}

	formattedRows := formatRows(table, columnWidths)
	headerSeparator := buildHeaderSeparator(columnWidths, 0)
	formattedRows = insertHeaderSeparator(formattedRows, headerSeparator)

	return strings.Join(formattedRows, "\n"), nil
}

// ToMDTable converts tab-delimited text into a Markdown table without
// adjusting column widths. The column alignment remains as-is, and the
// header separator uses a uniform width.
//
// The function returns an error if the input has inconsistent column counts.
func ToMDTable(text string) (string, error) {
	table, _, err := parseAndAnalyze(text)
	if err != nil {
		return "", err
	}

	formattedRows := formatRows(table, nil)
	if len(formattedRows) > 1 {
		headerSeparator := buildHeaderSeparator(make([]int, len(table[0])), DefaultHeaderSeparatorWidth)
		formattedRows = insertHeaderSeparator(formattedRows, headerSeparator)
	}

	return strings.Join(formattedRows, "\n"), nil
}

// parseAndAnalyze parses tab-delimited text into rows and columns, while
// calculating the maximum width of each column. It ensures that all rows
// have the same number of columns.
//
// Returns the parsed table, column widths, and an error if the data is inconsistent.
func parseAndAnalyze(text string) ([][]string, []int, error) {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) == 0 {
		return nil, nil, nil
	}

	tableData := make([][]string, len(lines))

	var columnWidths []int

	for lineIndex, line := range lines {
		tableData[lineIndex] = strings.Split(line, "\t")
		if lineIndex == 0 {
			columnWidths = make([]int, len(tableData[lineIndex]))
		} else if len(tableData[lineIndex]) != len(columnWidths) {
			return nil, nil, fmt.Errorf("%w: row %d", ErrInconsistentColumnCount, lineIndex+1)
		}

		for columnIndex, column := range tableData[lineIndex] {
			if len(column) > columnWidths[columnIndex] {
				columnWidths[columnIndex] = len(column)
			}
		}
	}

	return tableData, columnWidths, nil
}

// formatRows formats the rows of the Markdown table. When columnWidths are
// provided, each column is padded to the specified width for alignment.
// Otherwise, columns are left as-is.
//
// Returns a slice of formatted Markdown rows.
func formatRows(table [][]string, columnWidths []int) []string {
	formattedRows := make([]string, len(table))

	for rowIndex, row := range table {
		formattedRow := make([]string, len(row))

		for columnIndex, column := range row {
			if columnWidths != nil {
				formattedRow[columnIndex] = fmt.Sprintf("%-*s", columnWidths[columnIndex], column)
			} else {
				formattedRow[columnIndex] = column
			}
		}

		formattedRows[rowIndex] = "| " + strings.Join(formattedRow, " | ") + " |"
	}

	return formattedRows
}

// buildHeaderSeparator creates a Markdown header separator row using hyphens.
// If uniformWidth is greater than 0, all columns use the specified width
// instead of individual column widths.
//
// Returns the formatted header separator string.
func buildHeaderSeparator(columnWidths []int, uniformWidth int) string {
	var builder strings.Builder

	builder.WriteString("| ")

	for columnIndex := range columnWidths {
		width := columnWidths[columnIndex]
		if uniformWidth > 0 {
			width = uniformWidth
		}

		builder.WriteString(strings.Repeat("-", width))

		if columnIndex < len(columnWidths)-1 {
			builder.WriteString(" | ")
		}
	}

	builder.WriteString(" |")

	return builder.String()
}

// insertHeaderSeparator inserts the header separator row into the table
// after the header row.
//
// Returns the updated slice of table rows.
func insertHeaderSeparator(rows []string, separator string) []string {
	return append([]string{rows[0], separator}, rows[1:]...)
}
