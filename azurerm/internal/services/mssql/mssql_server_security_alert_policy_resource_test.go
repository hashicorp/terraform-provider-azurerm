package mssql_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMssqlServerSecurityAlertPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server_security_alert_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMssqlServerSecurityAlertPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMssqlServerSecurityAlertPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMssqlServerSecurityAlertPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "Enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "disabled_alerts.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "email_account_admins", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_days", "20"),
					resource.TestCheckResourceAttr(data.ResourceName, "email_addresses.#", "0"),
				),
			},
			data.ImportStep("storage_account_access_key"),
		},
	})
}

func TestAccAzureRMMssqlServerSecurityAlertPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_server_security_alert_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMssqlServerSecurityAlertPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMssqlServerSecurityAlertPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMssqlServerSecurityAlertPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "Enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "disabled_alerts.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "email_account_admins", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_days", "20"),
					resource.TestCheckResourceAttr(data.ResourceName, "email_addresses.#", "0"),
				),
			},
			data.ImportStep("storage_account_access_key"),
			{
				Config: testAccAzureRMMssqlServerSecurityAlertPolicy_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMssqlServerSecurityAlertPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "Enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "disabled_alerts.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "email_account_admins", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_days", "30"),
					resource.TestCheckResourceAttr(data.ResourceName, "email_addresses.#", "0"),
				),
			},
			data.ImportStep("storage_account_access_key"),
		},
	})
}

func testCheckAzureRMMssqlServerSecurityAlertPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.ServerSecurityAlertPoliciesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]

		resp, err := client.Get(ctx, resourceGroup, serverName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("security alert policy was not found for resource group %q, sql server %q",
					resourceGroup, serverName)
			}

			return err
		}

		return nil
	}
}

func testCheckAzureRMMssqlServerSecurityAlertPolicyDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.ServerSecurityAlertPoliciesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mssql_server_security_alert_policy" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]

		resp, err := client.Get(ctx, resourceGroup, serverName)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Security Alert Policy still exists:\n%#v", resp.SecurityAlertPolicyProperties)
		}
	}

	return nil
}

func testAccAzureRMMssqlServerSecurityAlertPolicy_basic(data acceptance.TestData) string {
	server := testAccAzureRMMssqlServerSecurityAlertPolicy_server(data)

	return fmt.Sprintf(`
%s

resource "azurerm_mssql_server_security_alert_policy" "test" {
  resource_group_name        = azurerm_resource_group.test.name
  server_name                = azurerm_sql_server.test.name
  state                      = "Enabled"
  storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  disabled_alerts = [
    "Sql_Injection",
    "Data_Exfiltration"
  ]
  retention_days = 20
}
`, server)
}

func testAccAzureRMMssqlServerSecurityAlertPolicy_update(data acceptance.TestData) string {
	server := testAccAzureRMMssqlServerSecurityAlertPolicy_server(data)

	return fmt.Sprintf(`
%s

resource "azurerm_mssql_server_security_alert_policy" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  server_name          = azurerm_sql_server.test.name
  state                = "Enabled"
  email_account_admins = true
  retention_days       = 30
}
`, server)
}

func testAccAzureRMMssqlServerSecurityAlertPolicy_server(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
