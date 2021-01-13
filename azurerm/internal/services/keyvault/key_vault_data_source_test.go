package keyvault_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type KeyVaultDataSource struct {
}

func TestAccDataSourceKeyVault_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault", "test")
	r := KeyVaultDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tenant_id").Exists(),
				check.That(data.ResourceName).Key("sku_name").Exists(),
				check.That(data.ResourceName).Key("access_policy.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("access_policy.0.object_id").Exists(),
				check.That(data.ResourceName).Key("access_policy.0.key_permissions.0").HasValue("create"),
				check.That(data.ResourceName).Key("access_policy.0.secret_permissions.0").HasValue("set"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceKeyVault_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault", "test")
	r := KeyVaultDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tenant_id").Exists(),
				check.That(data.ResourceName).Key("sku_name").Exists(),
				check.That(data.ResourceName).Key("access_policy.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("access_policy.0.object_id").Exists(),
				check.That(data.ResourceName).Key("access_policy.0.key_permissions.0").HasValue("get"),
				check.That(data.ResourceName).Key("access_policy.0.secret_permissions.0").HasValue("get"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
			),
		},
	})
}

func TestAccDataSourceKeyVault_networkAcls(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault", "test")
	r := KeyVaultDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.networkAcls(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tenant_id").Exists(),
				check.That(data.ResourceName).Key("sku_name").Exists(),
				check.That(data.ResourceName).Key("access_policy.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("access_policy.0.object_id").Exists(),
				check.That(data.ResourceName).Key("access_policy.0.key_permissions.0").HasValue("create"),
				check.That(data.ResourceName).Key("access_policy.0.secret_permissions.0").HasValue("set"),
				check.That(data.ResourceName).Key("network_acls.#").HasValue("1"),
				check.That(data.ResourceName).Key("network_acls.0.default_action").HasValue("Allow"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceKeyVault_softDelete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault", "test")
	r := KeyVaultDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.enableSoftDelete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("soft_delete_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("purge_protection_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("sku_name").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func (KeyVaultDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault" "test" {
  name                = azurerm_key_vault.test.name
  resource_group_name = azurerm_key_vault.test.resource_group_name
}
`, KeyVaultResource{}.basic(data))
}

func (KeyVaultDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault" "test" {
  name                = azurerm_key_vault.test.name
  resource_group_name = azurerm_key_vault.test.resource_group_name
}
`, KeyVaultResource{}.complete(data))
}

func (KeyVaultDataSource) networkAcls(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault" "test" {
  name                = azurerm_key_vault.test.name
  resource_group_name = azurerm_key_vault.test.resource_group_name
}
`, KeyVaultResource{}.networkAclsUpdated(data))
}

func (KeyVaultDataSource) enableSoftDelete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault" "test" {
  name                = azurerm_key_vault.test.name
  resource_group_name = azurerm_key_vault.test.resource_group_name
}
`, KeyVaultResource{}.softDelete(data))
}
