package eventhub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccEventHubConsumerGroupDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_consumer_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckEventHubConsumerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventHubConsumerGroupDataSource_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckEventHubConsumerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "user_metadata", "some-meta-data"),
				),
			},
		},
	})
}

func TestAccEventHubConsumerGroupDataSourceDefault_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_consumer_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckEventHubConsumerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventHubConsumerGroupDataSourceDefault_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckEventHubConsumerGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "user_metadata", "some-meta-data"),
				),
			},
		},
	})
}

func testAccEventHubConsumerGroupDataSource_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 7
}

resource "azurerm_eventhub_consumer_group" "test" {
  name                = "acctesteventhubcg-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
  user_metadata       = "some-meta-data"
}

data "azurerm_eventhub_consumer_group" "test" {
  name                = azurerm_eventhub_consumer_group.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccEventHubConsumerGroupDataSourceDefault_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 7
}

data "azurerm_eventhub_consumer_group" "test" {
  name                = "$Default"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
