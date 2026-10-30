// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package windows_web_app

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.ResourceWithCustomImporter = Resource{}
	_ sdk.ResourceWithCustomizeDiff  = Resource{}
	_ sdk.ResourceWithStateMigration = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_windows_web_app"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateAppServiceID
}
