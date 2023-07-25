// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity

import "strings"

type Type string

const (
	TypeNone                       Type = "None"
	TypeSystemAssigned             Type = "SystemAssigned"
	TypeUserAssigned               Type = "UserAssigned"
	TypeSystemAssignedUserAssigned Type = "SystemAssigned, UserAssigned"

	// this is an internal-only type to transform the legacy API value to the type we want to expose
	typeLegacySystemAssignedUserAssigned Type = "SystemAssigned,UserAssigned"
)

func normalizeType(input Type) Type {
	// switch out the legacy API value (no space) for the value used in the Schema (w/space for consistency)
	if strings.EqualFold(string(input), string(typeLegacySystemAssignedUserAssigned)) {
		return TypeSystemAssignedUserAssigned
	}

	vals := []Type{
		TypeNone,
		TypeSystemAssigned,
		TypeUserAssigned,
		TypeSystemAssignedUserAssigned,
	}
	for _, v := range vals {
		if strings.EqualFold(string(input), string(v)) {
			return v
		}
	}
	return input
}
