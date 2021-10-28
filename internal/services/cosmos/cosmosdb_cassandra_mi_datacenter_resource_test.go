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

type CassandraMIDatacenterResource struct {
}

func TestAccCassandraMIDatacenter_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_cassandra_managed_instance_datacenter", "test")
	r := CassandraMIDatacenterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("initial_cassandra_admin_password"),
	})
}

func TestAccCassandraMIDatacenter_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_cassandra_managed_instance_datacenter", "test")
	r := CassandraMIDatacenterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.update(data, 3),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("initial_cassandra_admin_password"),
		{
			Config: r.update(data, 5),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("initial_cassandra_admin_password"),
	})
}

//Basic test case also covers the Complete test case because all fields are required by the resource

func (t CassandraMIDatacenterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CassandraDatacenterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cosmos.CassandraDatacentersClient.Get(ctx, id.ResourceGroup, id.CassandraClusterName, id.DataCenterName)
	if err != nil {
		return nil, fmt.Errorf("reading Cassandra MI Datacenter (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (CassandraMIDatacenterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ca-%d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]

  tags = {
    environment = "Test"
  }
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
  principal_id         = "e5007d2c-4b13-4a74-9b6a-605d99f03501"
}

resource "azurerm_cosmosdb_cassandra_managed_instance_cluster" "test" {
  name                     = "acctca-mi-cluster-%[1]d"
  resource_group_name              = azurerm_resource_group.test.name
  location                         = azurerm_resource_group.test.location
  delegated_management_subnet_id   = azurerm_subnet.test.id
  initial_cassandra_admin_password = "Password1234"
}

resource "azurerm_cosmosdb_cassandra_managed_instance_datacenter" "test" {
  name                   = azurerm_cosmosdb_cassandra_managed_instance_cluster.test.name
  datacenter_name                = "acctca-mi-dc-%[1]d"
  resource_group_name            = azurerm_resource_group.test.name
  location                       = azurerm_resource_group.test.location
  delegated_management_subnet_id = azurerm_subnet.test.id
  node_count                     = 3
}
`, data.RandomInteger, data.Locations.Secondary)
}

func (CassandraMIDatacenterResource) update(data acceptance.TestData, nodeCount int) string {
	return fmt.Sprintf(`

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ca-%d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]

  tags = {
    environment = "Test"
  }
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
  principal_id         = "e5007d2c-4b13-4a74-9b6a-605d99f03501"
}

resource "azurerm_cosmosdb_cassandra_managed_instance_cluster" "test" {
  name                     = "acctca-mi-cluster-%[1]d"
  resource_group_name              = azurerm_resource_group.test.name
  location                         = azurerm_resource_group.test.location
  delegated_management_subnet_id   = azurerm_subnet.test.id
  initial_cassandra_admin_password = "Password1234"
}

resource "azurerm_cosmosdb_cassandra_managed_instance_datacenter" "test" {
  name                   = azurerm_cosmosdb_cassandra_managed_instance_cluster.test.name
  datacenter_name                = "acctca-mi-dc-%[1]d"
  resource_group_name            = azurerm_resource_group.test.name
  location                       = azurerm_resource_group.test.location
  delegated_management_subnet_id = azurerm_subnet.test.id
  node_count                     = %[3]s
}
`, data.RandomInteger, data.Locations.Secondary, fmt.Sprint(nodeCount))
}
