package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

// TODO: add import tests
func TestAccAzureRMMsSqlElasticPool_basic_DTU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlElasticPool_basic_DTU(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "BasicPool"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "50"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_database_settings.0.min_capacity", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_database_settings.0.max_capacity", "5"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "max_size_gb"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "zone_redundant"),
				),
			},
		},
	})
}

func TestAccAzureRMMsSqlElasticPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlElasticPool_basic_DTU(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMMsSqlElasticPool_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_mssql_elasticpool"),
			},
		},
	})
}

func TestAccAzureRMMsSqlElasticPool_standard_DTU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlElasticPool_standard_DTU(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "StandardPool"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "50"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_database_settings.0.min_capacity", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_database_settings.0.max_capacity", "50"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "max_size_gb"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "zone_redundant"),
				),
			},
			data.ImportStep("max_size_gb"),
		},
	})
}

func TestAccAzureRMMsSqlElasticPool_premium_DTU_zone_redundant(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlElasticPool_premium_DTU_zone_redundant(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "PremiumPool"),
					resource.TestCheckResourceAttr(data.ResourceName, "zone_redundant", "true"),
				),
			},
			data.ImportStep("max_size_gb"),
		},
	})
}

func TestAccAzureRMMsSqlElasticPool_basic_vCore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlElasticPool_basic_vCore(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "GP_Gen5"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_database_settings.0.min_capacity", "0.25"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_database_settings.0.max_capacity", "4"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "max_size_gb"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "zone_redundant"),
				),
			},
			data.ImportStep("max_size_gb"),
		},
	})
}

func TestAccAzureRMMsSqlElasticPool_basic_vCore_MaxSizeBytes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlElasticPool_basic_vCore_MaxSizeBytes(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "GP_Gen5"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_database_settings.0.min_capacity", "0.25"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_database_settings.0.max_capacity", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "max_size_bytes", "214748364800"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "max_size_gb"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "zone_redundant"),
				),
			},
			data.ImportStep("max_size_gb"),
		},
	})
}

func TestAccAzureRMMsSqlElasticPool_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlElasticPool_standard_DTU(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "StandardPool"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "50"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_database_settings.0.min_capacity", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_database_settings.0.max_capacity", "50"),
					testCheckAzureRMMsSqlElasticPoolDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMMsSqlElasticPool_resize_DTU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlElasticPool_standard_DTU(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "StandardPool"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "50"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_database_settings.0.min_capacity", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_database_settings.0.max_capacity", "50"),
				),
			},
			{
				Config: testAccAzureRMMsSqlElasticPool_resize_DTU(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "StandardPool"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_database_settings.0.min_capacity", "50"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_database_settings.0.max_capacity", "100"),
				),
			},
		},
	})
}

func TestAccAzureRMMsSqlElasticPool_resize_vCore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlElasticPool_basic_vCore(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "GP_Gen5"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_database_settings.0.min_capacity", "0.25"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_database_settings.0.max_capacity", "4"),
				),
			},
			{
				Config: testAccAzureRMMsSqlElasticPool_resize_vCore(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "GP_Gen5"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "8"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_database_settings.0.min_capacity", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_database_settings.0.max_capacity", "8"),
				),
			},
		},
	})
}

func TestAccAzureRMMsSqlElasticPool_licenseType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlElasticPool_licenseType_Template(data, "LicenseIncluded"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "license_type", "LicenseIncluded"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMsSqlElasticPool_licenseType_Template(data, "BasePrice"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "license_type", "BasePrice"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMMsSqlElasticPoolExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.ElasticPoolsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		poolName := rs.Primary.Attributes["name"]

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
		client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.ElasticPoolsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		poolName := rs.Primary.Attributes["name"]

		if _, err := client.Delete(ctx, resourceGroup, serverName, poolName); err != nil {
			return fmt.Errorf("Bad: Delete on msSqlElasticPoolsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMMsSqlElasticPool_basic_DTU(data acceptance.TestData) string {
	return testAccAzureRMMsSqlElasticPool_DTU_Template(data, "BasicPool", "Basic", 50, 4.8828125, 0, 5, false)
}

func testAccAzureRMMsSqlElasticPool_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlElasticPool_DTU_Template(data, "BasicPool", "Basic", 50, 4.8828125, 0, 5, false)
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_elasticpool" "import" {
  name                = azurerm_mssql_elasticpool.test.name
  resource_group_name = azurerm_mssql_elasticpool.test.resource_group_name
  location            = azurerm_mssql_elasticpool.test.location
  server_name         = azurerm_mssql_elasticpool.test.server_name
  max_size_gb         = 4.8828125

  sku {
    name     = "BasicPool"
    tier     = "Basic"
    capacity = 50
  }

  per_database_settings {
    min_capacity = 0
    max_capacity = 5
  }
}
`, template)
}

func testAccAzureRMMsSqlElasticPool_premium_DTU_zone_redundant(data acceptance.TestData) string {
	return testAccAzureRMMsSqlElasticPool_DTU_Template(data, "PremiumPool", "Premium", 125, 50, 0, 50, true)
}

func testAccAzureRMMsSqlElasticPool_standard_DTU(data acceptance.TestData) string {
	return testAccAzureRMMsSqlElasticPool_DTU_Template(data, "StandardPool", "Standard", 50, 50, 0, 50, false)
}

func testAccAzureRMMsSqlElasticPool_resize_DTU(data acceptance.TestData) string {
	return testAccAzureRMMsSqlElasticPool_DTU_Template(data, "StandardPool", "Standard", 100, 100, 50, 100, false)
}

func testAccAzureRMMsSqlElasticPool_basic_vCore(data acceptance.TestData) string {
	return testAccAzureRMMsSqlElasticPool_vCore_Template(data, "GP_Gen5", "GeneralPurpose", 4, "Gen5", 0.25, 4)
}

func testAccAzureRMMsSqlElasticPool_basic_vCore_MaxSizeBytes(data acceptance.TestData) string {
	return testAccAzureRMMsSqlElasticPool_vCore_MaxSizeBytes_Template(data, "GP_Gen5", "GeneralPurpose", 4, "Gen5", 0.25, 4)
}

func testAccAzureRMMsSqlElasticPool_resize_vCore(data acceptance.TestData) string {
	return testAccAzureRMMsSqlElasticPool_vCore_Template(data, "GP_Gen5", "GeneralPurpose", 8, "Gen5", 0, 8)
}

func testAccAzureRMMsSqlElasticPool_vCore_Template(data acceptance.TestData, skuName string, skuTier string, skuCapacity int, skuFamily string, databaseSettingsMin float64, databaseSettingsMax float64) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctest%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-vcore-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_sql_server.test.name
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
`, data.RandomInteger, data.Locations.Primary, skuName, skuTier, skuCapacity, skuFamily, databaseSettingsMin, databaseSettingsMax)
}

func testAccAzureRMMsSqlElasticPool_vCore_MaxSizeBytes_Template(data acceptance.TestData, skuName string, skuTier string, skuCapacity int, skuFamily string, databaseSettingsMin float64, databaseSettingsMax float64) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctest%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-vcore-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_sql_server.test.name
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
`, data.RandomInteger, data.Locations.Primary, skuName, skuTier, skuCapacity, skuFamily, databaseSettingsMin, databaseSettingsMax)
}

func testAccAzureRMMsSqlElasticPool_DTU_Template(data acceptance.TestData, skuName string, skuTier string, skuCapacity int, maxSizeGB float64, databaseSettingsMin int, databaseSettingsMax int, zoneRedundant bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctest%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-dtu-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_sql_server.test.name
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
`, data.RandomInteger, data.Locations.Primary, skuName, skuTier, skuCapacity, maxSizeGB, databaseSettingsMin, databaseSettingsMax, zoneRedundant)
}

func testAccAzureRMMsSqlElasticPool_licenseType_Template(data acceptance.TestData, licenseType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctest%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-dtu-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_sql_server.test.name
  max_size_gb         = 50
  zone_redundant      = false
  license_type        = "%[3]s"

  sku {
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
    capacity = 4
    family   = "Gen5"
  }

  per_database_settings {
    min_capacity = 0
    max_capacity = 4
  }

}
`, data.RandomInteger, data.Locations.Primary, licenseType)
}
