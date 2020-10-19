package tests

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

func TestAccAzureRMMySQLServer_basicFiveSix(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServer_basic(data, "5.6"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccAzureRMMySQLServer_basicFiveSixWithIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServer_basicWithIdentity(data, "5.6"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccAzureRMMySQLServer_basicFiveSixWithIdentityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServer_basic(data, "5.6"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
			{
				Config: testAccAzureRMMySQLServer_basicWithIdentity(data, "5.6"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccAzureRMMySQLServer_basicFiveSixDeprecated(t *testing.T) { // remove in v3.0
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServer_basicDeprecated(data, "5.6"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccAzureRMMySQLServer_basicFiveSeven(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServer_basic(data, "5.7"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccAzureRMMySQLServer_basicEightZero(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServer_basic(data, "8.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccAzureRMMySQLServer_autogrowOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")
	mysqlVersion := "5.7"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServer_autogrow(data, mysqlVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
			{
				Config: testAccAzureRMMySQLServer_basic(data, mysqlVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccAzureRMMySQLServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServer_basic(data, "5.7"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMMySQLServer_requiresImport),
		},
	})
}

func TestAccAzureRMMySQLServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServer_complete(data, "8.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccAzureRMMySQLServer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")
	mysqlVersion := "8.0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServer_basic(data, mysqlVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
			{
				Config: testAccAzureRMMySQLServer_complete(data, mysqlVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
			{
				Config: testAccAzureRMMySQLServer_complete2(data, mysqlVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password", "threat_detection_policy.0.storage_account_access_key"),
			{
				Config: testAccAzureRMMySQLServer_complete3(data, mysqlVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password", "threat_detection_policy.0.storage_account_access_key"),
			{
				Config: testAccAzureRMMySQLServer_basic(data, mysqlVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccAzureRMMySQLServer_completeDeprecatedMigrate(t *testing.T) { // remove in v3.0
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")
	mysqlVersion := "5.6"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServer_completeDeprecated(data, mysqlVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
			{
				Config: testAccAzureRMMySQLServer_complete(data, mysqlVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccAzureRMMySQLServer_updateDeprecated(t *testing.T) { // remove in v3.0
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")
	mysqlVersion := "5.6"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServer_basicDeprecated(data, mysqlVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
			{
				Config: testAccAzureRMMySQLServer_completeDeprecated(data, mysqlVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
			{
				Config: testAccAzureRMMySQLServer_basicDeprecated(data, mysqlVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccAzureRMMySQLServer_updateSKU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServer_sku(data, "GP_Gen5_2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
			{
				Config: testAccAzureRMMySQLServer_sku(data, "MO_Gen5_16"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"), // not returned as sensitive
		},
	})
}

func TestAccAzureRMMySQLServer_createReplica(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")
	mysqlVersion := "8.0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServer_basic(data, mysqlVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				Config: testAccAzureRMMySQLServer_createReplica(data, mysqlVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
					testCheckAzureRMMySQLServerExists("azurerm_mysql_server.replica"),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMMySQLServer_createPointInTimeRestore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")
	restoreTime := time.Now().Add(11 * time.Minute)
	mysqlVersion := "8.0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMySQLServer_basic(data, mysqlVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				PreConfig: func() { time.Sleep(restoreTime.Sub(time.Now().Add(-7 * time.Minute))) },
				Config:    testAccAzureRMMySQLServer_createPointInTimeRestore(data, mysqlVersion, restoreTime.Format(time.RFC3339)),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
					testCheckAzureRMMySQLServerExists("azurerm_mysql_server.restore"),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func testCheckAzureRMMySQLServerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MySQL.ServersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMMySQLServer_basic(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                             = "acctestmysqlsvr-%d"
  location                         = azurerm_resource_group.test.location
  resource_group_name              = azurerm_resource_group.test.name
  sku_name                         = "GP_Gen5_2"
  administrator_login              = "acctestun"
  administrator_login_password     = "H@Sh1CoR3!"
  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "TLS1_1"
  storage_mb                       = 51200
  version                          = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func testAccAzureRMMySQLServer_basicDeprecated(data acceptance.TestData, version string) string { // remove in v3.0
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb = 51200
  }

  administrator_login              = "acctestun"
  administrator_login_password     = "H@Sh1CoR3!"
  version                          = "%s"
  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "TLS1_1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func testAccAzureRMMySQLServer_basicWithIdentity(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                             = "acctestmysqlsvr-%d"
  location                         = azurerm_resource_group.test.location
  resource_group_name              = azurerm_resource_group.test.name
  sku_name                         = "GP_Gen5_2"
  administrator_login              = "acctestun"
  administrator_login_password     = "H@Sh1CoR3!"
  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "TLS1_1"
  storage_mb                       = 51200
  version                          = "%s"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func testAccAzureRMMySQLServer_complete(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mysql_server" "test" {
  name                             = "acctestmysqlsvr-%[1]d"
  location                         = azurerm_resource_group.test.location
  resource_group_name              = azurerm_resource_group.test.name
  sku_name                         = "GP_Gen5_2"
  administrator_login              = "acctestun"
  administrator_login_password     = "H@Sh1CoR3!"
  auto_grow_enabled                = true
  backup_retention_days            = 7
  create_mode                      = "Default"
  geo_redundant_backup_enabled     = false
  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "TLS1_2"
  storage_mb                       = 51200
  version                          = "%[3]s"
  threat_detection_policy {
    enabled              = true
    disabled_alerts      = ["Sql_Injection", "Data_Exfiltration"]
    email_account_admins = true
    email_addresses      = ["pearcec@example.com", "admin@example.com"]
    retention_days       = 7
  }
}
`, data.RandomInteger, data.Locations.Primary, version)
}

func testAccAzureRMMySQLServer_complete2(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mysql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_mysql_server" "test" {
  name                         = "acctestmysqlsvr-%[1]d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  sku_name                     = "GP_Gen5_2"
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!updated"
  auto_grow_enabled            = true
  backup_retention_days        = 7
  create_mode                  = "Default"
  geo_redundant_backup_enabled = false
  ssl_enforcement_enabled      = false
  storage_mb                   = 51200
  version                      = "%[3]s"
  threat_detection_policy {
    enabled                    = true
    disabled_alerts            = ["Sql_Injection"]
    email_account_admins       = true
    email_addresses            = ["pearcec@example.com"]
    retention_days             = 7
    storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
    storage_account_access_key = azurerm_storage_account.test.primary_access_key
  }
}
`, data.RandomInteger, data.Locations.Primary, version)
}

func testAccAzureRMMySQLServer_complete3(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mysql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_mysql_server" "test" {
  name                         = "acctestmysqlsvr-%[1]d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  sku_name                     = "GP_Gen5_2"
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!updated"
  auto_grow_enabled            = true
  backup_retention_days        = 7
  create_mode                  = "Default"
  geo_redundant_backup_enabled = false
  ssl_enforcement_enabled      = false
  storage_mb                   = 51200
  version                      = "%[3]s"
  threat_detection_policy {
    enabled                    = true
    email_account_admins       = true
    retention_days             = 7
    storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
    storage_account_access_key = azurerm_storage_account.test.primary_access_key
  }
}
`, data.RandomInteger, data.Locations.Primary, version)
}

func testAccAzureRMMySQLServer_completeDeprecated(data acceptance.TestData, version string) string { // remove in v3.0
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
    auto_grow             = "Enabled"
  }

  administrator_login              = "acctestun"
  administrator_login_password     = "H@Sh1CoR3!"
  version                          = "%s"
  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "TLS1_2"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func testAccAzureRMMySQLServer_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_server" "import" {
  name                = azurerm_mysql_server.test.name
  location            = azurerm_mysql_server.test.location
  resource_group_name = azurerm_mysql_server.test.resource_group_name
  sku_name            = "GP_Gen5_2"
  version             = "5.7"

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  ssl_enforcement_enabled      = true
  storage_mb                   = 51200
}
`, testAccAzureRMMySQLServer_basic(data, "5.7"))
}

func testAccAzureRMMySQLServer_sku(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "%s"
  version             = "5.7"

  storage_mb                   = 4194304
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  ssl_enforcement_enabled      = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, sku)
}

func testAccAzureRMMySQLServer_autogrow(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "GP_Gen5_2"
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

func testAccAzureRMMySQLServer_createReplica(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_server" "replica" {
  name                = "acctestmysqlsvr-%d-replica"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "GP_Gen5_2"
  version             = "%s"
  storage_mb          = 51200

  create_mode                      = "Replica"
  creation_source_server_id        = azurerm_mysql_server.test.id
  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "TLS1_1"
}
`, testAccAzureRMMySQLServer_basic(data, version), data.RandomInteger, version)
}

func testAccAzureRMMySQLServer_createPointInTimeRestore(data acceptance.TestData, version, restoreTime string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_server" "restore" {
  name                = "acctestmysqlsvr-%d-restore"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "GP_Gen5_2"
  version             = "%s"

  create_mode                      = "PointInTimeRestore"
  creation_source_server_id        = azurerm_mysql_server.test.id
  restore_point_in_time            = "%s"
  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "TLS1_1"
  storage_mb                       = 51200
}
`, testAccAzureRMMySQLServer_basic(data, version), data.RandomInteger, version, restoreTime)
}
