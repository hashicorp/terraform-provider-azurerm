package postgres_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPostgreSQLVirtualNetworkRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_virtual_network_rule", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLVirtualNetworkRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPostgreSQLVirtualNetworkRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_virtual_network_rule", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLVirtualNetworkRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMPostgreSQLVirtualNetworkRule_requiresImport),
		},
	})
}

func TestAccAzureRMPostgreSQLVirtualNetworkRule_switchSubnets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_virtual_network_rule", "test")

	// Create regex strings that will ensure that one subnet name exists, but not the other
	preConfigRegex := regexp.MustCompile(fmt.Sprintf("(acctest-SN1-%d)$|(acctest-SN[^2]-%d)$", data.RandomInteger, data.RandomInteger))  // subnet 1 but not 2
	postConfigRegex := regexp.MustCompile(fmt.Sprintf("(acctest-SN2-%d)$|(acctest-SN-[^1]%d)$", data.RandomInteger, data.RandomInteger)) // subnet 2 but not 1

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLVirtualNetworkRule_subnetSwitchPre(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(data.ResourceName),
					resource.TestMatchResourceAttr(data.ResourceName, "subnet_id", preConfigRegex),
				),
			},
			{
				Config: testAccAzureRMPostgreSQLVirtualNetworkRule_subnetSwitchPost(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(data.ResourceName),
					resource.TestMatchResourceAttr(data.ResourceName, "subnet_id", postConfigRegex),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLVirtualNetworkRule_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_virtual_network_rule", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLVirtualNetworkRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(data.ResourceName),
					testCheckAzureRMPostgreSQLVirtualNetworkRuleDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMPostgreSQLVirtualNetworkRule_multipleSubnets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_virtual_network_rule", "rule1")
	resourceName2 := "azurerm_postgresql_virtual_network_rule.rule2"
	resourceName3 := "azurerm_postgresql_virtual_network_rule.rule3"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLVirtualNetworkRule_multipleSubnets(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(data.ResourceName),
					testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(resourceName2),
					testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(resourceName3),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLVirtualNetworkRule_IgnoreEndpointValid(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_virtual_network_rule", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLVirtualNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLVirtualNetworkRule_ignoreEndpointValid(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "ignore_missing_vnet_service_endpoint", "true"),
				),
			},
		},
	})
}

func testCheckAzureRMPostgreSQLVirtualNetworkRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Postgres.VirtualNetworkRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		ruleName := rs.Primary.Attributes["name"]

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
	client := acceptance.AzureProvider.Meta().(*clients.Client).Postgres.VirtualNetworkRulesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_postgresql_virtual_network_rule" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		ruleName := rs.Primary.Attributes["name"]

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
		client := acceptance.AzureProvider.Meta().(*clients.Client).Postgres.VirtualNetworkRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		ruleName := rs.Primary.Attributes["name"]

		future, err := client.Delete(ctx, resourceGroup, serverName, ruleName)
		if err != nil {
			// If the error is that the resource we want to delete does not exist in the first
			// place (404), then just return with no error.
			if response.WasNotFound(future.Response()) {
				return nil
			}

			return fmt.Errorf("Error deleting PostgreSQL Virtual Network Rule: %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			// Same deal as before. Just in case.
			if response.WasNotFound(future.Response()) {
				return nil
			}

			return fmt.Errorf("Error deleting PostgreSQL Virtual Network Rule: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMPostgreSQLVirtualNetworkRule_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%d"
  address_space       = ["10.7.29.0/29"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-SN-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.0/29"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_postgresql_server" "test" {
  name                         = "acctest-psql-server-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }
}

resource "azurerm_postgresql_virtual_network_rule" "test" {
  name                = "acctest-PSQL-VNET-rule-%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_postgresql_server.test.name
  subnet_id           = azurerm_subnet.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPostgreSQLVirtualNetworkRule_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMPostgreSQLVirtualNetworkRule_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_virtual_network_rule" "import" {
  name                = azurerm_postgresql_virtual_network_rule.test.name
  resource_group_name = azurerm_postgresql_virtual_network_rule.test.resource_group_name
  server_name         = azurerm_postgresql_virtual_network_rule.test.server_name
  subnet_id           = azurerm_postgresql_virtual_network_rule.test.subnet_id
}
`, template)
}

func testAccAzureRMPostgreSQLVirtualNetworkRule_subnetSwitchPre(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%d"
  address_space       = ["10.7.29.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test1" {
  name                 = "acctest-SN1-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.0/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "test2" {
  name                 = "acctest-SN2-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.128/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_postgresql_server" "test" {
  name                         = "acctest-psql-server-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }
}

resource "azurerm_postgresql_virtual_network_rule" "test" {
  name                = "acctest-PSQL-VNET-rule-%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_postgresql_server.test.name
  subnet_id           = azurerm_subnet.test1.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPostgreSQLVirtualNetworkRule_subnetSwitchPost(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%d"
  address_space       = ["10.7.29.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test1" {
  name                 = "acctest-SN1-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.0/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "test2" {
  name                 = "acctest-SN2-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.128/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_postgresql_server" "test" {
  name                         = "acctest-psql-server-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }
}

resource "azurerm_postgresql_virtual_network_rule" "test" {
  name                = "acctest-PSQL-VNET-rule-%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_postgresql_server.test.name
  subnet_id           = azurerm_subnet.test2.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPostgreSQLVirtualNetworkRule_multipleSubnets(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "vnet1" {
  name                = "acctestvnet1%d"
  address_space       = ["10.7.29.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "vnet2" {
  name                = "acctestvnet2%d"
  address_space       = ["10.1.29.0/29"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "vnet1_subnet1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.vnet1.name
  address_prefix       = "10.7.29.0/29"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "vnet1_subnet2" {
  name                 = "acctestsubnet2%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.vnet1.name
  address_prefix       = "10.7.29.128/29"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "vnet2_subnet1" {
  name                 = "acctestsubnet3%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.vnet2.name
  address_prefix       = "10.1.29.0/29"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_postgresql_server" "test" {
  name                         = "acctest-psql-server-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }
}

resource "azurerm_postgresql_virtual_network_rule" "rule1" {
  name                                 = "acctestpostgresqlvnetrule1%d"
  resource_group_name                  = azurerm_resource_group.test.name
  server_name                          = azurerm_postgresql_server.test.name
  subnet_id                            = azurerm_subnet.vnet1_subnet1.id
  ignore_missing_vnet_service_endpoint = false
}

resource "azurerm_postgresql_virtual_network_rule" "rule2" {
  name                                 = "acctestpostgresqlvnetrule2%d"
  resource_group_name                  = azurerm_resource_group.test.name
  server_name                          = azurerm_postgresql_server.test.name
  subnet_id                            = azurerm_subnet.vnet1_subnet2.id
  ignore_missing_vnet_service_endpoint = false
}

resource "azurerm_postgresql_virtual_network_rule" "rule3" {
  name                                 = "acctestpostgresqlvnetrule3%d"
  resource_group_name                  = azurerm_resource_group.test.name
  server_name                          = azurerm_postgresql_server.test.name
  subnet_id                            = azurerm_subnet.vnet2_subnet1.id
  ignore_missing_vnet_service_endpoint = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPostgreSQLVirtualNetworkRule_ignoreEndpointValid(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%d"
  address_space       = ["10.7.29.0/29"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-SN-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.0/29"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctestpostgresqlsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true
}

resource "azurerm_postgresql_virtual_network_rule" "test" {
  name                                 = "acctestpostgresqlvnetrule%d"
  resource_group_name                  = azurerm_resource_group.test.name
  server_name                          = azurerm_postgresql_server.test.name
  subnet_id                            = azurerm_subnet.test.id
  ignore_missing_vnet_service_endpoint = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
