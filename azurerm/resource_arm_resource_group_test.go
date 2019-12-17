package azurerm

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
	testData := acceptance.BuildTestData(t, "azurerm_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceGroup_basic(testData),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(testData.ResourceName),
				),
			},
			testData.ImportStep(),
		},
	})
}

func TestAccAzureRMResourceGroup_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	testData := acceptance.BuildTestData(t, "azurerm_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceGroup_basic(testData),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(testData.ResourceName),
				),
			},
			testData.RequiresImportErrorStep(testAccAzureRMResourceGroup_requiresImport),
		},
	})
}

func TestAccAzureRMResourceGroup_disappears(t *testing.T) {
	testData := acceptance.BuildTestData(t, "azurerm_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			testData.DisappearsStep(acceptance.DisappearsStepData{
				Config:      testAccAzureRMResourceGroup_basic,
				CheckExists: testCheckAzureRMResourceGroupExists,
				Destroy:     testCheckAzureRMResourceGroupDisappears,
			}),
		},
	})
}

func TestAccAzureRMResourceGroup_withTags(t *testing.T) {
	testData := acceptance.BuildTestData(t, "azurerm_resource_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceGroup_withTags(testData),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(testData.ResourceName),
					resource.TestCheckResourceAttr(testData.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(testData.ResourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(testData.ResourceName, "tags.cost_center", "MSFT"),
				),
			},
			testData.ImportStep(),
			{
				Config: testAccAzureRMResourceGroup_withTagsUpdated(testData),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(testData.ResourceName),
					resource.TestCheckResourceAttr(testData.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(testData.ResourceName, "tags.environment", "staging"),
				),
			},
			testData.ImportStep(),
		},
	})
}

func testCheckAzureRMResourceGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["name"]

		// Ensure resource group exists in API
		client := acceptance.AzureProvider.Meta().(*clients.Client).Resource.GroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["name"]

		// Ensure resource group exists in API
		client := acceptance.AzureProvider.Meta().(*clients.Client).Resource.GroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMResourceGroup_basic(testData acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
`, testData.RandomInteger, testData.Locations.Primary)
}

func testAccAzureRMResourceGroup_requiresImport(testData acceptance.TestData) string {
	template := testAccAzureRMResourceGroup_basic(testData)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group" "import" {
  name     = "${azurerm_resource_group.test.name}"
  location = "${azurerm_resource_group.test.location}"
}
`, template)
}

func testAccAzureRMResourceGroup_withTags(testData acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, testData.RandomInteger, testData.Locations.Primary)
}

func testAccAzureRMResourceGroup_withTagsUpdated(testData acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  tags = {
    environment = "staging"
  }
}
`, testData.RandomInteger, testData.Locations.Primary)
}
