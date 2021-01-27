package cosmos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type CosmosDBCassandraTableResource struct {
}

func (t CosmosDBCassandraTableResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.CassandraTableID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cosmos.CassandraClient.GetCassandraTable(ctx, id.ResourceGroup, id.DatabaseAccountName, id.CassandraKeyspaceName, id.TableName)
	if err != nil {
		return nil, fmt.Errorf("reading Cosmos Cassandra Table (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func TestAccCosmosDbCassandraTable_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_cassandra_table", "test")
	r := CosmosDBCassandraTableResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{

			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (CosmosDBCassandraTableResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_cassandra_table" "test" {
  name                = "acctest-CCASST-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  keyspace_name       = azurerm_cosmosdb_cassandra_keyspace.test.name
}
`, CosmosDbCassandraKeyspaceResource{}.basic(data), data.RandomInteger)
}
