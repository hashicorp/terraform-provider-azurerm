package resource

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMResources_ByName(t *testing.T) {
	dataSourceName := "data.azurerm_resources.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMResources_template(ri, rs, location),
			},
			{
				Config: testAccDataSourceAzureRMResources_ByName(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "resources.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMResources_ByResourceGroup(t *testing.T) {
	dataSourceName := "data.azurerm_resources.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMResources_template(ri, rs, location),
			},
			{
				Config: testAccDataSourceAzureRMResources_ByResourceGroup(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "resources.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMResources_ByResourceType(t *testing.T) {
	dataSourceName := "data.azurerm_resources.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := acceptance.Location()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMResources_template(ri, rs, location),
			},
			{
				Config: testAccDataSourceAzureRMResources_ByResourceType(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "resources.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMResources_FilteredByTags(t *testing.T) {
	dataSourceName := "data.azurerm_resources.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMResources_template(ri, rs, location),
			},
			{
				Config: testAccDataSourceAzureRMResources_FilteredByTags(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "resources.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMResources_ByName(rInt int, rString string, location string) string {
	r := testAccDataSourceAzureRMResources_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_resources" "test" {
  name = "${azurerm_storage_account.test.name}"
}
`, r)
}

func testAccDataSourceAzureRMResources_ByResourceGroup(rInt int, rString string, location string) string {
	r := testAccDataSourceAzureRMResources_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_resources" "test" {
  resource_group_name = "${azurerm_storage_account.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMResources_ByResourceType(rInt int, rString string, location string) string {
	r := testAccDataSourceAzureRMResources_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_resources" "test" {
  resource_group_name = "${azurerm_storage_account.test.resource_group_name}"
  type                = "Microsoft.Storage/storageAccounts"
}
`, r)
}

func testAccDataSourceAzureRMResources_FilteredByTags(rInt int, rString string, location string) string {
	r := testAccDataSourceAzureRMResources_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_resources" "test" {
  name                = "${azurerm_storage_account.test.name}"
  resource_group_name = "${azurerm_storage_account.test.resource_group_name}"

  required_tags = {
    environment = "production"
  }
}
`, r)
}

func testAccDataSourceAzureRMResources_template(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsads%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}
