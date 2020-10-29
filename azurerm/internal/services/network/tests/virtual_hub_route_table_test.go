package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualHubRouteTable_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_route_table", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubRouteTable_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubRouteTableExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualHubRouteTable_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_route_table", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubRouteTable_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubRouteTableExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMVirtualHubRouteTable_requiresImport),
		},
	})
}

func TestAccAzureRMVirtualHubRouteTable_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_route_table", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubRouteTable_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubRouteTableExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualHubRouteTable_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_route_table", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubRouteTable_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubRouteTableExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualHubRouteTable_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubRouteTableExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualHubRouteTable_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubRouteTableExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMVirtualHubRouteTableExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.HubRouteTableClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("network HubRouteTable not found: %s", resourceName)
		}

		id, err := parse.VirtualHubRouteTableID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Network HubRouteTable %q does not exist", id.Name)
			}

			return fmt.Errorf("bad: Get on Network.HubRouteTableClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMVirtualHubRouteTableDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.HubRouteTableClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_network_hub_route_table" {
			continue
		}

		id, err := parse.VirtualHubRouteTableID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Network.HubRouteTableClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMVirtualHubRouteTable_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-VHUB-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%d"
  address_space       = ["10.5.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test" {
  name                = "acctest-NSG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-SUBNET-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.1.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctest-VWAN-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-VHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.2.0/24"
}

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctest-VHUBCONN-%d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualHubRouteTable_basic(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHubRouteTable_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_route_table" "test" {
  name           = "acctest-RouteTable-%d"
  virtual_hub_id = azurerm_virtual_hub.test.id
  labels         = ["Label1"]
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualHubRouteTable_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMVirtualHubRouteTable_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_route_table" "import" {
  name           = azurerm_virtual_hub_route_table.test.name
  virtual_hub_id = azurerm_virtual_hub_route_table.test.virtual_hub_id
  labels         = azurerm_virtual_hub_route_table.test.labels
}
`, config)
}

func testAccAzureRMVirtualHubRouteTable_complete(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHubRouteTable_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_route_table" "test" {
  name           = "acctest-RouteTable-%d"
  virtual_hub_id = azurerm_virtual_hub.test.id
  labels         = ["labeL1", "AnotherLabel"]

  route {
    name              = "VHub-Route-Test"
    destinations_type = "CIDR"
    destinations      = ["10.0.0.0/16"]
    next_hop_type     = "ResourceId"
    next_hop          = azurerm_virtual_hub_connection.test.id
  }
}
`, template, data.RandomInteger)
}
