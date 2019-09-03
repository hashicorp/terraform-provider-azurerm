package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSearchService_basic(t *testing.T) {
	resourceName := "azurerm_search_service.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSearchServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSearchService_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSearchServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
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

func TestAccAzureRMSearchService_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_search_service.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSearchServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSearchService_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSearchServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
				),
			},
			{
				Config:      testAccAzureRMSearchService_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_search_service"),
			},
		},
	})
}

func TestAccAzureRMSearchService_complete(t *testing.T) {
	resourceName := "azurerm_search_service.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSearchServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSearchService_complete(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSearchServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "replica_count", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_key"),
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

func TestAccAzureRMSearchService_tagUpdate(t *testing.T) {
	resourceName := "azurerm_search_service.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSearchServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSearchService_withCustomTagValue(ri, testLocation(), "staging"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSearchServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "staging"),
				),
			},
			{
				Config: testAccAzureRMSearchService_withCustomTagValue(ri, testLocation(), "production"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSearchServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "production"),
				),
			},
		},
	})
}

func testCheckAzureRMSearchServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		searchName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).search.ServicesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, searchName, nil)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Search Service %q (resource group %q) was not found: %+v", searchName, resourceGroup, err)
			}

			return fmt.Errorf("Bad: GetSearchService: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSearchServiceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_search_service" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		searchName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).search.ServicesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, searchName, nil)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Bad: Search Service %q (resource group %q) still exists: %+v", searchName, resourceGroup, resp)
	}

	return nil
}

func testAccAzureRMSearchService_withCustomTagValue(rInt int, location string, tagValue string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_search_service" "test" {
  name                = "acctestsearchservice%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "standard"

  tags = {
    environment = "%s"
  }
}
`, rInt, location, rInt, tagValue)
}

func testAccAzureRMSearchService_basic(rInt int, location string) string {
	return testAccAzureRMSearchService_withCustomTagValue(rInt, location, "staging")
}

func testAccAzureRMSearchService_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s
resource "azurerm_search_service" "import" {
  name                = "${azurerm_search_service.test.name}"
  resource_group_name = "${azurerm_search_service.test.resource_group_name}"
  location            = "${azurerm_search_service.test.location}"
  sku                 = "${azurerm_search_service.test.sku}"

  tags = {
    environment = "staging"
  }
}
`, testAccAzureRMSearchService_basic(rInt, location))
}

func testAccAzureRMSearchService_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_search_service" "test" {
  name                = "acctestsearchservice%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "standard"
  replica_count       = 2

  tags = {
    environment = "production"
  }
}
`, rInt, location, rInt)
}
