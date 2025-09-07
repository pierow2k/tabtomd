// Package mdalign_test provides black-box tests and runnable examples for
// the public API of the mdalign package.
package mdalign_test

import (
	"reflect"
	"testing"

	"github.com/pierow2k/tabtomd/internal/mdalign"
)

// TestAlign tests the Align function with various input scenarios.
//
//nolint:funlen,gosmopolitan
func TestAlign(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		rows    []string
		want    []string
		wantErr bool
	}{
		{
			name: "standard table alignment",
			rows: []string{
				"| Name | Age |",
				"|---|---|",
				"| Alice | 25 |",
				"| Bob | 300 |",
			},
			want: []string{
				"| Name  | Age |",
				"| ----- | --- |",
				"| Alice | 25  |",
				"| Bob   | 300 |",
			},
			wantErr: false,
		},
		{
			name:    "inconsistent column count",
			rows:    []string{"| h1 | h2 |", "|---|---|", "| d1 |"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty input",
			rows:    []string{},
			want:    []string{},
			wantErr: false,
		},
		{
			name: "table with only header",
			rows: []string{"| Header1 | Header2 |"},
			want: []string{
				"| Header1 | Header2 |",
				"| ------- | ------- |",
			},
			wantErr: false,
		},
		{
			name: "table with multi-byte characters",
			rows: []string{"| 你好 | a |", "|---|---|", "| 世界 | b |"},
			want: []string{
				"| 你好 | a |",
				"| --- | --- |",
				"| 世界 | b |",
			},
			wantErr: false,
		},
	}
	for _, testTable := range tests {
		t.Run(testTable.name, func(t *testing.T) {
			t.Parallel()

			got, err := mdalign.Align(testTable.rows)
			if (err != nil) != testTable.wantErr {
				t.Errorf("Align() error = %v, wantErr %v", err, testTable.wantErr)

				return
			}

			if !reflect.DeepEqual(got, testTable.want) {
				t.Errorf("Align() = %v, want %v", got, testTable.want)
			}
		})
	}
}
