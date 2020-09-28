package tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"testing"
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
				),
			},
			data.ImportStep(),
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
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMpostgresqlflexibleServer_update(t *testing.T) {
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
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMpostgresqlflexibleServer_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMpostgresqlflexibleServer_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
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
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMpostgresqlflexibleServer_updateMaintenanceWindow(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMpostgresqlflexibleServer_updateMaintenanceWindowUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMpostgresqlflexibleServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
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
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMpostgresqlflexibleServer_updateSku(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMpostgresqlflexibleServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMpostgresqlflexibleServer_updateStorageProfile(t *testing.T) {
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
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMpostgresqlflexibleServer_updateStorageProfile(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMpostgresqlflexibleServer_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
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
  name     = "acctest-postgresql-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-s-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMpostgresqlflexibleServer_basic(data acceptance.TestData) string {
	template := testAccAzureRMpostgresqlflexibleServer_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                = "acctest-FS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
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
  name                = azurerm_postgresql_flexible_server.test.name
  resource_group_name = azurerm_postgresql_flexible_server.test.resource_group_name
  location            = azurerm_postgresql_flexible_server.test.location
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
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                         = "acctest-FS-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_login_password = "QAZwsx123"
  availability_zone            = "1"
  display_name                 = "fsTerraform"
  version                      = "12"
  ha_enabled                   = false

  delegated_subnet_argument {
    subnet_arm_resource_id = azurerm_subnet.test.id
  }

  identity {
    type = "SystemAssigned"
  }

  maintenance_window {
    day_of_week  = 0
    start_hour   = 8
    start_minute = 0
  }

  sku {
    name = "Standard_D2s_v3"
    tier = "GeneralPurpose"
  }

  storage_profile {
    backup_retention_days = 7
    storage_mb            = 32768
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMpostgresqlflexibleServer_update(data acceptance.TestData) string {
	template := testAccAzureRMpostgresqlflexibleServer_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                         = "acctest-FS-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_login_password = "123wsxQAZ"
  availability_zone            = "1"
  display_name                 = "fsTerraform"
  version                      = "12"
  ha_enabled                   = true

  delegated_subnet_argument {
    subnet_arm_resource_id = azurerm_subnet.test.id
  }

  identity {
    type = "SystemAssigned"
  }

  maintenance_window {
    day_of_week  = 0
    start_hour   = 8
    start_minute = 0
  }

  sku {
    name = "Standard_D2s_v3"
    tier = "GeneralPurpose"
  }

  storage_profile {
    backup_retention_days = 7
    storage_mb            = 32768
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
  name                = "acctest-FS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
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
  name                = "acctest-FS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
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
  name                = "acctest-FS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku {
    name = "Standard_E2s_v3"
    tier = "MemoryOptimized"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMpostgresqlflexibleServer_updateStorageProfile(data acceptance.TestData) string {
	template := testAccAzureRMpostgresqlflexibleServer_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                         = "acctest-FS-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_login_password = "QAZwsx123"
  availability_zone            = "1"
  display_name                 = "fsTerraform"
  version                      = "12"
  ha_enabled                   = false

  delegated_subnet_argument {
    subnet_arm_resource_id = azurerm_subnet.test.id
  }

  identity {
    type = "SystemAssigned"
  }

  maintenance_window {
    day_of_week  = 0
    start_hour   = 8
    start_minute = 0
  }

  sku {
    name = "Standard_D2s_v3"
    tier = "GeneralPurpose"
  }

  storage_profile {
    backup_retention_days = 3
    storage_mb            = 65536
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}
