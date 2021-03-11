package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type PostgreSQLServerKeyResource struct {
}

func TestAccPostgreSQLServerKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server_key", "test")
	r := PostgreSQLServerKeyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPostgreSQLServerKey_updateKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server_key", "test")
	r := PostgreSQLServerKeyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPostgreSQLServerKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server_key", "test")
	r := PostgreSQLServerKeyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t PostgreSQLServerKeyResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ServerKeyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Postgres.ServerKeysClient.Get(ctx, id.ResourceGroup, id.ServerName, id.KeyName)
	if err != nil {
		return nil, fmt.Errorf("reading Postgresql Server Key (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (PostgreSQLServerKeyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = false
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
  soft_delete_enabled      = true
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "server" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_postgresql_server.test.identity.0.principal_id

  key_permissions    = ["get", "unwrapkey", "wrapkey"]
  secret_permissions = ["get"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions    = ["get", "create", "delete", "list", "restore", "recover", "unwrapkey", "wrapkey", "purge", "encrypt", "decrypt", "sign", "verify"]
  secret_permissions = ["get"]
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

func (r PostgreSQLServerKeyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_server_key" "import" {
  server_id        = azurerm_postgresql_server_key.test.server_id
  key_vault_key_id = azurerm_postgresql_server_key.test.key_vault_key_id
}
`, r.basic(data))
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
