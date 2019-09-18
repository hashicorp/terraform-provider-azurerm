package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSqlAdministrator_basic(t *testing.T) {
	resourceName := "azurerm_sql_active_directory_administrator.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlAdministratorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlAdministrator_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlAdministratorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "login", "sqladmin"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMSqlAdministrator_withUpdates(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlAdministratorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "login", "sqladmin2"),
				),
			},
		},
	})
}
func TestAccAzureRMSqlAdministrator_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	resourceName := "azurerm_sql_active_directory_administrator.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlAdministratorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlAdministrator_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlAdministratorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "login", "sqladmin"),
				),
			},
			{
				Config:      testAccAzureRMSqlAdministrator_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_sql_active_directory_administrator"),
			},
		},
	})
}

func TestAccAzureRMSqlAdministrator_disappears(t *testing.T) {
	resourceName := "azurerm_sql_active_directory_administrator.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMSqlAdministrator_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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

func testCheckAzureRMSqlAdministratorExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]

		client := testAccProvider.Meta().(*ArmClient).Sql.ServerAzureADAdministratorsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		_, err := client.Get(ctx, resourceGroup, serverName)
		return err
	}
}

func testCheckAzureRMSqlAdministratorDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]

		client := testAccProvider.Meta().(*ArmClient).Sql.ServerAzureADAdministratorsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if _, err := client.Delete(ctx, resourceGroup, serverName); err != nil {
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

		client := testAccProvider.Meta().(*ArmClient).Sql.ServerAzureADAdministratorsClient
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
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_active_directory_administrator" "test" {
  server_name         = "${azurerm_sql_server.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  login               = "sqladmin"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"
  object_id           = "${data.azurerm_client_config.current.client_id}"
}
`, rInt, location, rInt)
}

func testAccAzureRMSqlAdministrator_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sql_active_directory_administrator" "import" {
  server_name         = "${azurerm_sql_active_directory_administrator.test.server_name}"
  resource_group_name = "${azurerm_sql_active_directory_administrator.test.resource_group_name}"
  login               = "${azurerm_sql_active_directory_administrator.test.login}"
  tenant_id           = "${azurerm_sql_active_directory_administrator.test.tenant_id}"
  object_id           = "${azurerm_sql_active_directory_administrator.test.object_id}"
}
`, testAccAzureRMSqlAdministrator_basic(rInt, location))
}

func testAccAzureRMSqlAdministrator_withUpdates(rInt int, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_active_directory_administrator" "test" {
  server_name         = "${azurerm_sql_server.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  login               = "sqladmin2"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"
  object_id           = "${data.azurerm_client_config.current.client_id}"
}
`, rInt, location, rInt)
}
