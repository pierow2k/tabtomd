// Package tabto_test provides black-box tests and runnable examples
// for the public API of the tabto package.
package tabto_test

import (
	"fmt"

	"github.com/pierow2k/tabtomd/internal/tabto"
)

// HTML converts tab delimited input into an HTML table.
func ExampleHTML() {
	input := "Name\tAge\tCity\nAlice\t30\tNew York\nBob\t25\tLos Angeles"

	htmlTable, err := tabto.HTML(input)
	if err != nil {
		panic(err)
	}

	fmt.Println(htmlTable)

	// Output:
	// <table class="go-pretty-table">
	//   <thead>
	//   <tr>
	//     <th>Name</th>
	//     <th>Age</th>
	//     <th>City</th>
	//   </tr>
	//   </thead>
	//   <tbody>
	//   <tr>
	//     <td>Alice</td>
	//     <td>30</td>
	//     <td>New York</td>
	//   </tr>
	//   <tr>
	//     <td>Bob</td>
	//     <td>25</td>
	//     <td>Los Angeles</td>
	//   </tr>
	//   </tbody>
	// </table>
}

// The Markdown function converts tab delimited input into Markdown table
// format.
func ExampleMarkdown() {
	input := "Name\tAge\tCity\nAlice\t30\tNew York\nBob\t25\tLos Angeles"

	markdownTable, err := tabto.Markdown(input)
	if err != nil {
		panic(err)
	}

	fmt.Println(markdownTable)

	// Output:
	// | Name | Age | City |
	// | --- | --- | --- |
	// | Alice | 30 | New York |
	// | Bob | 25 | Los Angeles |
}

// Empty input results in an empty string result.
func ExampleMarkdown_emptyInput() {
	input := ""

	markdownTable, err := tabto.Markdown(input)
	if err != nil {
		panic(err)
	}

	if markdownTable == "" {
		fmt.Println("No data to convert.")
	} else {
		fmt.Println(markdownTable)
	}

	// Output:
	// No data to convert.
}

// Special characters such as pipes and backslashes are escaped within the table data.
func ExampleMarkdown_specialCharacters() {
	input := "Name|Pipe\tAge\tCity\nAlice|test\t30\tNew York\nBob\\back\t25\tLos Angeles"

	markdownTable, err := tabto.Markdown(input)
	if err != nil {
		panic(err)
	}

	fmt.Println(markdownTable)
	// Output:
	// | Name\|Pipe | Age | City |
	// | --- | --- | --- |
	// | Alice\|test | 30 | New York |
	// | Bob\back | 25 | Los Angeles |
}

// A single row of tab delimited characters produces a header separator of
// `---:`.
func ExampleMarkdown_singleRow() {
	input := "Header1\tHeader2\tHeader3"

	markdownTable, err := tabto.Markdown(input)
	if err != nil {
		panic(err)
	}

	fmt.Println(markdownTable)
	// Output:
	// | Header1 | Header2 | Header3 |
	// | ---:| ---:| ---:|
}

// multi-byte characters are handled correctly when determining column
// widths for basic table generation.
//
//nolint:gosmopolitan
func ExampleMarkdown_multiByteCharacters() {
	input := "你好\ta\n世界\tb"

	markdownTable, err := tabto.Markdown(input)
	if err != nil {
		panic(err)
	}

	fmt.Println(markdownTable)
	// Output:
	// | 你好 | a |
	// | --- | --- |
	// | 世界 | b |
}

// Markdown returns an error when the input data has an inconsistent
// number of columns. The error message includes the row number and
// expected/actual counts.
func ExampleMarkdown_inconsistentColumns() {
	input := "h1\th2\nh1" // Second row has fewer columns

	_, err := tabto.Markdown(input)
	if err != nil {
		fmt.Println(err)

		return
	}
	// Output:
	// failed to parse table: row has inconsistent column count: row 2 (expected 2 columns, got 1)
}

// Leading and trailing empty lines are trimmed, but internal empty lines
// are preserved as empty rows in the resulting table.
func ExampleMarkdown_withEmptyLines() {
	input := "\n\nName\tAge\n\nAlice\t30\n\n"

	markdownTable, err := tabto.Markdown(input)
	if err != nil {
		panic(err)
	}

	fmt.Println(markdownTable)
	// Output:
	// | Name | Age |
	// | --- | --- |
	// |  |  |
	// | Alice | 30 |
}
