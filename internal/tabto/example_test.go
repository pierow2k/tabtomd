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
// characters. Most notably, it shows that the table does not align
// properly with characters that occupy more than one byte, such as Chinese
// characters.
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
