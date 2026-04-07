// Package mdalign_test provides black-box tests and runnable examples
// for the public API of the mdalign package.
package mdalign_test

import (
	"errors"
	"fmt"

	"github.com/pierow2k/tabtomd/internal/mdalign"
)

func ExampleAlign() {
	rows := []string{
		"| Name | Age | City         |",
		"|----|-----|----|",
		"| Alice   |   30  |   New York     |",
		"|  Bob     | 25  |  Los Angeles  |",
		"| Charlie   | 35  | San Francisco|",
	}

	alignedRows, err := mdalign.Align(rows)
	if err != nil {
		panic(err)
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

// Align returns an empty output for an empty input.
func ExampleAlign_emptyInput() {
	rows := []string{}

	alignedRows, err := mdalign.Align(rows)
	if err != nil {
		panic(err)
	}

	if len(alignedRows) == 0 {
		fmt.Println("No rows to align.")
	}

	// Output:
	// No rows to align.
}

// Align returns the sentinel error ErrAlignmentFailed when the input data
// has an inconsistent number of columns.
func ExampleAlign_inconsistentColumns() {
	rows := []string{
		"| Name    | Age | City         |",
		"|---------|-----|--------------|",
		"| Alice   | 30  | New York     |",
		"| Bob     | 25  | Los Angeles  |",
		"| Charlie | 35  |", // Inconsistent column count
	}

	_, err := mdalign.Align(rows)
	if errors.Is(err, mdalign.ErrAlignmentFailed) {
		fmt.Println("Error:", err)
	}

	// Output:
	// Error: alignment failed: row 5 has inconsistent column count
}

// Align properly handles escaped characters in the input.
func ExampleAlign_escapedCharacters() {
	rows := []string{
		"| Name    | Age | City         |",
		"|---------|-----|--------------|",
		"| Alice\\A  | 30  | New York     |",
		"| Bob\\B  | 25  | Los Angeles  |",
		"| Charlie | 35  | San Francisco|",
	}

	alignedRows, err := mdalign.Align(rows)
	if err != nil {
		panic(err)
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

// Align returns a formatted header row when given a single row to process.
func ExampleAlign_singleRow() {
	rows := []string{
		"| Header1 | Header2 | Header3 |",
	}

	alignedRows, err := mdalign.Align(rows)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	for _, row := range alignedRows {
		fmt.Println(row)
	}

	// Output:
	// | Header |
	// | ------ |
}

// Align handles multi-byte characters. Notably, the table align properly
// with characters that occupy more than one byte, such as Chinese
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
		panic(err)
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
		panic(err)
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
