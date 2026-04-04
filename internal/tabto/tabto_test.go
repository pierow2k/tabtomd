// Package tabto_test provides black-box tests and testable examples for
// the public API of the tabto package.
//
//nolint:funlen
package tabto_test

import (
	"strings"
	"testing"

	"github.com/pierow2k/tabtomd/internal/tabto"
)

func TestMarkdown(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		text    string
		want    string
		wantErr bool
	}{
		{
			name:    "normal table",
			text:    "Name\tAge\tCity\nAlice\t30\tNew York\nBob\t25\tLos Angeles",
			wantErr: false,
			want:    "| Name | Age | City |",
		},
		{
			name:    "empty input",
			text:    "",
			wantErr: false,
			want:    "",
		},
		{
			name:    "special characters",
			text:    "Name|Pipe\tAge\tCity\nAlice|test\t30\tNew York\nBob\\back\t25\tLos Angeles",
			wantErr: false,
			want:    "Alice\\|test", // Should escape the pipe
		},
		{
			name:    "single row",
			text:    "Header1\tHeader2\tHeader3",
			wantErr: false,
			want:    "| Header1 | Header2 | Header3 |",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result, err := tabto.Markdown(testCase.text)
			if (err != nil) != testCase.wantErr {
				t.Errorf("Markdown() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if !testCase.wantErr && !strings.Contains(result, testCase.want) {
				t.Errorf(
					"Markdown() result does not contain expected text:\ngot:\n%s\nexpected to contain: %s",
					result,
					testCase.want)
			}
		})
	}
}

func TestHTML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		text    string
		want    string
		wantErr bool
	}{
		{
			name:    "normal table",
			text:    "Name\tAge\tCity\nAlice\t30\tNew York\nBob\t25\tLos Angeles",
			wantErr: false,
			want: `<table class="go-pretty-table">
  <thead>
  <tr>
    <th>Name</th>
    <th>Age</th>
    <th>City</th>
  </tr>
  </thead>
  <tbody>
  <tr>
    <td>Alice</td>
    <td>30</td>
    <td>New York</td>
  </tr>
  <tr>
    <td>Bob</td>
    <td>25</td>
    <td>Los Angeles</td>
  </tr>
  </tbody>
</table>
`,
		},
		{
			name:    "empty input",
			text:    "",
			wantErr: false,
			want:    "",
		},
		{
			name:    "special characters",
			text:    "Name|Pipe\tAge\tCity\nAlice|test\t30\tNew York\nBob\\back\t25\tLos Angeles",
			wantErr: false,
			want: `<table class="go-pretty-table">
  <thead>
  <tr>
    <th>Name|Pipe</th>
    <th>Age</th>
    <th>City</th>
  </tr>
  </thead>
  <tbody>
  <tr>
    <td>Alice|test</td>
    <td>30</td>
    <td>New York</td>
  </tr>
  <tr>
    <td>Bob\back</td>
    <td>25</td>
    <td>Los Angeles</td>
  </tr>
  </tbody>
</table>
`, // Should escape the pipe
		},
		{
			name:    "single row",
			text:    "Header1\tHeader2\tHeader3",
			wantErr: false,
			want: `<table class="go-pretty-table">
  <thead>
  <tr>
    <th align="right">Header1</th>
    <th align="right">Header2</th>
    <th align="right">Header3</th>
  </tr>
  </thead>
</table>
`,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := tabto.HTML(testCase.text)
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("HTML() failed: %v", gotErr)
				}

				return
			}

			if testCase.wantErr {
				t.Fatal("HTML() succeeded unexpectedly")
			}

			if !testCase.wantErr && got != testCase.want {
				t.Errorf(
					"HTML() result does not contain expected text:\ngot:\n'%s'\nexpected:\n'%s'",
					got,
					testCase.want)
			}
		})
	}
}
