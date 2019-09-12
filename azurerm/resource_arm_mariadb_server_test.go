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

func TestAccAzureRMMariaDbServer_basic(t *testing.T) {
	resourceName := "azurerm_mariadb_server.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMariaDbServer_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(resourceName, "version", "10.2"),
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

func TestAccAzureRMMariaDbServer_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_mariadb_server.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMariaDbServer_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMMariaDbServer_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_mariadb_server"),
			},
		},
	})
}

func TestAccAzureRMMariaDbServer_basicMaxStorage(t *testing.T) {
	resourceName := "azurerm_mariadb_server.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMariaDbServer_basicMaxStorage(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(resourceName, "version", "10.2"),
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

func TestAccAzureRMMariaDbServer_generalPurpose(t *testing.T) {
	resourceName := "azurerm_mariadb_server.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMariaDbServer_generalPurpose(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(resourceName),
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

func TestAccAzureRMMariaDbServer_memoryOptimized(t *testing.T) {
	resourceName := "azurerm_mariadb_server.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMariaDbServer_memoryOptimizedGeoRedundant(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(resourceName),
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

func TestAccAzureRMMariaDbServer_updatePassword(t *testing.T) {
	resourceName := "azurerm_mariadb_server.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMMariaDbServer_basic(ri, location)
	updatedConfig := testAccAzureRMMariaDbServer_basicUpdatedPassword(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(resourceName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMMariaDbServer_updated(t *testing.T) {
	resourceName := "azurerm_mariadb_server.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMMariaDbServer_basic(ri, location)
	updatedConfig := testAccAzureRMMariaDbServer_basicUpdated(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "B_Gen5_2"),
					resource.TestCheckResourceAttr(resourceName, "version", "10.2"),
					resource.TestCheckResourceAttr(resourceName, "storage_profile.0.storage_mb", "51200"),
					resource.TestCheckResourceAttr(resourceName, "administrator_login", "acctestun"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "B_Gen5_1"),
					resource.TestCheckResourceAttr(resourceName, "version", "10.2"),
					resource.TestCheckResourceAttr(resourceName, "storage_profile.0.storage_mb", "640000"),
					resource.TestCheckResourceAttr(resourceName, "administrator_login", "acctestun"),
				),
			},
		},
	})
}

func TestAccAzureRMMariaDbServer_updateSKU(t *testing.T) {
	resourceName := "azurerm_mariadb_server.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMMariaDbServer_generalPurpose(ri, location)
	updatedConfig := testAccAzureRMMariaDbServer_memoryOptimized(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(resourceName),
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
					testCheckAzureRMMariaDbServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "MO_Gen5_16"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "16"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.tier", "MemoryOptimized"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.family", "Gen5"),
					resource.TestCheckResourceAttr(resourceName, "storage_profile.0.storage_mb", "4096000"),
					resource.TestCheckResourceAttr(resourceName, "administrator_login", "acctestun"),
				),
			},
		},
	})
}

func TestAccAzureRMMariaDbServer_storageAutogrow(t *testing.T) {
	resourceName := "azurerm_mariadb_server.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMMariaDbServer_basic(ri, location)
	updatedConfig := testAccAzureRMMariaDbServer_storageAutogrowUpdated(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "storage_profile.0.auto_grow", "Enabled"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "storage_profile.0.auto_grow", "Disabled"),
				),
			},
		},
	})
}

func testCheckAzureRMMariaDbServerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for MariaDB Server: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).mariadb.ServersClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: MariaDB Server %q (resource group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on mariadbServersClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMMariaDbServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).mariadb.ServersClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mariadb_server" {
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

		return fmt.Errorf("MariaDB Server still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMMariaDbServer_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "B_Gen5_2"
    capacity = 2
    tier     = "Basic"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}

func testAccAzureRMMariaDbServer_requiresImport(rInt int, location string) string {
	template := testAccAzureRMMariaDbServer_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_mariadb_server" "import" {
  name                = "${azurerm_mariadb_server.test.name}"
  location            = "${azurerm_mariadb_server.test.location}"
  resource_group_name = "${azurerm_mariadb_server.test.resource_group_name}"

  sku {
    name     = "B_Gen5_2"
    capacity = 2
    tier     = "Basic"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, template)
}

func testAccAzureRMMariaDbServer_basicUpdatedPassword(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "B_Gen5_2"
    capacity = 2
    tier     = "Basic"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "R3dH0TCh1l1P3pp3rs!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}

func testAccAzureRMMariaDbServer_basicUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "B_Gen5_1"
    capacity = 1
    tier     = "Basic"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 640000
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}

func testAccAzureRMMariaDbServer_basicMaxStorage(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "B_Gen5_2"
    capacity = 2
    tier     = "Basic"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 947200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}

func testAccAzureRMMariaDbServer_generalPurpose(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
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
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}

func testAccAzureRMMariaDbServer_memoryOptimized(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "MO_Gen5_16"
    capacity = 16
    tier     = "MemoryOptimized"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 4096000
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}

func testAccAzureRMMariaDbServer_memoryOptimizedGeoRedundant(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "MO_Gen5_16"
    capacity = 16
    tier     = "MemoryOptimized"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 4096000
    backup_retention_days = 7
    geo_redundant_backup  = "Enabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}

func testAccAzureRMMariaDbServer_storageAutogrowUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "B_Gen5_2"
    capacity = 2
    tier     = "Basic"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
	auto_grow      		  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}
