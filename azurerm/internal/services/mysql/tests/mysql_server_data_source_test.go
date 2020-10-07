package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMMySQLServer_basicFiveSix(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mysql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMySQLServer_basic(data, "5.6"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "GP_Gen5_2"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_grow_enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_minimal_tls_version_enforced", "TLS1_1"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_mb", "51200"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "5.6"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMMySQLServer_basicFiveSixWithIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mysql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMySQLServer_basicWithIdentity(data, "5.6"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "GP_Gen5_2"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_grow_enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_minimal_tls_version_enforced", "TLS1_1"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_mb", "51200"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "5.6"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.tenant_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMMySQLServer_basicFiveSeven(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mysql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMySQLServer_basic(data, "5.7"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "GP_Gen5_2"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_grow_enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_minimal_tls_version_enforced", "TLS1_1"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_mb", "51200"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "5.7"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMMySQLServer_basicEightZero(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mysql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMySQLServer_basic(data, "8.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "GP_Gen5_2"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_grow_enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_minimal_tls_version_enforced", "TLS1_1"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_mb", "51200"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "8.0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMMySQLServer_autogrowOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mysql_server", "test")
	mysqlVersion := "5.7"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMySQLServer_autogrow(data, mysqlVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "GP_Gen5_2"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_grow_enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_minimal_tls_version_enforced", "TLS1_1"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_mb", "51200"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "5.7"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMMySQLServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mysql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMySQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMySQLServer_complete(data, "8.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMySQLServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "GP_Gen5_2"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_grow_enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_minimal_tls_version_enforced", "TLS1_1"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_mb", "51200"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "8.0"),
					resource.TestCheckResourceAttr(data.ResourceName, "threat_detection_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "threat_detection_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "threat_detection_policy.0.email_account_admins", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "threat_detection_policy.0.retention_days", "7"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMySQLServer_basic(data acceptance.TestData, version string) string {
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

data "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version, data.RandomInteger)
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

data "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version, data.RandomInteger)
}

func testAccDataSourceAzureRMMySQLServer_autogrow(data acceptance.TestData, version string) string {
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

data "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version, data.RandomInteger)
}

func testAccDataSourceAzureRMMySQLServer_complete(data acceptance.TestData, version string) string {
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

data "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, version, data.RandomInteger)
}
