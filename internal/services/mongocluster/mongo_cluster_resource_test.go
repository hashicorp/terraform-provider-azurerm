// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongocluster_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2025-09-01/mongoclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MongoClusterResource struct{}

func TestAccMongoClusterFreeTier(t *testing.T) {
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"freeTier": { // Run tests in sequence since each subscription is limited to one free tier cluster per region and free tier is currently only available in South India.
			"basic":  testAccMongoCluster_basic,
			"update": testAccMongoCluster_update,
			"import": testAccMongoCluster_requiresImport,
		},
	})
}

func testAccMongoCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mongo_cluster", "test")
	r := MongoClusterResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("connection_strings.0.value").HasValue(
					fmt.Sprintf(`mongodb+srv://adminTerraform:QAZwsx123basic@acctest-mc%d.global.mongocluster.cosmos.azure.com/?tls=true&authMechanism=SCRAM-SHA-256&retrywrites=false&maxIdleTimeMS=120000`,
						data.RandomInteger)),
				check.That(data.ResourceName).Key("connection_strings.1.value").HasValue(
					fmt.Sprintf(`mongodb+srv://adminTerraform:QAZwsx123basic@acctest-mc%d.mongocluster.cosmos.azure.com/?tls=true&authMechanism=SCRAM-SHA-256&retrywrites=false&maxIdleTimeMS=120000`,
						data.RandomInteger),
				),
			),
		},
		data.ImportStep("administrator_password", "create_mode", "connection_strings.0.value", "connection_strings.1.value"),
	})
}

func testAccMongoCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mongo_cluster", "test")
	r := MongoClusterResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode", "connection_strings.0.value", "connection_strings.1.value"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode", "connection_strings.0.value", "connection_strings.1.value"),
	})
}

func testAccMongoCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mongo_cluster", "test")
	r := MongoClusterResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMongoCluster_previewFeature(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mongo_cluster", "test")
	r := MongoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.previewFeature(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode", "connection_strings.0.value", "connection_strings.1.value"),
		{
			Config: r.geoReplica(data, r.previewFeature(data)),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode", "source_location", "connection_strings.0.value", "connection_strings.1.value"),
	})
}

func TestAccMongoCluster_geoReplica(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mongo_cluster", "test")
	r := MongoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.geoReplica(data, r.source(data)),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode", "source_location", "connection_strings.0.value", "connection_strings.1.value"),
	})
}

func TestAccMongoCluster_cmkKeyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mongo_cluster", "test")
	r := MongoClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cmkKeyVault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode", "source_location", "connection_strings.0.value", "connection_strings.1.value"),
		{
			Config: r.cmkKeyVaultUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode", "source_location", "connection_strings.0.value", "connection_strings.1.value"),
		{
			Config:      r.cmkKeyVaultUpdateError(data),
			ExpectError: regexp.MustCompile("the value specified in `customer_managed_key.user_assigned_identity_id` must also be specified in `identity.identity_ids`"),
		},
	})
}

func (r MongoClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := mongoclusters.ParseMongoClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MongoCluster.MongoClustersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r MongoClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mongo_cluster" "test" {
  name                   = "acctest-mc%s"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_username = "adminTerraform"
  administrator_password = "QAZwsx123basic"
  shard_count            = "1"
  compute_tier           = "Free"
  high_availability_mode = "Disabled"
  storage_size_in_gb     = "32"
  version                = "7.0"
}
`, r.template(data, data.Locations.Ternary), data.RandomString)
}

func (r MongoClusterResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mongo_cluster" "test" {
  name                   = "acctest-mc%s"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_username = "adminTerraform"
  administrator_password = "QAZwsx123update"
  shard_count            = "1"
  compute_tier           = "M30"
  high_availability_mode = "ZoneRedundantPreferred"
  public_network_access  = "Disabled"
  storage_size_in_gb     = "64"
  version                = "8.0"

  tags = {
    environment = "test"
  }
}
`, r.template(data, data.Locations.Ternary), data.RandomString)
}

func (r MongoClusterResource) cmkKeyVault(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mongo_cluster" "test" {
  name                = "acctest-mc%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  administrator_username = "adminTerraform"
  administrator_password = "QAZwsx123"

  shard_count            = "1"
  compute_tier           = "M30"
  high_availability_mode = "ZoneRedundantPreferred"
  public_network_access  = "Disabled"
  storage_size_in_gb     = "64"
  version                = "8.0"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  customer_managed_key {
    key_vault_key_id          = azurerm_key_vault_key.test.versionless_id
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }

  tags = {
    environment = "test"
  }
}
`, r.templateCMKKeyVault(data), data.RandomString)
}

func (r MongoClusterResource) cmkKeyVaultUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mongo_cluster" "test" {
  name                = "acctest-mc%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  administrator_username = "adminTerraform"
  administrator_password = "QAZwsx123"

  shard_count            = "1"
  compute_tier           = "M30"
  high_availability_mode = "ZoneRedundantPreferred"
  public_network_access  = "Disabled"
  storage_size_in_gb     = "64"
  version                = "8.0"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test2.id]
  }

  customer_managed_key {
    key_vault_key_id          = azurerm_key_vault_key.test2.versionless_id
    user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  }

  tags = {
    environment = "test"
  }
}
`, r.templateCMKKeyVault(data), data.RandomString)
}

func (r MongoClusterResource) cmkKeyVaultUpdateError(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test3" {
  name                = "acctestUAI3-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}


resource "azurerm_mongo_cluster" "test" {
  name                = "acctest-mc%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  administrator_username = "adminTerraform"
  administrator_password = "QAZwsx123"

  shard_count            = "1"
  compute_tier           = "M30"
  high_availability_mode = "ZoneRedundantPreferred"
  public_network_access  = "Disabled"
  storage_size_in_gb     = "64"
  version                = "8.0"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test3.id]
  }

  customer_managed_key {
    key_vault_key_id          = azurerm_key_vault_key.test2.versionless_id
    user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  }

  tags = {
    environment = "test"
  }
}
`, r.templateCMKKeyVault(data), data.RandomInteger, data.RandomString)
}

func (r MongoClusterResource) source(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mongo_cluster" "test" {
  name                   = "acctest-mc%s"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_username = "adminTerraform"
  administrator_password = "QAZwsx123update"
  high_availability_mode = "ZoneRedundantPreferred"
  shard_count            = "1"
  compute_tier           = "M30"
  storage_size_in_gb     = "64"
  version                = "8.0"
}
`, r.template(data, data.Locations.Primary), data.RandomString)
}

func (r MongoClusterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mongo_cluster" "import" {
  name                   = azurerm_mongo_cluster.test.name
  resource_group_name    = azurerm_mongo_cluster.test.resource_group_name
  location               = azurerm_mongo_cluster.test.location
  administrator_username = azurerm_mongo_cluster.test.administrator_username
  administrator_password = azurerm_mongo_cluster.test.administrator_password
  shard_count            = azurerm_mongo_cluster.test.shard_count
  compute_tier           = azurerm_mongo_cluster.test.compute_tier
  high_availability_mode = azurerm_mongo_cluster.test.high_availability_mode
  storage_size_in_gb     = azurerm_mongo_cluster.test.storage_size_in_gb
  version                = azurerm_mongo_cluster.test.version
}
`, r.basic(data))
}

func (r MongoClusterResource) previewFeature(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mongo_cluster" "test" {
  name                   = "acctest-mc%s"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_username = "adminTerraform"
  administrator_password = "testQAZwsx123"
  shard_count            = "1"
  compute_tier           = "M30"
  high_availability_mode = "ZoneRedundantPreferred"
  storage_size_in_gb     = "64"
  preview_features       = ["GeoReplicas"]
  version                = "8.0"
}
`, r.template(data, data.Locations.Primary), data.RandomString)
}

func (r MongoClusterResource) geoReplica(data acceptance.TestData, source string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mongo_cluster" "geo_replica" {
  name                = "acctest-mc-replica%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  source_server_id    = azurerm_mongo_cluster.test.id
  source_location     = azurerm_mongo_cluster.test.location
  create_mode         = "GeoReplica"

  lifecycle {
    ignore_changes = ["administrator_username", "high_availability_mode", "preview_features", "shard_count", "storage_size_in_gb", "compute_tier", "version"]
  }
}
`, source, data.RandomString, data.Locations.Secondary)
}

func (r MongoClusterResource) template(data acceptance.TestData, location string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctestUAI2-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, location)
}

func (r MongoClusterResource) templateKeyVaultBase(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault" "test" {
  name                       = "acctestKV-%[2]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  purge_protection_enabled   = true
  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions         = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy", "SetRotationPolicy"]
  secret_permissions      = ["Delete", "Get", "Set"]
  certificate_permissions = ["Create", "Delete", "DeleteIssuers", "Get", "Purge", "Update"]
}
`, r.template(data, data.Locations.Secondary), data.RandomString)
}

func (r MongoClusterResource) templateCMKKeyVault(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault_access_policy" "uai" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id

  key_permissions = ["Get", "List", "WrapKey", "UnwrapKey", "GetRotationPolicy", "SetRotationPolicy"]
}

resource "azurerm_key_vault_access_policy" "uai2" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_user_assigned_identity.test2.principal_id

  key_permissions = ["Get", "List", "WrapKey", "UnwrapKey", "GetRotationPolicy", "SetRotationPolicy"]
}

resource "azurerm_key_vault_key" "test" {
  name         = "key1-%[2]s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.uai,
    azurerm_key_vault_access_policy.uai2,
  ]
}

resource "azurerm_key_vault_key" "test2" {
  name         = "key2-%[2]s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.uai,
    azurerm_key_vault_access_policy.uai2,
  ]
}
`, r.templateKeyVaultBase(data), data.RandomString, data.RandomInteger)
}
