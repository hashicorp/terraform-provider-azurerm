package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSearchService_basic(t *testing.T) {
	resourceName := "azurerm_search_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMSearchService_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSearchServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSearchServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMSearchService_complete(t *testing.T) {
	resourceName := "azurerm_search_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMSearchService_complete(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSearchServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSearchServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "replica_count", "2"),
				),
			},
		},
	})
}

func testCheckAzureRMSearchServiceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		searchName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).searchServicesClient
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

		client := testAccProvider.Meta().(*ArmClient).searchServicesClient
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

func testAccAzureRMSearchService_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_search_service" "test" {
    name = "acctestsearchservice%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    sku = "standard"

    tags {
    	environment = "staging"
    }
}
`, rInt, location, rInt)
}

func testAccAzureRMSearchService_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}
resource "azurerm_search_service" "test" {
    name = "acctestsearchservice%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    sku = "standard"
    replica_count = 2

    tags {
    	environment = "production"
    }
}
`, rInt, location, rInt)
}
