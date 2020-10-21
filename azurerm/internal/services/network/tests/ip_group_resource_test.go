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

func TestAccAzureRMIpGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ip_group", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIpGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIpGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIpGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIpGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ip_group", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIpGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIpGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIpGroupExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMIpGroup_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_ip_group"),
			},
		},
	})
}

func TestAccAzureRMIpGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ip_group", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIpGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIpGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIpGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIpGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ip_group", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIpGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIpGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIpGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "cidrs.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMIpGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIpGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "cidrs.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMIpGroup_completeUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIpGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "cidrs.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMIpGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIpGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "cidrs.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMIpGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIpGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "cidrs.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMIpGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.IPGroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		ipGroupName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for IP Group: %q", ipGroupName)
		}

		resp, err := client.Get(ctx, resourceGroup, ipGroupName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: IP Group %q (resource group: %q) does not exist", ipGroupName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on IPGroupsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMIpGroupDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.IPGroupsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_ip_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name, "")

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return fmt.Errorf("IP Group still exists:\n%#v", resp.IPGroupPropertiesFormat)
	}

	return nil
}

func testAccAzureRMIpGroup_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_ip_group" "test" {
  name                = "acceptanceTestIpGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMIpGroup_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMIpGroup_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_ip_group" "import" {
  name                = azurerm_ip_group.test.name
  location            = azurerm_ip_group.test.location
  resource_group_name = azurerm_ip_group.test.resource_group_name
}
`, template)
}

func testAccAzureRMIpGroup_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_ip_group" "test" {
  name                = "acceptanceTestIpGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cidrs = ["192.168.0.1", "172.16.240.0/20", "10.48.0.0/12"]

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMIpGroup_completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_ip_group" "test" {
  name                = "acceptanceTestIpGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cidrs = ["172.16.240.0/20"]

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
