// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2020-01-01/serverkeys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PostgreSQLServerKeyResource struct{}

func TestAccPostgreSQLServerKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server_key", "test")
	r := PostgreSQLServerKeyResource{}

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

func TestAccPostgreSQLServerKey_updateKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server_key", "test")
	r := PostgreSQLServerKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPostgreSQLServerKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server_key", "test")
	r := PostgreSQLServerKeyResource{}

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

func TestAccPostgreSQLServerKey_replica(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server_key", "test")
	r := PostgreSQLServerKeyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.replica(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t PostgreSQLServerKeyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := serverkeys.ParseKeyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Postgres.ServerKeysClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Postgresql Server Key (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (PostgreSQLServerKeyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "server" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_postgresql_server.test.identity.0.principal_id

  key_permissions    = ["Get", "UnwrapKey", "WrapKey", "GetRotationPolicy", "SetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy", "SetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_key" "first" {
  name         = "first"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.server,
  ]
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-postgre-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name   = "GP_Gen5_2"
  version    = "11"
  storage_mb = 51200

  ssl_enforcement_enabled = true

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r PostgreSQLServerKeyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_server_key" "test" {
  server_id        = azurerm_postgresql_server.test.id
  key_vault_key_id = azurerm_key_vault_key.first.id
}
`, r.template(data))
}

func (r PostgreSQLServerKeyResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_key" "second" {
  name         = "second"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.server,
  ]
}

resource "azurerm_postgresql_server_key" "test" {
  server_id        = azurerm_postgresql_server.test.id
  key_vault_key_id = azurerm_key_vault_key.second.id
}
`, r.template(data))
}

func (r PostgreSQLServerKeyResource) replica(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_server_key" "test" {
  server_id        = azurerm_postgresql_server.test.id
  key_vault_key_id = azurerm_key_vault_key.first.id
}

resource "azurerm_key_vault_access_policy" "replica" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_postgresql_server.replica.identity.0.principal_id

  key_permissions    = ["Get", "UnwrapKey", "WrapKey", "GetRotationPolicy", "SetRotationPolicy"]
  secret_permissions = ["Get"]
}

resource "azurerm_resource_group" "replica" {
  name     = "acctestRG-psql-%d-replica"
  location = "%s"
}

resource "azurerm_postgresql_server" "replica" {
  name                = "acctest-postgre-replica-%d"
  location            = azurerm_resource_group.replica.location
  resource_group_name = azurerm_resource_group.replica.name

  create_mode               = "Replica"
  creation_source_server_id = azurerm_postgresql_server.test.id

  sku_name   = "GP_Gen5_2"
  version    = "11"
  storage_mb = 51200

  ssl_enforcement_enabled = true

  identity {
    type = "SystemAssigned"
  }
}


`, r.template(data), data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
