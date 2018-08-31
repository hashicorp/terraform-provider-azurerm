package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func init() {
	resource.AddTestSweepers("azurerm_resource_group", &resource.Sweeper{
		Name: "azurerm_resource_group",
		F:    testSweepResourceGroups,
	})
}

func testSweepResourceGroups(region string) error {
	armClient, err := buildConfigForSweepers()
	if err != nil {
		return err
	}

	client := (*armClient).resourceGroupsClient
	ctx := (*armClient).StopContext

	log.Printf("Retrieving the Resource Groups..")
	results, err := client.List(ctx, "", utils.Int32(int32(1000)))
	if err != nil {
		return fmt.Errorf("Error Listing on Resource Groups: %+v", err)
	}

	for _, resourceGroup := range results.Values() {
		if !shouldSweepAcceptanceTestResource(*resourceGroup.Name, *resourceGroup.Location, region) {
			continue
		}

		name := *resourceGroup.Name
		log.Printf("Deleting Resource Group %q", name)
		deleteFuture, err := client.Delete(ctx, name)
		if err != nil {
			return err
		}

		err = deleteFuture.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestAccAzureRMResourceGroup_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMResourceGroup_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists("azurerm_resource_group.test"),
				),
			},
		},
	})
}

func TestAccAzureRMResourceGroup_requiresImport(t *testing.T) {
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceGroup_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists("azurerm_resource_group.test"),
				),
			},
			{
				Config:      testAccAzureRMResourceGroup_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_resource_group"),
			},
		},
	})
}

func TestAccAzureRMResourceGroup_disappears(t *testing.T) {
	resourceName := "azurerm_resource_group.test"
	ri := acctest.RandInt()
	config := testAccAzureRMResourceGroup_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(resourceName),
					testCheckAzureRMResourceGroupDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMResourceGroup_withTags(t *testing.T) {
	resourceName := "azurerm_resource_group.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMResourceGroup_withTags(ri, location)
	postConfig := testAccAzureRMResourceGroup_withTagsUpdated(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(resourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "staging"),
				),
			},
		},
	})
}

func testCheckAzureRMResourceGroupExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["name"]

		// Ensure resource group exists in API
		client := testAccProvider.Meta().(*ArmClient).resourceGroupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup)
		if err != nil {
			return fmt.Errorf("Bad: Get on resourceGroupClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Virtual Network %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMResourceGroupDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["name"]

		// Ensure resource group exists in API
		client := testAccProvider.Meta().(*ArmClient).resourceGroupsClient
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
	client := testAccProvider.Meta().(*ArmClient).resourceGroupsClient
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
    name = "acctestRG-%d"
    location = "%s"
}
`, rInt, location)
}

func testAccAzureRMResourceGroup_requiresImport(rInt int, location string) string {
	template := testAccAzureRMResourceGroup_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group" "import" {
  name     = "${azurerm_resource_group.test.name}"
  location = "${azurerm_resource_group.test.location}"
}
`, template)
}

func testAccAzureRMResourceGroup_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"

    tags {
	environment = "Production"
	cost_center = "MSFT"
    }
}
`, rInt, location)
}

func testAccAzureRMResourceGroup_withTagsUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"

    tags {
	environment = "staging"
    }
}
`, rInt, location)
}
