package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPostgreSQLConfiguration_backslashQuote(t *testing.T) {
	resourceName := "azurerm_postgresql_configuration.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMPostgreSQLConfiguration_backslashQuote(ri, location)
	serverOnlyConfig := testAccAzureRMPostgreSQLConfiguration_empty(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLConfigurationValue(resourceName, "on"),
				),
			},
			{
				Config: serverOnlyConfig,
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckAzureRMPostgreSQLConfigurationValueReset(ri, "backslash_quote"),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLConfiguration_clientMinMessages(t *testing.T) {
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMPostgreSQLConfiguration_clientMinMessages(ri, location)
	serverOnlyConfig := testAccAzureRMPostgreSQLConfiguration_empty(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLConfigurationValue("azurerm_postgresql_configuration.test", "DEBUG5"),
				),
			},
			{
				Config: serverOnlyConfig,
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckAzureRMPostgreSQLConfigurationValueReset(ri, "client_min_messages"),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLConfiguration_deadlockTimeout(t *testing.T) {
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMPostgreSQLConfiguration_deadlockTimeout(ri, location)
	serverOnlyConfig := testAccAzureRMPostgreSQLConfiguration_empty(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLConfigurationValue("azurerm_postgresql_configuration.test", "5000"),
				),
			},
			{
				Config: serverOnlyConfig,
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckAzureRMPostgreSQLConfigurationValueReset(ri, "deadlock_timeout"),
				),
			},
		},
	})
}

func testCheckAzureRMPostgreSQLConfigurationValue(resourceName string, value string) resource.TestCheckFunc {
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
			return fmt.Errorf("Bad: no resource group found in state for PostgreSQL Configuration: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).postgresqlConfigurationsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: PostgreSQL Configuration %q (server %q resource group: %q) does not exist", name, serverName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on postgresqlConfigurationsClient: %+v", err)
		}

		if *resp.Value != value {
			return fmt.Errorf("PostgreSQL Configuration wasn't set. Expected '%s' - got '%s': \n%+v", value, *resp.Value, resp)
		}

		return nil
	}
}

func testCheckAzureRMPostgreSQLConfigurationValueReset(rInt int, configurationName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		resourceGroup := fmt.Sprintf("acctestRG-%d", rInt)
		serverName := fmt.Sprintf("acctestpsqlsvr-%d", rInt)

		client := testAccProvider.Meta().(*ArmClient).postgresqlConfigurationsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, configurationName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: PostgreSQL Configuration %q (server %q resource group: %q) does not exist", configurationName, serverName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on postgresqlConfigurationsClient: %+v", err)
		}

		actualValue := *resp.Value
		defaultValue := *resp.DefaultValue

		if defaultValue != actualValue {
			return fmt.Errorf("PostgreSQL Configuration wasn't set to the default value. Expected '%s' - got '%s': \n%+v", defaultValue, actualValue, resp)
		}

		return nil
	}
}

func testCheckAzureRMPostgreSQLConfigurationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).postgresqlConfigurationsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_postgresql_configuration" {
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

func testAccAzureRMPostgreSQLConfiguration_backslashQuote(rInt int, location string) string {
	return testAccAzureRMPostgreSQLConfiguration_template(rInt, location, "backslash_quote", "on")
}

func testAccAzureRMPostgreSQLConfiguration_clientMinMessages(rInt int, location string) string {
	return testAccAzureRMPostgreSQLConfiguration_template(rInt, location, "client_min_messages", "DEBUG5")
}

func testAccAzureRMPostgreSQLConfiguration_deadlockTimeout(rInt int, location string) string {
	return testAccAzureRMPostgreSQLConfiguration_template(rInt, location, "deadlock_timeout", "5000")
}

func testAccAzureRMPostgreSQLConfiguration_template(rInt int, location string, name string, value string) string {
	server := testAccAzureRMPostgreSQLConfiguration_empty(rInt, location)
	config := fmt.Sprintf(`
resource "azurerm_postgresql_configuration" "test" {
  name                = "%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_postgresql_server.test.name}"
  value               = "%s"
}
`, name, value)
	return server + config
}

func testAccAzureRMPostgreSQLConfiguration_empty(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctestpsqlsvr-%d"
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
  version                      = "9.6"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}
