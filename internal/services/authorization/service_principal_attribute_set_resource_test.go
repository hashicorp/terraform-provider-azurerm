// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package authorization

import "testing"

func TestUpsertManagedAttributeSet_preservesOtherSets(t *testing.T) {
	t.Parallel()

	existing := map[string]map[string]interface{}{
		"Finance": {
			"CostCenter": "1234",
		},
	}
	attributes := map[string]interface{}{
		"Environment": "Prod",
	}

	out, err := upsertManagedAttributeSet(existing, "Engineering", attributes)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	if got := out["Finance"]["CostCenter"]; got != "1234" {
		t.Fatalf("expected Finance.CostCenter to be preserved as 1234 but got %v", got)
	}
	if got := out["Engineering"]["Environment"]; got != "Prod" {
		t.Fatalf("expected Engineering.Environment to be Prod but got %v", got)
	}
}

func TestUpsertManagedAttributeSet_requiresStringValues(t *testing.T) {
	t.Parallel()

	attributes := map[string]interface{}{
		"Environment": 123,
	}

	if _, err := upsertManagedAttributeSet(map[string]map[string]interface{}{}, "Engineering", attributes); err == nil {
		t.Fatal("expected an error for non-string attribute value but got nil")
	}
}

func TestFlattenManagedAttributeSet_skipsODataType(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"@odata.type": "#Microsoft.DirectoryServices.CustomSecurityAttributeValue",
		"Environment": "Prod",
	}

	out, err := flattenManagedAttributeSet(input)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	if _, exists := out["@odata.type"]; exists {
		t.Fatal("expected @odata.type to be removed from flattened attributes")
	}
	if got := out["Environment"]; got != "Prod" {
		t.Fatalf("expected Environment=Prod but got %v", got)
	}
}
