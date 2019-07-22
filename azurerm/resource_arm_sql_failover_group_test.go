package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSqlFailoverGroup_basic(t *testing.T) {
	resourceName := "azurerm_sql_failover_group.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlFailoverGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlFailoverGroup_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlFailoverGroupExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSqlFailoverGroup_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skiiping since resources aren't required to be imported")
		return
	}
	resourceName := "azurerm_sql_failover_group.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlFailoverGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlFailoverGroup_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlFailoverGroupExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMSqlFailoverGroup_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_sql_failover_group"),
			},
		},
	})
}

func TestAccAzureRMSqlFailoverGroup_disappears(t *testing.T) {
	resourceName := "azurerm_sql_failover_group.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlFailoverGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlFailoverGroup_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlFailoverGroupExists(resourceName),
					testCheckAzureRMSqlFailoverGroupDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMSqlFailoverGroup_withTags(t *testing.T) {
	resourceName := "azurerm_sql_failover_group.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMSqlFailoverGroup_withTags(ri, location)
	postConfig := testAccAzureRMSqlFailoverGroup_withTagsUpdate(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlFailoverGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlFailoverGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlFailoverGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func testCheckAzureRMSqlFailoverGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		name := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sqlFailoverGroupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("SQL Failover Group %q (server %q / resource group %q) was not found", name, serverName, resourceGroup)
			}

			return err
		}

		return nil
	}
}

func testCheckAzureRMSqlFailoverGroupDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		name := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sqlFailoverGroupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if _, err := client.Delete(ctx, resourceGroup, serverName, name); err != nil {
			return fmt.Errorf("Bad: Delete on sqlFailoverGroupsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSqlFailoverGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azure_sql_failover_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]

		client := testAccProvider.Meta().(*ArmClient).sqlFailoverGroupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("SQL Failover Group %q (server %q / resource group %q) still exists: %+v", name, serverName, resourceGroup, resp)
	}

	return nil
}

func testAccAzureRMSqlFailoverGroup_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
	
resource "azurerm_sql_server" "test_primary" {
  name                         = "acctestsqlserver%d-primary"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_server" "test_secondary" {
  name                         = "acctestsqlserver%d-secondary"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}
  
resource "azurerm_sql_database" "test" {
  name                             = "acctestdb%d"
  resource_group_name              = "${azurerm_resource_group.test.name}"
  server_name                      = "${azurerm_sql_server.test_primary.name}"
  location                         = "${azurerm_resource_group.test.location}"
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"
}

resource "azurerm_sql_failover_group" "test" {
  name                = "acctest%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_sql_server.test_primary.name}"
  databases           = ["${azurerm_sql_database.test.id}"]
  
  partner_servers {
    id = "${azurerm_sql_server.test_secondary.id}"
  }

  read_write_endpoint_failover_policy = {
    mode          = "Automatic"
    grace_minutes = 60
  }
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMSqlFailoverGroup_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sql_failover_group" "import" {
	name                                = "${azurerm_sql_failover_group.test.name}"
	resource_group_name                 = "${azurerm_sql_failover_group.test.resource_group_name}"
	server_name                         = "${azurerm_sql_failover_group.test.server_name}"
	databases                           = "${azurerm_sql_failover_group.test.databases}"
	partner_servers                     = "${azurerm_sql_failover_group.test.partner_servers}"
	read_write_endpoint_failover_policy = "${azurerm_sql_failover_group.test.read_write_endpoint_failover_policy}"
  }
`, testAccAzureRMSqlFailoverGroup_basic(rInt, location))
}

func testAccAzureRMSqlFailoverGroup_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
	
resource "azurerm_sql_server" "test_primary" {
  name                         = "acctestsqlserver%d-primary"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_server" "test_secondary" {
  name                         = "acctestsqlserver%d-secondary"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}
  
resource "azurerm_sql_database" "test" {
  name                             = "acctestdb%d"
  resource_group_name              = "${azurerm_resource_group.test.name}"
  server_name                      = "${azurerm_sql_server.test_primary.name}"
  location                         = "${azurerm_resource_group.test.location}"
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"
}

resource "azurerm_sql_failover_group" "test" {
  name                = "acctest%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_sql_server.test_primary.name}"
  databases           = ["${azurerm_sql_database.test.id}"]
  
  partner_servers {
    id = "${azurerm_sql_server.test_secondary.id}"
  }

  read_write_endpoint_failover_policy = {
    mode          = "Automatic"
    grace_minutes = 60
  }

  tags = {
    environment = "staging"
    database    = "test"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMSqlFailoverGroup_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
	
resource "azurerm_sql_server" "test_primary" {
  name                         = "acctestsqlserver%d-primary"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_server" "test_secondary" {
  name                         = "acctestsqlserver%d-secondary"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}
  
resource "azurerm_sql_database" "test" {
  name                             = "acctestdb%d"
  resource_group_name              = "${azurerm_resource_group.test.name}"
  server_name                      = "${azurerm_sql_server.test_primary.name}"
  location                         = "${azurerm_resource_group.test.location}"
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"
}

resource "azurerm_sql_failover_group" "test" {
  name                = "acctest%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_sql_server.test_primary.name}"
  databases           = ["${azurerm_sql_database.test.id}"]
  
  partner_servers {
    id = "${azurerm_sql_server.test_secondary.id}"
  }

  read_write_endpoint_failover_policy = {
    mode          = "Automatic"
    grace_minutes = 60
  }

  tags = {
    environment = "production"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt)
}
