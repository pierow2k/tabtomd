// White-box tests for unexported functions in the mdalign package.
package mdalign

import (
	"reflect"
	"testing"
)

// Test_isSeparatorRow tests the isSeparatorRow function with various input scenarios.
func Test_isSeparatorRow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		row  string
		want bool
	}{
		{name: "simple separator", row: "|---|---|", want: true},
		{name: "separator with alignment", row: "|:--|--:|", want: true},
		{name: "separator with center alignment", row: "|:---:|:---:|", want: true},
		{name: "separator with extra whitespace", row: " | --- | --- | ", want: true},
		{name: "separator with single column", row: "| --- |", want: true},
		{name: "separator with short hyphens", row: "|-|-|", want: true},
		{name: "header row", row: "| Header 1 | Header 2 |", want: false},
		{name: "data row", row: "| Data 1 | Data 2 |", want: false},
		{name: "row with numbers", row: "| 1 | 2 |", want: false},
		{name: "empty row", row: "", want: false},
		{name: "whitespace row", row: "   ", want: false},
		{name: "row with only pipes", row: "|||", want: false},
		{name: "row with mixed characters", row: "|a--|--b|", want: false},
		{name: "separator without outer pipes", row: "--- | ---", want: true},
	}
	for _, testTable := range tests {
		t.Run(testTable.name, func(t *testing.T) {
			t.Parallel()

			got := isSeparatorRow(testTable.row)
			if got != testTable.want {
				t.Errorf("isSeparatorRow() = %v, want %v", got, testTable.want)
			}
		})
	}
}

// Test_parseMarkdownTable tests the parseMarkdownTable function with various input scenarios.
//
//nolint:gosmopolitan,funlen
func Test_parseMarkdownTable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		rows             []string
		wantTable        [][]string
		wantColumnWidths []int
		wantErr          bool
	}{
		{
			name:             "standard table",
			rows:             []string{"| h1 | h2 |", "|---|---|", "| d1 | d22 |"},
			wantTable:        [][]string{{"h1", "h2"}, {"d1", "d22"}},
			wantColumnWidths: []int{2, 3},
			wantErr:          false,
		},
		{
			name:             "inconsistent column count",
			rows:             []string{"| h1 | h2 |", "|---|---|", "| d1 |"},
			wantTable:        nil,
			wantColumnWidths: nil,
			wantErr:          true,
		},
		{
			name:             "empty input",
			rows:             []string{},
			wantTable:        [][]string{},
			wantColumnWidths: nil,
			wantErr:          false,
		},
		{
			name:             "table with only header and separator",
			rows:             []string{"| Header |", "| --- |"},
			wantTable:        [][]string{{"Header"}},
			wantColumnWidths: []int{6},
			wantErr:          false,
		},
		{
			name:             "table with multi-byte characters",
			rows:             []string{"| 你好 | a |", "|---|---|", "| 世界 | b |"},
			wantTable:        [][]string{{"你好", "a"}, {"世界", "b"}},
			wantColumnWidths: []int{2, 1},
			wantErr:          false,
		},
		{
			name:             "malformed row (no pipes)",
			rows:             []string{"h1", "h2"},
			wantTable:        [][]string{},
			wantColumnWidths: nil,
			wantErr:          false,
		},
	}
	for _, testTable := range tests {
		t.Run(testTable.name, func(t *testing.T) {
			t.Parallel()

			gotTable, gotColumnWidths, err := parseMarkdownTable(testTable.rows)
			if (err != nil) != testTable.wantErr {
				t.Errorf("parseMarkdownTable() error = %v, wantErr %v", err, testTable.wantErr)

				return
			}

			if !reflect.DeepEqual(gotTable, testTable.wantTable) {
				t.Errorf("parseMarkdownTable() gotTable = %v, want %v", gotTable, testTable.wantTable)
			}

			if !reflect.DeepEqual(gotColumnWidths, testTable.wantColumnWidths) {
				t.Errorf("parseMarkdownTable() gotColumnWidths = %v, want %v", gotColumnWidths, testTable.wantColumnWidths)
			}
		})
	}
}

// Test_formatRows tests the formatRows function with various input scenarios.
//
//nolint:gosmopolitan
func Test_formatRows(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		table        [][]string
		columnWidths []int
		want         []string
	}{
		{
			name:         "standard table",
			table:        [][]string{{"h1", "h2"}, {"d1", "d2"}},
			columnWidths: []int{2, 2},
			want:         []string{"| h1 | h2 |", "| d1 | d2 |"},
		},
		{
			name:         "table with padding",
			table:        [][]string{{"a", "b"}, {"long", "longer"}},
			columnWidths: []int{4, 6},
			want:         []string{"| a    | b      |", "| long | longer |"},
		},
		{
			name:         "table with multi-byte characters",
			table:        [][]string{{"你好", "世界"}, {"a", "b"}},
			columnWidths: []int{2, 2},
			want:         []string{"| 你好 | 世界 |", "| a  | b  |"},
		},
		{
			name:         "single row table",
			table:        [][]string{{"header1", "header2"}},
			columnWidths: []int{7, 7},
			want:         []string{"| header1 | header2 |"},
		},
		{
			name:         "empty table",
			table:        [][]string{},
			columnWidths: []int{},
			want:         []string{},
		},
	}
	for _, testTable := range tests {
		t.Run(testTable.name, func(t *testing.T) {
			t.Parallel()

			got := formatRows(testTable.table, testTable.columnWidths)
			if !reflect.DeepEqual(got, testTable.want) {
				t.Errorf("formatRows() = %v, want %v", got, testTable.want)
			}
		})
	}
}

// Test_buildHeaderSeparator tests the buildHeaderSeparator function with various column width scenarios.
func Test_buildHeaderSeparator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		columnWidths []int
		want         string
	}{
		{
			name:         "standard widths",
			columnWidths: []int{5, 6, 4},
			want:         "| ----- | ------ | ---- |",
		},
		{
			name:         "widths below minimum",
			columnWidths: []int{1, 2, 1},
			want:         "| --- | --- | --- |",
		},
		{
			name:         "mixed widths",
			columnWidths: []int{5, 2, 8},
			want:         "| ----- | --- | -------- |",
		},
		{
			name:         "single column",
			columnWidths: []int{7},
			want:         "| ------- |",
		},
		{
			name:         "empty input",
			columnWidths: []int{},
			want:         "|  |",
		},
	}
	for _, testTable := range tests {
		t.Run(testTable.name, func(t *testing.T) {
			t.Parallel()

			got := buildHeaderSeparator(testTable.columnWidths)
			if got != testTable.want {
				t.Errorf("buildHeaderSeparator() = %v, want %v", got, testTable.want)
			}
		})
	}
}

// Test_insertHeaderSeparator tests the insertHeaderSeparator function with various input scenarios.
func Test_insertHeaderSeparator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		rows      []string
		separator string
		want      []string
	}{
		{
			name:      "insert into multi-row table",
			rows:      []string{"| header1 | header2 |", "| data1   | data2   |"},
			separator: "|---------|---------|",
			want:      []string{"| header1 | header2 |", "|---------|---------|", "| data1   | data2   |"},
		},
		{
			name:      "append to single-row table (header only)",
			rows:      []string{"| header1 | header2 |"},
			separator: "|---------|---------|",
			want:      []string{"| header1 | header2 |", "|---------|---------|"},
		},
		{
			name:      "empty input rows",
			rows:      []string{},
			separator: "|---|",
			want:      []string{},
		},
		{
			name:      "table with three rows",
			rows:      []string{"| h1 |", "| r1 |", "| r2 |"},
			separator: "|----|",
			want:      []string{"| h1 |", "|----|", "| r1 |", "| r2 |"},
		},
	}
	for _, testTable := range tests {
		t.Run(testTable.name, func(t *testing.T) {
			t.Parallel()

			got := insertHeaderSeparator(testTable.rows, testTable.separator)
			// Use reflect.DeepEqual for slice comparison.
			if !reflect.DeepEqual(got, testTable.want) {
				t.Errorf("insertHeaderSeparator() = %v, want %v", got, testTable.want)
			}
		})
	}
}
