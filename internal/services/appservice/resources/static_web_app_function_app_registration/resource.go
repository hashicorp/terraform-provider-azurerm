// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package static_web_app_function_app_registration

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_static_web_app_function_app_registration"
}

func (r Resource) Description() string {
	return "Manages an App Service Static Web App Function App Registration."
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return staticsites.ValidateUserProvidedFunctionAppID
}
