// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SqlMiAdministratorResource struct{}

func TestAccSqlMiAdministrator_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_managed_instance_active_directory_administrator", "test")
	r := SqlMiAdministratorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.template(data),
		},
		{
			PreConfig: func() { time.Sleep(5 * time.Minute) },
			Config:    r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithAadAuthOnlyEqualTo(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithAadAuthOnlyEqualTo(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSqlMiAdministrator_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_managed_instance_active_directory_administrator", "test")
	r := SqlMiAdministratorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.template(data),
		},
		{
			PreConfig: func() { time.Sleep(5 * time.Minute) },
			Config:    r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r SqlMiAdministratorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ManagedInstanceAzureActiveDirectoryAdministratorID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Sql.ManagedInstanceAdministratorsClient.Get(ctx, id.ResourceGroup, id.ManagedInstanceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r SqlMiAdministratorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sql_managed_instance_active_directory_administrator" "test" {
  managed_instance_name = azurerm_sql_managed_instance.test.name
  resource_group_name   = azurerm_resource_group.test.name
  login                 = data.azuread_service_principal.test.display_name
  tenant_id             = data.azurerm_client_config.current.tenant_id
  object_id             = data.azurerm_client_config.current.client_id

  depends_on = [azuread_directory_role_member.test]
}
`, r.template(data))
}

func (r SqlMiAdministratorResource) basicWithAadAuthOnlyEqualTo(data acceptance.TestData, aadAuthOnly bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sql_managed_instance_active_directory_administrator" "test" {
  managed_instance_name       = azurerm_sql_managed_instance.test.name
  resource_group_name         = azurerm_resource_group.test.name
  login                       = data.azuread_service_principal.test.display_name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  object_id                   = data.azurerm_client_config.current.client_id
  azuread_authentication_only = %t
}
`, r.template(data), aadAuthOnly)
}

func (r SqlMiAdministratorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sql_managed_instance_active_directory_administrator" "import" {
  managed_instance_name = azurerm_sql_managed_instance_active_directory_administrator.test.managed_instance_name
  resource_group_name   = azurerm_sql_managed_instance_active_directory_administrator.test.resource_group_name
  login                 = azurerm_sql_managed_instance_active_directory_administrator.test.login
  tenant_id             = azurerm_sql_managed_instance_active_directory_administrator.test.tenant_id
  object_id             = azurerm_sql_managed_instance_active_directory_administrator.test.object_id
}
`, r.basic(data))
}

func (r SqlMiAdministratorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azuread_directory_role" "reader" {
  display_name = "Directory Readers"
}

data "azurerm_client_config" "current" {}

data "azuread_service_principal" "test" {
  object_id = data.azurerm_client_config.current.object_id
}

resource "azuread_directory_role_member" "test" {
  role_object_id   = azuread_directory_role.reader.object_id
  member_object_id = azurerm_sql_managed_instance.test.identity.0.principal_id
}
`, SqlManagedInstanceResource{}.identity(data))
}
