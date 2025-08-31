// White-box tests for unexported functions in the config package.
package convert

import (
	"reflect"
	"testing"
)

// Test_parseAndAnalyze tests the parseAndAnalyze function to ensure it correctly
// parses tab-delimited text into a 2D slice and calculates column widths.
// It also checks for error handling with inconsistent column counts.
func Test_parseAndAnalyze(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    [][]string
		want1   []int
		wantErr bool
	}{
		{
			name:    "Empty input",
			args:    args{text: ""},
			want:    nil,
			want1:   nil,
			wantErr: false,
		},
		{
			name:    "Whitespace-only input",
			args:    args{text: "  \n\t  "},
			want:    nil,
			want1:   nil,
			wantErr: false,
		},
		{
			name: "Standard table",
			args: args{text: "Header 1\tHeader 2\nRow 1 Col 1\tRow 1 Col 2"},
			want: [][]string{
				{"Header 1", "Header 2"},
				{"Row 1 Col 1", "Row 1 Col 2"},
			},
			want1:   []int{11, 11},
			wantErr: false,
		},
		{
			name:    "Inconsistent column count",
			args:    args{text: "a\tb\nc"},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Table with multi-byte characters",
			args: args{text: "你好\tWorld\n世界\tHello"},
			want: [][]string{
				{"你好", "World"},
				{"世界", "Hello"},
			},
			want1:   []int{2, 5},
			wantErr: false,
		},
		{
			name: "Single row table",
			args: args{text: "a\tb\tc"},
			want: [][]string{
				{"a", "b", "c"},
			},
			want1:   []int{1, 1, 1},
			wantErr: false,
		},
		{
			name: "Table with leading and trailing whitespace",
			args: args{text: "\n  Header\tValue  \n  Row1\tVal1  \n"},
			want: [][]string{
				{"  Header", "Value  "},
				{"  Row1", "Val1  "},
			},
			want1:   []int{8, 7},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseAndAnalyze(tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAndAnalyze() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAndAnalyze() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseAndAnalyze() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
