package azurerm

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSqlDatabase_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMSqlDatabase_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists("azurerm_sql_database.test"),
				),
			},
		},
	})
}

func TestAccAzureRMSqlDatabase_disappears(t *testing.T) {
	resourceName := "azurerm_sql_database.test"
	ri := acctest.RandInt()
	config := testAccAzureRMSqlDatabase_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists(resourceName),
					testCheckAzureRMSqlDatabaseDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMSqlDatabase_elasticPool(t *testing.T) {
	resourceName := "azurerm_sql_database.test"
	ri := acctest.RandInt()
	config := testAccAzureRMSqlDatabase_elasticPool(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "elastic_pool_name", fmt.Sprintf("acctestep%d", ri)),
				),
			},
		},
	})
}

func TestAccAzureRMSqlDatabase_withTags(t *testing.T) {
	resourceName := "azurerm_sql_database.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMSqlDatabase_withTags(ri, location)
	postConfig := testAccAzureRMSqlDatabase_withTagsUpdate(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMSqlDatabase_dataWarehouse(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMSqlDatabase_dataWarehouse(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists("azurerm_sql_database.test"),
				),
			},
		},
	})
}

func TestAccAzureRMSqlDatabase_restorePointInTime(t *testing.T) {
	resourceName := "azurerm_sql_database.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMSqlDatabase_basic(ri, location)
	timeToRestore := time.Now().Add(15 * time.Minute)
	formattedTime := string(timeToRestore.UTC().Format(time.RFC3339))
	postCongif := testAccAzureRMSqlDatabase_restorePointInTime(ri, formattedTime, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				PreventPostDestroyRefresh: true,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists(resourceName),
				),
			},
			{
				PreConfig: func() { time.Sleep(timeToRestore.Sub(time.Now().Add(-1 * time.Minute))) },
				Config:    postCongif,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists(resourceName),
					testCheckAzureRMSqlDatabaseExists("azurerm_sql_database.test_restore"),
				),
			},
		},
	})
}

func TestAccAzureRMSqlDatabase_collation(t *testing.T) {
	resourceName := "azurerm_sql_database.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMSqlDatabase_basic(ri, location)
	postConfig := testAccAzureRMSqlDatabase_collationUpdate(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "collation", "SQL_Latin1_General_CP1_CI_AS"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "collation", "Japanese_Bushu_Kakusu_100_CS_AS_KS_WS"),
				),
			},
		},
	})
}

func testCheckAzureRMSqlDatabaseExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		databaseName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sqlDatabasesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, databaseName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("SQL Database %q (server %q / resource group %q) was not found", databaseName, serverName, resourceGroup)
			}

			return err
		}

		return nil
	}
}

func testCheckAzureRMSqlDatabaseDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_sql_database" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		databaseName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sqlDatabasesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, databaseName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("SQL Database %q (server %q / resource group %q) still exists: %+v", databaseName, serverName, resourceGroup, resp)
	}

	return nil
}

func testCheckAzureRMSqlDatabaseDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		databaseName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sqlDatabasesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		_, err := client.Delete(ctx, resourceGroup, serverName, databaseName)
		if err != nil {
			return fmt.Errorf("Bad: Delete on sqlDatabasesClient: %+v", err)
		}

		return nil
	}
}

func TestAccAzureRMSqlDatabase_bacpac(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMSqlDatabase_bacpac(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists("azurerm_sql_database.test"),
				),
			},
		},
	})
}

func testAccAzureRMSqlDatabase_basic(rInt int, location string) string {
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

resource "azurerm_sql_database" "test" {
    name = "acctestdb%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    server_name = "${azurerm_sql_server.test.name}"
    location = "${azurerm_resource_group.test.location}"
    edition = "Standard"
    collation = "SQL_Latin1_General_CP1_CI_AS"
    max_size_bytes = "1073741824"
    requested_service_objective_name = "S0"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMSqlDatabase_withTags(rInt int, location string) string {
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

resource "azurerm_sql_database" "test" {
    name = "acctestdb%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    server_name = "${azurerm_sql_server.test.name}"
    location = "${azurerm_resource_group.test.location}"
    edition = "Standard"
    collation = "SQL_Latin1_General_CP1_CI_AS"
    max_size_bytes = "1073741824"
    requested_service_objective_name = "S0"

    tags {
    	environment = "staging"
    	database = "test"
    }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMSqlDatabase_withTagsUpdate(rInt int, location string) string {
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

resource "azurerm_sql_database" "test" {
    name = "acctestdb%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    server_name = "${azurerm_sql_server.test.name}"
    location = "${azurerm_resource_group.test.location}"
    edition = "Standard"
    collation = "SQL_Latin1_General_CP1_CI_AS"
    max_size_bytes = "1073741824"
    requested_service_objective_name = "S0"

    tags {
    	environment = "production"
    }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMSqlDatabase_dataWarehouse(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctest_rg_%d"
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

resource "azurerm_sql_database" "test" {
    name = "acctestdb%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    server_name = "${azurerm_sql_server.test.name}"
    location = "${azurerm_resource_group.test.location}"
    edition = "DataWarehouse"
    collation = "SQL_Latin1_General_CP1_CI_AS"
    requested_service_objective_name = "DW400"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMSqlDatabase_restorePointInTime(rInt int, formattedTime string, location string) string {
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

resource "azurerm_sql_database" "test" {
    name = "acctestdb%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    server_name = "${azurerm_sql_server.test.name}"
    location = "${azurerm_resource_group.test.location}"
    edition = "Standard"
    collation = "SQL_Latin1_General_CP1_CI_AS"
    max_size_bytes = "1073741824"
    requested_service_objective_name = "S0"
}

resource "azurerm_sql_database" "test_restore" {
    name = "acctestdb_restore%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    server_name = "${azurerm_sql_server.test.name}"
    location = "${azurerm_resource_group.test.location}"
    create_mode = "PointInTimeRestore"
    source_database_id = "${azurerm_sql_database.test.id}"
    restore_point_in_time = "%s"
}
`, rInt, location, rInt, rInt, rInt, formattedTime)
}

func testAccAzureRMSqlDatabase_elasticPool(rInt int, location string) string {
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

resource "azurerm_sql_elasticpool" "test" {
    name = "acctestep%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    server_name = "${azurerm_sql_server.test.name}"
    edition = "Basic"
    dtu = 50
    pool_size = 5000
}

resource "azurerm_sql_database" "test" {
    name = "acctestdb%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    server_name = "${azurerm_sql_server.test.name}"
    location = "${azurerm_resource_group.test.location}"
    edition = "${azurerm_sql_elasticpool.test.edition}"
    collation = "SQL_Latin1_General_CP1_CI_AS"
    max_size_bytes = "1073741824"
    elastic_pool_name = "${azurerm_sql_elasticpool.test.name}"
    requested_service_objective_name = "ElasticPool"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMSqlDatabase_collationUpdate(rInt int, location string) string {
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

resource "azurerm_sql_database" "test" {
    name = "acctestdb%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    server_name = "${azurerm_sql_server.test.name}"
    location = "${azurerm_resource_group.test.location}"
    edition = "Standard"
    collation = "Japanese_Bushu_Kakusu_100_CS_AS_KS_WS"
    max_size_bytes = "1073741824"
    requested_service_objective_name = "S0"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMSqlDatabase_bacpac(rInt int, location string) string {
	return fmt.Sprintf(`
		resource "azurerm_resource_group" "test" {
			name     = "acctestRG_%d"
			location = "%s"
		}
		
		resource "azurerm_storage_account" "test" {
			name                     = "accsa%d"
			resource_group_name      = "${azurerm_resource_group.test.name}"
			location                 = "${azurerm_resource_group.test.location}"
			account_tier             = "Standard"
			account_replication_type = "LRS"
		}
		
		resource "azurerm_storage_container" "test" {
			name                  = "bacpac"
			resource_group_name   = "${azurerm_resource_group.test.name}"
			storage_account_name  = "${azurerm_storage_account.test.name}"
			container_access_type = "private"
		}
		
		resource "azurerm_storage_blob" "test" {
			name                   = "test.bacpac"
			resource_group_name    = "${azurerm_resource_group.test.name}"
			storage_account_name   = "${azurerm_storage_account.test.name}"
			storage_container_name = "${azurerm_storage_container.test.name}"
			type                   = "block"
			source                 = "testdata/sql_import.bacpac"
		}
		
		resource "azurerm_sql_server" "test" {
			name                         = "acctestsqlserver%d"
			resource_group_name          = "${azurerm_resource_group.test.name}"
			location                     = "${azurerm_resource_group.test.location}"
			version                      = "12.0"
			administrator_login          = "mradministrator"
			administrator_login_password = "thisIsDog11"
		}
		
		resource "azurerm_sql_firewall_rule" "test" {
			name                = "allowazure"
			resource_group_name = "${azurerm_resource_group.test.name}"
			server_name         = "${azurerm_sql_server.test.name}"
			start_ip_address    = "0.0.0.0"
			end_ip_address      = "0.0.0.0"
		}
		
		resource "azurerm_sql_database" "test" {
			name                             = "acctestdb%d"
			resource_group_name              = "${azurerm_resource_group.test.name}"
			server_name                      = "${azurerm_sql_server.test.name}"
			location                         = "${azurerm_resource_group.test.location}"
			edition                          = "Standard"
			collation                        = "SQL_Latin1_General_CP1_CI_AS"
			max_size_bytes                   = "1073741824"
			requested_service_objective_name = "S0"
		
			import {
				storage_uri                  = "${azurerm_storage_blob.test.url}"
				storage_key                  = "${azurerm_storage_account.test.primary_access_key}"
				storage_key_type             = "StorageAccessKey"
				administrator_login          = "${azurerm_sql_server.test.administrator_login}"
				administrator_login_password = "${azurerm_sql_server.test.administrator_login_password}"
				authentication_type          = "SQL"
			}
		}
		`, rInt, location, rInt, rInt, rInt)
}
