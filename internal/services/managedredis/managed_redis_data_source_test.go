// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedredis_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ManagedRedisDataSource struct{}

func TestAccManagedRedisDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_redis", "test")
	r := ManagedRedisDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.dataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("customer_managed_key.#").HasValue("1"),
				check.That(data.ResourceName).Key("customer_managed_key.0.key_vault_key_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("customer_managed_key.0.user_assigned_identity_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("default_database.#").HasValue("1"),
				check.That(data.ResourceName).Key("default_database.0.access_keys_authentication_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("default_database.0.client_protocol").HasValue(string(redisenterprise.ProtocolEncrypted)),
				check.That(data.ResourceName).Key("default_database.0.clustering_policy").HasValue(string(redisenterprise.ClusteringPolicyOSSCluster)),
				check.That(data.ResourceName).Key("default_database.0.eviction_policy").HasValue(string(redisenterprise.EvictionPolicyVolatileLRU)),
				check.That(data.ResourceName).Key("default_database.0.geo_replication_group_name").HasValue(fmt.Sprintf("acctest-geo-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("default_database.0.geo_replication_linked_database_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("default_database.0.module.#").HasValue("0"),
				check.That(data.ResourceName).Key("default_database.0.persistence.#").HasValue("0"),
				check.That(data.ResourceName).Key("default_database.0.port").HasValue("10000"),
				check.That(data.ResourceName).Key("default_database.0.primary_access_key").IsNotEmpty(),
				check.That(data.ResourceName).Key("default_database.0.secondary_access_key").IsNotEmpty(),
				check.That(data.ResourceName).Key("high_availability_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("hostname").IsNotEmpty(),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
				check.That(data.ResourceName).Key("sku_name").HasValue(string(redisenterprise.SkuNameBalancedBThree)),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.env").HasValue("testing"),
			),
		},
	})
}

func TestAccManagedRedisDataSource_dbPersistence(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_redis", "test")
	r := ManagedRedisDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.dataSourceDbPersistence(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("default_database.0.persistence.0.redis_database_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("default_database.0.persistence.0.redis_database_backup_frequency").HasValue("1h"),
			),
		},
	})
}

func (ManagedRedisDataSource) dataSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-managedRedis-%[1]d"
  location = "%[2]s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_key_vault" "test" {
  name                = "acctestAmr%[3]s"
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

resource "azurerm_managed_redis" "test" {
  name                = "acctest-amr-%[1]d"
  resource_group_name = azurerm_resource_group.test.name

  location = "%[2]s"
  sku_name = "Balanced_B3"

  customer_managed_key {
    key_vault_key_id          = azurerm_key_vault_key.test.id
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }

  default_database {
    access_keys_authentication_enabled = true
    geo_replication_group_name         = "acctest-geo-%[1]d"
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    env = "testing"
  }
}

data "azurerm_managed_redis" "test" {
  name                = azurerm_managed_redis.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (ManagedRedisDataSource) dataSourceDbPersistence(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-managedRedis-%[1]d"
  location = "%[2]s"
}

resource "azurerm_managed_redis" "test" {
  name                = "acctest-amr-%[1]d"
  resource_group_name = azurerm_resource_group.test.name

  location = "%[2]s"
  sku_name = "Balanced_B0"

  default_database {
    persistence {
      redis_database_enabled          = true
      redis_database_backup_frequency = "1h"
    }
  }
}

data "azurerm_managed_redis" "test" {
  name                = azurerm_managed_redis.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
