package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Function to find the header and return its index
func findHeaderIndex(lines []string, header string) int {
	for i, line := range lines {
		if strings.HasPrefix(line, header) {
			return i
		}
	}
	return -1
}

// Function to append the new entry under the appropriate header in alphabetical order
func appendUnderHeader(filePath string, newEntry, header string) error {
	// Open the file for reading and appending
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file content
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Find the correct header section
	headerIndex := findHeaderIndex(lines, header)

	// Now append the new entry under the correct header
	// Check if the next line is empty or not for proper formatting
	insertIndex := headerIndex + 1
	for i := headerIndex + 1; i < len(lines); i++ {
		// Look for the next header to break the section
		if strings.HasPrefix(lines[i], "[") {
			insertIndex = i
			break
		}
	}

	// Remove the header prefix from the new entry
	// Trim the header prefix based on which one it matches
	var trimmedEntry string
	if strings.HasPrefix(newEntry, "[BUG]") {
		trimmedEntry = strings.TrimPrefix(newEntry, "[BUG] ")
	} else if strings.HasPrefix(newEntry, "[ENHANCEMENT]") {
		trimmedEntry = strings.TrimPrefix(newEntry, "[ENHANCEMENT] ")
	} else if strings.HasPrefix(newEntry, "[FEATURE]") {
		trimmedEntry = strings.TrimPrefix(newEntry, "[FEATURE] ")
	} else {
		// If the entry doesn't match one of the expected headers, print an error
		fmt.Println("Error: New entry must start with one of the headers [BUG], [ENHANCEMENT], or [FEATURE].")
		return nil
	}

	// Insert the new entry under the header
	var section []string
	for i := headerIndex + 1; i < insertIndex; i++ {
		section = append(section, lines[i])
	}
	section = append(section, trimmedEntry)

	// Rebuild the file content
	lines = append(lines[:headerIndex+1], append(section, lines[insertIndex:]...)...)

	// Open the file for writing and overwrite the content
	file, err = os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the updated content back to the file
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run update_changelog.go CHANGELOG.md <new_entry>")
		return
	}

	filePath := os.Args[1]
	newEntry := os.Args[2]

	// Validate and determine the correct header for the new entry
	var selectedHeader string
	if strings.HasPrefix(newEntry, "[BUG]") {
		selectedHeader = "BUG FIXES:"
	} else if strings.HasPrefix(newEntry, "[ENHANCEMENT]") {
		selectedHeader = "ENHANCEMENTS:"
	} else if strings.HasPrefix(newEntry, "[FEATURE]") {
		selectedHeader = "FEATURES:"
	} else {
		// If the entry doesn't match one of the expected headers, print an error
		fmt.Println("Error: New entry must start with one of the headers [BUG], [ENHANCEMENT], or [FEATURE].")
		return
	}

	// Call the function to append under the appropriate header
	if err := appendUnderHeader(filePath, newEntry, selectedHeader); err != nil {
		fmt.Println("Error appending to file:", err)
		return
	}

	fmt.Println("Successfully appended the new entry under the", selectedHeader, "header.")
}
