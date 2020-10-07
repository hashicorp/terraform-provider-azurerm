package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualWan_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_wan", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualWanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualWan_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualWanExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualWan_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_wan", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualWanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualWan_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualWanExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualWan_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_virtual_wan"),
			},
		},
	})
}

func TestAccAzureRMVirtualWan_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_wan", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualWanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualWan_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualWanExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMVirtualWanDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VirtualWanClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_wan" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Virtual WAN still exists:\n%+v", resp)
		}
	}

	return nil
}

func testCheckAzureRMVirtualWanExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VirtualWanClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		virtualWanName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Virtual WAN: %s", virtualWanName)
		}

		resp, err := client.Get(ctx, resourceGroup, virtualWanName)
		if err != nil {
			return fmt.Errorf("Bad: Get on virtualWanClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Virtual WAN %q (resource group: %q) does not exist", virtualWanName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMVirtualWan_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMVirtualWan_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMVirtualWan_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_wan" "import" {
  name                = azurerm_virtual_wan.test.name
  resource_group_name = azurerm_virtual_wan.test.resource_group_name
  location            = azurerm_virtual_wan.test.location
}
`, template)
}

func testAccAzureRMVirtualWan_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  disable_vpn_encryption            = false
  allow_branch_to_branch_traffic    = true
  office365_local_breakout_category = "All"
  type                              = "Standard"

  tags = {
    Hello = "There"
    World = "Example"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
