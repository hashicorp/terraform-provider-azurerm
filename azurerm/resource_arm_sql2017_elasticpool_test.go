package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMSqlElasticPool_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMSqlElasticPool_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlElasticPoolExists("azurerm_sql_elasticpool.test"),
				),
			},
		},
	})
}

func TestAccAzureRMSqlElasticPool_disappears(t *testing.T) {
	resourceName := "azurerm_sql_elasticpool.test"
	ri := acctest.RandInt()
	config := testAccAzureRMSqlElasticPool_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlElasticPoolExists(resourceName),
					testCheckAzureRMSqlElasticPoolDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMSqlElasticPool_resizeDtu(t *testing.T) {
	resourceName := "azurerm_sql_elasticpool.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMSqlElasticPool_basic(ri, location)
	postConfig := testAccAzureRMSqlElasticPool_resizedDtu(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlElasticPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "dtu", "50"),
					resource.TestCheckResourceAttr(resourceName, "pool_size", "5000"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlElasticPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "dtu", "100"),
					resource.TestCheckResourceAttr(resourceName, "pool_size", "10000"),
				),
			},
		},
	})
}

func testCheckAzureRMSqlElasticPoolExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		poolName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sqlElasticPoolsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, poolName)
		if err != nil {
			return fmt.Errorf("Bad: Get on sqlElasticPoolsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: SQL Elastic Pool %q on server: %q (resource group: %q) does not exist", name, serverName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMSqlElasticPoolDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).sqlElasticPoolsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_sql_elasticpool" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		poolName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, serverName, poolName)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("SQL Elastic Pool still exists:\n%#v", resp.ElasticPoolProperties)
		}
	}

	return nil
}

func testCheckAzureRMSqlElasticPoolDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		poolName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sqlElasticPoolsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		_, err := client.Delete(ctx, resourceGroup, serverName, poolName)
		if err != nil {
			return fmt.Errorf("Bad: Delete on sqlElasticPoolsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMSqlElasticPool_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctest-%[1]d"
    location = "%s"
}

resource "azurerm_sql_server" "test" {
    name = "acctest%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    version = "12.0"
    administrator_login = "4dm1n157r470r"
    administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_sql_elasticpool" "test" {
    name = "acctest-pool-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    server_name = "${azurerm_sql_server.test.name}"
    edition = "Basic"
    dtu = 50
    pool_size = 5000
}
`, rInt, location)
}

func testAccAzureRMSqlElasticPool_resizedDtu(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctest-%[1]d"
    location = "%s"
}

resource "azurerm_sql_server" "test" {
    name = "acctest%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    version = "12.0"
    administrator_login = "4dm1n157r470r"
    administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_sql_elasticpool" "test" {
    name = "acctest-pool-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    server_name = "${azurerm_sql_server.test.name}"
    edition = "Basic"
    dtu = 100
    pool_size = 10000
}
`, rInt, location)
}
