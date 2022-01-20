package cosmos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CassandraClusterResource struct {
}

func TestAccCassandraCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_cassandra_cluster", "test")
	r := CassandraClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("default_admin_password"),
	})
}

func TestAccCassandraCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_cassandra_cluster", "test")
	r := CassandraClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("default_admin_password"),
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_cosmosdb_cassandra_cluster"),
		},
	})
}

func (t CassandraClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CassandraClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cosmos.CassandraClustersClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading %q: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (CassandraClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-ca-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = "255f3c8e-0c3d-4f06-ba9d-2fb68af0faed"
}

resource "azurerm_cosmosdb_cassandra_cluster" "test" {
  name                           = "acctca-mi-cluster-%[1]d"
  resource_group_name            = azurerm_resource_group.test.name
  location                       = azurerm_resource_group.test.location
  delegated_management_subnet_id = azurerm_subnet.test.id
  default_admin_password         = "Password1234"
}
`, data.RandomInteger, data.Locations.Secondary)
}

func (CassandraClusterResource) requiresImport(data acceptance.TestData) string {
	template := CassandraClusterResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_cassandra_cluster" "import" {
  name                           = azurerm_cosmosdb_cassandra_cluster.test.name
  resource_group_name            = azurerm_cosmosdb_cassandra_cluster.test.resource_group_name
  location                       = azurerm_cosmosdb_cassandra_cluster.test.location
  delegated_management_subnet_id = azurerm_subnet.test.id
  default_admin_password         = "Password1234"
}
`, template)
}
