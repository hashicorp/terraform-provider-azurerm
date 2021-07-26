package cosmos_test

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-01-15/documentdb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type CosmosDBAccountDataSourceResource struct {
}

func TestAccDataSourceCosmosDBAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountDataSourceResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.BoundedStaleness, 1),
			),
		},
	})
}

func TestAccDataSourceCosmosDBAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cosmosdb_account", "test")
	r := CosmosDBAccountDataSourceResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkAccCosmosDBAccount_basic(data, documentdb.BoundedStaleness, 3),
				check.That(data.ResourceName).Key("geo_location.0.location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("geo_location.1.location").HasValue(data.Locations.Secondary),
				check.That(data.ResourceName).Key("geo_location.2.location").HasValue(data.Locations.Ternary),
				check.That(data.ResourceName).Key("geo_location.0.failover_priority").HasValue("0"),
				check.That(data.ResourceName).Key("geo_location.1.failover_priority").HasValue("1"),
				check.That(data.ResourceName).Key("geo_location.2.failover_priority").HasValue("2"),
			),
		},
	})
}

func (CosmosDBAccountDataSourceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cosmosdb_account" "test" {
  name                = azurerm_cosmosdb_account.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, CosmosDBAccountResource{}.basic(data, documentdb.GlobalDocumentDB, documentdb.BoundedStaleness))
}

func (CosmosDBAccountDataSourceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cosmosdb_account" "test" {
  name                = azurerm_cosmosdb_account.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, CosmosDBAccountResource{}.complete(data, documentdb.GlobalDocumentDB, documentdb.BoundedStaleness))
}
