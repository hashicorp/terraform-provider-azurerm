// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlServerResource struct{}

func TestAccMsSqlServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlServer_minimumTLSVersionDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithMinimumTLSVersionDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

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

func TestAccMsSqlServer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.basicWithMinimumTLSVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlServer_systemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlServer_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlServer_systemAndUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAndUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlServer_azureadAdmin(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.aadAdmin(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlServer_azureadAdminUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.aadAdmin(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.aadAdminWithAADAuthOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlServer_azureadAdminWithAADAuthOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.aadAdminWithAADAuthOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlServer_azureadAuthenticationOnlyWithIdentityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.updateAzureadAuthenticationOnlyWithIdentity(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.updateAzureadAuthenticationOnlyWithIdentity(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.updateAzureadAuthenticationOnlyWithIdentity(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlServer_TDECMKServerDeployment(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tdeCMKServerDeployment(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlServer_CMKServerTagsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.CMKServerTags(data, "Sandbox"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.CMKServerTags(data, "Production"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.CMKServerNoTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func (MsSqlServerResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ServerID(state.ID)
	if err != nil {
		return nil, err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.Name)

	resp, err := client.MSSQL.ServersClient.Get(ctx, serverId, servers.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("SQL Server %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
		}
		return nil, fmt.Errorf("reading SQL Server %q (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (MsSqlServerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"

  outbound_network_restriction_enabled = true
}
`, data.RandomInteger, data.Locations.Primary)
}

func (MsSqlServerResource) basicWithMinimumTLSVersionDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"
  minimum_tls_version          = "Disabled"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (MsSqlServerResource) basicWithMinimumTLSVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"
  minimum_tls_version          = "1.1"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r MsSqlServerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_server" "import" {
  name                         = azurerm_mssql_server.test.name
  resource_group_name          = azurerm_mssql_server.test.resource_group_name
  location                     = azurerm_mssql_server.test.location
  version                      = azurerm_mssql_server.test.version
  administrator_login          = azurerm_mssql_server.test.administrator_login
  administrator_login_password = azurerm_mssql_server.test.administrator_login_password
}
`, r.basic(data))
}

func (MsSqlServerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "service" {
  name                 = "acctestsnetservice-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.1.0/24"]

  private_link_service_network_policies_enabled = true
}

resource "azurerm_subnet" "endpoint" {
  name                 = "acctestsnetendpoint-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  private_endpoint_network_policies = "Disabled"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctesta%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"
  minimum_tls_version          = "1.2"

  public_network_access_enabled     = true
  primary_user_assigned_identity_id = azurerm_user_assigned_identity.test.id

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    ENV      = "Staging"
    database = "NotProd"
  }
}

resource "azurerm_private_dns_zone" "finance" {
  name                = "privatelink.sql.database.azure.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = "acctest-privatelink-mssc-%[1]d"
    private_connection_resource_id = azurerm_mssql_server.test.id
    subresource_names              = ["sqlServer"]
    is_manual_connection           = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(15))
}

func (MsSqlServerResource) completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "service" {
  name                 = "acctestsnetservice-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.1.0/24"]

  private_link_service_network_policies_enabled = true
}

resource "azurerm_subnet" "endpoint" {
  name                 = "acctestsnetendpoint-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  private_endpoint_network_policies = "Disabled"
}

resource "azurerm_storage_account" "testb" {
  name                     = "acctestb%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"
  minimum_tls_version          = "1.0"

  public_network_access_enabled     = false
  primary_user_assigned_identity_id = azurerm_user_assigned_identity.test.id

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    DB = "NotProd"
  }
}

resource "azurerm_private_dns_zone" "finance" {
  name                = "privatelink.sql.database.azure.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = "acctest-privatelink-mssc-%[1]d"
    private_connection_resource_id = azurerm_mssql_server.test.id
    subresource_names              = ["sqlServer"]
    is_manual_connection           = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(15))
}

func (MsSqlServerResource) systemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (MsSqlServerResource) userAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test1" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity_1"
}

resource "azurerm_user_assigned_identity" "test2" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity_2"
}

resource "azurerm_mssql_server" "test" {
  name                              = "acctestsqlserver%[1]d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = azurerm_resource_group.test.location
  version                           = "12.0"
  administrator_login               = "missadministrator"
  administrator_login_password      = "thisIsKat11"
  primary_user_assigned_identity_id = azurerm_user_assigned_identity.test1.id

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test1.id, azurerm_user_assigned_identity.test2.id]
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (MsSqlServerResource) systemAndUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (MsSqlServerResource) aadAdmin(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azurerm_client_config" "test" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"

  azuread_administrator {
    login_username = "AzureAD Admin"
    object_id      = data.azurerm_client_config.test.object_id
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (MsSqlServerResource) aadAdminWithAADAuthOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azurerm_client_config" "test" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                = "acctestsqlserver%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  version             = "12.0"

  azuread_administrator {
    login_username              = "AzureAD Admin2"
    object_id                   = data.azurerm_client_config.test.object_id
    azuread_authentication_only = true
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (MsSqlServerResource) updateAzureadAuthenticationOnlyWithIdentity(data acceptance.TestData, enableAzureadAuthenticationOnly bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azurerm_client_config" "test" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity_1"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  minimum_tls_version          = "1.2"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"

  azuread_administrator {
    login_username              = "AzureAD Admin"
    object_id                   = data.azurerm_client_config.test.object_id
    azuread_authentication_only = %[3]t
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  primary_user_assigned_identity_id = azurerm_user_assigned_identity.test.id
}
`, data.RandomInteger, data.Locations.Primary, enableAzureadAuthenticationOnly)
}

func (MsSqlServerResource) tdeCMKServerDeployment(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "test" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity_2112"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "DaveLister"
  administrator_login_password = "7h1515K4711"
  minimum_tls_version          = "1.2"

  azuread_administrator {
    login_username = azurerm_user_assigned_identity.test.name
    object_id      = azurerm_user_assigned_identity.test.principal_id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  primary_user_assigned_identity_id            = azurerm_user_assigned_identity.test.id
  transparent_data_encryption_key_vault_key_id = azurerm_key_vault_key.test.id
}

resource "azurerm_key_vault" "test" {
  name                        = "vault%[1]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  enabled_for_disk_encryption = true
  tenant_id                   = azurerm_user_assigned_identity.test.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = true

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    object_id = data.azurerm_client_config.test.object_id

    key_permissions = ["Get", "List", "Create", "Delete", "Update", "Recover", "Purge", "GetRotationPolicy"]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.test.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    key_permissions = ["Get", "WrapKey", "UnwrapKey"]
  }
}

resource "azurerm_key_vault_key" "test" {
  depends_on = [azurerm_key_vault.test]

  name         = "key-%[3]s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = ["unwrapKey", "wrapKey"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (MsSqlServerResource) CMKServerTags(data acceptance.TestData, tag string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "test" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity_2112"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "DaveLister"
  administrator_login_password = "7h1515K4711"
  minimum_tls_version          = "1.2"

  azuread_administrator {
    login_username = azurerm_user_assigned_identity.test.name
    object_id      = azurerm_user_assigned_identity.test.principal_id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  primary_user_assigned_identity_id            = azurerm_user_assigned_identity.test.id
  transparent_data_encryption_key_vault_key_id = azurerm_key_vault_key.test.id

  tags = {
    DB = "%[4]s"
  }
}

resource "azurerm_key_vault" "test" {
  name                        = "vault%[1]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  enabled_for_disk_encryption = true
  tenant_id                   = azurerm_user_assigned_identity.test.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = true

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    object_id = data.azurerm_client_config.test.object_id

    key_permissions = ["Get", "List", "Create", "Delete", "Update", "Recover", "Purge", "GetRotationPolicy"]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.test.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    key_permissions = ["Get", "WrapKey", "UnwrapKey"]
  }
}

resource "azurerm_key_vault_key" "test" {
  depends_on = [azurerm_key_vault.test]

  name         = "key-%[3]s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = ["unwrapKey", "wrapKey"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, tag)
}

func (MsSqlServerResource) CMKServerNoTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "test" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity_2112"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "DaveLister"
  administrator_login_password = "7h1515K4711"
  minimum_tls_version          = "1.2"

  azuread_administrator {
    login_username = azurerm_user_assigned_identity.test.name
    object_id      = azurerm_user_assigned_identity.test.principal_id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  primary_user_assigned_identity_id            = azurerm_user_assigned_identity.test.id
  transparent_data_encryption_key_vault_key_id = azurerm_key_vault_key.test.id
}

resource "azurerm_key_vault" "test" {
  name                        = "vault%[1]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  enabled_for_disk_encryption = true
  tenant_id                   = azurerm_user_assigned_identity.test.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = true

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    object_id = data.azurerm_client_config.test.object_id

    key_permissions = ["Get", "List", "Create", "Delete", "Update", "Recover", "Purge", "GetRotationPolicy"]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.test.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    key_permissions = ["Get", "WrapKey", "UnwrapKey"]
  }
}

resource "azurerm_key_vault_key" "test" {
  depends_on = [azurerm_key_vault.test]

  name         = "key-%[3]s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = ["unwrapKey", "wrapKey"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
