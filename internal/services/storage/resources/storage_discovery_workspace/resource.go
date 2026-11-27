// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_discovery_workspace

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagediscovery/2025-09-01/storagediscoveryworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

var (
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithIdentity      = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
)

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_storage_discovery_workspace"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return storagediscoveryworkspaces.ValidateProviderStorageDiscoveryWorkspaceID
}

func (r Resource) Identity() resourceids.ResourceId {
	return &storagediscoveryworkspaces.ProviderStorageDiscoveryWorkspaceId{}
}
