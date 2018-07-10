package azurerm

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMCosmosDBAccount_boundedStaleness_complete(t *testing.T) {
	ri := acctest.RandInt()
	dataSourceName := "data.azurerm_cosmosdb_account.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDAtaSourceAzureRMCosmosDBAccount_boundedStaleness_complete(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(dataSourceName, testLocation(), string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(dataSourceName, "consistency_policy.0.max_interval_in_seconds", "10"),
					resource.TestCheckResourceAttr(dataSourceName, "consistency_policy.0.max_staleness_prefix", "200"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMCosmosDBAccount_geoReplicated_customId(t *testing.T) {
	ri := acctest.RandInt()
	dataSourceName := "data.azurerm_cosmosdb_account.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMCosmosDBAccount_geoReplicated_customId(ri, testLocation(), testAltLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(dataSourceName, testLocation(), string(documentdb.BoundedStaleness), 2),
					resource.TestCheckResourceAttr(dataSourceName, "geo_location.0.location", testLocation()),
					resource.TestCheckResourceAttr(dataSourceName, "geo_location.1.location", testAltLocation()),
					resource.TestCheckResourceAttr(dataSourceName, "geo_location.0.failover_priority", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "geo_location.1.failover_priority", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMCosmosDBAccount_complete(t *testing.T) {
	ri := acctest.RandInt()
	dataSourceName := "data.azurerm_cosmosdb_account.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMCosmosDBAccount_complete(ri, testLocation(), testAltLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(dataSourceName, testLocation(), string(documentdb.BoundedStaleness), 2),
					resource.TestCheckResourceAttr(dataSourceName, "ip_range_filter", "104.42.195.92,40.76.54.131,52.176.6.30,52.169.50.45,52.187.184.26"),
					resource.TestCheckResourceAttr(dataSourceName, "enable_automatic_failover", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "geo_location.0.location", testLocation()),
					resource.TestCheckResourceAttr(dataSourceName, "geo_location.1.location", testAltLocation()),
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
