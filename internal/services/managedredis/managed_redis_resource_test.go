// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedredis_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagedRedisResource struct{}

func TestAccManagedRedis_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis", "test")
	r := ManagedRedisResource{}
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

func TestAccManagedRedis_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis", "test")
	r := ManagedRedisResource{}
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

func TestAccManagedRedis_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis", "test")
	r := ManagedRedisResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedRedis_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis", "test")
	r := ManagedRedisResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		// Create B3 SKU without db
		{
			Config: r.template(data, "Balanced_B3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		// Create the db and update all non force-new prop
		{
			Config: r.update(data, true, "false"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		// Update the db
		{
			Config: r.update(data, true, "true"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		// remove the db
		{
			Config: r.update(data, false, "true"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedRedis_withPrivateEndpoint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis", "test")
	r := ManagedRedisResource{}
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

func TestAccManagedRedis_skuDoesNotSupportGeoReplication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis", "test")
	r := ManagedRedisResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.skuDoesNotSupportGeoReplication(),
			ExpectError: regexp.MustCompile(`SKU .* does not support geo-replication`),
		},
	})
}

func TestAccManagedRedis_moduleDoesNotSupportGeoReplication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis", "test")
	r := ManagedRedisResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.moduleDoesNotSupportGeoReplication(),
			ExpectError: regexp.MustCompile(`invalid module .*, only following modules are supported`),
		},
	})
}

func TestAccManagedRedis_hasToUseNoEvictionWithRediSearch(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis", "test")
	r := ManagedRedisResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.hasToUseNoEvictionWithRediSearch(),
			ExpectError: regexp.MustCompile(`invalid eviction_policy .*, when using RediSearch module, eviction_policy must be set to NoEviction`),
		},
	})
}

func TestAccManagedRedis_hasToUseEnterpriseClusteringWithRediSearch(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis", "test")
	r := ManagedRedisResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.hasToUseEnterpriseClusteringWithRediSearch(),
			ExpectError: regexp.MustCompile(`invalid clustering_policy .*, when using RediSearch module, clustering_policy must be set to EnterpriseCluster`),
		},
	})
}

func (r ManagedRedisResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r ManagedRedisResource) template(data acceptance.TestData, skuName string) string {
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
  sku_name = "%[3]s"
}
`, data.RandomInteger, data.Locations.Primary, skuName)
}

func (r ManagedRedisResource) basic(data acceptance.TestData) string {
	return r.template(data, "Balanced_B0")
}

func (r ManagedRedisResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_redis" "import" {
  name                = azurerm_managed_redis.test.name
  resource_group_name = azurerm_managed_redis.test.resource_group_name

  location = azurerm_managed_redis.test.location
  sku_name = azurerm_managed_redis.test.sku_name
}
`, r.basic(data))
}

func (r ManagedRedisResource) complete(data acceptance.TestData) string {
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

// Wait for key vault permissions to propagate
resource "time_sleep" "wait_key_vault_key" {
  depends_on = [azurerm_key_vault_key.test]

  create_duration = "30s"
}

resource "azurerm_managed_redis" "test" {
  name                = "acctest-amr-%[1]d"
  resource_group_name = azurerm_resource_group.test.name

  location = "%[2]s"
  sku_name = "Balanced_B3"

  customer_managed_key {
    encryption_key_url        = azurerm_key_vault_key.test.id
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }

  default_database {
    access_keys_authentication_enabled = true
    client_protocol                    = "Encrypted"
    clustering_policy                  = "EnterpriseCluster"
    eviction_policy                    = "NoEviction"
    geo_replication_group_name         = "acctest-amr-georep-%[1]d"

    module {
      name = "RediSearch"
      args = ""
    }

    module {
      name = "RedisJSON"
      args = ""
    }

    port = 10000
  }

  high_availability_enabled = true

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    ENV = "Test"
  }

  depends_on = [time_sleep.wait_key_vault_key]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(5))
}

func (r ManagedRedisResource) update(data acceptance.TestData, withDb bool, accessKeyAuthEnabled string) string {
	db := ""
	if withDb {
		db = fmt.Sprintf(`
  default_database {
    access_keys_authentication_enabled = %[1]s
    geo_replication_group_name         = "acctest-amr-georep-%[2]d"
  }
`, accessKeyAuthEnabled, data.RandomInteger)
	}

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

// Wait for key vault permissions to propagate
resource "time_sleep" "wait_key_vault_key" {
  depends_on = [azurerm_key_vault_key.test]

  create_duration = "30s"
}

resource "azurerm_managed_redis" "test" {
  name                = "acctest-amr-%[1]d"
  resource_group_name = azurerm_resource_group.test.name

  location = "%[2]s"
  sku_name = "Balanced_B3"

  customer_managed_key {
    encryption_key_url        = azurerm_key_vault_key.test.id
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }

%[4]s

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    ENV = "Test"
  }

  depends_on = [time_sleep.wait_key_vault_key]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(5), db)
}

func (r ManagedRedisResource) withPrivateEndpoint(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-managedRedis-%[1]d"
  location = "%[2]s"
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

resource "azurerm_managed_redis" "test" {
  name     = "acctest-amr-%[1]d"
  location = azurerm_virtual_network.test.location

  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Balanced_B3"

  default_database {
    geo_replication_group_name = "acctest-amr-georep-%[1]d"
  }
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-redis-pe-%[1]d"
  location            = azurerm_managed_redis.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test.id

  private_service_connection {
    name                           = "acctest-redis-psc-%[1]d"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_managed_redis.test.id
    subresource_names              = ["redisEnterprise"]
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ManagedRedisResource) skuDoesNotSupportGeoReplication() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_managed_redis" "test" {
  name     = "acctest-invalid"
  location = "eastus"

  resource_group_name = "my-rg"
  sku_name            = "Balanced_B0"

  default_database {
    geo_replication_group_name = "acctest-amr"
  }
}
`
}

func (r ManagedRedisResource) moduleDoesNotSupportGeoReplication() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_managed_redis" "test" {
  name     = "acctest-invalid"
  location = "eastus"

  resource_group_name = "my-rg"
  sku_name            = "Balanced_B3"

  default_database {
    geo_replication_group_name = "acctest-amr"

    module {
      name = "RedisTimeSeries"
    }
  }
}
`
}

func (r ManagedRedisResource) hasToUseNoEvictionWithRediSearch() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_managed_redis" "test" {
  name     = "acctest-invalid"
  location = "eastus"

  resource_group_name = "my-rg"
  sku_name            = "Balanced_B0"

  default_database {
    clustering_policy = "EnterpriseCluster"
    module {
      name = "RediSearch"
    }
    eviction_policy = "AllKeysLRU"
  }
}
`
}

func (r ManagedRedisResource) hasToUseEnterpriseClusteringWithRediSearch() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_managed_redis" "test" {
  name     = "acctest-invalid"
  location = "eastus"

  resource_group_name = "my-rg"
  sku_name            = "Balanced_B0"

  default_database {
    clustering_policy = "OSSCluster"
    module {
      name = "RediSearch"
    }
    eviction_policy = "NoEviction"
  }
}
`
}
