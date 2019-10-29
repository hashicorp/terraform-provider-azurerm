package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMResourceGroup_basic(t *testing.T) {
	resourceName := "azurerm_resource_group.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceGroup_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMResourceGroup_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_resource_group.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceGroup_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMResourceGroup_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_resource_group"),
			},
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
		client := testAccProvider.Meta().(*ArmClient).Resource.GroupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
		client := testAccProvider.Meta().(*ArmClient).Resource.GroupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
	client := testAccProvider.Meta().(*ArmClient).Resource.GroupsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_resource_group" {
			continue
		}

		resourceGroup := rs.Primary.ID

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

func testAccAzureRMResourceGroup_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_advanced_threat_protection" "test" {
  resource_id = azurerm_resource_group.test.id
  enable      = true
}
`, rInt, location)
}

func testAccAzureRMResourceGroup_storageAccount(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ATP-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "acctest%d"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}

resource "azurerm_advanced_threat_protection" "test" {
  resource_id = azurerm_resource_group.azurerm_storage_account.id
  enable      = true
}
`, rInt, location)
}
