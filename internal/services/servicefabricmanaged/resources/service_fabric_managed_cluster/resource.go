// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package service_fabric_managed_cluster

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2024-04-01/managedcluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithCustomizeDiff = Resource{}

func (k Resource) ResourceType() string {
	return "azurerm_service_fabric_managed_cluster"
}

func (k Resource) Description() string {
	return "Manages a Resource Group."
}

func (k Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return managedcluster.ValidateManagedClusterID
}
