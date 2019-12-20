package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMKeyVault_basic(t *testing.T) {
	dataSourceName := "data.azurerm_key_vault.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccDataSourceAzureRMKeyVault_basic(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "tenant_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "sku_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "access_policy.0.tenant_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "access_policy.0.object_id"),
					resource.TestCheckResourceAttr(dataSourceName, "access_policy.0.key_permissions.0", "create"),
					resource.TestCheckResourceAttr(dataSourceName, "access_policy.0.secret_permissions.0", "set"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKeyVault_basicClassic(t *testing.T) {
	dataSourceName := "data.azurerm_key_vault.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccDataSourceAzureRMKeyVault_basic(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "tenant_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "sku.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "access_policy.0.tenant_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "access_policy.0.object_id"),
					resource.TestCheckResourceAttr(dataSourceName, "access_policy.0.key_permissions.0", "create"),
					resource.TestCheckResourceAttr(dataSourceName, "access_policy.0.secret_permissions.0", "set"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKeyVault_complete(t *testing.T) {
	dataSourceName := "data.azurerm_key_vault.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccDataSourceAzureRMKeyVault_complete(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "tenant_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "sku.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "access_policy.0.tenant_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "access_policy.0.object_id"),
					resource.TestCheckResourceAttr(dataSourceName, "access_policy.0.key_permissions.0", "get"),
					resource.TestCheckResourceAttr(dataSourceName, "access_policy.0.secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.environment", "Production"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKeyVault_networkAcls(t *testing.T) {
	dataSourceName := "data.azurerm_key_vault.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccDataSourceAzureRMKeyVault_networkAcls(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "tenant_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "sku.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "access_policy.0.tenant_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "access_policy.0.object_id"),
					resource.TestCheckResourceAttr(dataSourceName, "access_policy.0.key_permissions.0", "create"),
					resource.TestCheckResourceAttr(dataSourceName, "access_policy.0.secret_permissions.0", "set"),
					resource.TestCheckResourceAttr(dataSourceName, "network_acls.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "network_acls.0.default_action", "Allow"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKeyVault_basic(rInt int, location string) string {
	r := testAccAzureRMKeyVault_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_key_vault" "test" {
  name                = "${azurerm_key_vault.test.name}"
  resource_group_name = "${azurerm_key_vault.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKeyVault_complete(rInt int, location string) string {
	r := testAccAzureRMKeyVault_complete(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_key_vault" "test" {
  name                = "${azurerm_key_vault.test.name}"
  resource_group_name = "${azurerm_key_vault.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKeyVault_networkAcls(rInt int, location string) string {
	r := testAccAzureRMKeyVault_networkAclsUpdated(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_key_vault" "test" {
  name                = "${azurerm_key_vault.test.name}"
  resource_group_name = "${azurerm_key_vault.test.resource_group_name}"
}
`, r)
}
