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

func TestAccAzureRMVirtualHubIP_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_ip", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubIP_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubIPExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualHubIP_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_ip", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubIP_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubIPExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualHubIP_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_virtual_hub_ip"),
			},
		},
	})
}

func TestAccAzureRMVirtualHubIP_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_ip", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubIP_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubIPExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualHubIP_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_ip", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubIP_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubIPExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualHubIP_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubIPExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMVirtualHubIPExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VirtualHubIPClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Virtual Hub IP not found: %s", resourceName)
		}

		id, err := parse.VirtualHubIpConfigurationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.IpConfigurationName); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Virtual Hub IP %q (Resource Group %q) does not exist", id.IpConfigurationName, id.ResourceGroup)
			}
			return fmt.Errorf("Bad: Get on network.VirtualHubIPClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMVirtualHubIPDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VirtualHubIPClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_hub_ip" {
			continue
		}

		id, err := parse.VirtualHubIpConfigurationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.IpConfigurationName); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on network.VirtualHubIPClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMVirtualHubIP_basic(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHubIP_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_ip" "test" {
  name           = "acctest-vhubipconfig-%d"
  virtual_hub_id = azurerm_virtual_hub.test.id
  subnet_id      = azurerm_subnet.test.id

  depends_on = [azurerm_subnet.gateway]
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualHubIP_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHubIP_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_ip" "import" {
  name           = azurerm_virtual_hub_ip.test.name
  virtual_hub_id = azurerm_virtual_hub_ip.test.virtual_hub_id
  subnet_id      = azurerm_virtual_hub_ip.test.subnet_id

  depends_on = [azurerm_subnet.gateway]
}
`, template)
}

func testAccAzureRMVirtualHubIP_complete(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHubIP_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_ip" "test" {
  name                         = "acctest-vhubipconfig-%d"
  virtual_hub_id               = azurerm_virtual_hub.test.id
  private_ip_address           = "10.5.1.18"
  private_ip_allocation_method = "Static"
  public_ip_address_id         = azurerm_public_ip.test.id
  subnet_id                    = azurerm_subnet.test.id

  depends_on = [azurerm_subnet.gateway]
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualHubIP_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vhub-%d"
  location = "%s"
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-vhub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-pip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  address_space       = ["10.5.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%d"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
