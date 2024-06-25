package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".go") && !strings.HasSuffix(info.Name(), "_test.go") {
			found := processGoFile(path)
			if found {
				testFilePath := strings.TrimSuffix(path, ".go") + "_test.go"
				processTestFile(testFilePath)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path: %v\n", err)
	}
}

func processGoFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", path, err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var fieldStack [][]string       // Stack to track field names and their indentation
	var lastRelevantFields []string // To store the output fields

	for scanner.Scan() {
		line := scanner.Text()
		indentLevel := countLeadingTabs(line)

		// Ensure the fieldStack is large enough for the current indent level
		for len(fieldStack) <= indentLevel {
			fieldStack = append(fieldStack, []string{}) // Expand fieldStack to accommodate new level
		}

		// Extract field names
		if fieldMatch := regexp.MustCompile(`"(\w+)": {`).FindStringSubmatch(line); fieldMatch != nil {
			currentField := fmt.Sprintf("%s%s", strings.Repeat(" ", indentLevel*4), fieldMatch[1])
			fieldStack[indentLevel] = []string{currentField}
			lastRelevantFields = fieldStack[indentLevel]
		}

		if strings.Contains(line, "CaseDifferenceV2Only") {
			fmt.Printf("File: %s\n", path)
			for _, field := range lastRelevantFields {
				fmt.Println(field)
			}
			fmt.Println("") // Print a blank line for readability
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", path, err)
	}
	return false
}

func processTestFile(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return // No corresponding test file exists
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", path, err)
		return
	}
	defer file.Close()

	fmt.Println("Processing Test File:", path)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "acceptance.BuildTestData") {
			fmt.Println(line) // Output each matching line directly
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", path, err)
	}
}

func countLeadingTabs(s string) int {
	return len(s) - len(strings.TrimLeft(s, "\t"))
}
