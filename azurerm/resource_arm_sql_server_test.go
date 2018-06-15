package azurerm

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func init() {
	resource.AddTestSweepers("azurerm_sql_server", &resource.Sweeper{
		Name: "azurerm_sql_server",
		F:    testSweepSQLServer,
	})
}

func testSweepSQLServer(region string) error {
	armClient, err := buildConfigForSweepers()
	if err != nil {
		return err
	}

	client := (*armClient).sqlServersClient
	ctx := (*armClient).StopContext

	log.Printf("Retrieving the SQL Servers..")
	results, err := client.List(ctx)
	if err != nil {
		return fmt.Errorf("Error Listing on SQL Servers: %+v", err)
	}

	for _, server := range results.Values() {
		if !shouldSweepAcceptanceTestResource(*server.Name, *server.Location, region) {
			continue
		}

		resourceId, err := parseAzureResourceID(*server.ID)
		if err != nil {
			return err
		}

		resourceGroup := resourceId.ResourceGroup
		name := resourceId.Path["servers"]

		log.Printf("Deleting SQL Server '%s' in Resource Group '%s'", name, resourceGroup)
		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			return err
		}

		err = future.WaitForCompletion(ctx, client.Client)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestAccAzureRMSqlServer_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMSqlServer_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlServerExists("azurerm_sql_server.test"),
				),
			},
		},
	})
}

func TestAccAzureRMSqlServer_disappears(t *testing.T) {
	resourceName := "azurerm_sql_server.test"
	ri := acctest.RandInt()
	config := testAccAzureRMSqlServer_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlServerExists(resourceName),
					testCheckAzureRMSqlServerDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMSqlServer_withTags(t *testing.T) {
	resourceName := "azurerm_sql_server.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMSqlServer_withTags(ri, location)
	postConfig := testAccAzureRMSqlServer_withTagsUpdated(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func testCheckAzureRMSqlServerExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		sqlServerName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for SQL Server: %s", sqlServerName)
		}

		conn := testAccProvider.Meta().(*ArmClient).sqlServersClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

func testCheckAzureRMSqlServerDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).sqlServersClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_sql_server" {
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

func testCheckAzureRMSqlServerDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sqlServersClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		future, err := client.Delete(ctx, resourceGroup, serverName)
		if err != nil {
			return err
		}

		err = future.WaitForCompletion(ctx, client.Client)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccAzureRMSqlServer_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_sql_server" "test" {
    name = "acctestsqlserver%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    version = "12.0"
    administrator_login = "mradministrator"
    administrator_login_password = "thisIsDog11"
}
`, rInt, location, rInt)
}

func testAccAzureRMSqlServer_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_sql_server" "test" {
    name = "acctestsqlserver%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    version = "12.0"
    administrator_login = "mradministrator"
    administrator_login_password = "thisIsDog11"

    tags {
    	environment = "staging"
    	database = "test"
    }
}
`, rInt, location, rInt)
}

func testAccAzureRMSqlServer_withTagsUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_sql_server" "test" {
    name = "acctestsqlserver%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    version = "12.0"
    administrator_login = "mradministrator"
    administrator_login_password = "thisIsDog11"

    tags {
    	environment = "production"
    }
}
`, rInt, location, rInt)
}
