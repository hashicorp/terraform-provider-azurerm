// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_api_version_set

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiversionset"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

var _ sdk.ResourceWithCustomizeDiff = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_api_management_workspace_api_version_set"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return apiversionset.ValidateWorkspaceApiVersionSetID
}
