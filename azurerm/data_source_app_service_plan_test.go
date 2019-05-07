package azurerm

import (
	"fmt"

	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMAppServicePlan_basic(t *testing.T) {
	dataSourceName := "data.azurerm_app_service_plan.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppServicePlan_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "kind", "Windows"),
					resource.TestCheckResourceAttr(dataSourceName, "sku.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "sku.0.tier", "Basic"),
					resource.TestCheckResourceAttr(dataSourceName, "sku.0.size", "B1"),
					resource.TestCheckResourceAttr(dataSourceName, "properties.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "properties.0.per_site_scaling", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAppServicePlan_complete(t *testing.T) {
	dataSourceName := "data.azurerm_app_service_plan.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppServicePlan_complete(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "kind", "Windows"),
					resource.TestCheckResourceAttr(dataSourceName, "sku.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "sku.0.tier", "Standard"),
					resource.TestCheckResourceAttr(dataSourceName, "sku.0.size", "S1"),
					resource.TestCheckResourceAttr(dataSourceName, "properties.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "properties.0.per_site_scaling", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.environment", "Test"),
				),
			},
		},
	})
}

func testAccDataSourceAppServicePlan_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Basic"
    size = "B1"
  }
}

data "azurerm_app_service_plan" "test" {
  name                = "${azurerm_app_service_plan.test.name}"
  resource_group_name = "${azurerm_app_service_plan.test.resource_group_name}"
}
`, rInt, location, rInt)
}

func testAccDataSourceAppServicePlan_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "Windows"

  sku {
    tier = "Standard"
    size = "S1"
  }

  properties {
    per_site_scaling = true
  }

  tags = {
    environment = "Test"
  }
}

data "azurerm_app_service_plan" "test" {
  name                = "${azurerm_app_service_plan.test.name}"
  resource_group_name = "${azurerm_app_service_plan.test.resource_group_name}"
}
`, rInt, location, rInt)
}
