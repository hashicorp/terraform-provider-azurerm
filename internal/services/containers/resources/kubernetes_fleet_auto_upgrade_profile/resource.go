// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_auto_upgrade_profile

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/autoupgradeprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "name" -compare-values "subscription_id:kubernetes_fleet_manager_id,resource_group_name:kubernetes_fleet_manager_id,fleet_name:kubernetes_fleet_manager_id"

var (
	_ sdk.ResourceWithIdentity = Resource{}
	_ sdk.ResourceWithUpdate   = Resource{}
)

type Resource struct{}

func (r Resource) Identity() resourceids.ResourceId {
	return &autoupgradeprofiles.AutoUpgradeProfileId{}
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autoupgradeprofiles.ValidateAutoUpgradeProfileID
}

func (r Resource) ResourceType() string {
	return "azurerm_kubernetes_fleet_auto_upgrade_profile"
}
