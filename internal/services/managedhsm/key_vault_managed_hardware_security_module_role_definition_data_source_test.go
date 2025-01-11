// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedhsm_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

// TODO: check the UUIDs

type KeyVaultManagedHardwareSecurityModuleRoleDefinitionDataSource struct{}

func testAccDataSourceKeyVaultManagedHardwareSecurityModuleRoleDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_managed_hardware_security_module_role_definition", "test")
	r := KeyVaultManagedHardwareSecurityModuleRoleDefinitionDataSource{}

	data.DataSourceTestInSequence(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("managed_hsm_id").Exists(),
				check.That(data.ResourceName).Key("role_name").HasValue(fmt.Sprintf("myRole%s", data.RandomString)),
				check.That(data.ResourceName).Key("description").HasValue("desc foo"),
				check.That(data.ResourceName).Key("permission.%").HasValue("1"),
				check.That(data.ResourceName).Key("permission.0.data_actions.%").HasValue("5"),
				check.That(data.ResourceName).Key("permission.0.not_data_actions.%").HasValue("1"),
			),
		},
	})
}

func (KeyVaultManagedHardwareSecurityModuleRoleDefinitionDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "test" {
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  name           = "21dbd100-6940-42c2-9190-5d6cb909625b"
}
`, KeyVaultMHSMRoleDefinitionResource{}.basic(data))
}
