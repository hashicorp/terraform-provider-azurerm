package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMySQLConfiguration_characterSetServer(t *testing.T) {
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMMySQLConfiguration_characterSetServer(ri, location)
	serverOnlyConfig := testAccAzureRMMySQLConfiguration_empty(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMySQLConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLConfigurationValue("azurerm_mysql_configuration.test", "hebrew"),
				),
			},
			{
				Config: serverOnlyConfig,
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckAzureRMMySQLConfigurationValueReset(ri, "character_set_server"),
				),
			},
		},
	})
}

func TestAccAzureRMMySQLConfiguration_interactiveTimeout(t *testing.T) {
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMMySQLConfiguration_interactiveTimeout(ri, location)
	serverOnlyConfig := testAccAzureRMMySQLConfiguration_empty(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMySQLConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLConfigurationValue("azurerm_mysql_configuration.test", "30"),
				),
			},
			{
				Config: serverOnlyConfig,
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckAzureRMMySQLConfigurationValueReset(ri, "interactive_timeout"),
				),
			},
		},
	})
}

func TestAccAzureRMMySQLConfiguration_logSlowAdminStatements(t *testing.T) {
	resourceName := "azurerm_mysql_configuration.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMMySQLConfiguration_logSlowAdminStatements(ri, location)
	serverOnlyConfig := testAccAzureRMMySQLConfiguration_empty(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMySQLConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLConfigurationValue(resourceName, "on"),
				),
			},
			{
				Config: serverOnlyConfig,
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckAzureRMMySQLConfigurationValueReset(ri, "log_slow_admin_statements"),
				),
			},
		},
	})
}

func testCheckAzureRMMySQLConfigurationValue(resourceName string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for MySQL Configuration: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).mysqlConfigurationsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MySQL Configuration %q (server %q resource group: %q) does not exist", name, serverName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on mysqlConfigurationsClient: %+v", err)
		}

		if *resp.Value != value {
			return fmt.Errorf("MySQL Configuration wasn't set. Expected '%s' - got '%s': \n%+v", value, *resp.Value, resp)
		}

		return nil
	}
}

func testCheckAzureRMMySQLConfigurationValueReset(rInt int, configurationName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		resourceGroup := fmt.Sprintf("acctestRG-%d", rInt)
		serverName := fmt.Sprintf("acctestmysqlsvr-%d", rInt)

		client := testAccProvider.Meta().(*ArmClient).mysqlConfigurationsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, configurationName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MySQL Configuration %q (server %q resource group: %q) does not exist", configurationName, serverName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on mysqlConfigurationsClient: %+v", err)
		}

		actualValue := *resp.Value
		defaultValue := *resp.DefaultValue

		if defaultValue != actualValue {
			return fmt.Errorf("MySQL Configuration wasn't set to the default value. Expected '%s' - got '%s': \n%+v", defaultValue, actualValue, resp)
		}

		return nil
	}
}

func testCheckAzureRMMySQLConfigurationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).mysqlConfigurationsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mysql_configuration" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, serverName, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}
	}

	return nil
}

func testAccAzureRMMySQLConfiguration_characterSetServer(rInt int, location string) string {
	return testAccAzureRMMySQLConfiguration_template(rInt, location, "character_set_server", "hebrew")
}

func testAccAzureRMMySQLConfiguration_interactiveTimeout(rInt int, location string) string {
	return testAccAzureRMMySQLConfiguration_template(rInt, location, "interactive_timeout", "30")
}

func testAccAzureRMMySQLConfiguration_logSlowAdminStatements(rInt int, location string) string {
	return testAccAzureRMMySQLConfiguration_template(rInt, location, "log_slow_admin_statements", "on")
}

func testAccAzureRMMySQLConfiguration_template(rInt int, location string, name string, value string) string {
	server := testAccAzureRMMySQLConfiguration_empty(rInt, location)
	config := fmt.Sprintf(`
resource "azurerm_mysql_configuration" "test" {
  name                = "%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mysql_server.test.name}"
  value               = "%s"
}
`, name, value)
	return server + config
}

func testAccAzureRMMySQLConfiguration_empty(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "B_Gen4_2"
    capacity = 2
    tier     = "Basic"
    family   = "Gen4"
  }

  storage_profile {
    storage_mb = 51200
    backup_retention_days = 7
    geo_redundant_backup = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "5.7"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}
