// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package edgezones

import "github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"

func FlattenEdgeZone(input *edgezones.Model) string {
	if input == nil || input.Name == "" {
		return ""
	}
	return edgezones.Normalize(input.Name)
}
