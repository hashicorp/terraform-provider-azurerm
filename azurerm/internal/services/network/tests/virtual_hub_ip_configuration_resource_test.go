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

func TestAccAzureRMVirtualHubIPConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_ip_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubIPConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubIPConfiguration_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubIPConfigurationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualHubIPConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_ip_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubIPConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubIPConfiguration_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubIPConfigurationExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualHubIPConfiguration_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_virtual_hub_ip_configuration"),
			},
		},
	})
}

func TestAccAzureRMVirtualHubIPConfiguration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_ip_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubIPConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubIPConfiguration_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubIPConfigurationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualHubIPConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_ip_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubIPConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubIPConfiguration_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubIPConfigurationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualHubIPConfiguration_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubIPConfigurationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMVirtualHubIPConfigurationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VirtualHubIPConfigurationClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Virtual Hub IP Configuration not found: %s", resourceName)
		}

		id, err := parse.VirtualHubIPConfigurationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Virtual Hub IP Configuration %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Bad: Get on network.VirtualHubIPConfigurationClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMVirtualHubIPConfigurationDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VirtualHubIPConfigurationClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_hub_ip_configuration" {
			continue
		}

		id, err := parse.VirtualHubIPConfigurationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on network.VirtualHubIPConfigurationClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMVirtualHubIPConfiguration_basic(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHubIPConfiguration_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_ip_configuration" "test" {
  name                         = "acctest-vhubipconfig-%d"
  virtual_hub_id               = azurerm_virtual_hub.test.id
  subnet_id                    = azurerm_subnet.test.id
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualHubIPConfiguration_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHubIPConfiguration_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_ip_configuration" "import" {
  name           = azurerm_virtual_hub_ip_configuration.test.name
  virtual_hub_id = azurerm_virtual_hub_ip_configuration.test.virtual_hub_id
}
`, template)
}

func testAccAzureRMVirtualHubIPConfiguration_complete(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHubIPConfiguration_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_ip_configuration" "test" {
  name                         = "acctest-vhubipconfig-%d"
  virtual_hub_id               = azurerm_virtual_hub.test.id
  private_ip_address           = "10.5.1.18"
  private_ip_allocation_method = "Static"
  public_ip_address_id         = azurerm_public_ip.test.id
  subnet_id                    = azurerm_subnet.test.id
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualHubIPConfiguration_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vhub-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-virtnet%d"
  address_space       = ["10.5.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.1.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctest-vwan-%d"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
