// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package edgezones

import "github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"

func Expand(input string) *edgezones.Model {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &edgezones.Model{
		Name: normalized,
	}
}
