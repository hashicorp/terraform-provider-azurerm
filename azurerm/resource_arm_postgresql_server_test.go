package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPostgreSQLServer_basicNinePointFive(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPostgreSQLServer_basicNinePointFive(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(resourceName, "version", "9.5"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enforcement", "Enabled"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"administrator_login_password", // not returned as sensitive
				},
			},
		},
	})
}

func TestAccAzureRMPostgreSQLServer_basicNinePointSix(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPostgreSQLServer_basicNinePointSix(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(resourceName, "version", "9.6"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enforcement", "Enabled"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"administrator_login_password", // not returned as sensitive
				},
			},
		},
	})
}

func TestAccAzureRMPostgreSQLServer_basicTenPointZero(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_basicTenPointZero(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(resourceName, "version", "10.0"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enforcement", "Enabled"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"administrator_login_password", // not returned as sensitive
				},
			},
		},
	})
}

func TestAccAzureRMPostgreSQLServer_basicEleven(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_basicEleven(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(resourceName, "version", "11"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enforcement", "Enabled"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"administrator_login_password", // not returned as sensitive
				},
			},
		},
	})
}

func TestAccAzureRMPostgreSQLServer_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_postgresql_server.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_basicTenPointZero(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(resourceName, "version", "10.0"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enforcement", "Enabled"),
				),
			},
			{
				Config:      testAccAzureRMPostgreSQLServer_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_postgresql_server"),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLServer_basicMaxStorage(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPostgreSQLServer_basicMaxStorage(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(resourceName, "version", "9.6"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enforcement", "Enabled"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"administrator_login_password", // not returned as sensitive
				},
			},
		},
	})
}

func TestAccAzureRMPostgreSQLServer_generalPurpose(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPostgreSQLServer_generalPurpose(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"administrator_login_password", // not returned as sensitive
				},
			},
		},
	})
}

func TestAccAzureRMPostgreSQLServer_memoryOptimized(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPostgreSQLServer_memoryOptimizedGeoRedundant(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"administrator_login_password", // not returned as sensitive
				},
			},
		},
	})
}

func TestAccAzureRMPostgreSQLServer_updatePassword(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMPostgreSQLServer_basicNinePointSix(ri, location)
	updatedConfig := testAccAzureRMPostgreSQLServer_basicNinePointSixUpdatedPassword(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(resourceName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLServer_updated(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMPostgreSQLServer_basicNinePointSix(ri, location)
	updatedConfig := testAccAzureRMPostgreSQLServer_basicNinePointSixUpdated(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "GP_Gen5_2"),
					resource.TestCheckResourceAttr(resourceName, "version", "9.6"),
					resource.TestCheckResourceAttr(resourceName, "storage_profile.0.storage_mb", "51200"),
					resource.TestCheckResourceAttr(resourceName, "storage_profile.0.auto_grow", "Disabled"),
					resource.TestCheckResourceAttr(resourceName, "administrator_login", "acctestun"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "GP_Gen5_4"),
					resource.TestCheckResourceAttr(resourceName, "version", "9.6"),
					resource.TestCheckResourceAttr(resourceName, "storage_profile.0.storage_mb", "640000"),
					resource.TestCheckResourceAttr(resourceName, "storage_profile.0.auto_grow", "Enabled"),
					resource.TestCheckResourceAttr(resourceName, "administrator_login", "acctestun"),
				),
			},
		},
	})
}

func TestAccAzureRMPostgreSQLServer_updateSKU(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMPostgreSQLServer_generalPurpose(ri, location)
	updatedConfig := testAccAzureRMPostgreSQLServer_memoryOptimized(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "GP_Gen5_32"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "32"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "GeneralPurpose"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.family", "Gen5"),
					resource.TestCheckResourceAttr(resourceName, "storage_profile.0.storage_mb", "640000"),
					resource.TestCheckResourceAttr(resourceName, "administrator_login", "acctestun"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "MO_Gen5_16"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "16"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "MemoryOptimized"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.family", "Gen5"),
					resource.TestCheckResourceAttr(resourceName, "storage_profile.0.storage_mb", "4194304"),
					resource.TestCheckResourceAttr(resourceName, "administrator_login", "acctestun"),
				),
			},
		},
	})
}

//

func testCheckAzureRMPostgreSQLServerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
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

		client := testAccProvider.Meta().(*ArmClient).postgres.ServersClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
	client := testAccProvider.Meta().(*ArmClient).postgres.ServersClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMPostgreSQLServer_basic(rInt int, location string, version string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctestpsqlsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "GP_Gen5_2"
    capacity = 2
    tier     = "GeneralPurpose"
    family   = "Gen5"
  }

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
`, rInt, location, rInt, version)
}

func testAccAzureRMPostgreSQLServer_basicNinePointFive(rInt int, location string) string {
	return testAccAzureRMPostgreSQLServer_basic(rInt, location, "9.5")
}

func testAccAzureRMPostgreSQLServer_basicNinePointSix(rInt int, location string) string {
	return testAccAzureRMPostgreSQLServer_basic(rInt, location, "9.6")
}

func testAccAzureRMPostgreSQLServer_basicTenPointZero(rInt int, location string) string {
	return testAccAzureRMPostgreSQLServer_basic(rInt, location, "10.0")
}

func testAccAzureRMPostgreSQLServer_basicEleven(rInt int, location string) string {
	return testAccAzureRMPostgreSQLServer_basic(rInt, location, "11")
}

func testAccAzureRMPostgreSQLServer_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_server" "import" {
  name                = "${azurerm_postgresql_server.test.name}"
  location            = "${azurerm_postgresql_server.test.location}"
  resource_group_name = "${azurerm_postgresql_server.test.resource_group_name}"

  sku {
    name     = "GP_Gen5_2"
    capacity = 2
    tier     = "GeneralPurpose"
    family   = "Gen5"
  }

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
`, testAccAzureRMPostgreSQLServer_basicTenPointZero(rInt, location))
}

func testAccAzureRMPostgreSQLServer_basicNinePointSixUpdatedPassword(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctestpsqlsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "GP_Gen5_2"
    capacity = 2
    tier     = "GeneralPurpose"
    family   = "Gen5"
  }

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
`, rInt, location, rInt)
}

func testAccAzureRMPostgreSQLServer_basicNinePointSixUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctestpsqlsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "GP_Gen5_4"
    capacity = 4
    tier     = "GeneralPurpose"
    family   = "Gen5"
  }

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
`, rInt, location, rInt)
}

func testAccAzureRMPostgreSQLServer_basicMaxStorage(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctestpsqlsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "GP_Gen5_2"
    capacity = 2
    tier     = "GeneralPurpose"
    family   = "Gen5"
  }

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
`, rInt, location, rInt)
}

func testAccAzureRMPostgreSQLServer_generalPurpose(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctestpsqlsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "GP_Gen5_32"
    capacity = 32
    tier     = "GeneralPurpose"
    family   = "Gen5"
  }

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
`, rInt, location, rInt)
}

func testAccAzureRMPostgreSQLServer_memoryOptimized(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctestpsqlsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "MO_Gen5_16"
    capacity = 16
    tier     = "MemoryOptimized"
    family   = "Gen5"
  }

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
`, rInt, location, rInt)
}

func testAccAzureRMPostgreSQLServer_memoryOptimizedGeoRedundant(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctestpsqlsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "MO_Gen5_16"
    capacity = 16
    tier     = "MemoryOptimized"
    family   = "Gen5"
  }

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
`, rInt, location, rInt)
}
