package azurerm

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMCosmosDBAccount_boundedStaleness_complete(t *testing.T) {
	ri := tf.AccRandTimeInt()
	dataSourceName := "data.azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDAtaSourceAzureRMCosmosDBAccount_boundedStaleness_complete(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(dataSourceName, acceptance.Location(), string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(dataSourceName, "consistency_policy.0.max_interval_in_seconds", "10"),
					resource.TestCheckResourceAttr(dataSourceName, "consistency_policy.0.max_staleness_prefix", "200"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMCosmosDBAccount_geoReplicated_customId(t *testing.T) {
	ri := tf.AccRandTimeInt()
	dataSourceName := "data.azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMCosmosDBAccount_geoReplicated_customId(ri, acceptance.Location(), acceptance.AltLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(dataSourceName, acceptance.Location(), string(documentdb.BoundedStaleness), 2),
					resource.TestCheckResourceAttr(dataSourceName, "geo_location.0.location", acceptance.Location()),
					resource.TestCheckResourceAttr(dataSourceName, "geo_location.1.location", acceptance.AltLocation()),
					resource.TestCheckResourceAttr(dataSourceName, "geo_location.0.failover_priority", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "geo_location.1.failover_priority", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMCosmosDBAccount_virtualNetworkFilter(t *testing.T) {
	ri := tf.AccRandTimeInt()
	dataSourceName := "data.azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMCosmosDBAccount_virtualNetworkFilter(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(dataSourceName, acceptance.Location(), string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(dataSourceName, "is_virtual_network_filter_enabled", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "virtual_network_rule.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "virtual_network_rule.1.id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMCosmosDBAccount_complete(t *testing.T) {
	ri := tf.AccRandTimeInt()
	dataSourceName := "data.azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMCosmosDBAccount_complete(ri, acceptance.Location(), acceptance.AltLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(dataSourceName, acceptance.Location(), string(documentdb.BoundedStaleness), 2),
					resource.TestCheckResourceAttr(dataSourceName, "ip_range_filter", "104.42.195.92,40.76.54.131,52.176.6.30,52.169.50.45/32,52.187.184.26,10.20.0.0/16"),
					resource.TestCheckResourceAttr(dataSourceName, "enable_automatic_failover", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "enable_multiple_write_locations", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "geo_location.0.location", acceptance.Location()),
					resource.TestCheckResourceAttr(dataSourceName, "geo_location.1.location", acceptance.AltLocation()),
					resource.TestCheckResourceAttr(dataSourceName, "geo_location.0.failover_priority", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "geo_location.1.failover_priority", "1"),
				),
			},
		},
	})
}

func testAccDAtaSourceAzureRMCosmosDBAccount_boundedStaleness_complete(rInt int, location string) string {
	return fmt.Sprintf(`
%s

data "azurerm_cosmosdb_account" "test" {
  name                = "${azurerm_cosmosdb_account.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, testAccAzureRMCosmosDBAccount_boundedStaleness_complete(rInt, location))
}

func testAccDataSourceAzureRMCosmosDBAccount_geoReplicated_customId(rInt int, location string, altLocation string) string {
	return fmt.Sprintf(`
%s

data "azurerm_cosmosdb_account" "test" {
  name                = "${azurerm_cosmosdb_account.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, testAccAzureRMCosmosDBAccount_geoReplicated_customId(rInt, location, altLocation))
}

func testAccDataSourceAzureRMCosmosDBAccount_complete(rInt int, location string, altLocation string) string {
	return fmt.Sprintf(`
%s

data "azurerm_cosmosdb_account" "test" {
  name                = "${azurerm_cosmosdb_account.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, testAccAzureRMCosmosDBAccount_complete(rInt, location, altLocation))
}

func testAccDataSourceAzureRMCosmosDBAccount_virtualNetworkFilter(rInt int, location string) string {
	return fmt.Sprintf(`
%s

data "azurerm_cosmosdb_account" "test" {
  name                = "${azurerm_cosmosdb_account.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, testAccAzureRMCosmosDBAccount_virtualNetworkFilter(rInt, location))
}
