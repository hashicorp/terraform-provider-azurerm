// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package linux_function_app_slot

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/linux_function_app_slot/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

var _ sdk.ResourceWithStateMigration = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_linux_function_app_slot"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return webapps.ValidateSlotID
}

func (r Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.LinuxFunctionAppSlotV0toV1{},
		},
	}
}
