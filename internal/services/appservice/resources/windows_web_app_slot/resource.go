// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package windows_web_app_slot

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.ResourceWithCustomizeDiff  = Resource{}
	_ sdk.ResourceWithUpdate         = Resource{}
	_ sdk.ResourceWithStateMigration = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_windows_web_app_slot"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return webapps.ValidateSlotID
}
