package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
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
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"administrator_login_password"},
			},
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

func TestAccAzureRMMsSqlServer_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlServer_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
			{
				Config: testAccAzureRMMsSqlServer_withTagsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"administrator_login_password"},
			},
		},
	})
}

func TestAccAzureRMMsSqlServer_withIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlServer_withIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.tenant_id", validate.UUIDRegExp),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"administrator_login_password"},
			},
		},
	})
}

func TestAccAzureRMMsSqlServer_updateWithIdentityAdded(t *testing.T) {
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
			{
				Config: testAccAzureRMMsSqlServer_withIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.tenant_id", validate.UUIDRegExp),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"administrator_login_password"},
			},
		},
	})
}

func TestAccAzureRMMsSqlServer_updateWithBlobAuditingPolices(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlServer_withBlobAuditingPolices(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "extended_auditing_policy.0.storage_account_access_key_is_secondary", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "extended_auditing_policy.0.retention_in_days", "6"),
				),
			},
			data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key"),
			{
				Config: testAccAzureRMMsSqlServer_withBlobAuditingPolicesUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "extended_auditing_policy.0.storage_account_access_key_is_secondary", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "extended_auditing_policy.0.retention_in_days", "11"),
				),
			},
			data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key"),
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
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
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

func testAccAzureRMMsSqlServer_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"

  tags = {
    environment = "staging"
    database    = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMsSqlServer_withTagsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMsSqlServer_withIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMMsSqlServer_withBlobAuditingPolices(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[1]d"
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
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"

  extended_auditing_policy {
    storage_account_access_key              = azurerm_storage_account.test.primary_access_key
    storage_endpoint                        = azurerm_storage_account.test.primary_blob_endpoint
    storage_account_access_key_is_secondary = true
    retention_in_days                       = 6
  }
}
`, data.RandomIntOfLength(15), data.Locations.Primary)
}

func testAccAzureRMMsSqlServer_withBlobAuditingPolicesUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[1]d"
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
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"

  extended_auditing_policy {
    storage_account_access_key              = azurerm_storage_account.test2.primary_access_key
    storage_endpoint                        = azurerm_storage_account.test2.primary_blob_endpoint
    storage_account_access_key_is_secondary = false
    retention_in_days                       = 11
  }
}
`, data.RandomIntOfLength(15), data.Locations.Primary)
}
