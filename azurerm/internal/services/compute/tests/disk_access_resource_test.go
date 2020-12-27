package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMDiskAccess_empty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_access", "test")
	var d compute.DiskAccess

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDiskAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDiskAccess_empty(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDiskAccessExists(data.ResourceName, &d, true),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDiskAccess_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_access", "test")
	var d compute.DiskAccess

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDiskAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDiskAccess_empty(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDiskAccessExists(data.ResourceName, &d, true),
				),
			},
			{
				Config:      testAccAzureRMDiskAccess_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_disk_access"),
			},
		},
	})
}

func TestAccAzureRMDiskAccess_import(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_access", "test")
	var d compute.DiskAccess

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDiskAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDiskAccess_import(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDiskAccessExists(data.ResourceName, &d, true),
				),
			},
		},
	})
}

func testCheckAzureRMDiskAccessExists(resourceName string, d *compute.DiskAccess, shouldExist bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.DiskAccessClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		daName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for disk access: %s", daName)
		}
		resp, err := client.Get(ctx, resourceGroup, daName)
		if err != nil {
			return fmt.Errorf("Bad: Get on diskAccessClient: %v", err)
		}
		if resp.StatusCode == http.StatusNotFound && shouldExist {
			return fmt.Errorf("Bad: DiskAccess %q (resource group %q) does not exist", daName, resourceGroup)
		}
		if resp.StatusCode != http.StatusNotFound && !shouldExist {
			return fmt.Errorf("Bad: DiskAccess %q (resource group %q) still exists", daName, resourceGroup)
		}

		*d = resp

		return nil
	}
}

func testAccAzureRMDiskAccess_empty(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_disk_access" "test" {
  name                = "acctestda-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}

	`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMDiskAccess_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDiskAccess_empty(data)
	return fmt.Sprintf(`
	%s

resource "azurerm_disk_access" "import" {
  name                = azurerm_disk_access.test.name
  location            = azurerm_disk_access.test.location
  resource_group_name = azurerm_disk_access.test.resource_group_name

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
		`, template)
}

func testAccAzureRMDiskAccess_import(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_disk_access" "test" {
  name                = "accda%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location


  tags = {
    environment = "staging"
  }
}


	`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testCheckAzureRMDiskAccessDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.DiskAccessClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_disk_access" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Disk Access still exists: \n%#v", resp.DiskAccessProperties)
		}
	}

	return nil
}
