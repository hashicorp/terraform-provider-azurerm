package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAppServiceVirtualNetworkSwiftConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_virtual_network_swift_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceVirtualNetworkSwiftConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceVirtualNetworkSwiftConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_virtual_network_swift_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceVirtualNetworkSwiftConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAppServiceVirtualNetworkSwiftConnection_requiresImport),
		},
	})
}

func TestAccAzureRMAppServiceVirtualNetworkSwiftConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_virtual_network_swift_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceVirtualNetworkSwiftConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
				),
			},
			{
				Config: testAccAzureRMAppServiceVirtualNetworkSwiftConnection_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceVirtualNetworkSwiftConnection_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_virtual_network_swift_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceVirtualNetworkSwiftConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
					testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.VirtualNetworkSwiftConnectionID(rs.Primary.Attributes["id"])
		if err != nil {
			return fmt.Errorf("Error parsing Azure Resource ID %q", id)
		}

		resp, err := client.GetSwiftVirtualNetworkConnection(ctx, id.ResourceGroup, id.SiteName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: App Service Virtual Network Association %q (Resource Group: %q) does not exist", id.SiteName, id.ResourceGroup)
			}

			return fmt.Errorf("Bad: Get on appServicesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.VirtualNetworkSwiftConnectionID(rs.Primary.Attributes["id"])
		if err != nil {
			return fmt.Errorf("Error parsing Azure Resource ID %q", id)
		}

		resp, err := client.DeleteSwiftVirtualNetwork(ctx, id.ResourceGroup, id.SiteName)
		if err != nil {
			if !utils.ResponseWasNotFound(resp) {
				return fmt.Errorf("Bad: Delete on appServicesClient: %+v", err)
			}
		}

		return nil
	}
}

func testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_virtual_network_swift_connection" {
			continue
		}

		id, err := parse.VirtualNetworkSwiftConnectionID(rs.Primary.Attributes["id"])
		if err != nil {
			return fmt.Errorf("Error parsing Azure Resource ID %q", id)
		}

		resp, err := client.GetSwiftVirtualNetworkConnection(ctx, id.ResourceGroup, id.SiteName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return nil
	}

	return nil
}

func testAccAzureRMAppServiceVirtualNetworkSwiftConnection_base(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appservice-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  lifecycle {
    ignore_changes = [ddos_protection_plan]
  }
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestSubnet1"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"

  delegation {
    name = "acctestdelegation"

    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_subnet" "test2" {
  name                 = "acctestSubnet2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"

  delegation {
    name = "acctestdelegation"

    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctest-ASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctest-AS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceVirtualNetworkSwiftConnection_basic(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceVirtualNetworkSwiftConnection_base(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_virtual_network_swift_connection" "test" {
  app_service_id = azurerm_app_service.test.id
  subnet_id      = azurerm_subnet.test1.id
}
`, template)
}

func testAccAzureRMAppServiceVirtualNetworkSwiftConnection_update(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceVirtualNetworkSwiftConnection_base(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_virtual_network_swift_connection" "test" {
  app_service_id = azurerm_app_service.test.id
  subnet_id      = azurerm_subnet.test2.id
}
`, template)
}

func testAccAzureRMAppServiceVirtualNetworkSwiftConnection_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceVirtualNetworkSwiftConnection_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_virtual_network_swift_connection" "import" {
  app_service_id = azurerm_app_service_virtual_network_swift_connection.test.app_service_id
  subnet_id      = azurerm_app_service_virtual_network_swift_connection.test.subnet_id
}
`, template)
}
