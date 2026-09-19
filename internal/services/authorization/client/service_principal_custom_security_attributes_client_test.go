// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import "testing"

func TestNormalizeCustomSecurityAttributesForPatch_addsODataType(t *testing.T) {
	t.Parallel()

	input := map[string]map[string]interface{}{
		"Engineering": {
			"Environment": "Prod",
		},
	}

	out := normalizeCustomSecurityAttributesForPatch(input)
	if got := out["Engineering"]["@odata.type"]; got != customSecurityAttributeValueODataType {
		t.Fatalf("expected @odata.type to be %q but got %v", customSecurityAttributeValueODataType, got)
	}
	if got := out["Engineering"]["Environment"]; got != "Prod" {
		t.Fatalf("expected Environment to be preserved but got %v", got)
	}
}

func TestNormalizeCustomSecurityAttributesForPatch_replacesExistingODataType(t *testing.T) {
	t.Parallel()

	input := map[string]map[string]interface{}{
		"Engineering": {
			"@odata.type": "unexpected",
			"Environment": "Prod",
		},
	}

	out := normalizeCustomSecurityAttributesForPatch(input)
	if got := out["Engineering"]["@odata.type"]; got != customSecurityAttributeValueODataType {
		t.Fatalf("expected @odata.type to be %q but got %v", customSecurityAttributeValueODataType, got)
	}
}
