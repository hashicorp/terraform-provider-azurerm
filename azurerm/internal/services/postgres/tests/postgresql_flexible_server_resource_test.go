package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMpostgresqlflexibleServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMpostgresqlflexibleServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMpostgresqlflexibleServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "byok_enforcement"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ha_state"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_network_access"),
				),
			},
			data.ImportStep("administrator_login_password", "create_mode"),
		},
	})
}

func TestAccAzureRMpostgresqlflexibleServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMpostgresqlflexibleServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMpostgresqlflexibleServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMpostgresqlflexibleServer_requiresImport),
		},
	})
}

func TestAccAzureRMpostgresqlflexibleServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMpostgresqlflexibleServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMpostgresqlflexibleServer_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.tenant_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "byok_enforcement"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ha_state"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_network_access"),
				),
			},
			data.ImportStep("administrator_login_password", "create_mode"),
		},
	})
}

func TestAccAzureRMpostgresqlflexibleServer_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMpostgresqlflexibleServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMpostgresqlflexibleServer_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.tenant_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "byok_enforcement"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ha_state"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_network_access"),
				),
			},
			data.ImportStep("administrator_login_password", "create_mode"),
			{
				Config: testAccAzureRMpostgresqlflexibleServer_completeUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.tenant_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "byok_enforcement"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ha_state"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_network_access"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "standby_availability_zone"),
				),
			},
			data.ImportStep("administrator_login_password", "create_mode"),
		},
	})
}

func TestAccAzureRMpostgresqlflexibleServer_updateMaintenanceWindow(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMpostgresqlflexibleServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMpostgresqlflexibleServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "byok_enforcement"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ha_state"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_network_access"),
				),
			},
			data.ImportStep("administrator_login_password", "create_mode"),
			{
				Config: testAccAzureRMpostgresqlflexibleServer_updateMaintenanceWindow(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "byok_enforcement"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ha_state"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_network_access"),
				),
			},
			data.ImportStep("administrator_login_password", "create_mode"),
			{
				Config: testAccAzureRMpostgresqlflexibleServer_updateMaintenanceWindowUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "byok_enforcement"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ha_state"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_network_access"),
				),
			},
			data.ImportStep("administrator_login_password", "create_mode"),
			{
				Config: testAccAzureRMpostgresqlflexibleServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "byok_enforcement"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ha_state"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_network_access"),
				),
			},
			data.ImportStep("administrator_login_password", "create_mode"),
		},
	})
}

func TestAccAzureRMpostgresqlflexibleServer_updateSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMpostgresqlflexibleServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMpostgresqlflexibleServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "byok_enforcement"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ha_state"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_network_access"),
				),
			},
			data.ImportStep("administrator_login_password", "create_mode"),
			{
				Config: testAccAzureRMpostgresqlflexibleServer_updateSku(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "byok_enforcement"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ha_state"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_network_access"),
				),
			},
			data.ImportStep("administrator_login_password", "create_mode"),
			{
				Config: testAccAzureRMpostgresqlflexibleServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "byok_enforcement"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ha_state"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_network_access"),
				),
			},
			data.ImportStep("administrator_login_password", "create_mode"),
		},
	})
}

func TestAccAzureRMpostgresqlflexibleServer_pitr(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMpostgresqlflexibleServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMpostgresqlflexibleServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "byok_enforcement"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ha_state"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_network_access"),
				),
			},
			data.ImportStep("administrator_login_password", "create_mode"),
			{
				PreConfig: func() { time.Sleep(10 * time.Minute) },
				Config:    testAccAzureRMpostgresqlflexibleServer_pitr(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "byok_enforcement"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ha_state"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_network_access"),
				),
			},
			data.ImportStep("administrator_login_password", "create_mode"),
		},
	})
}

func testCheckAzureRMpostgresqlflexibleServerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Postgres.FlexibleServersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("postgresqlflexibleservers Server not found: %s", resourceName)
		}
		id, err := parse.FlexibleServerID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Postgresqlflexibleservers Server %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Postgresqlflexibleservers.ServerClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMpostgresqlflexibleServerDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Postgres.FlexibleServersClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_postgresql_flexible_server" {
			continue
		}
		id, err := parse.FlexibleServerID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Postgresqlflexibleservers.ServerClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMpostgresqlflexibleServer_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-postgresql-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMpostgresqlflexibleServer_basic(data acceptance.TestData) string {
	template := testAccAzureRMpostgresqlflexibleServer_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                         = "acctest-fs-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_login_password = "QAZwsx123"
  storage_mb                   = 32768
  version                      = "12"

  sku {
    name = "Standard_D2s_v3"
    tier = "GeneralPurpose"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMpostgresqlflexibleServer_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMpostgresqlflexibleServer_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "import" {
  name                         = azurerm_postgresql_flexible_server.test.name
  resource_group_name          = azurerm_postgresql_flexible_server.test.resource_group_name
  location                     = azurerm_postgresql_flexible_server.test.location
  administrator_login          = azurerm_postgresql_flexible_server.test.administrator_login
  administrator_login_password = azurerm_postgresql_flexible_server.test.administrator_login_password
  version                      = azurerm_postgresql_flexible_server.test.version
  storage_mb                   = azurerm_postgresql_flexible_server.test.storage_mb
  sku {
    name = azurerm_postgresql_flexible_server.test.sku.0.name
    tier = azurerm_postgresql_flexible_server.test.sku.0.tier
  }
}
`, config)
}

func testAccAzureRMpostgresqlflexibleServer_complete(data acceptance.TestData) string {
	template := testAccAzureRMpostgresqlflexibleServer_template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vn-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-sn-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  delegation {
    name = "fs"
    service_delegation {
      name = "Microsoft.DBforPostgreSQL/flexibleServers"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_postgresql_flexible_server" "test" {
  name                         = "acctest-fs-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_login_password = "QAZwsx123"
  availability_zone            = "1"
  version                      = "12"
  ha_enabled                   = false
  backup_retention_days        = 7
  storage_mb                   = 32768
  delegated_subnet_id          = azurerm_subnet.test.id

  identity {
    type = "SystemAssigned"
  }

  sku {
    name = "Standard_D2s_v3"
    tier = "GeneralPurpose"
  }

  maintenance_window {
    day_of_week  = 0
    start_hour   = 8
    start_minute = 0
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMpostgresqlflexibleServer_completeUpdate(data acceptance.TestData) string {
	template := testAccAzureRMpostgresqlflexibleServer_template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vn-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-sn-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  delegation {
    name = "fs"
    service_delegation {
      name = "Microsoft.DBforPostgreSQL/flexibleServers"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_postgresql_flexible_server" "test" {
  name                         = "acctest-fs-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_login_password = "123wsxQAZ"
  availability_zone            = "1"
  version                      = "12"
  ha_enabled                   = true
  backup_retention_days        = 10
  storage_mb                   = 65536
  delegated_subnet_id          = azurerm_subnet.test.id

  identity {
    type = "SystemAssigned"
  }

  sku {
    name = "Standard_D2s_v3"
    tier = "GeneralPurpose"
  }

  maintenance_window {
    day_of_week  = 0
    start_hour   = 8
    start_minute = 0
  }

  tags = {
    ENV = "Stage"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMpostgresqlflexibleServer_updateMaintenanceWindow(data acceptance.TestData) string {
	template := testAccAzureRMpostgresqlflexibleServer_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                         = "acctest-fs-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_login_password = "QAZwsx123"
  version                      = "12"
  storage_mb                   = 32768
  sku {
    name = "Standard_D2s_v3"
    tier = "GeneralPurpose"
  }

  maintenance_window {
    day_of_week  = 0
    start_hour   = 8
    start_minute = 0
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMpostgresqlflexibleServer_updateMaintenanceWindowUpdated(data acceptance.TestData) string {
	template := testAccAzureRMpostgresqlflexibleServer_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                         = "acctest-fs-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_login_password = "QAZwsx123"
  version                      = "12"
  storage_mb                   = 32768
  sku {
    name = "Standard_D2s_v3"
    tier = "GeneralPurpose"
  }

  maintenance_window {
    day_of_week  = 3
    start_hour   = 7
    start_minute = 15
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMpostgresqlflexibleServer_updateSku(data acceptance.TestData) string {
	template := testAccAzureRMpostgresqlflexibleServer_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                         = "acctest-fs-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_login_password = "QAZwsx123"
  version                      = "12"
  storage_mb                   = 32768
  sku {
    name = "Standard_E2s_v3"
    tier = "MemoryOptimized"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMpostgresqlflexibleServer_pitr(data acceptance.TestData) string {
	template := testAccAzureRMpostgresqlflexibleServer_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "pitr" {
  name                = "acctest-fs-pitr-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  create_mode         = "PointInTimeRestore"
  source_server_name  = azurerm_postgresql_flexible_server.test.name
  point_in_time_utc   = "%s"
}
`, template, data.RandomInteger, time.Now().Add(time.Duration(20)*time.Minute).UTC().Format(time.RFC3339))
}
