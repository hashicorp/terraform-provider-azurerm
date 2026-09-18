// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strconv"
	"testing"
)

func TestPortNumber(t *testing.T) {
	cases := []struct {
		Port   int
		Errors int
	}{
		{
			Port:   -1,
			Errors: 1,
		},
		{
			Port:   0,
			Errors: 1,
		},
		{
			Port:   1,
			Errors: 0,
		},
		{
			Port:   8477,
			Errors: 0,
		},
		{
			Port:   65535,
			Errors: 0,
		},
		{
			Port:   65536,
			Errors: 1,
		},
		{
			Port:   7000000,
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(strconv.Itoa(tc.Port), func(t *testing.T) {
			_, errors := PortNumber(tc.Port, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected PortNumber to return %d error(s) not %d", len(errors), tc.Errors)
			}
		})
	}
}

func TestPortNumberOrZero(t *testing.T) {
	cases := []struct {
		Port   int
		Errors int
	}{
		{
			Port:   -1,
			Errors: 1,
		},
		{
			Port:   0,
			Errors: 0,
		},
		{
			Port:   1,
			Errors: 0,
		},
		{
			Port:   8477,
			Errors: 0,
		},
		{
			Port:   65535,
			Errors: 0,
		},
		{
			Port:   65536,
			Errors: 1,
		},
		{
			Port:   7000000,
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(strconv.Itoa(tc.Port), func(t *testing.T) {
			_, errors := PortNumberOrZero(tc.Port, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected PortNumberOrZero to return %d error(s) not %d", len(errors), tc.Errors)
			}
		})
	}
}
