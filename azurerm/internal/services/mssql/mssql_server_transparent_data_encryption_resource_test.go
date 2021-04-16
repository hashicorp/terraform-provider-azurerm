package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MsSqlServerTransparentDataEncryptionResource struct{}

func TestAccMsSqlServerTransparentDataEncryption_keyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server_transparent_data_encryption", "test")
	r := MsSqlServerTransparentDataEncryptionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.keyVault(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlServerTransparentDataEncryption_systemManaged(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server_transparent_data_encryption", "test")
	r := MsSqlServerTransparentDataEncryptionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.systemManaged(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_vault_key_id").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlServerTransparentDataEncryption_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server_transparent_data_encryption", "test")
	r := MsSqlServerTransparentDataEncryptionResource{}

	// Test going from systemManaged to keyVault and back
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.keyVault(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.systemManaged(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_vault_key_id").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func (MsSqlServerTransparentDataEncryptionResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.EncryptionProtectorID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.EncryptionProtectorClient.Get(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil, fmt.Errorf("Encryption protector for server %q (Resource Group %q) does not exist", id.ServerName, id.ResourceGroup)
		}

		return nil, fmt.Errorf("reading Encryption Protector for server %q (Resource Group %q): %v", id.ServerName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r MsSqlServerTransparentDataEncryptionResource) keyVault(data acceptance.TestData) string {
	return fmt.Sprintf(
		`
%s

resource "azurerm_key_vault" "test" {
	name                        = "acctestsqlserver%[2]s"
	location                    = azurerm_resource_group.test.location
	resource_group_name         = azurerm_resource_group.test.name
	enabled_for_disk_encryption = true
	tenant_id                   = data.azurerm_client_config.current.tenant_id
	soft_delete_retention_days  = 7
	purge_protection_enabled    = false

	sku_name = "standard"

	access_policy {
	  tenant_id    = data.azurerm_client_config.current.tenant_id
	  object_id    = data.azurerm_client_config.current.object_id
  
	  key_permissions = [
		"Get",  "List", "Create", "Delete", "Update", "Purge", 
	  ]
	}

	access_policy {
	tenant_id = azurerm_mssql_server.test.identity[0].tenant_id
	object_id = azurerm_mssql_server.test.identity[0].principal_id
  
	key_permissions = [
		"Get", "WrapKey", "UnwrapKey", "List", "Create", 
	  ]
	}
  }

  resource "azurerm_key_vault_key" "generated" {
	name         = "keyVault"
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
	  azurerm_key_vault.test,
	]
  }

  resource "azurerm_mssql_server_transparent_data_encryption" "test" {
	server_id = azurerm_mssql_server.test.id
	key_vault_key_id = azurerm_key_vault_key.generated.id
  }
`, r.server(data), data.RandomStringOfLength(5))
}

func (r MsSqlServerTransparentDataEncryptionResource) systemManaged(data acceptance.TestData) string {
	return fmt.Sprintf(
		`
%s

  resource "azurerm_mssql_server_transparent_data_encryption" "test" {
	server_id = azurerm_mssql_server.test.id
  }
`, r.server(data))
}

func (MsSqlServerTransparentDataEncryptionResource) server(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"

  identity {
    type = "SystemAssigned"
  }
}


`, data.RandomInteger, data.Locations.Primary)
}
