package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMMsSqlElasticPool_basic_DTU(t *testing.T) {
	resourceName := "azurerm_mssql_elasticpool.test"
	ri := acctest.RandInt()
	config := testAccAzureRMMsSqlElasticPool_basic_DTU(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "StandardPool"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "50"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "0"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.max_capacity", "50"),
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

func TestAccAzureRMMsSqlElasticPool_basic_vCore(t *testing.T) {
	resourceName := "azurerm_mssql_elasticpool.test"
	ri := acctest.RandInt()
	config := testAccAzureRMMsSqlElasticPool_basic_vCore(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "GP_Gen5"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "GeneralPurpose"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.family", "Gen5"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "0.25"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.max_capacity", "4"),
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

func TestAccAzureRMMsSqlElasticPool_disappears(t *testing.T) {
	resourceName := "azurerm_mssql_elasticpool.test"
	ri := acctest.RandInt()
	config := testAccAzureRMMsSqlElasticPool_basic_DTU(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "StandardPool"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "50"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "0"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.max_capacity", "50"),
					testCheckAzureRMMsSqlElasticPoolDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMMsSqlElasticPool_resize_DTU(t *testing.T) {
	resourceName := "azurerm_mssql_elasticpool.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMMsSqlElasticPool_basic_DTU(ri, location)
	postConfig := testAccAzureRMMsSqlElasticPool_resize_DTU(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "StandardPool"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "50"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "0"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.max_capacity", "50"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "StandardPool"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "100"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "50"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.max_capacity", "100"),
				),
			},
		},
	})
}

func TestAccAzureRMMsSqlElasticPool_resize_vCore(t *testing.T) {
	resourceName := "azurerm_mssql_elasticpool.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMMsSqlElasticPool_basic_vCore(ri, location)
	postConfig := testAccAzureRMMsSqlElasticPool_resize_vCore(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "GP_Gen5"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "GeneralPurpose"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.family", "Gen5"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "0.25"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.max_capacity", "4"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "GP_Gen5"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "GeneralPurpose"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "8"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.family", "Gen5"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "0"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.max_capacity", "8"),
				),
			},
		},
	})
}

func testCheckAzureRMMsSqlElasticPoolExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		poolName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).msSqlElasticPoolsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, poolName)
		if err != nil {
			return fmt.Errorf("Bad: Get on msSqlElasticPoolsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: MsSql Elastic Pool %q on server: %q (resource group: %q) does not exist", name, serverName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMMsSqlElasticPoolDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).msSqlElasticPoolsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mssql_elasticpool" {
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
			return fmt.Errorf("MsSql Elastic Pool still exists:\n%#v", resp.ElasticPoolProperties)
		}
	}

	return nil
}

func testCheckAzureRMMsSqlElasticPoolDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		poolName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).msSqlElasticPoolsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if _, err := client.Delete(ctx, resourceGroup, serverName, poolName); err != nil {
			return fmt.Errorf("Bad: Delete on msSqlElasticPoolsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMMsSqlElasticPool_basic_DTU(rInt int, location string) string {
	return testAccAzureRMMsSqlElasticPool_DTU_Template(rInt, location, "StandardPool", "Standard", 50, 0, 50)
}

func testAccAzureRMMsSqlElasticPool_resize_DTU(rInt int, location string) string {
	return testAccAzureRMMsSqlElasticPool_DTU_Template(rInt, location, "StandardPool", "Standard", 100, 50, 100)
}

func testAccAzureRMMsSqlElasticPool_basic_vCore(rInt int, location string) string {
	return testAccAzureRMMsSqlElasticPool_vCore_Template(rInt, location, "GP_Gen5", "GeneralPurpose", 4, "Gen5", 0.25, 4)
}

func testAccAzureRMMsSqlElasticPool_resize_vCore(rInt int, location string) string {
	return testAccAzureRMMsSqlElasticPool_vCore_Template(rInt, location, "GP_Gen5", "GeneralPurpose", 8, "Gen5", 0, 8)
}

func testAccAzureRMMsSqlElasticPool_vCore_Template(rInt int, location string, skuName string, skuTier string, skuCapacity int, skuFamily string, databaseSettingsMin float64, databaseSettingsMax float64) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%[1]d"
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

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-vcore-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  server_name         = "${azurerm_sql_server.test.name}"

  sku {
    name     = "%[3]s"
    tier     = "%[4]s"
    capacity = %[5]d
    family   = "%[6]s"
  }

  per_database_settings {
    min_capacity = %.2[7]f
    max_capacity = %.2[8]f
  }
}
`, rInt, location, skuName, skuTier, skuCapacity, skuFamily, databaseSettingsMin, databaseSettingsMax)
}

func testAccAzureRMMsSqlElasticPool_DTU_Template(rInt int, location string, skuName string, skuTier string, skuCapacity int, databaseSettingsMin int, databaseSettingsMax int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%[1]d"
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

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-dtu-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  server_name         = "${azurerm_sql_server.test.name}"

  sku {
    name     = "%[3]s"
    tier     = "%[4]s"
    capacity = %[5]d
  }

  per_database_settings {
    min_capacity = %[6]d
    max_capacity = %[7]d
  }
}
`, rInt, location, skuName, skuTier, skuCapacity, databaseSettingsMin, databaseSettingsMax)
}
