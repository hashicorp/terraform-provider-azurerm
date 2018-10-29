package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMSqlElasticPool2017_basic_DTU(t *testing.T) {
	resourceName := "azurerm_sql2017_elasticpool.test"
	ri := acctest.RandInt()
	config := testAccAzureRMSqlElasticPool2017_basic_DTU(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlElasticPool2017Destroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlElasticPool2017Exists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "BasicPool"),
				),
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		},
	})
}

func TestAccAzureRMSqlElasticPool2017_basic_vCore(t *testing.T) {
	resourceName := "azurerm_sql2017_elasticpool.test"
	ri := acctest.RandInt()
	config := testAccAzureRMSqlElasticPool2017_basic_vCore(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlElasticPool2017Destroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlElasticPool2017Exists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "GP_Gen5"),
				),
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		},
	})
}

func TestAccAzureRMSqlElasticPool2017_disappears(t *testing.T) {
	resourceName := "azurerm_sql2017_elasticpool.test"
	ri := acctest.RandInt()
	config := testAccAzureRMSqlElasticPool2017_basic_DTU(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlElasticPool2017Destroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlElasticPool2017Exists(resourceName),
					testCheckAzureRMSqlElasticPool2017Disappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMSqlElasticPool2017_resize_DTU(t *testing.T) {
	resourceName := "azurerm_sql_elasticpool.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMSqlElasticPool2017_basic_DTU(ri, location)
	postConfig := testAccAzureRMSqlElasticPool2017_resize_DTU(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlElasticPool2017Destroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlElasticPool2017Exists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "0"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "100"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlElasticPool2017Exists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "50"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "1000"),
				),
			},
		},
	})
}

func TestAccAzureRMSqlElasticPool2017_resize_vCore(t *testing.T) {
	resourceName := "azurerm_sql_elasticpool.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMSqlElasticPool2017_basic_vCore(ri, location)
	postConfig := testAccAzureRMSqlElasticPool2017_resize_vCore(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlElasticPool2017Destroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlElasticPool2017Exists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "0.25"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "4"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlElasticPool2017Exists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "0.0"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "8"),
				),
			},
		},
	})
}

func testCheckAzureRMSqlElasticPool2017Exists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		poolName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sql2017ElasticPoolsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, poolName)
		if err != nil {
			return fmt.Errorf("Bad: Get on sql2017ElasticPoolsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: SQL2017 Elastic Pool %q on server: %q (resource group: %q) does not exist", name, serverName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMSqlElasticPool2017Destroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).sql2017ElasticPoolsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_sql2017_elasticpool" {
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
			return fmt.Errorf("SQL2017 Elastic Pool still exists:\n%#v", resp.ElasticPoolProperties)
		}
	}

	return nil
}

func testCheckAzureRMSqlElasticPool2017Disappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		poolName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sql2017ElasticPoolsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		_, err := client.Delete(ctx, resourceGroup, serverName, poolName)
		if err != nil {
			return fmt.Errorf("Bad: Delete on sql2017ElasticPoolsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMSqlElasticPool2017_basic_DTU(rInt int, location string) string {
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

resource "azurerm_sql2017_elasticpool" "test" {
  name                = "acctest-pool-dtu-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  server_name         = "${azurerm_sql_server.test.name}"

  sku {
    name     = "BasicPool"
    tier     = "Basic"
    capacity = 1000
  }

  per_database_settings {
    min_capacity = 0
    max_capacity = 100
  }
}
`, rInt, location)
}

func testAccAzureRMSqlElasticPool2017_basic_vCore(rInt int, location string) string {
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

resource "azurerm_sql2017_elasticpool" "test" {
  name                = "acctest-pool-vcore-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  server_name         = "${azurerm_sql_server.test.name}"

  sku {
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
    capacity = 4
    family   = "Gen5"
  }

  per_database_settings {
    min_capacity = 0.25
    max_capacity = 4
  }
}
`, rInt, location)
}

func testAccAzureRMSqlElasticPool2017_resize_DTU(rInt int, location string) string {
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

resource "azurerm_sql2017_elasticpool" "test" {
  name                = "acctest-pool-dtu-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  server_name         = "${azurerm_sql_server.test.name}"

  sku {
    name     = "BasicPool"
    tier     = "Basic"
    capacity = 1000
  }

  per_database_settings {
    min_capacity = 50
    max_capacity = 1000
  }
}
`, rInt, location)
}

func testAccAzureRMSqlElasticPool2017_resize_vCore(rInt int, location string) string {
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

resource "azurerm_sql2017_elasticpool" "test" {
  name                = "acctest-pool-vcore-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  server_name         = "${azurerm_sql_server.test.name}"

  sku {
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
    capacity = 8
    family   = "Gen5"
  }

  per_database_settings {
    min_capacity = 0
    max_capacity = 8
  }
}
`, rInt, location)
}
