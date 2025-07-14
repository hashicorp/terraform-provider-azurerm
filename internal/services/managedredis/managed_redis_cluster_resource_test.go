// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedredis_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagedRedisClusterResource struct{}

func TestAccManagedRedisCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_cluster", "test")
	r := ManagedRedisClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedRedisCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_cluster", "test")
	r := ManagedRedisClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccManagedRedisCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_cluster", "test")
	r := ManagedRedisClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("hostname").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedRedisCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_cluster", "test")
	r := ManagedRedisClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedRedisCluster_withCmk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_cluster", "test")
	r := ManagedRedisClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withCmk(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedRedisCluster_withPrivateEndpoint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_cluster", "test")
	r := ManagedRedisClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withPrivateEndpoint(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ManagedRedisClusterResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := redisenterprise.ParseRedisEnterpriseID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.ManagedRedis.Client.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r ManagedRedisClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-managedRedis-%[1]d"
  location = "eastus"
}

resource "azurerm_managed_redis_cluster" "test" {
  name                = "acctest-mrc-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku_name = "Balanced_B3"
}
`, data.RandomInteger)
}

func (r ManagedRedisClusterResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-managedRedis-%[1]d"
  location = "eastus"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_managed_redis_cluster" "test" {
  name                = "acctest-mrc-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku_name = "Balanced_B3"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger)
}

func (r ManagedRedisClusterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_redis_cluster" "import" {
  name                = azurerm_managed_redis_cluster.test.name
  resource_group_name = azurerm_managed_redis_cluster.test.resource_group_name
  location            = azurerm_managed_redis_cluster.test.location

  sku_name = azurerm_managed_redis_cluster.test.sku_name
}
`, r.basic(data))
}

func (r ManagedRedisClusterResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-managedRedis-%[1]d"
  location = "eastus"
}

resource "azurerm_managed_redis_cluster" "test" {
  name                = "acctest-mrc-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  minimum_tls_version       = "1.2"
  high_availability_enabled = true

  sku_name = "Balanced_B3"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger)
}

func (r ManagedRedisClusterResource) withCmk(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-managedRedis-%[1]d"
  location = "eastus"
}

data "azurerm_client_config" "current" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_key_vault" "test" {
  name                = "acctestMngdRedis%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  purge_protection_enabled   = true
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Delete",
      "Get",
      "List",
      "Purge",
      "Recover",
      "Update",
      "GetRotationPolicy",
      "SetRotationPolicy"
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    key_permissions = [
      "Get",
      "WrapKey",
      "UnwrapKey"
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctest-key-%[1]d"
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
}

resource "azurerm_managed_redis_cluster" "test" {
  name                = "acctest-mrc-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku_name = "Balanced_B3"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  customer_managed_key {
    encryption_key_url        = azurerm_key_vault_key.test.id
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }
}
`, data.RandomInteger, data.RandomStringOfLength(5))
}

func (r ManagedRedisClusterResource) withPrivateEndpoint(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-managedRedis-%[1]d"
  location = "eastus"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_managed_redis_cluster" "test" {
  name                = "acctest-mrc-%[1]d"
  location            = azurerm_virtual_network.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Balanced_B3"
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-redis-pe-%[1]d"
  location            = azurerm_managed_redis_cluster.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test.id

  private_service_connection {
    name                           = "acctest-redis-psc-%[1]d"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_managed_redis_cluster.test.id
    subresource_names              = ["redisEnterprise"]
  }
}

	`, data.RandomInteger)
}
