package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPostgreSQLServer_basicNinePointFive(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLServer_basicNinePointFive(ri, testLocation())

	resource.Test(t, resource.TestCase{
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
		},
	})
}

func TestAccAzureRMPostgreSQLServer_basicNinePointSix(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLServer_basicNinePointSix(ri, testLocation())

	resource.Test(t, resource.TestCase{
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
		},
	})
}

func TestAccAzureRMPostgreSQLServer_basicMaxStorage(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLServer_basicMaxStorage(ri, testLocation())

	resource.Test(t, resource.TestCase{
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
		},
	})
}

func TestAccAzureRMPostgreSQLServer_generalPurpose(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLServer_generalPurpose(ri, testLocation())

	resource.Test(t, resource.TestCase{
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
		},
	})
}

func TestAccAzureRMPostgreSQLServer_memoryOptimized(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLServer_memoryOptimized(ri, testLocation())

	resource.Test(t, resource.TestCase{
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
		},
	})
}

func TestAccAzureRMPostgreSQLServer_updatePassword(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMPostgreSQLServer_basicNinePointSix(ri, location)
	updatedConfig := testAccAzureRMPostgreSQLServer_basicNinePointSixUpdatedPassword(ri, location)

	resource.Test(t, resource.TestCase{
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

func TestAccAzureRMPostgreSQLServer_updateSKU(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMPostgreSQLServer_basicNinePointSix(ri, location)
	updatedConfig := testAccAzureRMPostgreSQLServer_basicNinePointSixUpdatedSKU(ri, location)

	resource.Test(t, resource.TestCase{
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

func testCheckAzureRMPostgreSQLServerExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for PostgreSQL Server: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).postgresqlServersClient
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
	client := testAccProvider.Meta().(*ArmClient).postgresqlServersClient
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

func testAccAzureRMPostgreSQLServer_basicNinePointFive(rInt int, location string) string {
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
    name     = "B_Gen4_2"
    capacity = 2
    tier     = "Basic"
    family   = "Gen4"
  }

  storage_profile {
    storage_mb = 51200
    backup_retention_days = 7
    geo_redundant_backup = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}

func testAccAzureRMPostgreSQLServer_basicNinePointSix(rInt int, location string) string {
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
    name     = "B_Gen4_2"
    capacity = 2
    tier     = "Basic"
    family   = "Gen4"
  }

  storage_profile {
    storage_mb = 51200
    backup_retention_days = 7
    geo_redundant_backup = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.6"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
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
    name     = "B_Gen4_2"
    capacity = 2
    tier     = "Basic"
    family   = "Gen4"
  }

  storage_profile {
    storage_mb = 51200
    backup_retention_days = 7
    geo_redundant_backup = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "R3dH0TCh1l1P3pp3rs!"
  version                      = "9.6"
  ssl_enforcement              = "Disabled"
}

`, rInt, location, rInt)
}

func testAccAzureRMPostgreSQLServer_basicNinePointSixUpdatedSKU(rInt int, location string) string {
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
    name     = "B_Gen4_1"
    capacity = 1
    tier     = "Basic"
    family   = "Gen4"
  }

  storage_profile {
    storage_mb = 51200
    backup_retention_days = 7
    geo_redundant_backup = "Disabled"
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
    name     = "B_Gen4_2"
    capacity = 2
    tier     = "Basic"
    family   = "Gen4"
  }

  storage_profile {
    storage_mb = 947200
    backup_retention_days = 7
    geo_redundant_backup = "Disabled"
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
    storage_mb = 640000
    backup_retention_days = 7
    geo_redundant_backup = "Disabled"
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
    storage_mb = 1048576
    backup_retention_days = 7
    geo_redundant_backup = "Enabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.6"
  ssl_enforcement              = "Enabled"
}
`, rInt, location, rInt)
}
