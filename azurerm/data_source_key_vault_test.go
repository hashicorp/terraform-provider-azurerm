package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMKeyVault_basic(t *testing.T) {
	dataSourceName := "data.azurerm_key_vault.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccDataSourceAzureRMKeyVault_basic(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccDataSourceAzureRMKeyVault_complete(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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

func testAccDataSourceAzureRMKeyVault_basic(rInt int, location string) string {
	resource := testAccAzureRMKeyVault_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_key_vault" "test" {
  name                = "${azurerm_key_vault.test.name}"
  resource_group_name = "${azurerm_key_vault.test.resource_group_name}"
}
`, resource)
}

func testAccDataSourceAzureRMKeyVault_complete(rInt int, location string) string {
	resource := testAccAzureRMKeyVault_complete(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_key_vault" "test" {
  name                = "${azurerm_key_vault.test.name}"
  resource_group_name = "${azurerm_key_vault.test.resource_group_name}"
}
`, resource)
}
