package mssql_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MsSqlServerResource struct{}

func TestAccMsSqlServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key"),
	})
}

func TestAccMsSqlServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

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

func TestAccMsSqlServer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key"),
		{
			Config: r.completeUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key"),
		{
			Config: r.basicWithMinimumTLSVersion(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlServer_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.identity(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key"),
	})
}

func TestAccMsSqlServer_azureadAdmin(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.aadAdmin(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.aadAdminUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMsSqlServer_blobAuditingPolicies_withFirewall(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key"),
		{
			Config: r.blobAuditingPoliciesWithFirewall(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key"),
	})
}

func TestAccMsSqlServer_customDiff(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")
	r := MsSqlServerResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicWithMinimumTLSVersion(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExpectError: regexp.MustCompile("`minimum_tls_version` cannot be removed once set, please set a valid value for this property"),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func (MsSqlServerResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ServerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.ServersClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil, fmt.Errorf("SQL Server %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
		}
		return nil, fmt.Errorf("reading SQL Server %q (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ID != nil), nil
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
  extended_auditing_policy     = []
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
  minimum_tls_version          = "1.2"
  extended_auditing_policy     = []
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

  enforce_private_link_service_network_policies = true
}

resource "azurerm_subnet" "endpoint" {
  name                 = "acctestsnetendpoint-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_storage_account" "test" {
  name                     = "acctesta%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"
  minimum_tls_version          = "1.2"

  public_network_access_enabled = true

  extended_auditing_policy {
    storage_account_access_key              = azurerm_storage_account.test.primary_access_key
    storage_endpoint                        = azurerm_storage_account.test.primary_blob_endpoint
    storage_account_access_key_is_secondary = true
    retention_in_days                       = 6
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

  enforce_private_link_service_network_policies = true
}

resource "azurerm_subnet" "endpoint" {
  name                 = "acctestsnetendpoint-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_storage_account" "testb" {
  name                     = "acctestb%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"
  minimum_tls_version          = "1.0"

  public_network_access_enabled = false

  extended_auditing_policy {
    storage_account_access_key              = azurerm_storage_account.testb.primary_access_key
    storage_endpoint                        = azurerm_storage_account.testb.primary_blob_endpoint
    storage_account_access_key_is_secondary = false
    retention_in_days                       = 11
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

func (MsSqlServerResource) identity(data acceptance.TestData) string {
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

func (MsSqlServerResource) aadAdmin(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

data "azuread_service_principal" "test" {
  application_id = "%[3]s"
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
    object_id      = data.azuread_service_principal.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, os.Getenv("ARM_CLIENT_ID"))
}

func (MsSqlServerResource) aadAdminUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

data "azuread_service_principal" "test" {
  application_id = "%[3]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"

  azuread_administrator {
    login_username = "AzureAD Admin2"
    object_id      = data.azuread_service_principal.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, os.Getenv("ARM_CLIENT_ID"))
}

func (MsSqlServerResource) blobAuditingPoliciesWithFirewall(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  network_rules {
    default_action             = "Allow"
    ip_rules                   = ["127.0.0.1"]
    virtual_network_subnet_ids = [azurerm_subnet.test.id]
  }
}

data "azuread_service_principal" "test" {
  application_id = "%[4]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"

  azuread_administrator {
    login_username = "AzureAD Admin2"
    object_id      = data.azuread_service_principal.test.id
  }

  extended_auditing_policy {
    storage_account_access_key              = azurerm_storage_account.test.primary_access_key
    storage_endpoint                        = azurerm_storage_account.test.primary_blob_endpoint
    storage_account_access_key_is_secondary = true
    retention_in_days                       = 6
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, os.Getenv("ARM_CLIENT_ID"))
}
