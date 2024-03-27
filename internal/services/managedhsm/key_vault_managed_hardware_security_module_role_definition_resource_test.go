// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedhsm_test

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KeyVaultMHSMRoleDefinitionResource struct{}

// real test nested in TestAccKeyVaultManagedHardwareSecurityModule, only provide Exists logic here
func (k KeyVaultMHSMRoleDefinitionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	baseURL := state.Attributes["vault_base_url"]
	id, err := parse.RoleNestedItemID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.ManagedHSMs.DataPlaneRoleDefinitionsClient.Get(ctx, baseURL, "/", id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Type %s: %+v", id, err)
	}
	return utils.Bool(resp.RoleDefinitionProperties != nil), nil
}

func (k KeyVaultMHSMRoleDefinitionResource) withRoleDefinition(data acceptance.TestData) string {
	hsm := KeyVaultManagedHardwareSecurityModuleResource{}.download(data, 3)
	return fmt.Sprintf(`


%s

locals {
  roleTestName = "c9562a52-2bd9-2671-3d89-cea5b4798a6b"
}

resource "azurerm_key_vault_managed_hardware_security_module_role_definition" "test" {
  name           = local.roleTestName
  vault_base_url = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  description    = "desc foo"
  permission {
    data_actions = [
      "Microsoft.KeyVault/managedHsm/keys/read/action",
      "Microsoft.KeyVault/managedHsm/keys/write/action",
      "Microsoft.KeyVault/managedHsm/keys/encrypt/action",
      "Microsoft.KeyVault/managedHsm/keys/create",
      "Microsoft.KeyVault/managedHsm/keys/delete",
    ]
    not_data_actions = [
      "Microsoft.KeyVault/managedHsm/roleAssignments/read/action",
    ]
  }
}
`, hsm)
}

func (k KeyVaultMHSMRoleDefinitionResource) withRoleDefinitionUpdate(data acceptance.TestData) string {
	hsm := KeyVaultManagedHardwareSecurityModuleResource{}.download(data, 3)
	return fmt.Sprintf(`


%s

locals {
  roleTestName = "c9562a52-2bd9-2671-3d89-cea5b4798a6b"
}

resource "azurerm_key_vault_managed_hardware_security_module_role_definition" "test" {
  name           = local.roleTestName
  vault_base_url = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  description    = "desc foo2"
  permission {
    data_actions = [
      "Microsoft.KeyVault/managedHsm/keys/read/action",
      "Microsoft.KeyVault/managedHsm/keys/write/action",
      "Microsoft.KeyVault/managedHsm/keys/encrypt/action",
      "Microsoft.KeyVault/managedHsm/keys/create",
    ]
    not_data_actions = [
      "Microsoft.KeyVault/managedHsm/roleAssignments/read/action",
      "Microsoft.KeyVault/managedHsm/keys/delete",
    ]
  }
}
`, hsm)
}
