package azurerm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPostgreSQLVirtualNetworkRule_basic(t *testing.T) {
	resourceName := "azurerm_postgresql_virtual_network_rule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLVirtualNetworkRule_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMPostgreSQLVirtualNetworkRule_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_postgresql_virtual_network_rule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLVirtualNetworkRule_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMPostgreSQLVirtualNetworkRule_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_postgresql_virtual_network_rule"),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLVirtualNetworkRule_switchSubnets(t *testing.T) {
	resourceName := "azurerm_postgresql_virtual_network_rule.test"
	ri := tf.AccRandTimeInt()

	preConfig := testAccAzureRMPostgreSQLVirtualNetworkRule_subnetSwitchPre(ri, testLocation())
	postConfig := testAccAzureRMPostgreSQLVirtualNetworkRule_subnetSwitchPost(ri, testLocation())

	// Create regex strings that will ensure that one subnet name exists, but not the other
	preConfigRegex := regexp.MustCompile(fmt.Sprintf("(subnet1%d)$|(subnet[^2]%d)$", ri, ri))  //subnet 1 but not 2
	postConfigRegex := regexp.MustCompile(fmt.Sprintf("(subnet2%d)$|(subnet[^1]%d)$", ri, ri)) //subnet 2 but not 1

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(resourceName),
					resource.TestMatchResourceAttr(resourceName, "subnet_id", preConfigRegex),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(resourceName),
					resource.TestMatchResourceAttr(resourceName, "subnet_id", postConfigRegex),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLVirtualNetworkRule_disappears(t *testing.T) {
	resourceName := "azurerm_postgresql_virtual_network_rule.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPostgreSQLVirtualNetworkRule_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(resourceName),
					testCheckAzureRMPostgreSQLVirtualNetworkRuleDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMPostgreSQLVirtualNetworkRule_multipleSubnets(t *testing.T) {
	resourceName1 := "azurerm_postgresql_virtual_network_rule.rule1"
	resourceName2 := "azurerm_postgresql_virtual_network_rule.rule2"
	resourceName3 := "azurerm_postgresql_virtual_network_rule.rule3"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPostgreSQLVirtualNetworkRule_multipleSubnets(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(resourceName1),
					testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(resourceName2),
					testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(resourceName3),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLVirtualNetworkRule_IgnoreEndpointValid(t *testing.T) {
	resourceName := "azurerm_postgresql_virtual_network_rule.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPostgreSQLVirtualNetworkRule_ignoreEndpointValid(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ignore_missing_vnet_service_endpoint", "true"),
				),
			},
		},
	})
}

func testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		ruleName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).Postgres.VirtualNetworkRulesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, ruleName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: PostgreSQL Virtual Network Rule %q (server %q / resource group %q) was not found", ruleName, serverName, resourceGroup)
			}

			return err
		}

		return nil
	}
}

func testCheckAzureRMPostgreSQLVirtualNetworkRuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_postgresql_virtual_network_rule" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		ruleName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).Postgres.VirtualNetworkRulesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, ruleName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Bad: PostgreSQL Virtual Network Rule %q (server %q / resource group %q) still exists: %+v", ruleName, serverName, resourceGroup, resp)
	}

	return nil
}

func testCheckAzureRMPostgreSQLVirtualNetworkRuleDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		ruleName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).Postgres.VirtualNetworkRulesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		future, err := client.Delete(ctx, resourceGroup, serverName, ruleName)
		if err != nil {
			//If the error is that the resource we want to delete does not exist in the first
			//place (404), then just return with no error.
			if response.WasNotFound(future.Response()) {
				return nil
			}

			return fmt.Errorf("Error deleting PostgreSQL Virtual Network Rule: %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			//Same deal as before. Just in case.
			if response.WasNotFound(future.Response()) {
				return nil
			}

			return fmt.Errorf("Error deleting PostgreSQL Virtual Network Rule: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMPostgreSQLVirtualNetworkRule_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/29"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.7.29.0/29"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_postgresql_server" "test" {
  name                         = "acctestpostgresqlsvr-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement              = "Enabled"

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
}

resource "azurerm_postgresql_virtual_network_rule" "test" {
  name                = "acctestpostgresqlvnetrule%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_postgresql_server.test.name}"
  subnet_id           = "${azurerm_subnet.test.id}"
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMPostgreSQLVirtualNetworkRule_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_virtual_network_rule" "import" {
  name                = "${azurerm_postgresql_virtual_network_rule.test.name}"
  resource_group_name = "${azurerm_postgresql_virtual_network_rule.test.resource_group_name}"
  server_name         = "${azurerm_postgresql_virtual_network_rule.test.server_name}"
  subnet_id           = "${azurerm_postgresql_virtual_network_rule.test.subnet_id}"
}
`, testAccAzureRMPostgreSQLVirtualNetworkRule_basic(rInt, location))
}

func testAccAzureRMPostgreSQLVirtualNetworkRule_subnetSwitchPre(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/24"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test1" {
  name                 = "subnet1%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.7.29.0/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "test2" {
  name                 = "subnet2%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.7.29.128/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_postgresql_server" "test" {
  name                         = "acctestpostgresqlsvr-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement              = "Enabled"

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
}

resource "azurerm_postgresql_virtual_network_rule" "test" {
  name                = "acctestpostgresqlvnetrule%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_postgresql_server.test.name}"
  subnet_id           = "${azurerm_subnet.test1.id}"
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMPostgreSQLVirtualNetworkRule_subnetSwitchPost(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/24"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test1" {
  name                 = "subnet1%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.7.29.0/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "test2" {
  name                 = "subnet2%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.7.29.128/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_postgresql_server" "test" {
  name                         = "acctestpostgresqlsvr-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement              = "Enabled"

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
}

resource "azurerm_postgresql_virtual_network_rule" "test" {
  name                = "acctestpostgresqlvnetrule%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_postgresql_server.test.name}"
  subnet_id           = "${azurerm_subnet.test2.id}"
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMPostgreSQLVirtualNetworkRule_multipleSubnets(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "vnet1" {
  name                = "acctestvnet1%d"
  address_space       = ["10.7.29.0/24"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_virtual_network" "vnet2" {
  name                = "acctestvnet2%d"
  address_space       = ["10.1.29.0/29"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "vnet1_subnet1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.vnet1.name}"
  address_prefix       = "10.7.29.0/29"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "vnet1_subnet2" {
  name                 = "acctestsubnet2%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.vnet1.name}"
  address_prefix       = "10.7.29.128/29"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "vnet2_subnet1" {
  name                 = "acctestsubnet3%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.vnet2.name}"
  address_prefix       = "10.1.29.0/29"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_postgresql_server" "test" {
  name                         = "acctestpostgresqlsvr-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement              = "Enabled"

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
}

resource "azurerm_postgresql_virtual_network_rule" "rule1" {
  name                                 = "acctestpostgresqlvnetrule1%d"
  resource_group_name                  = "${azurerm_resource_group.test.name}"
  server_name                          = "${azurerm_postgresql_server.test.name}"
  subnet_id                            = "${azurerm_subnet.vnet1_subnet1.id}"
  ignore_missing_vnet_service_endpoint = false
}

resource "azurerm_postgresql_virtual_network_rule" "rule2" {
  name                                 = "acctestpostgresqlvnetrule2%d"
  resource_group_name                  = "${azurerm_resource_group.test.name}"
  server_name                          = "${azurerm_postgresql_server.test.name}"
  subnet_id                            = "${azurerm_subnet.vnet1_subnet2.id}"
  ignore_missing_vnet_service_endpoint = false
}

resource "azurerm_postgresql_virtual_network_rule" "rule3" {
  name                                 = "acctestpostgresqlvnetrule3%d"
  resource_group_name                  = "${azurerm_resource_group.test.name}"
  server_name                          = "${azurerm_postgresql_server.test.name}"
  subnet_id                            = "${azurerm_subnet.vnet2_subnet1.id}"
  ignore_missing_vnet_service_endpoint = false
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMPostgreSQLVirtualNetworkRule_ignoreEndpointValid(rInt int, location string) string {
	return fmt.Sprintf(` 
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/29"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.7.29.0/29"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctestpostgresqlsvr-%d"
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
  version                      = "9.5"
  ssl_enforcement              = "Enabled"
}

resource "azurerm_postgresql_virtual_network_rule" "test" {
  name                                 = "acctestpostgresqlvnetrule%d"
  resource_group_name                  = "${azurerm_resource_group.test.name}"
  server_name                          = "${azurerm_postgresql_server.test.name}"
  subnet_id                            = "${azurerm_subnet.test.id}"
  ignore_missing_vnet_service_endpoint = true
}
`, rInt, location, rInt, rInt, rInt, rInt)
}
