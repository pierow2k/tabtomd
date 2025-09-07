// Package mdalign_test provides black-box tests and runnable examples
// for the public API of the mdalign package.
package mdalign_test

import (
	"fmt"

	"github.com/pierow2k/tabtomd/internal/mdalign"
)

func ExampleAlign() {
	rows := []string{
		"| Name    | Age | City         |",
		"|---------|-----|--------------|",
		"| Alice   | 30  | New York     |",
		"| Bob     | 25  | Los Angeles  |",
		"| Charlie | 35  | San Francisco|",
	}

	alignedRows, err := mdalign.Align(rows)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	for _, row := range alignedRows {
		fmt.Println(row)
	}

	// Output:
	// | Name    | Age | City          |
	// | ------- | --- | ------------- |
	// | Alice   | 30  | New York      |
	// | Bob     | 25  | Los Angeles   |
	// | Charlie | 35  | San Francisco |
}

func ExampleAlign_emptyInput() {
	rows := []string{}

	alignedRows, err := mdalign.Align(rows)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	if len(alignedRows) == 0 {
		fmt.Println("No rows to align.")

		return
	}

	for _, row := range alignedRows {
		fmt.Println(row)
	}

	// Output:
	// No rows to align.
}

func ExampleAlign_inconsistentColumns() {
	rows := []string{
		"| Name    | Age | City         |",
		"|---------|-----|--------------|",
		"| Alice   | 30  | New York     |",
		"| Bob     | 25  | Los Angeles  |",
		"| Charlie | 35  |", // Inconsistent column count
	}

	alignedRows, err := mdalign.Align(rows)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	for _, row := range alignedRows {
		fmt.Println(row)
	}

	// Output:
	// Error: row has inconsistent column count: row 5
}

func ExampleAlign_specialCharacters() {
	rows := []string{
		"| Name    | Age | City         |",
		"|---------|-----|--------------|",
		"| Alice\\A  | 30  | New York     |",
		"| Bob\\B  | 25  | Los Angeles  |",
		"| Charlie | 35  | San Francisco|",
	}

	alignedRows, err := mdalign.Align(rows)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	for _, row := range alignedRows {
		fmt.Println(row)
	}

	// Output:
	// | Name    | Age | City          |
	// | ------- | --- | ------------- |
	// | Alice\A | 30  | New York      |
	// | Bob\B   | 25  | Los Angeles   |
	// | Charlie | 35  | San Francisco |
}

func ExampleAlign_singleRow() {
	rows := []string{
		"| Header1 | Header2 | Header3 |",
	}

	alignedRows, err := mdalign.Align(rows)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	for _, row := range alignedRows {
		fmt.Println(row)
	}

	// Output:
	// | Header1 | Header2 | Header3 |
	// | ------- | ------- | ------- |
}

func ExampleAlign_onlyHeaderAndSeparator() {
	rows := []string{
		"| Header |",
		"| --- |",
	}

	alignedRows, err := mdalign.Align(rows)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	for _, row := range alignedRows {
		fmt.Println(row)
	}

	// Output:
	// | Header |
	// | ------ |
}

// ExampleAlign_multiByteCharacters demonstrates the handling of multi-byte
// characters. Most notably, it shows that the table does not align
// properly with characters that occupy more than one byte, such as Chinese
// characters.
//
//nolint:gosmopolitan
func ExampleAlign_multiByteCharacters() {
	rows := []string{
		"| 你好 | a |",
		"|---|---|",
		"| 世界 | b |",
	}

	alignedRows, err := mdalign.Align(rows)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	for _, row := range alignedRows {
		fmt.Println(row)
	}

	// Output:
	// | 你好 | a |
	// | --- | --- |
	// | 世界 | b |
}

func ExampleAlign_malformedRow() {
	rows := []string{
		"h1",
		"h2",
	}

	alignedRows, err := mdalign.Align(rows)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	if len(alignedRows) == 0 {
		fmt.Println("No valid table rows to align.")

		return
	}

	for _, row := range alignedRows {
		fmt.Println(row)
	}

	// Output:
	// h1
	// h2
}
