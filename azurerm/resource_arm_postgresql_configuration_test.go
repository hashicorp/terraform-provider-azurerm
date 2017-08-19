package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMPostgreSQLConfiguration_backslashQuote(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLConfiguration_backslashQuote(ri)
	serverOnlyConfig := testAccAzureRMPostgreSQLConfiguration_empty(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLConfigurationReset,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLConfigurationValue("azurerm_postgresql_configuration.test", "on"),
				),
			},
			{
				Config: serverOnlyConfig,
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckAzureRMPostgreSQLConfigurationValueReset(ri, "backslash_quote", "safe_encoding"),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLConfiguration_clientMinMessages(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLConfiguration_clientMinMessages(ri)
	serverOnlyConfig := testAccAzureRMPostgreSQLConfiguration_empty(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLConfigurationReset,
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
					testCheckAzureRMPostgreSQLConfigurationValueReset(ri, "client_min_messages", "NOTICE"),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLConfiguration_deadlockTimeout(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLConfiguration_deadlockTimeout(ri)
	serverOnlyConfig := testAccAzureRMPostgreSQLConfiguration_empty(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLConfigurationReset,
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
					testCheckAzureRMPostgreSQLConfigurationValueReset(ri, "deadlock_timeout", "1000"),
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

		resp, err := client.Get(resourceGroup, serverName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on postgresqlConfigurationsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: PostgreSQL Configuration %q (server %q resource group: %q) does not exist", name, serverName, resourceGroup)
		}

		if *resp.Value != value {
			return fmt.Errorf("PostgreSQL Configuration wasn't set. Expected '%s' - got '%s': \n%+v", value, *resp.Value, resp)
		}

		return nil
	}
}

func testCheckAzureRMPostgreSQLConfigurationValueReset(rInt int, configurationName string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		resourceGroup := fmt.Sprintf("acctestRG-%d", rInt)
		serverName := fmt.Sprintf("acctestpsqlsvr-%d", rInt)

		client := testAccProvider.Meta().(*ArmClient).postgresqlConfigurationsClient

		resp, err := client.Get(resourceGroup, serverName, configurationName)
		if err != nil {
			return fmt.Errorf("Bad: Get on postgresqlConfigurationsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: PostgreSQL Configuration %q (server %q resource group: %q) does not exist", configurationName, serverName, resourceGroup)
		}

		actualValue := *resp.Value
		defaultValue := *resp.DefaultValue

		if defaultValue != actualValue {
			return fmt.Errorf("PostgreSQL Configuration wasn't set to the default value. Expected '%s' - got '%s': \n%+v", value, *resp.Value, resp)
		}

		return nil
	}
}

func testCheckAzureRMPostgreSQLConfigurationReset(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).postgresqlConfigurationsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_postgresql_configuration" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(resourceGroup, serverName, name)

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return err
			}

			return nil
		}
	}

	return nil
}

func testAccAzureRMPostgreSQLConfiguration_backslashQuote(rInt int) string {
	return testAccAzureRMPostgreSQLConfiguration_template(rInt, "backslash_quote", "on")
}

func testAccAzureRMPostgreSQLConfiguration_clientMinMessages(rInt int) string {
	return testAccAzureRMPostgreSQLConfiguration_template(rInt, "client_min_messages", "DEBUG5")
}

func testAccAzureRMPostgreSQLConfiguration_deadlockTimeout(rInt int) string {
	return testAccAzureRMPostgreSQLConfiguration_template(rInt, "deadlock_timeout", "5000")
}

func testAccAzureRMPostgreSQLConfiguration_template(rInt int, name string, value string) string {
	server := testAccAzureRMPostgreSQLConfiguration_empty(rInt)
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

func testAccAzureRMPostgreSQLConfiguration_empty(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}
resource "azurerm_postgresql_server" "test" {
  name = "acctestpsqlsvr-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "PGSQLB50"
    capacity = 50
    tier = "Basic"
  }

  administrator_login = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version = "9.6"
  storage_mb = 51200
  ssl_enforcement = "Enabled"
}
`, rInt, rInt)
}
