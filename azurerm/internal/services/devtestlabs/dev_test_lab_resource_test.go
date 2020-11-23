package devtestlabs_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccDevTestLab_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_lab", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestLabDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestLab_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestLabExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_type", "Premium"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccDevTestLab_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_lab", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestLabDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestLab_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestLabExists(data.ResourceName),
				),
			},
			{
				Config:      testAccDevTestLab_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_dev_test_lab"),
			},
		},
	})
}

func TestAccDevTestLab_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_lab", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestLabDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestLab_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestLabExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_type", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Hello", "World"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckDevTestLabExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).DevTestLabs.LabsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		labName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for DevTest Lab: %s", labName)
		}

		resp, err := conn.Get(ctx, resourceGroup, labName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get devTestLabsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DevTest Lab %q (Resource Group: %q) does not exist", labName, resourceGroup)
		}

		return nil
	}
}

func testCheckDevTestLabDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).DevTestLabs.LabsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dev_test_lab" {
			continue
		}

		labName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, labName, "")
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("DevTest Lab still exists:\n%#v", resp)
	}

	return nil
}

func testAccDevTestLab_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccDevTestLab_requiresImport(data acceptance.TestData) string {
	template := testAccDevTestLab_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_lab" "import" {
  name                = azurerm_dev_test_lab.test.name
  location            = azurerm_dev_test_lab.test.location
  resource_group_name = azurerm_dev_test_lab.test.resource_group_name
}
`, template)
}

func testAccDevTestLab_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  storage_type        = "Standard"

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
