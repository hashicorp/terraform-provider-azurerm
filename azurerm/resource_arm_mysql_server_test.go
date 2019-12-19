package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMySQLServer_basicFiveSix(t *testing.T) {
	resourceName := "azurerm_mysql_server.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServer_basicFiveSix(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(resourceName),
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

func TestAccAzureRMMySQLServer_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_mysql_server.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServer_basicFiveSevenUpdated(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMMySQLServer_requiresImport(ri, acceptance.Location()),
				ExpectError: acceptance.RequiresImportError("azurerm_mysql_server"),
			},
		},
	})
}

func TestAccAzureRMMySQLServer_basicFiveSeven(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_mysql_server.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMySQLServer_basicFiveSeven(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(resourceName),
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

func TestAccAzureRMMySQLServer_basicEightZero(t *testing.T) {
	resourceName := "azurerm_mysql_server.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServer_basicEightZero(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(resourceName),
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

func TestAccAzureRMMySqlServer_generalPurpose(t *testing.T) {
	resourceName := "azurerm_mysql_server.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMySQLServer_generalPurpose(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(resourceName),
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

func TestAccAzureRMMySqlServer_memoryOptimized(t *testing.T) {
	resourceName := "azurerm_mysql_server.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMySQLServer_memoryOptimizedGeoRedundant(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(resourceName),
				),
			}, {
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

func TestAccAzureRMMySQLServer_basicFiveSevenUpdated(t *testing.T) {
	resourceName := "azurerm_mysql_server.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccAzureRMMySQLServer_basicFiveSeven(ri, location)
	updatedConfig := testAccAzureRMMySQLServer_basicFiveSevenUpdated(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "GP_Gen5_2"),
					resource.TestCheckResourceAttr(resourceName, "version", "5.7"),
					resource.TestCheckResourceAttr(resourceName, "storage_profile.0.storage_mb", "51200"),
					resource.TestCheckResourceAttr(resourceName, "administrator_login", "acctestun"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "GP_Gen5_4"),
					resource.TestCheckResourceAttr(resourceName, "version", "5.7"),
					resource.TestCheckResourceAttr(resourceName, "storage_profile.0.storage_mb", "640000"),
					resource.TestCheckResourceAttr(resourceName, "administrator_login", "acctestun"),
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

func TestAccAzureRMMySQLServer_updateSKU(t *testing.T) {
	resourceName := "azurerm_mysql_server.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccAzureRMMySQLServer_generalPurpose(ri, location)
	updatedConfig := testAccAzureRMMySQLServer_memoryOptimized(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(resourceName),
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
					testCheckAzureRMMySQLServerExists(resourceName),
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

func TestAccAzureRMMySQLServer_storageAutogrow(t *testing.T) {
	resourceName := "azurerm_mysql_server.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccAzureRMMySQLServer_basicFiveSeven(ri, location)
	updatedConfig := testAccAzureRMMySQLServer_autogrow(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "storage_profile.0.auto_grow", "Enabled"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "storage_profile.0.auto_grow", "Disabled"),
				),
			},
		},
	})
}

//

func testCheckAzureRMMySQLServerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for MySQL Server: %s", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).MySQL.ServersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MySQL Server %q (resource group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on mysqlServersClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMMySQLServerDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MySQL.ServersClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mysql_server" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("MySQL Server still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMMySQLServer_basicFiveSix(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
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
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "5.6"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}

func testAccAzureRMMySQLServer_basicFiveSeven(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
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
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "5.7"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}

func testAccAzureRMMySQLServer_basicEightZero(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
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
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "8.0"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}

func testAccAzureRMMySQLServer_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_server" "import" {
  name                = "${azurerm_mysql_server.test.name}"
  location            = "${azurerm_mysql_server.test.location}"
  resource_group_name = "${azurerm_mysql_server.test.name}"

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
  version                      = "5.7"
  ssl_enforcement              = "Enabled"
}
`, testAccAzureRMMySQLServer_basicFiveSevenUpdated(rInt, location))
}

func testAccAzureRMMySQLServer_basicFiveSevenUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
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
  version                      = "5.7"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}

func testAccAzureRMMySQLServer_generalPurpose(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
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
  version                      = "5.7"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}

func testAccAzureRMMySQLServer_memoryOptimized(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
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
  version                      = "5.7"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}

func testAccAzureRMMySQLServer_memoryOptimizedGeoRedundant(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
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
  version                      = "5.7"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}

func testAccAzureRMMySQLServer_autogrow(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
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
  version                      = "5.7"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}
