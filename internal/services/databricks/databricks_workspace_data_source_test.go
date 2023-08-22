// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databricks_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DatabricksWorkspaceDataSource struct{}

func TestAccDatabricksWorkspaceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestMatchResourceAttr(data.ResourceName, "workspace_url", regexp.MustCompile("azuredatabricks.net")),
				check.That(data.ResourceName).Key("workspace_id").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
			),
		},
	})
}

func TestAccDatabricksWorkspaceDataSource_storage_account_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.storageAccountIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestMatchResourceAttr(data.ResourceName, "workspace_url", regexp.MustCompile("azuredatabricks.net")),
				check.That(data.ResourceName).Key("workspace_id").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
			),
		},
	})
}

func (DatabricksWorkspaceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%d"
  location = "%s"
}

resource "azurerm_databricks_workspace" "test" {
  name                = "acctestDBW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "standard"
}

data "azurerm_databricks_workspace" "test" {
  name                = azurerm_databricks_workspace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (DatabricksWorkspaceDataSource) storageAccountIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_databricks_workspace" "test" {
  depends_on = [azurerm_key_vault_access_policy.managed]

  name                        = "acctestDBW-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku                         = "premium"
  managed_resource_group_name = "acctestRG-DBW-%[1]d-managed"

  customer_managed_key_enabled          = true
  managed_services_cmk_key_vault_key_id = azurerm_key_vault_key.test.id

  tags = {
    Environment = "Production"
    Pricing     = "Premium"
  }
}

resource "azurerm_databricks_workspace_root_dbfs_customer_managed_key" "test" {
  depends_on = [azurerm_key_vault_access_policy.databricks]

  workspace_id     = azurerm_databricks_workspace.test.id
  key_vault_key_id = azurerm_key_vault_key.test.id
}

resource "azurerm_key_vault" "test" {
  name                = "acctest-kv-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_key" "test" {
  depends_on = [azurerm_key_vault_access_policy.terraform]

  name         = "acctest-certificate"
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

resource "azurerm_key_vault_access_policy" "terraform" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_key_vault.test.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "List",
    "Create",
    "Decrypt",
    "Encrypt",
    "GetRotationPolicy",
    "Sign",
    "UnwrapKey",
    "Verify",
    "WrapKey",
    "Delete",
    "Restore",
    "Recover",
    "Update",
    "Purge",
  ]
}

resource "azurerm_key_vault_access_policy" "managed" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_key_vault.test.tenant_id
  object_id    = "%[4]s"

  key_permissions = [
    "Get",
    "GetRotationPolicy",
    "UnwrapKey",
    "WrapKey",
  ]
}

data "azurerm_databricks_workspace" "test" {
  name                = azurerm_databricks_workspace.test.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault_access_policy" "databricks" {
  depends_on = [azurerm_databricks_workspace.test]

  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_databricks_workspace.test.storage_account_identity.0.tenant_id
  object_id    = data.azurerm_databricks_workspace.test.storage_account_identity.0.principal_id

  key_permissions = [
    "Get",
    "GetRotationPolicy",
    "UnwrapKey",
    "WrapKey",
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, getDatabricksPrincipalId(data.Client().SubscriptionID))
}
