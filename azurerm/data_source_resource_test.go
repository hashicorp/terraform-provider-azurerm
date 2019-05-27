package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMResource_ByResourceID(t *testing.T) {
	dataSourceName := "data.azurerm_resource.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	config := testAccDataSourceAzureRMResource_ByResourceID(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "resources.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMResource_ByName(t *testing.T) {
	dataSourceName := "data.azurerm_resource.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	config := testAccDataSourceAzureRMResource_ByName(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "resources.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMResource_ByResourceGroup(t *testing.T) {
	dataSourceName := "data.azurerm_resource.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	config := testAccDataSourceAzureRMResource_ByResourceGroup(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "resources.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMResource_ByResourceType(t *testing.T) {
	dataSourceName := "data.azurerm_resource.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	config := testAccDataSourceAzureRMResource_ByResourceType(ri, rs, location)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "resources.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMResource_FilteredByTags(t *testing.T) {
	dataSourceName := "data.azurerm_resource.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	config := testAccDataSourceAzureRMResource_FilteredByTags(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "resources.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMResource_ByResourceID(rInt int, rString string, location string) string {
	r := testAccDataSourceAzureRMStorageAccount_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_resource" "test" {
  resource_id = "${azurerm_storage_account.test.id}"
}
`, r)
}

func testAccDataSourceAzureRMResource_ByName(rInt int, rString string, location string) string {
	r := testAccDataSourceAzureRMStorageAccount_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_resource" "test" {
  name	= "${azurerm_storage_account.test.name}"
}
`, r)
}

func testAccDataSourceAzureRMResource_ByResourceGroup(rInt int, rString string, location string) string {
	r := testAccDataSourceAzureRMStorageAccount_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_resource" "test" {
  resource_group_name = "${azurerm_storage_account.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMResource_ByResourceType(rInt int, rString string, location string) string {
	r := testAccDataSourceAzureRMStorageAccount_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_resource" "test" {
	resource_group_name = "${azurerm_storage_account.test.resource_group_name}"
  type 								= "Microsoft.Storage/storageAccounts"
}
`, r)
}

func testAccDataSourceAzureRMResource_FilteredByTags(rInt int, rString string, location string) string {
	r := testAccDataSourceAzureRMStorageAccount_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_resource" "test" {
  name                = "${azurerm_storage_account.test.name}"
	resource_group_name = "${azurerm_storage_account.test.resource_group_name}"
	
	required_tags = {
		environment = "production"
	}
}
`, r)
}
