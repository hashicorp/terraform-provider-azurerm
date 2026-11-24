// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_gateway_config_connection

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apigatewayconfigconnection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_api_management_gateway_config_connection"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return apigatewayconfigconnection.ValidateConfigConnectionID
}
