// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/validate"
)

func TestWorkspaceName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// empty
			Input: "",
			Valid: false,
		},
		{
			// basic example
			Input: "hello",
			Valid: true,
		},
		{
			// 63 chars
			Input: "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk",
			Valid: true,
		},
		{
			// 64 chars
			Input: "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkl",
			Valid: true,
		},
		{
			// 65 chars
			Input: "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklm",
			Valid: false,
		},
		{
			// may contain alphanumerics, dots, dashes and underscores
			Input: "hello_world7.goodbye-world4",
			Valid: true,
		},
		{
			// must begin with an alphanumeric
			Input: "_hello",
			Valid: false,
		},
		{
			// can't end with a period
			Input: "hello.",
			Valid: false,
		},
		{
			// can't end with a dash
			Input: "hello-",
			Valid: false,
		},
		{
			// can end with an underscore
			Input: "hello_",
			Valid: true,
		},
		{
			// can't contain an exclamation mark
			Input: "hello!",
			Valid: false,
		},
		{
			// can start with a number
			Input: "0abc",
			Valid: true,
		},
		{
			// can contain only numbers
			Input: "12345",
			Valid: true,
		},
		{
			// can start with upper case letter
			Input: "Test",
			Valid: true,
		},
		{
			// can end with upper case letter
			Input: "TEST",
			Valid: true,
		},
	}

	for _, tc := range cases {
		_, errs := validate.WorkspaceName(tc.Input, "name")
		valid := len(errs) == 0

		if valid != tc.Valid {
			t.Fatalf("expected %s to be %t, got %t", tc.Input, tc.Valid, valid)
		}
	}
}
