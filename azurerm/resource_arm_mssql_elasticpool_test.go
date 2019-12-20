package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

// TODO: add import tests
func TestAccAzureRMMsSqlElasticPool_basic_DTU(t *testing.T) {
	resourceName := "azurerm_mssql_elasticpool.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlElasticPool_basic_DTU(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "BasicPool"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "50"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "0"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.max_capacity", "5"),
					resource.TestCheckResourceAttrSet(resourceName, "max_size_gb"),
					resource.TestCheckResourceAttrSet(resourceName, "zone_redundant"),
				),
			},
		},
	})
}

func TestAccAzureRMMsSqlElasticPool_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_mssql_elasticpool.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlElasticPool_basic_DTU(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMMsSqlElasticPool_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_mssql_elasticpool"),
			},
		},
	})
}

func TestAccAzureRMMsSqlElasticPool_standard_DTU(t *testing.T) {
	resourceName := "azurerm_mssql_elasticpool.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMsSqlElasticPool_standard_DTU(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "StandardPool"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "50"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "0"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.max_capacity", "50"),
					resource.TestCheckResourceAttrSet(resourceName, "max_size_gb"),
					resource.TestCheckResourceAttrSet(resourceName, "zone_redundant"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"max_size_gb"},
			},
		},
	})
}

func TestAccAzureRMMsSqlElasticPool_premium_DTU_zone_redundant(t *testing.T) {
	resourceName := "azurerm_mssql_elasticpool.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMsSqlElasticPool_premium_DTU_zone_redundant(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "PremiumPool"),
					resource.TestCheckResourceAttr(resourceName, "zone_redundant", "true"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"max_size_gb"},
			},
		},
	})
}

func TestAccAzureRMMsSqlElasticPool_basic_vCore(t *testing.T) {
	resourceName := "azurerm_mssql_elasticpool.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMsSqlElasticPool_basic_vCore(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "GP_Gen5"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "0.25"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.max_capacity", "4"),
					resource.TestCheckResourceAttrSet(resourceName, "max_size_gb"),
					resource.TestCheckResourceAttrSet(resourceName, "zone_redundant"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"max_size_gb"},
			},
		},
	})
}

func TestAccAzureRMMsSqlElasticPool_basic_vCore_MaxSizeBytes(t *testing.T) {
	resourceName := "azurerm_mssql_elasticpool.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMsSqlElasticPool_basic_vCore_MaxSizeBytes(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "GP_Gen5"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "0.25"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.max_capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "max_size_bytes", "214748364800"),
					resource.TestCheckResourceAttrSet(resourceName, "max_size_gb"),
					resource.TestCheckResourceAttrSet(resourceName, "zone_redundant"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"max_size_gb"},
			},
		},
	})
}
func TestAccAzureRMMsSqlElasticPool_disappears(t *testing.T) {
	resourceName := "azurerm_mssql_elasticpool.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMsSqlElasticPool_standard_DTU(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "StandardPool"),
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
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	preConfig := testAccAzureRMMsSqlElasticPool_standard_DTU(ri, location)
	postConfig := testAccAzureRMMsSqlElasticPool_resize_DTU(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "StandardPool"),
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
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	preConfig := testAccAzureRMMsSqlElasticPool_basic_vCore(ri, location)
	postConfig := testAccAzureRMMsSqlElasticPool_resize_vCore(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "GP_Gen5"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "4"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "0.25"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.max_capacity", "4"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "GP_Gen5"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "8"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.min_capacity", "0"),
					resource.TestCheckResourceAttr(resourceName, "per_database_settings.0.max_capacity", "8"),
				),
			},
		},
	})
}

func testCheckAzureRMMsSqlElasticPoolExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		poolName := rs.Primary.Attributes["name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.ElasticPoolsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, poolName)
		if err != nil {
			return fmt.Errorf("Bad: Get on msSqlElasticPoolsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: MsSql Elastic Pool %q on server: %q (resource group: %q) does not exist", poolName, serverName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMMsSqlElasticPoolDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.ElasticPoolsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testCheckAzureRMMsSqlElasticPoolDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		poolName := rs.Primary.Attributes["name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.ElasticPoolsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if _, err := client.Delete(ctx, resourceGroup, serverName, poolName); err != nil {
			return fmt.Errorf("Bad: Delete on msSqlElasticPoolsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMMsSqlElasticPool_basic_DTU(rInt int, location string) string {
	return testAccAzureRMMsSqlElasticPool_DTU_Template(rInt, location, "BasicPool", "Basic", 50, 4.8828125, 0, 5, false)
}

func testAccAzureRMMsSqlElasticPool_requiresImport(rInt int, location string) string {
	template := testAccAzureRMMsSqlElasticPool_DTU_Template(rInt, location, "BasicPool", "Basic", 50, 5242880000, 0, 5, false)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_elasticpool" "import" {
  name                  = "${azurerm_mssql_elasticpool.test.name}"
  resource_group_name   = "${azurerm_mssql_elasticpool.test.resource_group_name}"
  location              = "${azurerm_mssql_elasticpool.test.location}"
  server_name           = "${azurerm_mssql_elasticpool.test.server_name}"
  max_size_bytes        = "${azurerm_mssql_elasticpool.test.max_size_bytes}"
  sku                   = "${azurerm_mssql_elasticpool.test.sku}"
  per_database_settings = "${azurerm_mssql_elasticpool.test.per_database_settings}"
}
`, template)
}

func testAccAzureRMMsSqlElasticPool_premium_DTU_zone_redundant(rInt int, location string) string {
	return testAccAzureRMMsSqlElasticPool_DTU_Template(rInt, location, "PremiumPool", "Premium", 125, 50, 0, 50, true)
}

func testAccAzureRMMsSqlElasticPool_standard_DTU(rInt int, location string) string {
	return testAccAzureRMMsSqlElasticPool_DTU_Template(rInt, location, "StandardPool", "Standard", 50, 50, 0, 50, false)
}

func testAccAzureRMMsSqlElasticPool_resize_DTU(rInt int, location string) string {
	return testAccAzureRMMsSqlElasticPool_DTU_Template(rInt, location, "StandardPool", "Standard", 100, 100, 50, 100, false)
}

func testAccAzureRMMsSqlElasticPool_basic_vCore(rInt int, location string) string {
	return testAccAzureRMMsSqlElasticPool_vCore_Template(rInt, location, "GP_Gen5", "GeneralPurpose", 4, "Gen5", 0.25, 4)
}

func testAccAzureRMMsSqlElasticPool_basic_vCore_MaxSizeBytes(rInt int, location string) string {
	return testAccAzureRMMsSqlElasticPool_vCore_MaxSizeBytes_Template(rInt, location, "GP_Gen5", "GeneralPurpose", 4, "Gen5", 0.25, 4)
}

func testAccAzureRMMsSqlElasticPool_resize_vCore(rInt int, location string) string {
	return testAccAzureRMMsSqlElasticPool_vCore_Template(rInt, location, "GP_Gen5", "GeneralPurpose", 8, "Gen5", 0, 8)
}

func testAccAzureRMMsSqlElasticPool_vCore_Template(rInt int, location string, skuName string, skuTier string, skuCapacity int, skuFamily string, databaseSettingsMin float64, databaseSettingsMax float64) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctest%[1]d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-vcore-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  server_name         = "${azurerm_sql_server.test.name}"
  max_size_gb         = 5

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

func testAccAzureRMMsSqlElasticPool_vCore_MaxSizeBytes_Template(rInt int, location string, skuName string, skuTier string, skuCapacity int, skuFamily string, databaseSettingsMin float64, databaseSettingsMax float64) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctest%[1]d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-vcore-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  server_name         = "${azurerm_sql_server.test.name}"
  max_size_bytes      = 214748364800

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

func testAccAzureRMMsSqlElasticPool_DTU_Template(rInt int, location string, skuName string, skuTier string, skuCapacity int, maxSizeGB float64, databaseSettingsMin int, databaseSettingsMax int, zoneRedundant bool) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctest%[1]d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-dtu-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  server_name         = "${azurerm_sql_server.test.name}"
  max_size_gb         = %.7[6]f
  zone_redundant      = %[9]t

  sku {
    name     = "%[3]s"
    tier     = "%[4]s"
    capacity = %[5]d
  }

  per_database_settings {
    min_capacity = %[7]d
    max_capacity = %[8]d
  }
}
`, rInt, location, skuName, skuTier, skuCapacity, maxSizeGB, databaseSettingsMin, databaseSettingsMax, zoneRedundant)
}
