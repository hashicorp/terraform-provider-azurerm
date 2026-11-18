// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action_version

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "name" -compare-values "subscription_id:cdn_frontdoor_edge_action_id,resource_group_name:cdn_frontdoor_edge_action_id,edge_action_name:cdn_frontdoor_edge_action_id"

type Resource struct{}

var (
	_ sdk.ResourceWithIdentity = Resource{}
	_ sdk.ResourceWithUpdate   = Resource{}
)

func (Resource) ResourceType() string {
	return "azurerm_cdn_frontdoor_edge_action_version"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return edgeactionversions.ValidateVersionID
}

func (Resource) Identity() resourceids.ResourceId {
	return &edgeactionversions.VersionId{}
}
