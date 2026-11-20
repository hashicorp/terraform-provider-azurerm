// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_test

import (
	"crypto/rand"
	"fmt"
	"os"
)

func populateTempFile(input *os.File) error {
	if err := input.Truncate(25*1024*1024 + 512); err != nil {
		return fmt.Errorf("Failed to truncate file to 25M")
	}

	for i := int64(0); i < 20; i += 2 {
		randomBytes := make([]byte, 1*1024*1024)
		if _, err := rand.Read(randomBytes); err != nil {
			return fmt.Errorf("Failed to read random bytes")
		}

		if _, err := input.WriteAt(randomBytes, i*1024*1024); err != nil {
			return fmt.Errorf("Failed to write random bytes to file")
		}
	}

	randomBytes := make([]byte, 5*1024*1024)
	if _, err := rand.Read(randomBytes); err != nil {
		return fmt.Errorf("Failed to read random bytes")
	}

	if _, err := input.WriteAt(randomBytes, 20*1024*1024); err != nil {
		return fmt.Errorf("Failed to write random bytes to file")
	}

	if err := input.Close(); err != nil {
		return fmt.Errorf("Failed to close source blob")
	}

	return nil
}
