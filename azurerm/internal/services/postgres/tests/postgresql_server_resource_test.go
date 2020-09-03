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

func TestAccAzureRMPostgreSQLServer_basicNinePointFive(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_basic(data, "9.5"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_basicNinePointFiveDeprecated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_basicDeprecated(data, "9.5"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				Config: testAccAzureRMPostgreSQLServer_gp(data, "9.5"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
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
				Config: testAccAzureRMPostgreSQLServer_basic(data, "9.6"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
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
				Config: testAccAzureRMPostgreSQLServer_basic(data, "10.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_gpTenPointZero(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_gp(data, "10.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_moTenPointZero(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_mo(data, "10.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
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
				Config: testAccAzureRMPostgreSQLServer_basic(data, "11"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_basicWithIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_basicWithIdentity(data, "11"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_autogrowOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_autogrow(data, "11"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				Config: testAccAzureRMPostgreSQLServer_gp(data, "11"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_basic(data, "10.0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMPostgreSQLServer_requiresImport),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_complete(data, "9.6"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_updatedDeprecated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_basicDeprecated(data, "9.6"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				Config: testAccAzureRMPostgreSQLServer_completeDeprecated(data, "9.6"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				Config: testAccAzureRMPostgreSQLServer_basicDeprecated(data, "9.6"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
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
				Config: testAccAzureRMPostgreSQLServer_gp(data, "9.6"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				Config: testAccAzureRMPostgreSQLServer_complete(data, "9.6"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				Config: testAccAzureRMPostgreSQLServer_complete2(data, "9.6"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				Config: testAccAzureRMPostgreSQLServer_gp(data, "9.6"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_completeDeprecatedUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_completeDeprecated(data, "9.6"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				Config: testAccAzureRMPostgreSQLServer_complete(data, "9.6"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
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
				Config: testAccAzureRMPostgreSQLServer_sku(data, "10.0", "B_Gen5_2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				Config: testAccAzureRMPostgreSQLServer_sku(data, "10.0", "GP_Gen5_2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				Config: testAccAzureRMPostgreSQLServer_sku(data, "10.0", "MO_Gen5_16"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_createReplica(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_gp(data, "11"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				Config: testAccAzureRMPostgreSQLServer_createReplica(data, "11"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
					testCheckAzureRMPostgreSQLServerExists("azurerm_postgresql_server.replica"),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

func TestAccAzureRMPostgreSQLServer_createPointInTimeRestore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_server", "test")
	restoreTime := time.Now().Add(11 * time.Minute)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPostgreSQLServer_gp(data, "11"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
				),
			},
			data.ImportStep("administrator_login_password"),
			{
				PreConfig: func() { time.Sleep(restoreTime.Sub(time.Now().Add(-7 * time.Minute))) },
				Config:    testAccAzureRMPostgreSQLServer_createPointInTimeRestore(data, "11", restoreTime.Format(time.RFC3339)),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(data.ResourceName),
					testCheckAzureRMPostgreSQLServerExists("azurerm_postgresql_server.restore"),
				),
			},
			data.ImportStep("administrator_login_password"),
		},
	})
}

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

func testAccAzureRMPostgreSQLServer_template(data acceptance.TestData, sku, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name   = "%s"
  version    = "%s"
  storage_mb = 51200

  ssl_enforcement_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, sku, version)
}

func testAccAzureRMPostgreSQLServer_basic(data acceptance.TestData, version string) string {
	return testAccAzureRMPostgreSQLServer_template(data, "B_Gen5_1", version)
}

func testAccAzureRMPostgreSQLServer_basicWithIdentity(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name   = "B_Gen5_1"
  version    = "%s"
  storage_mb = 51200

  ssl_enforcement_enabled = true

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func testAccAzureRMPostgreSQLServer_mo(data acceptance.TestData, version string) string {
	return testAccAzureRMPostgreSQLServer_template(data, "MO_Gen5_2", version)
}

func testAccAzureRMPostgreSQLServer_gp(data acceptance.TestData, version string) string {
	return testAccAzureRMPostgreSQLServer_template(data, "GP_Gen5_2", version)
}

func testAccAzureRMPostgreSQLServer_autogrow(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name          = "GP_Gen5_2"
  version           = "%s"
  auto_grow_enabled = true

  ssl_enforcement_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func testAccAzureRMPostgreSQLServer_basicDeprecated(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name = "GP_Gen5_2"
  version  = "%s"

  storage_profile {
    storage_mb = 51200
  }

  ssl_enforcement_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func testAccAzureRMPostgreSQLServer_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMPostgreSQLServer_basic(data, "10.0")
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_server" "import" {
  name                = azurerm_postgresql_server.test.name
  location            = azurerm_postgresql_server.test.location
  resource_group_name = azurerm_postgresql_server.test.resource_group_name

  administrator_login          = azurerm_postgresql_server.test.administrator_login
  administrator_login_password = azurerm_postgresql_server.test.administrator_login_password

  sku_name   = azurerm_postgresql_server.test.sku_name
  version    = azurerm_postgresql_server.test.version
  storage_mb = azurerm_postgresql_server.test.storage_mb

  ssl_enforcement_enabled = azurerm_postgresql_server.test.ssl_enforcement_enabled
}
`, template)
}

func testAccAzureRMPostgreSQLServer_completeDeprecated(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  version  = "%s"
  sku_name = "GP_Gen5_2"

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"

  infrastructure_encryption_enabled = true
  public_network_access_enabled     = false
  ssl_minimal_tls_version_enforced  = "TLS1_2"

  ssl_enforcement_enabled = true

  storage_profile {
    storage_mb            = 640000
    backup_retention_days = 7
    geo_redundant_backup  = "Enabled"
    auto_grow             = "Enabled"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func testAccAzureRMPostgreSQLServer_complete(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!updated"

  sku_name   = "GP_Gen5_4"
  version    = "%[3]s"
  storage_mb = 640000

  backup_retention_days        = 7
  geo_redundant_backup_enabled = true
  auto_grow_enabled            = true

  infrastructure_encryption_enabled = true
  public_network_access_enabled     = false
  ssl_enforcement_enabled           = true
  ssl_minimal_tls_version_enforced  = "TLS1_2"

  threat_detection_policy {
    enabled              = true
    disabled_alerts      = ["Sql_Injection", "Data_Exfiltration"]
    email_account_admins = true
    email_addresses      = ["kt@example.com", "admin@example.com"]

    retention_days = 7
  }
}
`, data.RandomInteger, data.Locations.Primary, version)
}

func testAccAzureRMPostgreSQLServer_complete2(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!updated"

  sku_name   = "GP_Gen5_4"
  version    = "%[3]s"
  storage_mb = 640000

  backup_retention_days        = 14
  geo_redundant_backup_enabled = false
  auto_grow_enabled            = false

  infrastructure_encryption_enabled = false
  public_network_access_enabled     = true
  ssl_enforcement_enabled           = false
  ssl_minimal_tls_version_enforced  = "TLS1_1"

  threat_detection_policy {
    enabled              = true
    disabled_alerts      = ["Sql_Injection"]
    email_account_admins = true
    email_addresses      = ["kt@example.com"]

    retention_days = 7

    storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
    storage_account_access_key = azurerm_storage_account.test.primary_access_key
  }
}
`, data.RandomInteger, data.Locations.Primary, version)
}

func testAccAzureRMPostgreSQLServer_sku(data acceptance.TestData, version, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-psql-server-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name   = "%s"
  storage_mb = 51200
  version    = "%s"

  ssl_enforcement_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, sku, version)
}

func testAccAzureRMPostgreSQLServer_createReplica(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_postgresql_server" "replica" {
  name                = "acctest-psql-server-%[2]d-replica"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_2"
  version  = "%[3]s"

  create_mode               = "Replica"
  creation_source_server_id = azurerm_postgresql_server.test.id

  ssl_enforcement_enabled = true
}
`, testAccAzureRMPostgreSQLServer_basic(data, version), data.RandomInteger, version)
}

func testAccAzureRMPostgreSQLServer_createPointInTimeRestore(data acceptance.TestData, version, restoreTime string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_postgresql_server" "restore" {
  name                = "acctest-psql-server-%[2]d-restore"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name   = "GP_Gen5_2"
  version    = "%[4]s"
  storage_mb = 51200

  create_mode               = "PointInTimeRestore"
  creation_source_server_id = azurerm_postgresql_server.test.id
  restore_point_in_time     = "%[3]s"

  ssl_enforcement_enabled = true
}
`, testAccAzureRMPostgreSQLServer_basic(data, version), data.RandomInteger, restoreTime, version)
}
