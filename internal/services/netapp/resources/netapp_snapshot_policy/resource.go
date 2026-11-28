// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_snapshot_policy

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/snapshotpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.Resource                  = Resource{}
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
)

func (Resource) ResourceType() string {
	return "azurerm_netapp_snapshot_policy"
}

func (Resource) Description() string {
	return "Manages a NetApp Snapshot Policy."
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return snapshotpolicies.ValidateSnapshotPolicyID
}
