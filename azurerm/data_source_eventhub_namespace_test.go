package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMEventHubNamespace_basic(t *testing.T) {
	dataSourceName := "data.azurerm_eventhub_namespace.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEventHubNamespace_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "sku", "Basic"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMEventHubNamespace_complete(t *testing.T) {
	dataSourceName := "data.azurerm_eventhub_namespace.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEventHubNamespace_complete(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "sku", "Standard"),
					resource.TestCheckResourceAttr(dataSourceName, "capacity", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "auto_inflate_enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "maximum_throughput_units", "20"),
				),
			},
		},
	})
}

func testAccDataSourceEventHubNamespace_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Basic"
}

data "azurerm_eventhub_namespace" "test" {
  name                = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_eventhub_namespace.test.resource_group_name}"
}
`, rInt, location, rInt)
}

func testAccDataSourceEventHubNamespace_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                     = "acctesteventhubnamespace-%d"
  location                 = "${azurerm_resource_group.test.location}"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  sku                      = "Standard"
  capacity                 = "2"
  auto_inflate_enabled     = true
  maximum_throughput_units = 20
}

data "azurerm_eventhub_namespace" "test" {
  name                = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_eventhub_namespace.test.resource_group_name}"
}
`, rInt, location, rInt)
}
