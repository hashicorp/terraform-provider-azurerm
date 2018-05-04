package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMDataLakeStore_payasyougo(t *testing.T) {
	resourceName := "azurerm_data_lake_store.test"
	ri := acctest.RandInt()
	config := testAccAzureRMDataLakeStore_payasyougo(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMDataLakeStore_monthlycommitment(t *testing.T) {
	resourceName := "azurerm_data_lake_store.test"
	ri := acctest.RandInt()
	config := TestAccAzureRMDataLakeStore_monthlycommitment(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMDataLakeStore_withTags(t *testing.T) {
	resourceName := "azurerm_data_lake_store.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMDataLakeStore_withTags(ri, location)
	postConfig := testAccAzureRMDataLakeStore_withTagsUpdate(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsZoneExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func testCheckAzureRMDataLakeStoreExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		accountName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for data lake store: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).dataLakeStoreAccountClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, accountName)
		if err != nil {
			return fmt.Errorf("Bad: Get on dataLakeStoreAccountClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Date Lake Store %q (resource group: %q) does not exist", accountName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDataLakeStoreDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).dataLakeStoreAccountClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_lake_store" {
			continue
		}

		accountName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, accountName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("Data Lake Store still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMDataLakeStore_payasyougo(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "%s"
}

resource "azurerm_data_lake_store" "test" {
    name = "acctestlake%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func TestAccAzureRMDataLakeStore_monthlycommitment(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "%s"
}

resource "azurerm_data_lake_store" "test" {
    name = "acctestlake%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    tier = "Commitment_1TB"
}
`, rInt, location, rInt)
}

func testAccAzureRMDataLakeStore_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "%s"
}

resource "azurerm_data_lake_store" "test" {
    name = "acctestlake%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    tags {
        environment = "Production"
        cost_center = "MSFT"
    }
}
`, rInt, location, rInt)
}

func testAccAzureRMDataLakeStore_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "%s"
}

resource "azurerm_data_lake_store" "test" {
    name = "acctestlake%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    tags {
        environment = "staging"
    }
}
`, rInt, location, rInt)
}

