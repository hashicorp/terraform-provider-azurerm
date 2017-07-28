package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestResourceAzureRMSqlDatabaseEdition_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "Random",
			ErrCount: 1,
		},
		{
			Value:    "Basic",
			ErrCount: 0,
		},
		{
			Value:    "Standard",
			ErrCount: 0,
		},
		{
			Value:    "Premium",
			ErrCount: 0,
		},
		{
			Value:    "DataWarehouse",
			ErrCount: 0,
		},
	}

	for _, tc := range cases {
		_, errors := validateArmSqlDatabaseEdition(tc.Value, "azurerm_sql_database")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM SQL Database edition to trigger a validation error")
		}
	}
}

func TestAccAzureRMSqlDatabase_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccAzureRMSqlDatabase_basic, ri, ri, ri)

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

func TestAccAzureRMSqlDatabase_elasticPool(t *testing.T) {
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccAzureRMSqlDatabase_elasticPool, ri, ri, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists("azurerm_sql_database.test"),
					resource.TestCheckResourceAttr("azurerm_sql_database.test", "elastic_pool_name", fmt.Sprintf("acctestep%d", ri)),
				),
			},
		},
	})
}

func TestAccAzureRMSqlDatabase_withTags(t *testing.T) {
	ri := acctest.RandInt()
	preConfig := fmt.Sprintf(testAccAzureRMSqlDatabase_withTags, ri, ri, ri)
	postConfig := fmt.Sprintf(testAccAzureRMSqlDatabase_withTagsUpdate, ri, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists("azurerm_sql_database.test"),
					resource.TestCheckResourceAttr(
						"azurerm_sql_database.test", "tags.%", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists("azurerm_sql_database.test"),
					resource.TestCheckResourceAttr(
						"azurerm_sql_database.test", "tags.%", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMSqlDatabase_datawarehouse(t *testing.T) {
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccAzureRMSqlDatabase_datawarehouse, ri, ri, ri)

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

func testCheckAzureRMSqlDatabaseExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for SQL Database: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).sqlDatabasesClient
		resp, err := conn.Get(resourceGroup, serverName, name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get SQL Database: %v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: SQL Database %s (resource group: %s) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMSqlDatabaseDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).sqlDatabasesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_sql_database_record" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		serverName := rs.Primary.Attributes["server_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, serverName, name, "")

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("SQL Database still exists:\n%#v", resp.DatabaseProperties)
		}

	}

	return nil
}

var testAccAzureRMSqlDatabase_elasticPool = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}

resource "azurerm_sql_server" "test" {
    name = "acctestsqlserver%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "West US"
    version = "12.0"
    administrator_login = "mradministrator"
    administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_elasticpool" "test" {
    name = "acctestep%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "West US"
    server_name = "${azurerm_sql_server.test.name}"
    edition = "Basic"
    dtu = 50
    pool_size = 5000
}

resource "azurerm_sql_database" "test" {
    name = "acctestdb%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    server_name = "${azurerm_sql_server.test.name}"
    location = "West US"
    edition = "${azurerm_sql_elasticpool.test.edition}"
    collation = "SQL_Latin1_General_CP1_CI_AS"
    max_size_bytes = "1073741824"
    elastic_pool_name = "${azurerm_sql_elasticpool.test.name}"
    requested_service_objective_name = "ElasticPool"
}
`

var testAccAzureRMSqlDatabase_basic = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_sql_server" "test" {
    name = "acctestsqlserver%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "West US"
    version = "12.0"
    administrator_login = "mradministrator"
    administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_database" "test" {
    name = "acctestdb%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    server_name = "${azurerm_sql_server.test.name}"
    location = "West US"
    edition = "Standard"
    collation = "SQL_Latin1_General_CP1_CI_AS"
    max_size_bytes = "1073741824"
    requested_service_objective_name = "S0"
}
`

var testAccAzureRMSqlDatabase_withTags = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_sql_server" "test" {
    name = "acctestsqlserver%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "West US"
    version = "12.0"
    administrator_login = "mradministrator"
    administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_database" "test" {
    name = "acctestdb%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    server_name = "${azurerm_sql_server.test.name}"
    location = "West US"
    edition = "Standard"
    collation = "SQL_Latin1_General_CP1_CI_AS"
    max_size_bytes = "1073741824"
    requested_service_objective_name = "S0"

    tags {
    	environment = "staging"
    	database = "test"
    }
}
`

var testAccAzureRMSqlDatabase_withTagsUpdate = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_sql_server" "test" {
    name = "acctestsqlserver%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "West US"
    version = "12.0"
    administrator_login = "mradministrator"
    administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_database" "test" {
    name = "acctestdb%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    server_name = "${azurerm_sql_server.test.name}"
    location = "West US"
    edition = "Standard"
    collation = "SQL_Latin1_General_CP1_CI_AS"
    max_size_bytes = "1073741824"
    requested_service_objective_name = "S0"

    tags {
    	environment = "production"
    }
}
`

var testAccAzureRMSqlDatabase_datawarehouse = `
resource "azurerm_resource_group" "test" {
    name = "acctest_rg_%d"
    location = "West US"
}
resource "azurerm_sql_server" "test" {
    name = "acctestsqlserver%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "West US"
    version = "12.0"
    administrator_login = "mradministrator"
    administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_database" "test" {
    name = "acctestdb%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    server_name = "${azurerm_sql_server.test.name}"
    location = "West US"
    edition = "DataWarehouse"
    collation = "SQL_Latin1_General_CP1_CI_AS"
    requested_service_objective_name = "DW400"
}
`
