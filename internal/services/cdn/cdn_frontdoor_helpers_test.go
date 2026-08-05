// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import "testing"

func TestNormalizeFrontDoorOriginGroupID_caseInsensitive(t *testing.T) {
	input := "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/example-rg/providers/microsoft.cdn/profiles/example-profile/origingroups/example-origin-group"
	expected := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Cdn/profiles/example-profile/originGroups/example-origin-group"

	actual, err := normalizeFrontDoorOriginGroupID(input)
	if err != nil {
		t.Fatalf("normalizing Front Door origin group ID: %+v", err)
	}

	if actual != expected {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}
