package eventhub_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type EventHubNamespaceDataSource struct {
}

func TestAccEventHubNamespaceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku").HasValue("Basic"),
			),
		},
	})
}

func TestAccEventHubNamespaceDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku").HasValue("Standard"),
				check.That(data.ResourceName).Key("capacity").HasValue("2"),
				check.That(data.ResourceName).Key("auto_inflate_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("maximum_throughput_units").HasValue("20"),
			),
		},
	})
}

func TestAccEventHubNamespaceDataSource_withAliasConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			// `default_primary_connection_string_alias` and `default_secondary_connection_string_alias` are still `nil` while `data.azurerm_eventhub_namespace` is retrieving acceptance. since `azurerm_eventhub_namespace_disaster_recovery_config` hasn't been created.
			// So these two properties should be checked in the second run.
			Config: r.withAliasConnectionString(data),
		},
		{
			Config: r.withAliasConnectionString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("default_primary_connection_string_alias").Exists(),
				check.That(data.ResourceName).Key("default_secondary_connection_string_alias").Exists(),
			),
		},
	})
}

func (EventHubNamespaceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

data "azurerm_eventhub_namespace" "test" {
  name                = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_eventhub_namespace.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                     = "acctesteventhubnamespace-%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  sku                      = "Standard"
  capacity                 = "2"
  auto_inflate_enabled     = true
  maximum_throughput_units = 20
}

data "azurerm_eventhub_namespace" "test" {
  name                = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_eventhub_namespace.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceDataSource) withAliasConnectionString(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_eventhub_namespace" "test" {
  name                = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_eventhub_namespace.test.resource_group_name
}
`, EventHubNamespaceResource{}.withAliasConnectionString(data))
}
