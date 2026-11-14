// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nginx_api_key

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxapikey"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = (*Resource)(nil)

func (m Resource) ResourceType() string {
	return "azurerm_nginx_api_key"
}

func (m Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return nginxapikey.ValidateApiKeyID
}
