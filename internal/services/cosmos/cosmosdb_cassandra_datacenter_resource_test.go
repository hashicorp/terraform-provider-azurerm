// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

type CassandraDatacenterResource struct{}

func testAccCassandraDatacenter_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_cassandra_datacenter", "test")
	r := CassandraDatacenterResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, 3),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccCassandraDatacenter_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_cassandra_datacenter", "test")
	r := CassandraDatacenterResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, 3),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, 5),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t CassandraDatacenterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CassandraDatacenterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cosmos.CassandraDatacentersClient.Get(ctx, id.ResourceGroup, id.CassandraClusterName, id.DataCenterName)
	if err != nil {
		return nil, fmt.Errorf("reading %q: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r CassandraDatacenterResource) basic(data acceptance.TestData, nodeCount int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_cassandra_datacenter" "test" {
  name                           = "acctca-mi-dc-%d"
  cassandra_cluster_id           = azurerm_cosmosdb_cassandra_cluster.test.id
  location                       = azurerm_cosmosdb_cassandra_cluster.test.location
  delegated_management_subnet_id = azurerm_subnet.test.id
  node_count                     = %d
  disk_count                     = 4
  sku_name                       = "Standard_DS14_v2"
  availability_zones_enabled     = false
}
`, r.template(data), data.RandomInteger, nodeCount)
}

func (r CassandraDatacenterResource) complete(data acceptance.TestData, nodeCount int) string {
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv-%s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  purge_protection_enabled   = true
}

resource "azurerm_key_vault_access_policy" "current_user" {
  key_vault_id = azurerm_key_vault.test.id

  tenant_id = azurerm_key_vault.test.tenant_id
  object_id = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy"
  ]
}

resource "azurerm_key_vault_access_policy" "system_identity" {
  key_vault_id = azurerm_key_vault.test.id

  tenant_id = azurerm_key_vault.test.tenant_id
  object_id = azurerm_cosmosdb_cassandra_cluster.test.identity.0.principal_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy"
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkey-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = [
    azurerm_key_vault_access_policy.current_user,
    azurerm_key_vault_access_policy.system_identity
  ]
}

resource "azurerm_cosmosdb_cassandra_datacenter" "test" {
  name                            = "acctca-mi-dc-%d"
  cassandra_cluster_id            = azurerm_cosmosdb_cassandra_cluster.test.id
  location                        = azurerm_cosmosdb_cassandra_cluster.test.location
  delegated_management_subnet_id  = azurerm_subnet.test.id
  node_count                      = %d
  disk_count                      = 4
  sku_name                        = "Standard_DS14_v2"
  availability_zones_enabled      = false
  disk_sku                        = "P30"
  backup_storage_customer_key_uri = azurerm_key_vault_key.test.id
  managed_disk_customer_key_uri   = azurerm_key_vault_key.test.id
  base64_encoded_yaml_fragment    = "Y29tcGFjdGlvbl90aHJvdWdocHV0X21iX3Blcl9zZWM6IDMyCmNvbXBhY3Rpb25fbGFyZ2VfcGFydGl0aW9uX3dhcm5pbmdfdGhyZXNob2xkX21iOiAxMDA="

  depends_on = [
    azurerm_key_vault_key.test
  ]
}
`, r.template(data), data.RandomString, data.RandomString, data.RandomInteger, nodeCount)
}

func (r CassandraDatacenterResource) update(data acceptance.TestData, nodeCount int) string {
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv-%s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  purge_protection_enabled   = true
}

resource "azurerm_key_vault_access_policy" "current_user" {
  key_vault_id = azurerm_key_vault.test.id

  tenant_id = azurerm_key_vault.test.tenant_id
  object_id = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy"
  ]
}

resource "azurerm_key_vault_access_policy" "system_identity" {
  key_vault_id = azurerm_key_vault.test.id

  tenant_id = azurerm_key_vault.test.tenant_id
  object_id = azurerm_cosmosdb_cassandra_cluster.test.identity.0.principal_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy"
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkey-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = [
    azurerm_key_vault_access_policy.current_user,
    azurerm_key_vault_access_policy.system_identity
  ]
}

resource "azurerm_key_vault_key" "test2" {
  name         = "acctestkey2-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = [azurerm_key_vault_key.test]
}

resource "azurerm_cosmosdb_cassandra_datacenter" "test" {
  name                            = "acctca-mi-dc-%d"
  cassandra_cluster_id            = azurerm_cosmosdb_cassandra_cluster.test.id
  location                        = azurerm_cosmosdb_cassandra_cluster.test.location
  delegated_management_subnet_id  = azurerm_subnet.test.id
  node_count                      = %d
  disk_count                      = 4
  sku_name                        = "Standard_DS14_v2"
  availability_zones_enabled      = false
  backup_storage_customer_key_uri = azurerm_key_vault_key.test2.id
  managed_disk_customer_key_uri   = azurerm_key_vault_key.test2.id
  base64_encoded_yaml_fragment    = "Z29tcGFjdGlvbl90aHJvdWdocHV0X21iX3Blcl9zZWM6IDMyCmNvbXBhY3Rpb25fbGFyZ2VfcGFydGl0aW9uX3dhcm5pbmdfdGhyZXNob2xkX21iOiAxMDA="

  depends_on = [
    azurerm_key_vault_key.test,
    azurerm_key_vault_key.test2
  ]
}
`, r.template(data), data.RandomString, data.RandomString, data.RandomString, data.RandomInteger, nodeCount)
}

func (CassandraDatacenterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ca-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

data "azuread_service_principal" "test" {
  display_name = "Azure Cosmos DB"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = data.azuread_service_principal.test.object_id
}

resource "azurerm_cosmosdb_cassandra_cluster" "test" {
  name                           = "acctca-mi-cluster-%d"
  resource_group_name            = azurerm_resource_group.test.name
  location                       = azurerm_resource_group.test.location
  delegated_management_subnet_id = azurerm_subnet.test.id
  default_admin_password         = "Password1234"

  identity {
    type = "SystemAssigned"
  }

  depends_on = [azurerm_role_assignment.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
