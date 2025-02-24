// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/administrators"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PostgresqlFlexibleServerAdministratorResource struct{}

func TestAccPostgresqlFlexibleServerAdministrator_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_active_directory_administrator", "test")
	r := PostgresqlFlexibleServerAdministratorResource{}

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

func TestAccPostgresqlFlexibleServerAdministrator_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_active_directory_administrator", "test")
	r := PostgresqlFlexibleServerAdministratorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_postgresql_flexible_server_active_directory_administrator"),
		},
	})
}

func TestAccPostgresqlFlexibleServerAdministrator_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_active_directory_administrator", "test")
	r := PostgresqlFlexibleServerAdministratorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func (r PostgresqlFlexibleServerAdministratorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := administrators.ParseAdministratorID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Postgres.FlexibleServerAdministratorsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Postgresql AAD Administrator (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r PostgresqlFlexibleServerAdministratorResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := administrators.ParseAdministratorID(state.ID)
	if err != nil {
		return nil, err
	}

	if _, err := client.Postgres.FlexibleServerAdministratorsClient.Delete(ctx, *id); err != nil {
		return nil, fmt.Errorf("deleting Postgresql AAD Administrator (%s): %+v", id.String(), err)
	}

	return utils.Bool(true), nil
}

func (PostgresqlFlexibleServerAdministratorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azurerm_client_config" "current" {}

data "azuread_service_principal" "test" {
  object_id = data.azurerm_client_config.current.object_id
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%[1]d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  storage_mb             = 32768
  version                = "12"
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "2"

  authentication {
    active_directory_auth_enabled = true
    tenant_id                     = data.azurerm_client_config.current.tenant_id
  }

}

resource "azurerm_postgresql_flexible_server_active_directory_administrator" "test" {
  server_name         = azurerm_postgresql_flexible_server.test.name
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  object_id           = data.azuread_service_principal.test.object_id
  principal_name      = data.azuread_service_principal.test.display_name
  principal_type      = "ServicePrincipal"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r PostgresqlFlexibleServerAdministratorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server_active_directory_administrator" "import" {
  server_name         = azurerm_postgresql_flexible_server_active_directory_administrator.test.server_name
  resource_group_name = azurerm_postgresql_flexible_server_active_directory_administrator.test.resource_group_name
  tenant_id           = azurerm_postgresql_flexible_server_active_directory_administrator.test.tenant_id
  object_id           = azurerm_postgresql_flexible_server_active_directory_administrator.test.object_id
  principal_name      = azurerm_postgresql_flexible_server_active_directory_administrator.test.principal_name
  principal_type      = azurerm_postgresql_flexible_server_active_directory_administrator.test.principal_type
}
`, r.basic(data))
}
