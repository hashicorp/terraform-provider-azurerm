package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMariaDbConfiguration_characterSetServer(t *testing.T) {
	resourceName := "azurerm_mariadb_configuration.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMariaDbConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMariaDbConfiguration_characterSetServer(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbConfigurationValue(resourceName, "hebrew"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMMariaDbConfiguration_empty(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckAzureRMMariaDbConfigurationValueReset(ri, "character_set_server"),
				),
			},
		},
	})
}

func TestAccAzureRMMariaDbConfiguration_interactiveTimeout(t *testing.T) {
	resourceName := "azurerm_mariadb_configuration.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMariaDbConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMariaDbConfiguration_interactiveTimeout(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbConfigurationValue(resourceName, "30"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMMariaDbConfiguration_empty(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckAzureRMMariaDbConfigurationValueReset(ri, "interactive_timeout"),
				),
			},
		},
	})
}

func TestAccAzureRMMariaDbConfiguration_logSlowAdminStatements(t *testing.T) {
	resourceName := "azurerm_mariadb_configuration.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMariaDbConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMariaDbConfiguration_logSlowAdminStatements(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbConfigurationValue(resourceName, "on"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMMariaDbConfiguration_empty(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					// "delete" resets back to the default value
					testCheckAzureRMMariaDbConfigurationValueReset(ri, "log_slow_admin_statements"),
				),
			},
		},
	})
}

func testCheckAzureRMMariaDbConfigurationValue(resourceName string, value string) resource.TestCheckFunc {
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
			return fmt.Errorf("Bad: no resource group found in state for MariaDb Configuration: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).mariadb.ConfigurationsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MariaDb Configuration %q (server %q resource group: %q) does not exist", name, serverName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on mariadbConfigurationsClient: %+v", err)
		}

		if *resp.Value != value {
			return fmt.Errorf("MariaDb Configuration wasn't set. Expected '%s' - got '%s': \n%+v", value, *resp.Value, resp)
		}

		return nil
	}
}

func testCheckAzureRMMariaDbConfigurationValueReset(rInt int, configurationName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		resourceGroup := fmt.Sprintf("acctestRG-%d", rInt)
		serverName := fmt.Sprintf("acctestmariadbsvr-%d", rInt)

		client := testAccProvider.Meta().(*ArmClient).mariadb.ConfigurationsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, configurationName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MariaDb Configuration %q (server %q resource group: %q) does not exist", configurationName, serverName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on mariadbConfigurationsClient: %+v", err)
		}

		actualValue := *resp.Value
		defaultValue := *resp.DefaultValue

		if defaultValue != actualValue {
			return fmt.Errorf("MariaDb Configuration wasn't set to the default value. Expected '%s' - got '%s': \n%+v", defaultValue, actualValue, resp)
		}

		return nil
	}
}

func testCheckAzureRMMariaDbConfigurationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).mariadb.ConfigurationsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mariadb_configuration" {
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

func testAccAzureRMMariaDbConfiguration_characterSetServer(rInt int, location string) string {
	return testAccAzureRMMariaDbConfiguration_template(rInt, location, "character_set_server", "hebrew")
}

func testAccAzureRMMariaDbConfiguration_interactiveTimeout(rInt int, location string) string {
	return testAccAzureRMMariaDbConfiguration_template(rInt, location, "interactive_timeout", "30")
}

func testAccAzureRMMariaDbConfiguration_logSlowAdminStatements(rInt int, location string) string {
	return testAccAzureRMMariaDbConfiguration_template(rInt, location, "log_slow_admin_statements", "on")
}

func testAccAzureRMMariaDbConfiguration_template(rInt int, location string, name string, value string) string {
	server := testAccAzureRMMariaDbConfiguration_empty(rInt, location)
	config := fmt.Sprintf(`
resource "azurerm_mariadb_configuration" "test" {
  name                = "%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mariadb_server.test.name}"
  value               = "%s"
}
`, name, value)
	return server + config
}

func testAccAzureRMMariaDbConfiguration_empty(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "GP_Gen5_2"
    capacity = 2
    tier     = "GeneralPurpose"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}
