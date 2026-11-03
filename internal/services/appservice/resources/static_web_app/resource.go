// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package static_web_app

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

var _ sdk.ResourceWithCustomizeDiff = Resource{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return staticsites.ValidateStaticSiteID
}

func (r Resource) ResourceType() string {
	return "azurerm_static_web_app"
}
