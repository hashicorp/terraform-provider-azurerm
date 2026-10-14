// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package schemaz

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperation"
)

func TestFlattenApiManagementOperationParameterExampleContract_stableOrder(t *testing.T) {
	input := map[string]apioperation.ParameterExampleContract{
		"UniqueMatch": {Summary: pointer.To("unique")},
		"NoMatch":     {Summary: pointer.To("none")},
		"Alpha":       {Summary: pointer.To("alpha")},
		"Zulu":        {Summary: pointer.To("zulu")},
		"Mike":        {Summary: pointer.To("mike")},
	}
	expected := []string{"Alpha", "Mike", "NoMatch", "UniqueMatch", "Zulu"}

	// map iteration order is randomised, so repeat to make an unsorted implementation fail reliably
	for range 50 {
		output, err := FlattenApiManagementOperationParameterExampleContract(input)
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if len(output) != len(expected) {
			t.Fatalf("expected %d examples but got %d", len(expected), len(output))
		}
		for j, item := range output {
			if name := item.(map[string]any)["name"].(string); name != expected[j] {
				t.Fatalf("expected example %d to be %q but got %q", j, expected[j], name)
			}
		}
	}
}
