package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMKeyVaultAccessPolicy(t *testing.T) {
	resourceName := "azurerm_key_vault.test"
	ri := acctest.RandInt()
	config := testAccAzureRMKeyVaultAccessPolicy(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "access_policy.0.key_permissions.0", "create"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVaultAccessPolicy_update(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_key_vault.test"
	preConfig := testAccAzureRMKeyVaultAccessPolicy(ri, testLocation())
	postConfig := testAccAzureRMKeyVaultAccessPolicy_update(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "access_policy.0.key_permissions.0", "create"),
					resource.TestCheckResourceAttr(resourceName, "access_policy.0.secret_permissions.0", "set"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "Production"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "access_policy.0.key_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "access_policy.0.secret_permissions.0", "get"),
				),
			},
		},
	})
}

func testCheckAzureRMKeyVaultAccessPolicyDestroy(s *terraform.State) error {
	// No op

	return nil
}

func testAccAzureRMKeyVaultAccessPolicy(rInt int, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "vault%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "premium"
  }
}

resource access_policy {
	tenant_id = "${data.azurerm_client_config.current.tenant_id}"
	object_id = "${data.azurerm_client_config.current.client_id}"

	key_permissions = [
		"create",
	]

	secret_permissions = [
		"set",
	]
}
`, rInt, location, rInt)
}

func testAccAzureRMKeyVaultAccessPolicy_update(rInt int, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "vault%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "premium"
  }
}

resource access_policy {
	tenant_id = "${data.azurerm_client_config.current.tenant_id}"
	object_id = "${data.azurerm_client_config.current.client_id}"

	key_permissions = [
		"get",
	]

	secret_permissions = [
		"get",
	]
}
`, rInt, location, rInt)
}
