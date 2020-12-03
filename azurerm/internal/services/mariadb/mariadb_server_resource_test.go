package mariadb_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccMariaDbServer_basicTenTwo(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")
	version := "10.2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariaDbServer_basic(data, version),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "version", version),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccMariaDbServer_basicTenTwoDeprecated(t *testing.T) { // remove in v3.0
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")
	version := "10.2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariaDbServer_basicDeprecated(data, version),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "version", version),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccMariaDbServer_basicTenThree(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")
	version := "10.3"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariaDbServer_basic(data, version),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "version", version),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccMariaDbServer_autogrowOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")
	version := "10.3"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariaDbServer_autogrow(data, version),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccMariaDbServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariaDbServer_basic(data, "10.3"),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccMariaDbServer_requiresImport),
		},
	})
}

func TestAccMariaDbServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariaDbServer_complete(data, "10.3"),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccMariaDbServer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")
	version := "10.3"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariaDbServer_basic(data, version),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
			{
				Config: testAccMariaDbServer_complete(data, version),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
			{
				Config: testAccMariaDbServer_basic(data, version),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccMariaDbServer_completeDeprecatedMigrate(t *testing.T) { // remove in v3.0
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")
	version := "10.3"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariaDbServer_completeDeprecated(data, version),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
			{
				Config: testAccMariaDbServer_complete(data, version),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccMariaDbServer_updateDeprecated(t *testing.T) { // remove in v3.0
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")
	version := "10.2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariaDbServer_basicDeprecated(data, version),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
			{
				Config: testAccMariaDbServer_completeDeprecated(data, version),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
			{
				Config: testAccMariaDbServer_basicDeprecated(data, version),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccMariaDbServer_updateSKU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariaDbServer_sku(data, "GP_Gen5_32"),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
			{
				Config: testAccMariaDbServer_sku(data, "MO_Gen5_16"),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccMariaDbServer_createReplica(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")
	version := "10.3"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariaDbServer_basic(data, version),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
			{
				Config: testAccMariaDbServer_createReplica(data, version),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
					testCheckMariaDbServerExists("azurerm_mariadb_server.replica"),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccMariaDbServer_createPointInTimeRestore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_server", "test")
	restoreTime := time.Now().Add(11 * time.Minute)
	version := "10.3"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariaDbServer_basic(data, version),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
			{
				PreConfig: func() { time.Sleep(restoreTime.Sub(time.Now().Add(-7 * time.Minute))) },
				Config:    testAccMariaDbServer_createPointInTimeRestore(data, version, restoreTime.Format(time.RFC3339)),
				Check: resource.ComposeTestCheckFunc(
					testCheckMariaDbServerExists(data.ResourceName),
					testCheckMariaDbServerExists("azurerm_mariadb_server.restore"),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func testCheckMariaDbServerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MariaDB.ServersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testCheckMariaDbServerDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MariaDB.ServersClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccMariaDbServer_basic(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "B_Gen5_2"
  version             = "%s"

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  ssl_enforcement_enabled      = true
  storage_mb                   = 51200
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func testAccMariaDbServer_basicDeprecated(data acceptance.TestData, version string) string { // remove in v3.0
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "B_Gen5_2"
  version             = "%s"

  storage_profile {
    storage_mb = 51200
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  ssl_enforcement_enabled      = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func testAccMariaDbServer_complete(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "B_Gen5_2"
  version             = "%s"

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  auto_grow_enabled            = true
  backup_retention_days        = 14
  create_mode                  = "Default"
  geo_redundant_backup_enabled = false
  ssl_enforcement_enabled      = true
  storage_mb                   = 51200
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func testAccMariaDbServer_completeDeprecated(data acceptance.TestData, version string) string { // remove in v3.0
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "B_Gen5_2"
  version             = "%s"

  storage_profile {
    auto_grow             = "Enabled"
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
    storage_mb            = 51200
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  create_mode                  = "Default"
  ssl_enforcement_enabled      = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func testAccMariaDbServer_autogrow(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "B_Gen5_2"
  version             = "%s"

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  auto_grow_enabled            = true
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  ssl_enforcement_enabled      = true
  storage_mb                   = 51200
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func testAccMariaDbServer_requiresImport(data acceptance.TestData) string {
	template := testAccMariaDbServer_basic(data, "10.3")
	return fmt.Sprintf(`
%s

resource "azurerm_mariadb_server" "import" {
  name                = azurerm_mariadb_server.test.name
  location            = azurerm_mariadb_server.test.location
  resource_group_name = azurerm_mariadb_server.test.resource_group_name
  sku_name            = "B_Gen5_2"
  version             = "10.3"

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  ssl_enforcement_enabled      = true
  storage_mb                   = 51200
}
`, template)
}

func testAccMariaDbServer_sku(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "%s"
  version             = "10.2"

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  ssl_enforcement_enabled      = true
  storage_mb                   = 640000
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, sku)
}

func testAccMariaDbServer_createReplica(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mariadb_server" "replica" {
  name                      = "acctestmariadbsvr-%d-replica"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  sku_name                  = "B_Gen5_2"
  version                   = "%s"
  create_mode               = "Replica"
  creation_source_server_id = azurerm_mariadb_server.test.id
  ssl_enforcement_enabled   = true
  storage_mb                = 51200
}
`, testAccMariaDbServer_basic(data, version), data.RandomInteger, version)
}

func testAccMariaDbServer_createPointInTimeRestore(data acceptance.TestData, version, restoreTime string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mariadb_server" "restore" {
  name                      = "acctestmariadbsvr-%d-restore"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  sku_name                  = "B_Gen5_2"
  version                   = "%s"
  create_mode               = "PointInTimeRestore"
  creation_source_server_id = azurerm_mariadb_server.test.id
  restore_point_in_time     = "%s"
  ssl_enforcement_enabled   = true
  storage_mb                = 51200
}
`, testAccMariaDbServer_basic(data, version), data.RandomInteger, version, restoreTime)
}
