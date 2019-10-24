package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSqlExtendedServerBlobAuditingPolicies_basic(t *testing.T) {
	resourceName := "azurerm_sql_extended_server_blob_auditing_policies.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlExtendedServerBlobAuditingPolicies_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlExtendedServerBlobAuditingPoliciesExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "state", "Enabled"),
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

func TestAccAzureRMSqlExtendedServerBlobAuditingPolicies_complete(t *testing.T) {
	resourceName := "azurerm_sql_extended_server_blob_auditing_policies.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlExtendedServerBlobAuditingPolicies_complete(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlExtendedServerBlobAuditingPoliciesExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "state", "Enabled"),
					resource.TestCheckResourceAttr(resourceName, "retention_days", "6"),
					resource.TestCheckResourceAttr(resourceName, "is_storage_secondary_key_in_use", "true"),
					resource.TestCheckResourceAttr(resourceName, "audit_actions_and_groups", "SUCCESSFUL_DATABASE_AUTHENTICATION_GROUP,FAILED_DATABASE_AUTHENTICATION_GROUP"),
					resource.TestCheckResourceAttr(resourceName, "is_azure_monitor_target_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "storage_account_subscription_id", "00000000-0000-0000-3333-000000000000"),
					resource.TestCheckResourceAttr(resourceName, "predicate_expression", "object_name = 'SensitiveData'"),
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

func testCheckAzureRMSqlExtendedServerBlobAuditingPoliciesExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		sqlServerName := rs.Primary.Attributes["server_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for SQL Server: %s Blob Auditing Policies", sqlServerName)
		}

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		conn := testAccProvider.Meta().(*ArmClient).Sql.ExtendedServerBlobAuditingPoliciesClient
		resp, err := conn.Get(ctx, resourceGroup, sqlServerName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: SQL Server %s Blob Auditing Policies(resource group: %s) does not exist", sqlServerName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get SQL Server: %v Blob Auditing Policies", err)
		}
		return nil
	}
}

func testAccAzureRMSqlExtendedServerBlobAuditingPolicies_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
 name     = "acctestRG-%d"
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
 name                     = "accstr%d"
 resource_group_name      = "${azurerm_resource_group.test.name}"
 location                 = "${azurerm_resource_group.test.location}"
 account_tier             = "Standard"
 account_replication_type = "GRS"
}

resource "azurerm_sql_extended_server_blob_auditing_policies" "test"{
resource_group_name           = "${azurerm_resource_group.test.name}"
server_name                   = "${azurerm_sql_server.test.name}"
state                         = "Enabled"
storage_endpoint              = "${azurerm_storage_account.test.primary_blob_endpoint}"
storage_account_access_key    = "${azurerm_storage_account.test.primary_access_key}"

}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMSqlExtendedServerBlobAuditingPolicies_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
 name     = "acctestRG-%d"
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
 name                     = "accstr%d"
 resource_group_name      = "${azurerm_resource_group.test.name}"
 location                 = "${azurerm_resource_group.test.location}"
 account_tier             = "Standard"
 account_replication_type = "GRS"
}

resource "azurerm_sql_extended_server_blob_auditing_policies" "test"{
resource_group_name               = "${azurerm_resource_group.test.name}"
server_name                       = "${azurerm_sql_server.test.name}"
state                             = "Enabled"
storage_endpoint                  = "${azurerm_storage_account.test.primary_blob_endpoint}"
storage_account_access_key        = "${azurerm_storage_account.test.primary_access_key}"
retention_days                    = 6
is_storage_secondary_key_in_use   = true
audit_actions_and_groups          = "SUCCESSFUL_DATABASE_AUTHENTICATION_GROUP,FAILED_DATABASE_AUTHENTICATION_GROUP"
is_azure_monitor_target_enabled   = true
storage_account_subscription_id   = "00000000-0000-0000-3333-000000000000"
predicate_expression              ="object_name = 'SensitiveData'"

}
`, rInt, location, rInt, rInt)
}
