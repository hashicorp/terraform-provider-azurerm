package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMySQLServerKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_server_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServerKey_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerKeyExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMMySQLServerKey_updateKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_server_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServerKey_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerKeyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMySQLServerKey_updated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerKeyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMySQLServerKey_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_server_key", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServerKey_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerKeyExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMMySQLServerKey_requiresImport),
		},
	})
}

func testCheckAzureRMMySQLServerKeyDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MySQL.ServerKeysClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mysql_server_key" {
			continue
		}

		id, err := parse.MySQLServerKeyID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("retrieving MySQL Server Key: %+v", err)
			}
			return nil
		}

		return fmt.Errorf("MySQL Server Key still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMMySQLServerKeyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MySQL.ServerKeysClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.MySQLServerKeyID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MySQL Server Key %q (Resource Group %q / Server %q) does not exist", id.Name, id.ResourceGroup, id.ServerName)
			}
			return fmt.Errorf("Bad: Get on MySQLServerKeysClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMMySQLServerKey_template(data acceptance.TestData) string {
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
  key_vault_id       = azurerm_key_vault.test.id
  tenant_id          = data.azurerm_client_config.current.tenant_id
  object_id          = azurerm_mysql_server.test.identity.0.principal_id
  key_permissions    = ["get", "unwrapkey", "wrapkey"]
  secret_permissions = ["get"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id       = azurerm_key_vault.test.id
  tenant_id          = data.azurerm_client_config.current.tenant_id
  object_id          = data.azurerm_client_config.current.object_id
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

resource "azurerm_mysql_server" "test" {
  name                             = "acctestmysqlsvr-%d"
  location                         = azurerm_resource_group.test.location
  resource_group_name              = azurerm_resource_group.test.name
  sku_name                         = "GP_Gen5_2"
  administrator_login              = "acctestun"
  administrator_login_password     = "H@Sh1CoR3!"
  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "TLS1_1"
  storage_mb                       = 51200
  version                          = "5.6"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func testAccAzureRMMySQLServerKey_basic(data acceptance.TestData) string {
	template := testAccAzureRMMySQLServerKey_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_server_key" "test" {
  server_id        = azurerm_mysql_server.test.id
  key_vault_key_id = azurerm_key_vault_key.first.id
}
`, template)
}

func testAccAzureRMMySQLServerKey_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMySQLServerKey_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_server_key" "import" {
  server_id        = azurerm_mysql_server_key.test.server_id
  key_vault_key_id = azurerm_mysql_server_key.test.key_vault_key_id
}
`, template)
}

func testAccAzureRMMySQLServerKey_updated(data acceptance.TestData) string {
	template := testAccAzureRMMySQLServerKey_template(data)
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
resource "azurerm_mysql_server_key" "test" {
  server_id        = azurerm_mysql_server.test.id
  key_vault_key_id = azurerm_key_vault_key.second.id
}
`, template)
}
