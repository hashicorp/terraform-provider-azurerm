package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPostgreSQLServerKeys_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server_key", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerKeysDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServerKeys_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerKeysExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "key_type", "AzureKeyVault"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPostgreSQLServerKeys_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server_key", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerKeysDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServerKeys_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerKeysExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "key_type", "AzureKeyVault"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMPostgreSQLServerKeys_requiresImport),
		},
	})
}

func testCheckAzureRMPostgreSQLServerKeysExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Postgres.ServerKeysClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for PostgreSQL server key: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: PostgreSQL server key %q (server %q resource group: %q) does not exist", name, serverName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on postgresqlServerKeyssClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMPostgreSQLServerKeysDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Postgres.DatabasesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_postgresql_server_key" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, serverName, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("PostgreSQL server key still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMPostgreSQLServerKeys_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                        = "acctestkv%s"
  location                    = azurerm_resource_group.test.location
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  resource_group_name         = azurerm_resource_group.test.name
  sku_name                    = "standard"
  enabled_for_disk_encryption = true
  soft_delete_enabled         = true
  purge_protection_enabled    = true

  #network_acls {}
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  object_id    = azurerm_postgresql_server.test.identity[0].principal_id
  tenant_id    = data.azurerm_client_config.current.tenant_id

  key_permissions = [
    "get", 
    "unwrapKey", 
    "wrapKey"
  ]
}

resource "azurerm_key_vault_access_policy" "test_tf" {
  key_vault_id = azurerm_key_vault.test.id
  object_id    = data.azurerm_client_config.current.object_id
  tenant_id    = data.azurerm_client_config.current.tenant_id

  key_permissions = [
    "get", 
    "create", 
    "delete"
  ]
}

resource "azurerm_key_vault_key" "generated" {
  name         = "acctest-generated-key-%d"
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

  depends_on = ["azurerm_key_vault_access_policy.test_tf"]
}

resource "azurerm_postgresql_server" "test" {
  name                              = "acctest-psql-server-%d"
  location                          = azurerm_resource_group.test.location
  resource_group_name               = azurerm_resource_group.test.name

  sku_name                          = "GP_Gen5_2"
  storage_mb                        = 51200
  backup_retention_days             = 7
  geo_redundant_backup_enabled      = false
  create_mode                       = "Default"
  public_network_access_enabled     = true
  ssl_minimal_tls_version_enforced  = "TLSEnforcementDisabled"
  infrastructure_encryption_enabled = false

  identity {
    type = "SystemAssigned"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.6"
  ssl_enforcement_enabled      = true

}

resource "azurerm_postgresql_server_key" "test" {
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_postgresql_server.test.name
  key_type            = "AzureKeyVault"
  key_url             = azurerm_key_vault_key.generated.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPostgreSQLServerKeys_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMPostgreSQLServerKeys_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_server_key" "import" {
  resource_group_name = azurerm_postgresql_server_key.test.resource_group_name
  server_name         = azurerm_postgresql_server_key.test.server_name
  key_type            = azurerm_postgresql_server_key.test.key_type
  key_url             = azurerm_postgresql_server_key.test.key_url
}
`, template)
}
