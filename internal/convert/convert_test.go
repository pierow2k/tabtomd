package convert_test

import (
	"strings"
	"testing"

	"github.com/pierow2k/tabtomd/internal/convert"
)

// Helper function to generate a large tab-delimited text block.
func generateTabDelimitedData(rows, cols, cellLength int) string {
	var stringBuilder strings.Builder

	for iterate := range rows {
		for j := range cols {
			stringBuilder.WriteString(strings.Repeat("x", cellLength))

			if j < cols-1 {
				stringBuilder.WriteString("\t")
			}
		}

		if iterate < rows-1 {
			stringBuilder.WriteString("\n")
		}
	}

	return stringBuilder.String()
}

// Benchmark for ToMDTablePretty.
func BenchmarkToMDTablePretty(b *testing.B) {
	// Generate synthetic data with 100 rows, 10 columns, and 20 characters per cell.
	data := generateTabDelimitedData(100, 10, 20)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := convert.ToMDTablePretty(data)
		if err != nil {
			b.Fatalf("Error during benchmark: %v", err)
		}
	}
}

// Benchmark for ToMDTable.
func BenchmarkToMDTable(b *testing.B) {
	// Generate synthetic data with 100 rows, 10 columns, and 20 characters per cell.
	data := generateTabDelimitedData(100, 10, 20)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := convert.ToMDTable(data)
		if err != nil {
			b.Fatalf("Error during benchmark: %v", err)
		}
	}
}
