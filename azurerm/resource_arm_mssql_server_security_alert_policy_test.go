package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMssqlServerSecurityAlertPolicy_basic(t *testing.T) {
	resourceName := "azurerm_mssql_server_security_alert_policy.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMssqlServerSecurityAlertPolicy_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMssqlServerSecurityAlertPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "state", "Enabled"),
					resource.TestCheckResourceAttr(resourceName, "disabled_alerts.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "email_account_admins", "false"),
					resource.TestCheckResourceAttr(resourceName, "retention_days", "20"),
					resource.TestCheckResourceAttr(resourceName, "email_addresses.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"storage_account_access_key"},
			},
		},
	})
}

func TestAccAzureRMMssqlServerSecurityAlertPolicy_update(t *testing.T) {
	resourceName := "azurerm_mssql_server_security_alert_policy.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMssqlServerSecurityAlertPolicy_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMssqlServerSecurityAlertPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "state", "Enabled"),
					resource.TestCheckResourceAttr(resourceName, "disabled_alerts.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "email_account_admins", "false"),
					resource.TestCheckResourceAttr(resourceName, "retention_days", "20"),
					resource.TestCheckResourceAttr(resourceName, "email_addresses.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"storage_account_access_key"},
			},
			{
				Config: testAccAzureRMMssqlServerSecurityAlertPolicy_update(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMssqlServerSecurityAlertPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "state", "Enabled"),
					resource.TestCheckResourceAttr(resourceName, "disabled_alerts.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "email_account_admins", "true"),
					resource.TestCheckResourceAttr(resourceName, "retention_days", "30"),
					resource.TestCheckResourceAttr(resourceName, "email_addresses.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"storage_account_access_key"},
			},
		},
	})
}

func testCheckAzureRMMssqlServerSecurityAlertPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.ServerSecurityAlertPoliciesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testCheckAzureRMStorageAccountSqlServerDestroy(s *terraform.State) error {
	err := testCheckAzureRMStorageAccountDestroy(s)
	if err != nil {
		return err
	}

	return testCheckAzureRMSqlServerDestroy(s)
}

func testAccAzureRMMssqlServerSecurityAlertPolicy_basic(rInt int, location string) string {
	server := testAccAzureRMMssqlServerSecurityAlertPolicy_server(rInt, location)

	return fmt.Sprintf(`
%s

resource "azurerm_mssql_server_security_alert_policy" "test" {
  resource_group_name        = azurerm_resource_group.test.name
  server_name                = azurerm_sql_server.test.name
  state                      = "Enabled"
  storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  disabled_alerts            = [
    "Sql_Injection",
    "Data_Exfiltration"
  ]
  retention_days             = 20
}
`, server)
}

func testAccAzureRMMssqlServerSecurityAlertPolicy_update(rInt int, location string) string {
	server := testAccAzureRMMssqlServerSecurityAlertPolicy_server(rInt, location)

	return fmt.Sprintf(`
%s

resource "azurerm_mssql_server_security_alert_policy" "test" {
  resource_group_name        = azurerm_resource_group.test.name
  server_name                = azurerm_sql_server.test.name
  state                      = "Enabled"
  email_account_admins       = true
  retention_days             = 30
}
`, server)
}

func testAccAzureRMMssqlServerSecurityAlertPolicy_server(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "%s"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}
`, rInt, location, rInt, rInt, location)
}
