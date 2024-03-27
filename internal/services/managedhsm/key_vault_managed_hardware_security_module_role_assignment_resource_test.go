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

type KeyVaultManagedHSMRoleAssignmentResource struct{}

// real test nested in TestAccKeyVaultManagedHardwareSecurityModule, only provide Exists logic here
func (k KeyVaultManagedHSMRoleAssignmentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.RoleNestedItemID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.ManagedHSMs.DataPlaneRoleAssignmentsClient.Get(ctx, id.VaultBaseUrl, id.Scope, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Type %s: %+v", id, err)
	}
	return utils.Bool(resp.Properties != nil), nil
}

func (k KeyVaultManagedHSMRoleAssignmentResource) withRoleAssignment(data acceptance.TestData) string {
	roleDef := KeyVaultMHSMRoleDefinitionResource{}.withRoleDefinition(data)

	return fmt.Sprintf(`


%s

locals {
  assignmentTestName = "1e243909-064c-6ac3-84e9-1c8bf8d6ad52"
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "test" {
  vault_base_url     = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  name               = local.assignmentTestName
  scope              = "/keys"
  role_definition_id = azurerm_key_vault_managed_hardware_security_module_role_definition.test.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}
`, roleDef)
}

func (k KeyVaultManagedHSMRoleAssignmentResource) withBuiltInRoleAssignment(data acceptance.TestData) string {
	roleDef := k.withRoleAssignment(data)

	return fmt.Sprintf(`


%s

locals {
  assignmentOfficerName = "706c03c7-69ad-33e5-2796-b3380d3a6e1a"
}

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "officer" {
  vault_base_url = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  name           = "515eb02d-2335-4d2d-92f2-b1cbdf9c3778"
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "officer" {
  vault_base_url     = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  name               = local.assignmentOfficerName
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.officer.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}
`, roleDef)
}
