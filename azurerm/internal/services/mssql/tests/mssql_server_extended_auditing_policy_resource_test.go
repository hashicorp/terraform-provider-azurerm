package tests

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMsSqlServerExtendedAuditingPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server_extended_auditing_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerExtendedAuditingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlServerExtendedAuditingPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExtendedAuditingPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep("storage_account_access_key"),
		},
	})
}

func TestAccAzureRMMsSqlServerExtendedAuditingPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server_extended_auditing_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerExtendedAuditingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlServerExtendedAuditingPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExtendedAuditingPolicyExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMMsSqlServerExtendedAuditingPolicy_requiresImport),
		},
	})
}

func TestAccAzureRMMsSqlServerExtendedAuditingPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server_extended_auditing_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerExtendedAuditingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlServerExtendedAuditingPolicy_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExtendedAuditingPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep("storage_account_access_key"),
		},
	})
}

func TestAccAzureRMMsSqlServerExtendedAuditingPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server_extended_auditing_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerExtendedAuditingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlServerExtendedAuditingPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExtendedAuditingPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep("storage_account_access_key"),
			{
				Config: testAccAzureRMMsSqlServerExtendedAuditingPolicy_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExtendedAuditingPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep("storage_account_access_key"),
			{
				Config: testAccAzureRMMsSqlServerExtendedAuditingPolicy_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExtendedAuditingPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep("storage_account_access_key"),
		},
	})
}

func TestAccAzureRMMsSqlServerExtendedAuditingPolicy_storageAccBehindFireWall(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server_extended_auditing_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerExtendedAuditingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlServerExtendedAuditingPolicy_storageAccountBehindFireWall(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExtendedAuditingPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep("storage_account_access_key"),
		},
	})
}

func testCheckAzureRMMsSqlServerExtendedAuditingPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.ServerExtendedBlobAuditingPoliciesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.MssqlServerExtendedAuditingPolicyID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.MsSqlServer)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("msSql Server ExtendedAuditingPolicy %q (resource group: %q) does not exist", id.MsSqlServer, id.ResourceGroup)
			}

			return fmt.Errorf("get on MsSql Server ExtendedAuditingPolicy Client: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMMsSqlServerExtendedAuditingPolicyDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.ServerExtendedBlobAuditingPoliciesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mssql_server_extended_auditing_policy" {
			continue
		}

		id, err := parse.MssqlServerExtendedAuditingPolicyID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.MsSqlServer); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("get on MsSql Server ExtendedAuditingPolicy Client: %+v", err)
			}

			if resp.ExtendedServerBlobAuditingPolicyProperties != nil && resp.ExtendedServerBlobAuditingPolicyProperties.State == sql.BlobAuditingPolicyStateEnabled {
				return fmt.Errorf("`azurerm_mssql_server_extended_auditing_policy` is still enabled")
			}
		}
		return nil
	}

	return nil
}

func testAccAzureRMMsSqlServerExtendedAuditingPolicy_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctest-sqlserver-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "AdminPassword123!"
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMMsSqlServerExtendedAuditingPolicy_basic(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlServerExtendedAuditingPolicy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_server_extended_auditing_policy" "test" {
  server_id                  = azurerm_mssql_server.test.id
  storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, template)
}

func testAccAzureRMMsSqlServerExtendedAuditingPolicy_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlServerExtendedAuditingPolicy_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_server_extended_auditing_policy" "import" {
  server_id                  = azurerm_mssql_server.test.id
  storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}
`, template)
}

func testAccAzureRMMsSqlServerExtendedAuditingPolicy_complete(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlServerExtendedAuditingPolicy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_server_extended_auditing_policy" "test" {
  server_id                               = azurerm_mssql_server.test.id
  storage_endpoint                        = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key              = azurerm_storage_account.test.primary_access_key
  storage_account_access_key_is_secondary = false
  retention_in_days                       = 6
}
`, template)
}

func testAccAzureRMMsSqlServerExtendedAuditingPolicy_update(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlServerExtendedAuditingPolicy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test2" {
  name                     = "unlikely23exst2acc2%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_server_extended_auditing_policy" "test" {
  server_id                               = azurerm_mssql_server.test.id
  storage_endpoint                        = azurerm_storage_account.test2.primary_blob_endpoint
  storage_account_access_key              = azurerm_storage_account.test2.primary_access_key
  storage_account_access_key_is_secondary = true
  retention_in_days                       = 3
}
`, template, data.RandomString)
}

func testAccAzureRMMsSqlServerExtendedAuditingPolicy_storageAccountBehindFireWall(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctest-sqlserver-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "AdminPassword123!"
  identity {
    type = "SystemAssigned"
  }
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
    default_action             = "Deny"
    ip_rules                   = ["127.0.0.1"]
    virtual_network_subnet_ids = [azurerm_subnet.test.id]
  }
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_mssql_server.test.identity.0.principal_id
}

resource "azurerm_mssql_server_extended_auditing_policy" "test" {
  server_id        = azurerm_mssql_server.test.id
  storage_endpoint = azurerm_storage_account.test.primary_blob_endpoint

  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
