// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlManagedInstanceActiveDirectoryAdministratorResource struct{}

func TestAccMsSqlManagedInstanceActiveDirectoryAdministrator_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_managed_instance_active_directory_administrator", "test")
	r := MsSqlManagedInstanceActiveDirectoryAdministratorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			PreConfig: func() { time.Sleep(5 * time.Minute) },
			Config:    r.basic(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func (r MsSqlManagedInstanceActiveDirectoryAdministratorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ManagedInstanceAzureActiveDirectoryAdministratorID(state.ID)
	if err != nil {
		return nil, err
	}

	instanceId := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

	resp, err := client.MSSQLManagedInstance.ManagedInstanceAdministratorsClient.Get(ctx, instanceId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r MsSqlManagedInstanceActiveDirectoryAdministratorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      /* Due to the creation of unmanaged Microsoft.Network/networkIntentPolicies in this service,
      prevent_deletion_if_contains_resources has been added here to allow the test resources to be
       deleted until this can be properly investigated
      */
      prevent_deletion_if_contains_resources = false
    }
  }
}

%[1]s

resource "azurerm_mssql_managed_instance" "test" {
  name                = "acctestsqlserver%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  license_type       = "BasePrice"
  sku_name           = "GP_Gen5"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.test.id
  vcores             = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_route_table_association.test,
  ]

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "staging"
    database    = "test"
  }
}

data "azuread_client_config" "test" {}

resource "azuread_application" "test" {
  display_name     = "acctest-ManagedInstance-%[2]d"
  sign_in_audience = "AzureADMyOrg"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azuread_directory_role" "reader" {
  display_name = "Directory Readers"
}

resource "azuread_directory_role_member" "test" {
  role_object_id   = azuread_directory_role.reader.object_id
  member_object_id = azurerm_mssql_managed_instance.test.identity.0.principal_id
}
`, MsSqlManagedInstanceResource{}.template(data, data.Locations.Primary), data.RandomInteger)
}

func (r MsSqlManagedInstanceActiveDirectoryAdministratorResource) basic(data acceptance.TestData, aadOnly bool) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_managed_instance_active_directory_administrator" "test" {
  managed_instance_id = azurerm_mssql_managed_instance.test.id
  login_username      = azuread_service_principal.test.display_name
  object_id           = azuread_service_principal.test.object_id
  tenant_id           = data.azuread_client_config.test.tenant_id

  azuread_authentication_only = %[2]t

  depends_on = [azuread_directory_role_member.test]
}
`, r.template(data), aadOnly)
}
