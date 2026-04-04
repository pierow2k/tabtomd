// White-box tests for unexported functions in the tabto package.
package tabto

import (
	"reflect"
	"testing"
)

// Test_parseRow tests the parseRow function with various input scenarios.
func Test_parseRow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		line string
		want []string
	}{
		{
			name: "normal row",
			line: "cell1\tcell2\tcell3",
			want: []string{"cell1", "cell2", "cell3"},
		},
		{
			name: "row with leading/trailing spaces in cells",
			line: "  cell1  \t  cell2  \t  cell3  ",
			want: []string{"cell1", "cell2", "cell3"},
		},
		{
			name: "row with empty cells",
			line: "cell1\t\tcell3",
			want: []string{"cell1", "", "cell3"},
		},
		{
			name: "row with only tabs",
			line: "\t\t",
			want: []string{"", "", ""},
		},
		{
			name: "empty line",
			line: "",
			want: nil,
		},
		{
			name: "single cell",
			line: "cell1",
			want: []string{"cell1"},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := parseRow(testCase.line); !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("parseRow() = %v, want %v", got, testCase.want)
			}
		})
	}
}

// Test_trimEmptyLines tests the trimEmptyLines function with various input scenarios.
//
//nolint:funlen
func Test_trimEmptyLines(t *testing.T) {
	t.Parallel()

	type args struct {
		lines []string
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "no empty lines",
			args: args{lines: []string{"line1", "line2", "line3"}},
			want: []string{"line1", "line2", "line3"},
		},
		{
			name: "leading empty lines",
			args: args{lines: []string{"", "", "line1", "line2"}},
			want: []string{"line1", "line2"},
		},
		{
			name: "trailing empty lines",
			args: args{lines: []string{"line1", "line2", "", ""}},
			want: []string{"line1", "line2"},
		},
		{
			name: "leading and trailing empty lines",
			args: args{lines: []string{"", "line1", "line2", ""}},
			want: []string{"line1", "line2"},
		},
		{
			name: "lines with only whitespace",
			args: args{lines: []string{"  ", "\t", "line1", "line2", " "}},
			want: []string{"line1", "line2"},
		},
		{
			name: "all empty lines",
			args: args{lines: []string{"", "  ", "\t"}},
			want: []string{},
		},
		{
			name: "empty input slice",
			args: args{lines: []string{}},
			want: []string{},
		},
		{
			name: "mixed content with internal empty lines",
			args: args{lines: []string{"", "line1", "", "line2", ""}},
			want: []string{"line1", "", "line2"},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := trimEmptyLines(testCase.args.lines); !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("trimEmptyLines() = %v, want %v", got, testCase.want)
			}
		})
	}
}

// Test_validateAndParseRow tests the validateAndParseRow function with various input scenarios.
//
//nolint:funlen
func Test_validateAndParseRow(t *testing.T) {
	t.Parallel()

	type args struct {
		line         string
		expectedCols int
		rowIndex     int
	}

	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "valid row",
			args:    args{line: "c1\tc2\tc3", expectedCols: 3, rowIndex: 0},
			want:    []string{"c1", "c2", "c3"},
			wantErr: false,
		},
		{
			name:    "inconsistent column count - too few",
			args:    args{line: "c1\tc2", expectedCols: 3, rowIndex: 1},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "inconsistent column count - too many",
			args:    args{line: "c1\tc2\tc3\tc4", expectedCols: 3, rowIndex: 2},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty line",
			args:    args{line: "", expectedCols: 3, rowIndex: 3},
			want:    []string{"", "", ""},
			wantErr: false,
		},
		{
			name:    "whitespace line",
			args:    args{line: "   ", expectedCols: 2, rowIndex: 4},
			want:    []string{"", ""},
			wantErr: false,
		},
		{
			name:    "row with extra whitespace in cells",
			args:    args{line: "  c1  \t  c2  ", expectedCols: 2, rowIndex: 5},
			want:    []string{"c1", "c2"},
			wantErr: false,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got, err := validateAndParseRow(testCase.args.line, testCase.args.expectedCols, testCase.args.rowIndex)
			if (err != nil) != testCase.wantErr {
				t.Errorf("validateAndParseRow() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("validateAndParseRow() = %v, want %v", got, testCase.want)
			}
		})
	}
}

// Test_determineColumnCount tests the determineColumnCount function.
func Test_determineColumnCount(t *testing.T) {
	t.Parallel()

	type args struct {
		lines []string
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "valid lines",
			args: args{lines: []string{"h1\th2", "r1\tr2"}},
			want: 2,
		},
		{
			name: "leading empty lines",
			args: args{lines: []string{"", "  ", "h1\th2\th3"}},
			want: 3,
		},
		{
			name: "no valid lines",
			args: args{lines: []string{"", "  ", "\t"}},
			want: 0,
		},
		{
			name: "empty input",
			args: args{lines: []string{}},
			want: 0,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := determineColumnCount(testCase.args.lines); got != testCase.want {
				t.Errorf("determineColumnCount() = %v, want %v", got, testCase.want)
			}
		})
	}
}

// Test_parseTable tests the parseTable function with various input scenarios.
//
//nolint:funlen
func Test_parseTable(t *testing.T) {
	t.Parallel()

	type args struct {
		text string
	}

	tests := []struct {
		name    string
		args    args
		want    [][]string
		wantErr bool
	}{
		{
			name: "valid table",
			args: args{text: "h1\th2\nr1\tr2"},
			want: [][]string{{"h1", "h2"}, {"r1", "r2"}},
		},
		{
			name: "empty input",
			args: args{text: ""},
			want: [][]string{},
		},
		{
			name: "input with only empty/whitespace lines",
			args: args{text: "\n  \n\t\n"},
			want: [][]string{},
		},
		{
			name:    "inconsistent columns",
			args:    args{text: "h1\th2\nr1"},
			wantErr: true,
		},
		{
			name: "leading and trailing empty lines",
			args: args{text: "\n\nh1\th2\nr1\tr2\n\n"},
			want: [][]string{{"h1", "h2"}, {"r1", "r2"}},
		},
		{
			name: "table with internal empty lines",
			args: args{text: "h1\th2\n\nr1\tr2"},
			want: [][]string{{"h1", "h2"}, {"", ""}, {"r1", "r2"}},
		},
		{
			name: "single row table",
			args: args{text: "h1\th2\th3"},
			want: [][]string{{"h1", "h2", "h3"}},
		},
		{
			name: "table with extra whitespace",
			args: args{text: "  h1 \t h2  \n r1\t r2 "},
			want: [][]string{{"h1", "h2"}, {"r1", "r2"}},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseTable(testCase.args.text)
			if (err != nil) != testCase.wantErr {
				t.Errorf("parseTable() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("parseTable() = %v, want %v", got, testCase.want)
			}
		})
	}
}
