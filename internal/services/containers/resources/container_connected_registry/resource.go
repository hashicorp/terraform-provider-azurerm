// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_connected_registry

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/connectedregistries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_container_connected_registry"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return connectedregistries.ValidateConnectedRegistryID
}
