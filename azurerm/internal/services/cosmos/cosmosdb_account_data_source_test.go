package cosmos_test

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2020-04-01/documentdb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMCosmosDBAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMCosmosDBAccount_basic(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, documentdb.BoundedStaleness, 1),
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
					checkAccAzureRMCosmosDBAccount_basic(data, documentdb.BoundedStaleness, 3),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_location.0.location", data.Locations.Primary),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_location.1.location", data.Locations.Secondary),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_location.2.location", data.Locations.Ternary),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_location.0.failover_priority", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_location.1.failover_priority", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "geo_location.2.failover_priority", "2"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMCosmosDBAccount_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cosmosdb_account" "test" {
  name                = azurerm_cosmosdb_account.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, testAccAzureRMCosmosDBAccount_basic(data, documentdb.GlobalDocumentDB, documentdb.BoundedStaleness))
}

func testAccDataSourceAzureRMCosmosDBAccount_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cosmosdb_account" "test" {
  name                = azurerm_cosmosdb_account.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, testAccAzureRMCosmosDBAccount_complete(data, documentdb.GlobalDocumentDB, documentdb.BoundedStaleness))
}
