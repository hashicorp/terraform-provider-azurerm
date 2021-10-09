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

type CassandraMIClusterResource struct {
}

func TestAccCassandraMICluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_cassandra_mi_cluster", "test")
	r := CassandraMIClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t CassandraMIClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CassandraClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cosmos.CassandraClustersClient.Get(ctx, id.ResourceGroup, id.ClusterName)
	if err != nil {
		return nil, fmt.Errorf("reading Cosmos Cassandra Keyspace (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (CassandraMIClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`

  provider "azurerm" {
	features {}
  }
  
  resource "azurerm_resource_group" "test" {
	name     = "cassandra-%[1]d"
	location = "East US"
  }
  
  
  resource "azurerm_virtual_network" "test" {
	name                = "vnet-tf"
	location            = azurerm_resource_group.test.location
	resource_group_name = azurerm_resource_group.test.name
	address_space       = ["10.0.0.0/16"]
  
	subnet {
	  name           = "default"
	  address_prefix = "10.0.1.0/24"
	}
  
	tags = {
	  environment = "Production"
	}
  }
  
  resource "azurerm_role_assignment" "test" {
	scope              = azurerm_virtual_network.test.id
	role_definition_name = "Network Contributor"
	principal_id       = "e5007d2c-4b13-4a74-9b6a-605d99f03501"
  }
  
  resource "azurerm_cosmosdb_cassandra_mi_cluster" "test" {
	cluster_name        			= "cassandra-mi-cluster-%[1]d"
	resource_group_name 			= azurerm_resource_group.test.name
	location	      				= "eastus"
	delegated_management_subnet_id 	= "${azurerm_role_assignment.test.scope}/subnets/default" 
	initial_cassandra_admin_password	= "Password1234"  
  }
`, data.RandomInteger)
}

// func (CassandraMIClusterResource) basic(data acceptance.TestData) string {
// 	return fmt.Sprintf(`

//   provider "azurerm" {
// 	features {}
//   }
//   variable "subscription_id" {
// 	default = "c8a23972-1b42-43fa-9bda-92e665014f30"
//   }

//   resource "azurerm_resource_group" "example" {
// 	name     = "cassandra-%[1]d"
// 	location = "East US"
//   }

//   resource "azurerm_virtual_network" "example" {
// 	name                = "vnet-tf"
// 	location            = azurerm_resource_group.example.location
// 	resource_group_name = azurerm_resource_group.example.name
// 	address_space       = ["10.0.0.0/16"]

// 	subnet {
// 	  name           = "default"
// 	  address_prefix = "10.0.1.0/24"
// 	}

// 	tags = {
// 	  environment = "Production"
// 	}
//   }

//   resource "azurerm_role_assignment" "example" {
// 	scope              = "/subscriptions/${var.subscription_id}/resourceGroups/${azurerm_resource_group.example.name}/providers/Microsoft.Network/virtualNetworks/${azurerm_virtual_network.example.name}"
// 	role_definition_name = "Network Contributor"
// 	principal_id       = "e5007d2c-4b13-4a74-9b6a-605d99f03501"
//   }

//   resource "azurerm_cosmosdb_cassandra_mi_cluster" "example" {
// 	cluster_name        			= "terraform-nova-cluster-revised-api-1"
// 	resource_group_name 			= azurerm_resource_group.example.name
// 	location	      			= "eastus"
// 	delegated_management_subnet_id 	= "/subscriptions/${var.subscription_id}/resourceGroups/${azurerm_resource_group.example.name}/providers/Microsoft.Network/virtualNetworks/${azurerm_virtual_network.example.name}/subnets/default"
// 	initial_cassandra_admin_password	= "Password1234"
//   }
// `, data.RandomInteger)
// }
