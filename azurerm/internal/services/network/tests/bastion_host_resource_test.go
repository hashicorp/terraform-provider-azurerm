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

func TestAccAzureRMBastionHost_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bastion_host", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBastionHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBastionHost_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBastionHostExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMBastionHost_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bastion_host", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBastionHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBastionHost_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBastionHostExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "production"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMBastionHost_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_bastion_host", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBastionHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBastionHost_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBastionHostExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMBastionHost_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_bastion_host"),
			},
		},
	})
}

func testAccAzureRMBastionHost_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-bastion-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVNet%s"
  address_space       = ["192.168.1.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name               = "AzureBastionSubnet"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "192.168.1.224/27"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestBastionPIP%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_bastion_host" "test" {
  name                = "acctestBastion%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                 = "ip-configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomString)
}

func testAccAzureRMBastionHost_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-bastion-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVNet%s"
  address_space       = ["192.168.1.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name               = "AzureBastionSubnet"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "192.168.1.224/27"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestBastionPIP%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_bastion_host" "test" {
  name                = "acctestBastion%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                 = "ip-configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomString)
}

func testAccAzureRMBastionHost_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMBastionHost_basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_bastion_host" "import" {
  name                = azurerm_bastion_host.test.name
  resource_group_name = azurerm_bastion_host.test.resource_group_name
  location            = azurerm_bastion_host.test.location

  ip_configuration {
    name                 = "ip-configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
`, template)
}

func testCheckAzureRMBastionHostExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.BastionHostsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Azure Bastion Host %q does not exist", rs.Primary.ID)
			}
			return fmt.Errorf("Bad: Get on Azure Bastion Host Client: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMBastionHostDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.BastionHostsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_bastion_host" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}
