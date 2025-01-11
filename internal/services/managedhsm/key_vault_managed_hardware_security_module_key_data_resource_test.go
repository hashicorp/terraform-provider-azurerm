// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedhsm_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type KeyVaultMHSMKeyTestDataSource struct{}

func testAccKeyVaultMHSMKeyDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module_key", "test")
	dataSourceName := "data." + data.ResourceName
	r := KeyVaultMHSMKeyTestDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(dataSourceName).Key("version").Exists(),
				check.That(dataSourceName).Key("key_type").Exists(),
			),
		},
	})
}

func (KeyVaultMHSMKeyTestDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_managed_hardware_security_module_key" "test" {
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  name           = azurerm_key_vault_managed_hardware_security_module_key.test.name
}
`, KeyVaultMHSMKeyTestResource{}.basic(data))
}
