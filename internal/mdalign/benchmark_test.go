package mdalign_test

import (
	"bufio"
	"compress/bzip2"
	"os"
	"testing"

	"github.com/pierow2k/tabtomd/internal/mdalign"
)

func loadRows(tb testing.TB, path string) []string {
	tb.Helper()

	f, err := os.Open(path)
	if err != nil {
		tb.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()

	bZipReader := bzip2.NewReader(f)

	var rows []string
	scanner := bufio.NewScanner(bZipReader)

	// Optional: raise this if you may have very long lines.
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		rows = append(rows, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		tb.Fatalf("scan %s: %v", path, err)
	}

	return rows
}

func BenchmarkAlign(b *testing.B) {
	small := []string{
		"| Name | Age | City         |",
		"|----|-----|----|",
		"| Alice   |   30  |   New York     |",
		"|  Bob     | 25  |  Los Angeles  |",
		"| Charlie   | 35  | San Francisco|",
	}

	large := loadRows(b, "testdata/2k.md.bz2")

	b.ReportAllocs()

	b.Run("small", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = mdalign.Align(small)
		}
	})

	b.Run("2k rows", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = mdalign.Align(large)
		}
	})
}
