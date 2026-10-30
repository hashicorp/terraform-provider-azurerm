// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_service_environment_v3

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// (@jackofallops) - Important property missing from the SDK / Swagger that will need to be added later: `upgrade_preference` https://docs.microsoft.com/azure/app-service/environment/using#upgrade-preference

type Resource struct{}

var (
	_ sdk.Resource                  = Resource{}
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_app_service_environment_v3"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateAppServiceEnvironmentID
}
