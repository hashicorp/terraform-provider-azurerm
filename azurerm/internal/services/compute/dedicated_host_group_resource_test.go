package compute_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDedicatedHostGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDedicatedHostGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDedicatedHostGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDedicatedHostGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDedicatedHostGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDedicatedHostGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostGroupExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMDedicatedHostGroup_requiresImport),
		},
	})
}

func TestAccAzureRMDedicatedHostGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDedicatedHostGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDedicatedHostGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "zones.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "zones.0", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "platform_fault_domain_count", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "prod"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDedicatedHostGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Dedicated Host Group not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.DedicatedHostGroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, resourceGroup, name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Dedicated Host Group %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on Compute.DedicatedHostGroupsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMDedicatedHostGroupDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.DedicatedHostGroupsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dedicated_host_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on Compute.DedicatedHostGroupsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMDedicatedHostGroup_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-compute-%d"
  location = "%s"
}

resource "azurerm_dedicated_host_group" "test" {
  name                        = "acctestDHG-compute-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  platform_fault_domain_count = 2
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMDedicatedHostGroup_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDedicatedHostGroup_basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_dedicated_host_group" "import" {
  resource_group_name         = azurerm_dedicated_host_group.test.resource_group_name
  name                        = azurerm_dedicated_host_group.test.name
  location                    = azurerm_dedicated_host_group.test.location
  platform_fault_domain_count = 2
}
`, template)
}

func testAccAzureRMDedicatedHostGroup_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-compute-%d"
  location = "%s"
}

resource "azurerm_dedicated_host_group" "test" {
  name                        = "acctestDHG-compute-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  platform_fault_domain_count = 2
  zones                       = ["1"]
  tags = {
    ENV = "prod"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
