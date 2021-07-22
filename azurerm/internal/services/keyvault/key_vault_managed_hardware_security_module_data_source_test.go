package keyvault_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type KeyVaultManagedHardwareSecurityModuleDataSource struct {
}

func testAccDataSourceKeyVaultManagedHardwareSecurityModule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_managed_hardware_security_module", "test")
	r := KeyVaultManagedHardwareSecurityModuleDataSource{}

	data.DataSourceTestInSequence(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku_name").HasValue("Standard_B1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func (KeyVaultManagedHardwareSecurityModuleDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_managed_hardware_security_module" "test" {
  name                = azurerm_key_vault_managed_hardware_security_module.test.name
  resource_group_name = azurerm_key_vault_managed_hardware_security_module.test.resource_group_name
}
`, KeyVaultManagedHardwareSecurityModuleResource{}.basic(data))
}
