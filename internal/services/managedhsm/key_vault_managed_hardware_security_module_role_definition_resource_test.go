// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedhsm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KeyVaultMHSMRoleDefinitionResource struct{}

func testAccKeyVaultManagedHardwareSecurityModuleRoleDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module_role_definition", "test")
	r := KeyVaultMHSMRoleDefinitionResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccKeyVaultManagedHardwareSecurityModuleRoleDefinition_legacyWithUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module_role_definition", "test")
	r := KeyVaultMHSMRoleDefinitionResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.legacy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.legacyUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// real test nested in TestAccKeyVaultManagedHardwareSecurityModule, only provide Exists logic here
func (r KeyVaultMHSMRoleDefinitionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	domainSuffix, ok := client.Account.Environment.ManagedHSM.DomainSuffix()
	if !ok {
		return nil, fmt.Errorf("this Environment doesn't specify the Domain Suffix for Managed HSM")
	}
	id, err := parse.ManagedHSMDataPlaneRoleDefinitionID(state.ID, domainSuffix)
	if err != nil {
		return nil, err
	}

	resp, err := client.ManagedHSMs.DataPlaneRoleDefinitionsClient.Get(ctx, id.BaseURI(), id.Scope, id.RoleDefinitionName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Type %s: %+v", id, err)
	}
	return utils.Bool(resp.RoleDefinitionProperties != nil), nil
}

func (r KeyVaultMHSMRoleDefinitionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

locals {
  roleTestName = "c9562a52-2bd9-2671-3d89-cea5b4798a6b"
}

resource "azurerm_key_vault_managed_hardware_security_module_role_definition" "test" {
  name           = local.roleTestName
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  role_name      = "myRole%s"
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
`, r.template(data), data.RandomString)
}

func (r KeyVaultMHSMRoleDefinitionResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

locals {
  roleTestName = "c9562a52-2bd9-2671-3d89-cea5b4798a6b"
}

resource "azurerm_key_vault_managed_hardware_security_module_role_definition" "test" {
  name           = local.roleTestName
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  role_name      = "myRole%s"
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
`, r.template(data), data.RandomString)
}

func (r KeyVaultMHSMRoleDefinitionResource) legacy(data acceptance.TestData) string {
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
`, r.template(data))
}

func (r KeyVaultMHSMRoleDefinitionResource) legacyUpdate(data acceptance.TestData) string {
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
`, r.template(data))
}

func (r KeyVaultMHSMRoleDefinitionResource) template(data acceptance.TestData) string {
	return KeyVaultManagedHardwareSecurityModuleResource{}.download(data, 3)
}
