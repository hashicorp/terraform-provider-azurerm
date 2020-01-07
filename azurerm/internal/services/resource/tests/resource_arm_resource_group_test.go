package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMResourceGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMResourceGroup_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMResourceGroup_requiresImport),
		},
	})
}

func TestAccAzureRMResourceGroup_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			data.DisappearsStep(acceptance.DisappearsStepData{
				Config:      testAccAzureRMResourceGroup_basic,
				CheckExists: testCheckAzureRMResourceGroupExists,
				Destroy:     testCheckAzureRMResourceGroupDisappears,
			}),
		},
	})
}

func TestAccAzureRMResourceGroup_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceGroup_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.cost_center", "MSFT"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMResourceGroup_withTagsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "staging"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMResourceGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Resource.GroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["name"]

		// Ensure resource group exists in API

		resp, err := client.Get(ctx, resourceGroup)
		if err != nil {
			return fmt.Errorf("Bad: Get on resourceGroupClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: resource group: %q does not exist", resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMResourceGroupDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Resource.GroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["name"]

		// Ensure resource group exists in API

		deleteFuture, err := client.Delete(ctx, resourceGroup)
		if err != nil {
			return fmt.Errorf("Failed deleting Resource Group %q: %+v", resourceGroup, err)
		}

		err = deleteFuture.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("Failed long polling for the deletion of Resource Group %q: %+v", resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMResourceGroupDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Resource.GroupsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_resource_group" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Resource Group still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccAzureRMResourceGroup_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMResourceGroup_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMResourceGroup_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group" "import" {
  name     = "${azurerm_resource_group.test.name}"
  location = "${azurerm_resource_group.test.location}"
}
`, template)
}

func testAccAzureRMResourceGroup_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMResourceGroup_withTagsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
