// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package loader

import (
	"fmt"
	"log"
	"os"
)

// DiffFileLoader loads changes from a diff file
type DiffFileLoader struct {
	filePath string
}

// Load loads changes from a diff file and returns a ChangeSet
func (l *DiffFileLoader) Load() (*ChangeSet, error) {
	cs := NewChangeSet()

	content, err := os.ReadFile(l.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read diff file: %w", err)
	}

	// Use the common parseDiffOutput function
	if err := cs.parseDiffOutput(string(content)); err != nil {
		return nil, err
	}

	if len(cs.changedFiles) == 0 {
		return nil, nil
	}

	log.Printf("✓ Found %d changed files with %d changed lines",
		len(cs.changedFiles), cs.getTotalChangedLines())

	return cs, nil
}
