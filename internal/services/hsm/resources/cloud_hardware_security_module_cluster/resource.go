// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cloud_hardware_security_module_cluster

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2025-03-31/cloudhsmclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.ResourceWithUpdate = Resource{}

type Resource struct{}

func (Resource) ResourceType() string {
	return "azurerm_cloud_hardware_security_module_cluster"
}

func (Resource) Description() string {
	return "Manages a Cloud Hardware Security Module Cluster."
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return cloudhsmclusters.ValidateCloudHsmClusterID
}
