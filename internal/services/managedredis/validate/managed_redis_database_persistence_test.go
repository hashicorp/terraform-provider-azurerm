// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"testing"
)

func TestAofBackupFrequency(t *testing.T) {
	tests := []struct {
		name           string
		input          interface{}
		expectedErrors []error
	}{
		{
			name:           "1 is a valid value",
			input:          1,
			expectedErrors: nil,
		},
		{
			name:           "-1 is not a valid value",
			input:          -1,
			expectedErrors: []error{fmt.Errorf(`expected "persistence_append_only_file_backup_frequency_in_seconds" to be one of [1], got -1`)},
		},
		{
			name:           "60 is not a valid value",
			input:          60,
			expectedErrors: []error{fmt.Errorf(`expected "persistence_append_only_file_backup_frequency_in_seconds" to be one of [1], got 60`)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, errors := AofBackupFrequency(tt.input, "persistence_append_only_file_backup_frequency_in_seconds")
			if !equalErrorSlices(errors, tt.expectedErrors) {
				t.Errorf("\nexpected errors: %v\ngot %v", tt.expectedErrors, errors)
			}
		})
	}
}

func TestRdbBackupFrequency(t *testing.T) {
	tests := []struct {
		name           string
		input          interface{}
		expectedErrors []error
	}{
		{
			name:           "1 is a valid value",
			input:          1,
			expectedErrors: nil,
		},
		{
			name:           "6 is a valid value",
			input:          6,
			expectedErrors: nil,
		},
		{
			name:           "12 is a valid value",
			input:          12,
			expectedErrors: nil,
		},
		{
			name:           "7 is not a valid value",
			input:          7,
			expectedErrors: []error{fmt.Errorf(`expected "persistence_redis_database_backup_frequency_in_hours" to be one of [1, 6, 12], got 7`)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, errors := RdbBackupFrequency(tt.input, "persistence_redis_database_backup_frequency_in_hours")
			if !equalErrorSlices(errors, tt.expectedErrors) {
				t.Errorf("\nexpected errors: %v\ngot %v", tt.expectedErrors, errors)
			}
		})
	}
}

func equalErrorSlices(a, b []error) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Error() != b[i].Error() {
			return false
		}
	}
	return true
}
