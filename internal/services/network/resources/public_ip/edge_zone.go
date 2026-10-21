// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package public_ip

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
)

// These will be renamed to expandEdgeZone when all calls to the former expandEdgeZone have been removed
func expandEdgeZoneNew(input string) *edgezones.Model {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &edgezones.Model{
		Name: normalized,
	}
}

func flattenEdgeZoneNew(input *edgezones.Model) string {
	if input == nil || input.Name == "" {
		return ""
	}
	return edgezones.Normalize(input.Name)
}
