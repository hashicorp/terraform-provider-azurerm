package eventhub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccEventHubDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEventHub_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "partition_count", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "partition_ids.#", "2"),
				),
			},
		},
	})
}

func testAccDataSourceEventHub_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest-eh-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  partition_count     = 2
  message_retention   = 1
}

data "azurerm_eventhub" "test" {
  name                = azurerm_eventhub.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}
