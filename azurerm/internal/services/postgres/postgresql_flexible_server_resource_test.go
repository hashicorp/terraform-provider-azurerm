package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type PostgresqlFlexibleServerResource struct {
}

func TestAccPostgresqlflexibleServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("cmk_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlflexibleServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccPostgresqlflexibleServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cmk_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlflexibleServer_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cmk_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.completeUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cmk_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlflexibleServer_updateMaintenanceWindow(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("cmk_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.updateMaintenanceWindow(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("cmk_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.updateMaintenanceWindowUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("cmk_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("cmk_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlflexibleServer_updateSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("cmk_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.updateSku(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("cmk_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("cmk_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlflexibleServer_pitr(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("cmk_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			PreConfig: func() { time.Sleep(15 * time.Minute) },
			Config:    r.pitr(data),
			Check: resource.ComposeTestCheckFunc(
				check.That("azurerm_postgresql_flexible_server.pitr").ExistsInAzure(r),
				check.That("azurerm_postgresql_flexible_server.pitr").Key("zone").Exists(),
				check.That("azurerm_postgresql_flexible_server.pitr").Key("cmk_enabled").Exists(),
				check.That("azurerm_postgresql_flexible_server.pitr").Key("fqdn").Exists(),
				check.That("azurerm_postgresql_flexible_server.pitr").Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode", "point_in_time_restore_time_in_utc"),
	})
}

func (PostgresqlFlexibleServerResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.FlexibleServerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Postgres.FlexibleServersClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Postgresql Flexible Server %q (resource group: %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ServerProperties != nil), nil
}

func (PostgresqlFlexibleServerResource) template(data acceptance.TestData) string {
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

func (r PostgresqlFlexibleServerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  storage_mb             = 32768
  version                = "12"
  sku_name               = "GP_Standard_D2s_v3"
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "import" {
  name                   = azurerm_postgresql_flexible_server.test.name
  resource_group_name    = azurerm_postgresql_flexible_server.test.resource_group_name
  location               = azurerm_postgresql_flexible_server.test.location
  administrator_login    = azurerm_postgresql_flexible_server.test.administrator_login
  administrator_password = azurerm_postgresql_flexible_server.test.administrator_password
  version                = azurerm_postgresql_flexible_server.test.version
  storage_mb             = azurerm_postgresql_flexible_server.test.storage_mb
  sku_name               = azurerm_postgresql_flexible_server.test.sku_name
}
`, r.basic(data))
}

func (r PostgresqlFlexibleServerResource) complete(data acceptance.TestData) string {
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
  service_endpoints    = ["Microsoft.Storage"]
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
  name                   = "acctest-fs-%[2]d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  zone                   = "1"
  version                = "12"
  backup_retention_days  = 7
  storage_mb             = 32768
  delegated_subnet_id    = azurerm_subnet.test.id
  sku_name               = "GP_Standard_D2s_v3"

  maintenance_window {
    day_of_week  = 0
    start_hour   = 8
    start_minute = 0
  }

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) completeUpdate(data acceptance.TestData) string {
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
  service_endpoints    = ["Microsoft.Storage"]
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
  name                   = "acctest-fs-%[2]d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "123wsxQAZ"
  zone                   = "1"
  version                = "12"
  backup_retention_days  = 10
  storage_mb             = 65536
  delegated_subnet_id    = azurerm_subnet.test.id
  sku_name               = "GP_Standard_D2s_v3"

  maintenance_window {
    day_of_week  = 0
    start_hour   = 8
    start_minute = 0
  }

  tags = {
    ENV = "Stage"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) updateMaintenanceWindow(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  version                = "12"
  storage_mb             = 32768
  sku_name               = "GP_Standard_D2s_v3"

  maintenance_window {
    day_of_week  = 0
    start_hour   = 8
    start_minute = 0
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) updateMaintenanceWindowUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  version                = "12"
  storage_mb             = 32768
  sku_name               = "GP_Standard_D2s_v3"

  maintenance_window {
    day_of_week  = 3
    start_hour   = 7
    start_minute = 15
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) updateSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  version                = "12"
  storage_mb             = 32768
  sku_name               = "MO_Standard_E2s_v3"
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) pitr(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "pitr" {
  name                              = "acctest-fs-pitr-%d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = azurerm_resource_group.test.location
  create_mode                       = "PointInTimeRestore"
  source_server_id                  = azurerm_postgresql_flexible_server.test.id
  point_in_time_restore_time_in_utc = "%s"
}
`, r.basic(data), data.RandomInteger, time.Now().Add(time.Duration(15)*time.Minute).UTC().Format(time.RFC3339))
}
