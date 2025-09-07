// Package tabto_test provides black-box tests and runnable examples for
// the public API of the convert package.
package tabto_test

import (
	"strings"
	"testing"

	"github.com/pierow2k/tabtomd/internal/tabto"
)

func TestMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		contains string
	}{
		{
			name:     "normal table",
			input:    "Name\tAge\tCity\nAlice\t30\tNew York\nBob\t25\tLos Angeles",
			wantErr:  false,
			contains: "| Name | Age | City |",
		},
		{
			name:     "empty input",
			input:    "",
			wantErr:  false,
			contains: "",
		},
		{
			name:     "special characters",
			input:    "Name|Pipe\tAge\tCity\nAlice|test\t30\tNew York\nBob\\back\t25\tLos Angeles",
			wantErr:  false,
			contains: "Alice\\|test", // Should escape the pipe
		},
		{
			name:     "single row",
			input:    "Header1\tHeader2\tHeader3",
			wantErr:  false,
			contains: "| Header1 | Header2 | Header3 |",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tabto.Markdown(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Markdown() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !strings.Contains(result, tt.contains) {
				t.Errorf("Markdown() result does not contain expected text:\ngot:\n%s\nexpected to contain: %s", result, tt.contains)
			}
		})
	}
}
