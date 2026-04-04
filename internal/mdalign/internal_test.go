// White-box tests for unexported functions in the mdalign package.
//
//nolint:gosmopolitan,funlen
package mdalign

import (
	"reflect"
	"testing"
)

// Test_isSeparatorRow tests the isSeparatorRow function with various input
// scenarios.
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

// Test_parseMarkdownTable tests the parseMarkdownTable function with
// various input scenarios.
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

// Test_buildHeaderSeparator tests the buildHeaderSeparator function with
// various column width scenarios.
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

func Test_formatRow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		cells        []string
		columnWidths []int
		want         string
	}{
		{
			name:         "standard row with padding",
			cells:        []string{"Alice", "25"},
			columnWidths: []int{10, 5},
			want:         "| Alice      | 25    |",
		},
		{
			name:         "exact width",
			cells:        []string{"Bob", "300"},
			columnWidths: []int{3, 3},
			want:         "| Bob | 300 |",
		},
		{
			name:         "empty cells",
			cells:        []string{"", "Data"},
			columnWidths: []int{3, 4},
			want:         "|     | Data |",
		},
		{
			name:         "multi-byte characters",
			cells:        []string{"你好", "a"},
			columnWidths: []int{2, 1},
			want:         "| 你好 | a |",
		},
		{
			name:         "single column",
			cells:        []string{"Header"},
			columnWidths: []int{8},
			want:         "| Header   |",
		},
		{
			name:         "empty row",
			cells:        []string{},
			columnWidths: []int{},
			want:         "|  |",
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := formatRow(testCase.cells, testCase.columnWidths); got != testCase.want {
				t.Errorf("formatRow() = %v, want %v", got, testCase.want)
			}
		})
	}
}
