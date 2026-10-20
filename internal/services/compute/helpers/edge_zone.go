// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
)

func ExpandEdgeZone(input string) *edgezones.Model {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &edgezones.Model{
		Name: normalized,
	}
}

func FlattenEdgeZone(input *edgezones.Model) string {
	if input == nil || input.Name == "" {
		return ""
	}
	return edgezones.Normalize(input.Name)
}

func ExpandManagedDiskEdgeZone(input string) *edgezones.Model {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &edgezones.Model{
		Name: normalized,
	}
}

func FlattenManagedDiskEdgeZone(input *edgezones.Model) string {
	if input == nil || input.Name == "" {
		return ""
	}
	return edgezones.NormalizeNilable(&input.Name)
}
