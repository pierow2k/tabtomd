// Package tabto_test provides black-box tests and runnable examples
// for the public API of the tabto package.
package tabto_test

import (
	"fmt"

	"github.com/pierow2k/tabtomd/internal/tabto"
)

// ExampleMarkdown demonstrates the use of the Markdown function to
// convert tab delimited input into Markdown table format.
func ExampleMarkdown() {
	input := "Name\tAge\tCity\nAlice\t30\tNew York\nBob\t25\tLos Angeles"

	markdownTable, err := tabto.Markdown(input)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	fmt.Println(markdownTable)

	// Output:
	// | Name | Age | City |
	// | --- | --- | --- |
	// | Alice | 30 | New York |
	// | Bob | 25 | Los Angeles |
}

// ExampleMarkdown_emptyInput demonstrates handling of empty input.
func ExampleMarkdown_emptyInput() {
	input := ""

	markdownTable, err := tabto.Markdown(input)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	if markdownTable == "" {
		fmt.Println("No data to convert.")
	} else {
		fmt.Println(markdownTable)
	}

	// Output:
	// No data to convert.
}

// ExampleMarkdown_specialCharacters demonstrates handling of special characters
// such as pipes and backslashes within the table data.
func ExampleMarkdown_specialCharacters() {
	input := "Name|Pipe\tAge\tCity\nAlice|test\t30\tNew York\nBob\\back\t25\tLos Angeles"

	markdownTable, err := tabto.Markdown(input)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	fmt.Println(markdownTable)
	// Output:
	// | Name\|Pipe | Age | City |
	// | --- | --- | --- |
	// | Alice\|test | 30 | New York |
	// | Bob\back | 25 | Los Angeles |
}

// ExampleMarkdown_singleRow demonstrates the handling of a single row of
// tab delimited characters. Note that the header separator defaults to
// `---:` in this case.
func ExampleMarkdown_singleRow() {
	input := "Header1\tHeader2\tHeader3"

	markdownTable, err := tabto.Markdown(input)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	fmt.Println(markdownTable)
	// Output:
	// | Header1 | Header2 | Header3 |
	// | ---:| ---:| ---:|
}

// ExampleMarkdown_multiByteCharacters demonstrates handling of multi-byte
// characters. The `go-pretty` library correctly handles these characters
// when determining column widths for basic table generation.
//
//nolint:gosmopolitan
func ExampleMarkdown_multiByteCharacters() {
	input := "你好\ta\n世界\tb"

	markdownTable, err := tabto.Markdown(input)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	fmt.Println(markdownTable)
	// Output:
	// | 你好 | a |
	// | --- | --- |
	// | 世界 | b |
}

// ExampleMarkdown_inconsistentColumns demonstrates how the function returns
// an error when the input data has an inconsistent number of columns.
func ExampleMarkdown_inconsistentColumns() {
	input := "h1\th2\nh1" // Second row has fewer columns

	_, err := tabto.Markdown(input)
	if err != nil {
		// The error message includes the row number and expected/actual counts.
		fmt.Println(err)

		return
	}
	// Output:
	// failed to parse table: row has inconsistent column count: row 2 (expected 2 columns, got 1): "h1"
}

// ExampleMarkdown_withEmptyLines demonstrates that leading and trailing
// empty lines are trimmed, but internal empty lines are preserved as empty rows
// in the resulting table.
func ExampleMarkdown_withEmptyLines() {
	input := "\n\nName\tAge\n\nAlice\t30\n\n"

	markdownTable, err := tabto.Markdown(input)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	fmt.Println(markdownTable)
	// Output:
	// | Name | Age |
	// | --- | --- |
	// |  |  |
	// | Alice | 30 |
}
