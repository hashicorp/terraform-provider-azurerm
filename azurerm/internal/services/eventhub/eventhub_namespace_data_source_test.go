package eventhub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccEventHubNamespaceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_namespace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEventHubNamespaceDataSource_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "Basic"),
				),
			},
		},
	})
}

func TestAccEventHubNamespaceDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_namespace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEventHubNamespaceDataSource_complete(data),
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

func TestAccEventHubNamespaceDataSource_withAliasConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_namespace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				// `default_primary_connection_string_alias` and `default_secondary_connection_string_alias` are still `nil` while `data.azurerm_eventhub_namespace` is retrieving resource. since `azurerm_eventhub_namespace_disaster_recovery_config` hasn't been created.
				// So these two properties should be checked in the second run.
				Config: testAccEventHubNamespace_withAliasConnectionString(data),
			},
			{
				Config: testAccEventHubNamespaceDataSource_withAliasConnectionString(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "default_primary_connection_string_alias"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "default_secondary_connection_string_alias"),
				),
			},
		},
	})
}

func testAccEventHubNamespaceDataSource_basic(data acceptance.TestData) string {
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
  sku                 = "Basic"
}

data "azurerm_eventhub_namespace" "test" {
  name                = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_eventhub_namespace.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccEventHubNamespaceDataSource_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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

func testAccEventHubNamespaceDataSource_withAliasConnectionString(data acceptance.TestData) string {
	template := testAccEventHubNamespace_withAliasConnectionString(data)
	return fmt.Sprintf(`
%s

data "azurerm_eventhub_namespace" "test" {
  name                = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_eventhub_namespace.test.resource_group_name
}
`, template)
}
