// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_role_definition

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/resources/key_vault_managed_hardware_security_module_role_definition/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.ResourceWithStateMigration = Resource{}
	_ sdk.ResourceWithUpdate         = Resource{}
)

func (r Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.ManagedHSMRoleDefinitionV0ToV1{},
		},
	}
}

func (r Resource) ResourceType() string {
	return "azurerm_key_vault_managed_hardware_security_module_role_definition"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ManagedHSMDataPlaneRoleDefinitionID
}
