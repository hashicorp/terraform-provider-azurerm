// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_named_value

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/namedvalue"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_api_management_workspace_named_value"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return namedvalue.ValidateWorkspaceNamedValueID
}
