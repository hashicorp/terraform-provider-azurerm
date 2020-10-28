package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPrivateEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateEndpoint_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_service_connection.0.private_ip_address"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPrivateEndpoint_updateTag(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateEndpoint_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPrivateEndpoint_withTag(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPrivateEndpoint_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPrivateEndpoint_requestMessage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateEndpoint_requestMessage(data, "CATS: ALL YOUR BASE ARE BELONG TO US."),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_service_connection.0.request_message", "CATS: ALL YOUR BASE ARE BELONG TO US."),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPrivateEndpoint_requestMessage(data, "CAPTAIN: WHAT YOU SAY!!"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_service_connection.0.request_message", "CAPTAIN: WHAT YOU SAY!!"),
				),
			},
			data.ImportStep(),
		},
	})
}

// The update and complete test cases had to be totally removed since there is a bug with tags and the support for
// tags has been removed, all other attributes are ForceNew.
// API Issue "Unable to remove Tags from Private Endpoint": https://github.com/Azure/azure-sdk-for-go/issues/6467

func TestAccAzureRMPrivateEndpoint_privateDnsZoneGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateEndpoint_privateDnsZoneGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_group.0.private_dns_zone_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_configs.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_group.#", "1"),
				),
			},
			data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
		},
	})
}

func TestAccAzureRMPrivateEndpoint_privateDnsZoneRename(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateEndpoint_privateDnsZoneGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_group.0.private_dns_zone_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_configs.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_group.#", "1"),
				),
			},
			data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
			{
				Config: testAccAzureRMPrivateEndpoint_privateDnsZoneGroupRename(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_group.0.private_dns_zone_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_configs.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_group.#", "1"),
				),
			},
			data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
		},
	})
}

func TestAccAzureRMPrivateEndpoint_privateDnsZoneUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateEndpoint_privateDnsZoneGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_group.0.private_dns_zone_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_configs.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_group.#", "1"),
				),
			},
			data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
			{
				Config: testAccAzureRMPrivateEndpoint_privateDnsZoneGroupUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_group.0.private_dns_zone_ids.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_configs.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_group.#", "1"),
				),
			},
			data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
		},
	})
}

func TestAccAzureRMPrivateEndpoint_privateDnsZoneRemove(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateEndpoint_privateDnsZoneGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_group.0.private_dns_zone_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_configs.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_group.#", "1"),
				),
			},
			data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
			{
				Config: testAccAzureRMPrivateEndpoint_privateDnsZoneGroupRemove(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(data.ResourceName),
				),
			},
			data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
			{
				Config: testAccAzureRMPrivateEndpoint_privateDnsZoneGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_group.0.private_dns_zone_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_configs.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_dns_zone_group.#", "1"),
				),
			},
			data.ImportStep("private_dns_zone_configs", "private_dns_zone_group"),
		},
	})
}

func testCheckAzureRMPrivateEndpointExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.PrivateEndpointClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Private Endpoint not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Private Endpoint %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on PrivateEndpointClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMPrivateEndpointDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.PrivateEndpointClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_private_endpoint" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on PrivateEndpointClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMPrivateEndpointTemplate_template(data acceptance.TestData, seviceCfg string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-privatelink-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "service" {
  name               = "acctestsnetservice-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefixes   = ["10.5.1.0/24"]

  enforce_private_link_service_network_policies = true
}

resource "azurerm_subnet" "endpoint" {
  name               = "acctestsnetendpoint-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefixes   = ["10.5.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  sku                 = "Standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  sku                 = "Standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  frontend_ip_configuration {
    name                 = azurerm_public_ip.test.name
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

%s
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, seviceCfg)
}

func testAccAzureRMPrivateEndpoint_serviceAutoApprove(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name      = "primaryIpConfiguration-%d"
    primary   = true
    subnet_id = azurerm_subnet.service.id
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateEndpoint_serviceManualApprove(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_private_link_service" "test" {
  name                = "acctestPLS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  nat_ip_configuration {
    name      = "primaryIpConfiguration-%d"
    primary   = true
    subnet_id = azurerm_subnet.service.id
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateEndpoint_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = azurerm_private_link_service.test.name
    is_manual_connection           = false
    private_connection_resource_id = azurerm_private_link_service.test.id
  }
}
`, testAccAzureRMPrivateEndpointTemplate_template(data, testAccAzureRMPrivateEndpoint_serviceAutoApprove(data)), data.RandomInteger)
}

func testAccAzureRMPrivateEndpoint_withTag(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = azurerm_private_link_service.test.name
    is_manual_connection           = false
    private_connection_resource_id = azurerm_private_link_service.test.id
  }

  tags = {
    env = "TEST"
  }
}
`, testAccAzureRMPrivateEndpointTemplate_template(data, testAccAzureRMPrivateEndpoint_serviceAutoApprove(data)), data.RandomInteger)
}

func testAccAzureRMPrivateEndpoint_requestMessage(data acceptance.TestData, msg string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = azurerm_private_link_service.test.name
    is_manual_connection           = true
    private_connection_resource_id = azurerm_private_link_service.test.id
    request_message                = %q
  }
}
`, testAccAzureRMPrivateEndpointTemplate_template(data, testAccAzureRMPrivateEndpoint_serviceManualApprove(data)), data.RandomInteger, msg)
}

func testAccAzureRMPrivateEndpoint_privateDnsZoneGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-privatelink-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "service" {
  name               = "acctestsnetservice-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefixes   = ["10.5.1.0/24"]

  enforce_private_link_service_network_policies = true
}

resource "azurerm_subnet" "endpoint" {
  name               = "acctestsnetendpoint-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefixes   = ["10.5.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-pe-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_4"

  storage_mb                   = 5120
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  auto_grow_enabled            = true

  administrator_login          = "psqladminun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true
}

resource "azurerm_private_dns_zone" "finance" {
  name                = "privatelink.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.endpoint.id

  private_dns_zone_group {
    name                 = "acctest-dzg-%d"
    private_dns_zone_ids = [azurerm_private_dns_zone.finance.id]
  }

  private_service_connection {
    name                           = "acctest-privatelink-psc-%d"
    private_connection_resource_id = azurerm_postgresql_server.test.id
    subresource_names              = ["postgresqlServer"]
    is_manual_connection           = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateEndpoint_privateDnsZoneGroupRemove(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-privatelink-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "service" {
  name               = "acctestsnetservice-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefixes   = ["10.5.1.0/24"]

  enforce_private_link_service_network_policies = true
}

resource "azurerm_subnet" "endpoint" {
  name               = "acctestsnetendpoint-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefixes   = ["10.5.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-pe-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_4"

  storage_mb                   = 5120
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  auto_grow_enabled            = true

  administrator_login          = "psqladminun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true
}

resource "azurerm_private_dns_zone" "finance" {
  name                = "privatelink.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = "acctest-privatelink-psc-%d"
    private_connection_resource_id = azurerm_postgresql_server.test.id
    subresource_names              = ["postgresqlServer"]
    is_manual_connection           = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateEndpoint_privateDnsZoneGroupUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-privatelink-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "service" {
  name               = "acctestsnetservice-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefixes   = ["10.5.1.0/24"]

  enforce_private_link_service_network_policies = true
}

resource "azurerm_subnet" "endpoint" {
  name               = "acctestsnetendpoint-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefixes   = ["10.5.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-pe-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_4"

  storage_mb                   = 5120
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  auto_grow_enabled            = true

  administrator_login          = "psqladminun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true
}

resource "azurerm_private_dns_zone" "finance" {
  name                = "privatelink.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone" "sales" {
  name                = "acctest.pdz.%d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest.privatelink.%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.endpoint.id

  private_dns_zone_group {
    name                 = "acctest-dzg-%d"
    private_dns_zone_ids = [azurerm_private_dns_zone.sales.id, azurerm_private_dns_zone.finance.id]
  }

  private_service_connection {
    name                           = "acctest-privatelink-psc-%d"
    private_connection_resource_id = azurerm_postgresql_server.test.id
    subresource_names              = ["postgresqlServer"]
    is_manual_connection           = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateEndpoint_privateDnsZoneGroupRename(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-privatelink-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "service" {
  name               = "acctestsnetservice-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefixes   = ["10.5.1.0/24"]

  enforce_private_link_service_network_policies = true
}

resource "azurerm_subnet" "endpoint" {
  name               = "acctestsnetendpoint-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefixes   = ["10.5.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-pe-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_4"

  storage_mb                   = 5120
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  auto_grow_enabled            = true

  administrator_login          = "psqladminun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true
}

resource "azurerm_private_dns_zone" "finance" {
  name                = "privatelink.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.endpoint.id

  private_dns_zone_group {
    name                 = "acctest-dzg-rn-%d"
    private_dns_zone_ids = [azurerm_private_dns_zone.finance.id]
  }

  private_service_connection {
    name                           = "acctest-privatelink-psc-%d"
    private_connection_resource_id = azurerm_postgresql_server.test.id
    subresource_names              = ["postgresqlServer"]
    is_manual_connection           = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
