// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/durabletask/validate"
)

func TestTaskHubName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{Input: "validhub", Valid: true},
		{Input: "valid-hub", Valid: true},
		{Input: "hub123", Valid: true},
		{Input: "my-task-hub", Valid: true},
		{Input: "abc", Valid: true},
		{Input: strings.Repeat("a", 63), Valid: true},
		{Input: "taskhub-1", Valid: true},
		{Input: "test123hub", Valid: true},
		{Input: "", Valid: false},
		{Input: "ab", Valid: false},
		{Input: "-invalid", Valid: false},
		{Input: "invalid-", Valid: false},
		{Input: strings.Repeat("a", 64), Valid: false},
	}

	for _, tc := range cases {
		_, errors := validate.TaskHubName(tc.Input, "task_hub_name")
		valid := len(errors) == 0
		if tc.Valid != valid {
			t.Fatalf("expected valid=%t for %q, got valid=%t", tc.Valid, tc.Input, valid)
		}
	}
}
