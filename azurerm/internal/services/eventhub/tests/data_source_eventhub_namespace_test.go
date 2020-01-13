package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMEventHubNamespace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_namespace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEventHubNamespace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "Basic"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMEventHubNamespace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_namespace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEventHubNamespace_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "capacity", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_inflate_enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "maximum_throughput_units", "20"),
				),
			},
		},
	})
}

func testAccDataSourceEventHubNamespace_basic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccDataSourceEventHubNamespace_complete(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
