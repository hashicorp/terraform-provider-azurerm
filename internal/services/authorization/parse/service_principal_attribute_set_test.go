// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package parse

import "testing"

func TestServicePrincipalAttributeSetID(t *testing.T) {
	t.Parallel()

	input := "00000000-0000-0000-0000-000000000000|EngineeringAttributes"

	id, err := ServicePrincipalAttributeSetID(input)
	if err != nil {
		t.Fatalf("expected no error but got: %+v", err)
	}
	if id.ServicePrincipalObjectId != "00000000-0000-0000-0000-000000000000" {
		t.Fatalf("unexpected service principal object ID: %q", id.ServicePrincipalObjectId)
	}
	if id.AttributeSetName != "EngineeringAttributes" {
		t.Fatalf("unexpected attribute set name: %q", id.AttributeSetName)
	}
	if got := id.ID(); got != input {
		t.Fatalf("unexpected ID() value: %q", got)
	}
}

func TestServicePrincipalAttributeSetID_requiresPipeSeparator(t *testing.T) {
	t.Parallel()

	if _, err := ServicePrincipalAttributeSetID("00000000-0000-0000-0000-000000000000"); err == nil {
		t.Fatal("expected an error for invalid ID format but got nil")
	}
}

func TestServicePrincipalAttributeSetID_requiresValidServicePrincipalObjectID(t *testing.T) {
	t.Parallel()

	if _, err := ServicePrincipalAttributeSetID("not-a-uuid|EngineeringAttributes"); err == nil {
		t.Fatal("expected an error for invalid service principal object ID but got nil")
	}
}

func TestServicePrincipalAttributeSetID_requiresAttributeSetName(t *testing.T) {
	t.Parallel()

	if _, err := ServicePrincipalAttributeSetID("00000000-0000-0000-0000-000000000000|"); err == nil {
		t.Fatal("expected an error for empty attribute set name but got nil")
	}
}
