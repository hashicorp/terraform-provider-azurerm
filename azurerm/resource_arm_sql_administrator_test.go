package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSqlAdministrator_basic(t *testing.T) {
	resourceName := "azurerm_sql_active_directory_administrator.test"
	ri := acctest.RandInt()
	preConfig := testAccAzureRMSqlAdministrator_basic(ri, testLocation())
	postConfig := testAccAzureRMSqlAdministrator_withUpdates(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlAdministratorDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlAdministratorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "login", "sqladmin"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlAdministratorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "login", "sqladmin2"),
				),
			},
		},
	})
}

func TestAccAzureRMSqlAdministrator_disappears(t *testing.T) {
	resourceName := "azurerm_sql_active_directory_administrator.test"
	ri := acctest.RandInt()
	config := testAccAzureRMSqlAdministrator_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlAdministratorDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlAdministratorExists(resourceName),
					testCheckAzureRMSqlAdministratorDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMSqlAdministratorExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]

		client := testAccProvider.Meta().(*ArmClient).sqlServerAzureADAdministratorsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		_, err := client.Get(ctx, resourceGroup, serverName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testCheckAzureRMSqlAdministratorDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]

		client := testAccProvider.Meta().(*ArmClient).sqlServerAzureADAdministratorsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		_, err := client.Delete(ctx, resourceGroup, serverName)
		if err != nil {
			return fmt.Errorf("Bad: Delete on sqlAdministratorClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSqlAdministratorDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_sql_active_directory_administrator" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]

		client := testAccProvider.Meta().(*ArmClient).sqlServerAzureADAdministratorsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("SQL Administrator (server %q / resource group %q) still exists: %+v", serverName, resourceGroup, resp)
	}

	return nil
}

func testAccAzureRMSqlAdministrator_basic(rInt int, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name = "acctestsqlserver%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location = "${azurerm_resource_group.test.location}"
  version = "12.0"
  administrator_login = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_active_directory_administrator" "test" {
    server_name = "${azurerm_sql_server.test.name}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    login = "sqladmin"
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.client_id}"
}
`, rInt, location, rInt)
}

func testAccAzureRMSqlAdministrator_withUpdates(rInt int, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name = "acctestsqlserver%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location = "${azurerm_resource_group.test.location}"
  version = "12.0"
  administrator_login = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_active_directory_administrator" "test" {
  server_name = "${azurerm_sql_server.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  login = "sqladmin2"
  tenant_id = "${data.azurerm_client_config.current.tenant_id}"
  object_id = "${data.azurerm_client_config.current.client_id}"
}
`, rInt, location, rInt)
}
