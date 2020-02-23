package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPostgreSQLServer_basicNinePointFive(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_basicNinePointFive(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "9.5"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_enforcement", "Enabled"),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_basicNinePointSix(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_basicNinePointSix(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "9.6"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_enforcement", "Enabled"),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_basicTenPointZero(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_basicTenPointZero(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "10.0"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_enforcement", "Enabled"),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_basicEleven(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_basicEleven(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "11"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_enforcement", "Enabled"),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_basicTenPointZero(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "10.0"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_enforcement", "Enabled"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMPostgreSQLServer_requiresImport),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_basicMaxStorage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_basicMaxStorage(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "9.6"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_enforcement", "Enabled"),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_generalPurpose(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_generalPurpose(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_memoryOptimized(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_memoryOptimizedGeoRedundant(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_updatePassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_basicNinePointSix(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMPostgreSQLServer_basicNinePointSixUpdatedPassword(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLServer_updated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_basicNinePointSix(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "GP_Gen5_2"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "9.6"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_profile.0.storage_mb", "51200"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_profile.0.auto_grow", "Disabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
				),
			},
			{
				Config: testAccAzureRMPostgreSQLServer_basicNinePointSixUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "GP_Gen5_4"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "9.6"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_profile.0.storage_mb", "640000"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_profile.0.auto_grow", "Enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLServer_updateSKU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_generalPurpose(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "GP_Gen5_32"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_profile.0.storage_mb", "640000"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
				),
			},
			{
				Config: testAccAzureRMPostgreSQLServer_memoryOptimized(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "MO_Gen5_16"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_profile.0.storage_mb", "4194304"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
				),
			},
		},
	})
}

//

func testCheckAzureRMPostgreSQLServerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Postgres.ServersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for PostgreSQL Server: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: PostgreSQL Server %q (resource group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on postgresqlServersClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMPostgreSQLServerDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Postgres.ServersClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_postgresql_server" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("PostgreSQL Server still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMPostgreSQLServer_basic(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
    auto_grow             = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "%s"
  ssl_enforcement              = "Enabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func testAccAzureRMPostgreSQLServer_basicNinePointFive(data acceptance.TestData) string {
	return testAccAzureRMPostgreSQLServer_basic(data, "9.5")
}

func testAccAzureRMPostgreSQLServer_basicNinePointSix(data acceptance.TestData) string {
	return testAccAzureRMPostgreSQLServer_basic(data, "9.6")
}

func testAccAzureRMPostgreSQLServer_basicTenPointZero(data acceptance.TestData) string {
	return testAccAzureRMPostgreSQLServer_basic(data, "10.0")
}

func testAccAzureRMPostgreSQLServer_basicEleven(data acceptance.TestData) string {
	return testAccAzureRMPostgreSQLServer_basic(data, "11")
}

func testAccAzureRMPostgreSQLServer_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMPostgreSQLServer_basicTenPointZero(data)
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_server" "import" {
  name                = "${azurerm_postgresql_server.test.name}"
  location            = "${azurerm_postgresql_server.test.location}"
  resource_group_name = "${azurerm_postgresql_server.test.resource_group_name}"

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.0"
  ssl_enforcement              = "Enabled"
}
`, template)
}

func testAccAzureRMPostgreSQLServer_basicNinePointSixUpdatedPassword(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "R3dH0TCh1l1P3pp3rs!"
  version                      = "9.6"
  ssl_enforcement              = "Disabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMPostgreSQLServer_basicNinePointSixUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "GP_Gen5_4"

  storage_profile {
    storage_mb            = 640000
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.6"
  ssl_enforcement              = "Enabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMPostgreSQLServer_basicMaxStorage(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 947200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
    auto_grow             = "Enabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.6"
  ssl_enforcement              = "Enabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMPostgreSQLServer_generalPurpose(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "GP_Gen5_32"

  storage_profile {
    storage_mb            = 640000
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.6"
  ssl_enforcement              = "Enabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMPostgreSQLServer_memoryOptimized(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "MO_Gen5_16"

  storage_profile {
    storage_mb            = 4194304
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.6"
  ssl_enforcement              = "Enabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMPostgreSQLServer_memoryOptimizedGeoRedundant(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku_name            = "MO_Gen5_16"

  storage_profile {
    storage_mb            = 4194304
    backup_retention_days = 7
    geo_redundant_backup  = "Enabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.6"
  ssl_enforcement              = "Enabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
