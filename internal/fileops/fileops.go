// Package fileops provides utilities for reading tab-delimited files
// and writing Markdown files. It is designed to support file I/O
// operations for table conversions and Markdown generation.
package fileops

import (
	"fmt"
	"io"
	"os"
)

// ReadTSV reads a tab-separated values (TSV) file and returns its contents
// as a string. If the file cannot be opened or read, it returns an error.
//
// The function ensures the file is properly closed after reading, even
// if an error occurs during the read operation.
func ReadTSV(filename string) (string, error) {
	// Open the specified file
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}

	// Ensure the file is closed after reading
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Printf("warning: failed to close file: %v\n", cerr)
		}
	}()

	// Read the entire content of the file
	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %w", err)
	}

	return string(content), nil
}

// WriteMD writes a string containing a Markdown table to the specified
// file. If the file does not exist, it is created. If the file exists, its
// content is overwritten.
//
// The function ensures the file is properly closed after writing, even
// if an error occurs during the write operation.
func WriteMD(filename, content string) error {
	// Create or overwrite the output file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	// Ensure the file is closed after writing
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Printf("warning: failed to close file: %v\n", cerr)
		}
	}()

	// Write the content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write content to file: %w", err)
	}

	return nil
}
