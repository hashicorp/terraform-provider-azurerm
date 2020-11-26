package tests

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMsSqlServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMMsSqlServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlServer_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key"),
		},
	})
}

func TestAccAzureRMMsSqlServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMMsSqlServer_requiresImport),
		},
	})
}

func TestAccAzureRMMsSqlServer_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
					testCheckAzureRMMsSqlServerDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMMsSqlServer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				Config: testAccAzureRMMsSqlServer_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key"),
			{
				Config: testAccAzureRMMsSqlServer_complete2(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key"),
			{
				Config: testAccAzureRMMsSqlServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMMsSqlServer_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlServer_identity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key"),
		},
	})
}

func TestAccAzureRMMsSqlServer_azureadAdmin(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				Config: testAccAzureRMMsSqlServer_azureadAdmin(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				Config: testAccAzureRMMsSqlServer_azureadAdminUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				Config: testAccAzureRMMsSqlServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMMsSqlServer_blobAuditingPolicies_withFirewall(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key"),
			{
				Config: testAccAzureRMMsSqlServer_blobAuditingPolicies_withFirewall(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key"),
			{
				Config: testAccAzureRMMsSqlServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key"),
		},
	})
}

func TestAccAzureRMMsSqlServer_customDiff(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlServer_basicWithMinimumTLSVersion(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				Config: testAccAzureRMMsSqlServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
				ExpectError: regexp.MustCompile("`minimum_tls_version` cannot be removed once set, please set a valid value for this property"),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func testCheckAzureRMMsSqlServerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Sql.ServersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		sqlServerName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for SQL Server: %s", sqlServerName)
		}

		resp, err := conn.Get(ctx, resourceGroup, sqlServerName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: SQL Server %s (resource group: %s) does not exist", sqlServerName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get SQL Server: %v", err)
		}

		return nil
	}
}

func testCheckAzureRMMsSqlServerDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.ServersClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mssql_server" {
			continue
		}

		sqlServerName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, sqlServerName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("Bad: Get Server: %+v", err)
		}

		return fmt.Errorf("SQL Server %s still exists", sqlServerName)
	}

	return nil
}

func testCheckAzureRMMsSqlServerDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.ServersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["name"]

		future, err := client.Delete(ctx, resourceGroup, serverName)
		if err != nil {
			return err
		}

		return future.WaitForCompletionRef(ctx, client.Client)
	}
}

func testAccAzureRMMsSqlServer_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%d"
  location = "%s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMsSqlServer_basicWithMinimumTLSVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%d"
  location = "%s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"
  minimum_tls_version          = "1.2"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMsSqlServer_requiresImport(data acceptance.TestData) string {
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
`, testAccAzureRMMsSqlServer_basic(data))
}

func testAccAzureRMMsSqlServer_complete(data acceptance.TestData) string {
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

func testAccAzureRMMsSqlServer_complete2(data acceptance.TestData) string {
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

func testAccAzureRMMsSqlServer_identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%d"
  location = "%s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMsSqlServer_azureadAdmin(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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

func testAccAzureRMMsSqlServer_azureadAdminUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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

func testAccAzureRMMsSqlServer_blobAuditingPolicies_withFirewall(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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
