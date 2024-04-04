// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

func expandEdgeZone(input string) *virtualmachines.Model {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &virtualmachines.Model{
		Name: utils.String(normalized),
		Type: compute.ExtendedLocationTypesEdgeZone,
	}
}

func flattenEdgeZone(input *compute.ExtendedLocation) string {
	if input == nil || input.Type != compute.ExtendedLocationTypesEdgeZone || input.Name == nil {
		return ""
	}
	return edgezones.NormalizeNilable(input.Name)
}

func expandManagedDiskEdgeZone(input string) *edgezones.Model {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &edgezones.Model{
		Name: normalized,
	}
}

func flattenManagedDiskEdgeZone(input *edgezones.Model) string {
	if input == nil || input.Name == "" {
		return ""
	}
	return edgezones.NormalizeNilable(&input.Name)
}
