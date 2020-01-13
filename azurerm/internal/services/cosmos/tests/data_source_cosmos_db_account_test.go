package tests

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMCosmosDBAccount_boundedStaleness_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDAtaSourceAzureRMCosmosDBAccount_boundedStaleness_complete(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(data.ResourceName, "consistency_policy.0.max_interval_in_seconds", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "consistency_policy.0.max_staleness_prefix", "200"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMCosmosDBAccount_geoReplicated_customId(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMCosmosDBAccount_geoReplicated_customId(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 2),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_location.0.location", data.Locations.Primary),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_location.1.location", data.Locations.Secondary),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_location.0.failover_priority", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_location.1.failover_priority", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMCosmosDBAccount_virtualNetworkFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMCosmosDBAccount_virtualNetworkFilter(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(data.ResourceName, "is_virtual_network_filter_enabled", "true"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_network_rule.0.id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_network_rule.1.id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMCosmosDBAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMCosmosDBAccount_complete(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 2),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_range_filter", "104.42.195.92,40.76.54.131,52.176.6.30,52.169.50.45/32,52.187.184.26,10.20.0.0/16"),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_automatic_failover", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_multiple_write_locations", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_location.0.location", data.Locations.Primary),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_location.1.location", data.Locations.Secondary),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_location.0.failover_priority", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_location.1.failover_priority", "1"),
				),
			},
		},
	})
}

func testAccDAtaSourceAzureRMCosmosDBAccount_boundedStaleness_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cosmosdb_account" "test" {
  name                = "${azurerm_cosmosdb_account.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, testAccAzureRMCosmosDBAccount_boundedStaleness_complete(data))
}

func testAccDataSourceAzureRMCosmosDBAccount_geoReplicated_customId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cosmosdb_account" "test" {
  name                = "${azurerm_cosmosdb_account.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, testAccAzureRMCosmosDBAccount_geoReplicated_customId(data))
}

func testAccDataSourceAzureRMCosmosDBAccount_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cosmosdb_account" "test" {
  name                = "${azurerm_cosmosdb_account.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, testAccAzureRMCosmosDBAccount_complete(data))
}

func testAccDataSourceAzureRMCosmosDBAccount_virtualNetworkFilter(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cosmosdb_account" "test" {
  name                = "${azurerm_cosmosdb_account.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, testAccAzureRMCosmosDBAccount_virtualNetworkFilter(data))
}
