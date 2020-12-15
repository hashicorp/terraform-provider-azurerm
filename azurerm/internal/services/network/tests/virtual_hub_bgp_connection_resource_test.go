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

func TestAccAzureRMVirtualHubBgpConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_bgp_connection", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubBgpConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubBgpConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubBgpConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualHubBgpConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_bgp_connection", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubBgpConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubBgpConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubBgpConnectionExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMVirtualHubBgpConnection_requiresImport),
		},
	})
}

func testCheckAzureRMVirtualHubBgpConnectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VirtualHubBgpConnectionClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("virtualHubBgpConnection not found: %s", resourceName)
		}

		id, err := parse.BgpConnectionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Network VirtualHubBgpConnection %q does not exist", id.Name)
			}

			return fmt.Errorf("bad: Get on Network.VirtualHubBgpConnectionClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMVirtualHubBgpConnectionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VirtualHubBgpConnectionClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_hub_bgp_connection" {
			continue
		}

		id, err := parse.BgpConnectionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Network.VirtualHubBgpConnectionClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMVirtualHubBgpConnection_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-VHub-%d"
  location = "%s"
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-VHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-PIP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNet-%d"
  address_space       = ["10.5.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-Subnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.5.1.0/24"
}

resource "azurerm_subnet" "gateway" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.5.0.0/24"
}

resource "azurerm_virtual_hub_ip" "test" {
  name                         = "acctest-VHub-IP-%d"
  virtual_hub_id               = azurerm_virtual_hub.test.id
  private_ip_address           = "10.5.1.18"
  private_ip_allocation_method = "Static"
  public_ip_address_id         = azurerm_public_ip.test.id
  subnet_id                    = azurerm_subnet.test.id

  depends_on = [azurerm_subnet.gateway]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualHubBgpConnection_basic(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHubBgpConnection_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_bgp_connection" "test" {
  name           = "acctest-VHub-BgpConnection-%d"
  virtual_hub_id = azurerm_virtual_hub.test.id
  peer_asn       = 65514
  peer_ip        = "169.254.21.5"

  depends_on = [azurerm_virtual_hub_ip.test]
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualHubBgpConnection_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMVirtualHubBgpConnection_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_bgp_connection" "import" {
  name           = azurerm_virtual_hub_bgp_connection.test.name
  virtual_hub_id = azurerm_virtual_hub_bgp_connection.test.virtual_hub_id
  peer_asn       = azurerm_virtual_hub_bgp_connection.test.peer_asn
  peer_ip        = azurerm_virtual_hub_bgp_connection.test.peer_ip
}
`, config)
}
