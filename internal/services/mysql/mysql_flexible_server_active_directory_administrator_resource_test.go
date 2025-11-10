// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/azureadadministrators"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MySQLFlexibleServerAdministratorResource struct{}

func TestAccMySQLFlexibleServerAdministrator_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server_active_directory_administrator", "test")
	r := MySQLFlexibleServerAdministratorResource{}

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

func TestAccMySQLFlexibleServerAdministrator_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server_active_directory_administrator", "test")
	r := MySQLFlexibleServerAdministratorResource{}

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

func TestAccMySQLFlexibleServerAdministrator_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server_active_directory_administrator", "test")
	r := MySQLFlexibleServerAdministratorResource{}

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
	})
}

func (r MySQLFlexibleServerAdministratorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FlexibleServerAzureActiveDirectoryAdministratorID(state.ID)
	if err != nil {
		return nil, err
	}

	flexibleServerId := azureadadministrators.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroup, id.FlexibleServerName)

	client := clients.MySQL.AzureADAdministratorsClient
	resp, err := client.Get(ctx, flexibleServerId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r MySQLFlexibleServerAdministratorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mysqlfsaad-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctestUAI2-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-mysqlfs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "_admin_Terraform_892123456789312"
  administrator_password = "QAZwsx123"
  sku_name               = "B_Standard_B1ms"
  zone                   = "2"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id, azurerm_user_assigned_identity.test2.id]
  }
}
`, data.RandomInteger, data.Locations.Ternary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r MySQLFlexibleServerAdministratorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server_active_directory_administrator" "test" {
  server_id   = azurerm_mysql_flexible_server.test.id
  identity_id = azurerm_user_assigned_identity.test.id
  login       = "sqladmin"
  object_id   = data.azurerm_client_config.current.client_id
  tenant_id   = data.azurerm_client_config.current.tenant_id
}
`, r.template(data))
}

func (r MySQLFlexibleServerAdministratorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server_active_directory_administrator" "import" {
  server_id   = azurerm_mysql_flexible_server_active_directory_administrator.test.server_id
  identity_id = azurerm_mysql_flexible_server_active_directory_administrator.test.identity_id
  login       = azurerm_mysql_flexible_server_active_directory_administrator.test.login
  object_id   = azurerm_mysql_flexible_server_active_directory_administrator.test.object_id
  tenant_id   = azurerm_mysql_flexible_server_active_directory_administrator.test.tenant_id
}
`, r.basic(data))
}

func (r MySQLFlexibleServerAdministratorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server_active_directory_administrator" "test" {
  server_id   = azurerm_mysql_flexible_server.test.id
  identity_id = azurerm_user_assigned_identity.test2.id
  login       = "sqladmin2"
  object_id   = data.azurerm_client_config.current.client_id
  tenant_id   = data.azurerm_client_config.current.tenant_id
}
`, r.template(data))
}
