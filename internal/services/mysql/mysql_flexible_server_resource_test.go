package mysql_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MySqlFlexibleServerResource struct {
}

func TestAccMySqlFlexibleServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("storage.#").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccMySqlFlexibleServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMySqlFlexibleServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccMySqlFlexibleServer_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccMySqlFlexibleServer_updateMaintenanceWindow(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.updateMaintenanceWindow(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.updateMaintenanceWindowUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccMySqlFlexibleServer_updateSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.updateSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccMySqlFlexibleServer_updateHA(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.updateHAZoneRedundant(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password"),

		{
			Config: r.updateHAZoneRedundantUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.updateHASameZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.updateHADisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccMySqlFlexibleServer_pitr(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
		{
			PreConfig: func() { time.Sleep(15 * time.Minute) },
			Config:    r.pitr(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_mysql_flexible_server.pitr").ExistsInAzure(r),
				check.That("azurerm_mysql_flexible_server.pitr").Key("zone").Exists(),
				check.That("azurerm_mysql_flexible_server.pitr").Key("fqdn").Exists(),
				check.That("azurerm_mysql_flexible_server.pitr").Key("public_network_access_enabled").Exists(),
				check.That("azurerm_mysql_flexible_server.pitr").Key("replica_capacity").Exists(),
				check.That("azurerm_mysql_flexible_server.pitr").Key("administrator_login").Exists(),
				check.That("azurerm_mysql_flexible_server.pitr").Key("sku_name").Exists(),
				check.That("azurerm_mysql_flexible_server.pitr").Key("version").Exists(),
				check.That("azurerm_mysql_flexible_server.pitr").Key("storage.#").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccMySqlFlexibleServer_replica(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.source(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
		{
			PreConfig: func() { time.Sleep(15 * time.Minute) },
			Config:    r.replica(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_mysql_flexible_server.replica").ExistsInAzure(r),
				check.That("azurerm_mysql_flexible_server.replica").Key("zone").Exists(),
				check.That("azurerm_mysql_flexible_server.replica").Key("fqdn").Exists(),
				check.That("azurerm_mysql_flexible_server.replica").Key("public_network_access_enabled").Exists(),
				check.That("azurerm_mysql_flexible_server.replica").Key("replica_capacity").Exists(),
				check.That("azurerm_mysql_flexible_server.replica").Key("replication_role").HasValue("Replica"),
				check.That("azurerm_mysql_flexible_server.replica").Key("administrator_login").Exists(),
				check.That("azurerm_mysql_flexible_server.replica").Key("sku_name").Exists(),
				check.That("azurerm_mysql_flexible_server.replica").Key("version").Exists(),
				check.That("azurerm_mysql_flexible_server.replica").Key("storage.#").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccMySqlFlexibleServer_geoRestore(t *testing.T) {
	if os.Getenv("ARM_GEO_RESTORE_LOCATION") == "" {
		t.Skip("Skipping as `ARM_GEO_RESTORE_LOCATION` is not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.geoRestoreSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
		{
			PreConfig: func() { time.Sleep(15 * time.Minute) },
			Config:    r.geoRestore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_mysql_flexible_server.geo_restore").ExistsInAzure(r),
				check.That("azurerm_mysql_flexible_server.geo_restore").Key("zone").Exists(),
				check.That("azurerm_mysql_flexible_server.geo_restore").Key("fqdn").Exists(),
				check.That("azurerm_mysql_flexible_server.geo_restore").Key("public_network_access_enabled").Exists(),
				check.That("azurerm_mysql_flexible_server.geo_restore").Key("replica_capacity").Exists(),
				check.That("azurerm_mysql_flexible_server.geo_restore").Key("administrator_login").Exists(),
				check.That("azurerm_mysql_flexible_server.geo_restore").Key("sku_name").Exists(),
				check.That("azurerm_mysql_flexible_server.geo_restore").Key("version").Exists(),
				check.That("azurerm_mysql_flexible_server.geo_restore").Key("storage.#").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccMySqlFlexibleServer_updateStorage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.updateStorage(data, 20, 360, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.updateStorage(data, 34, 402, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccMySqlFlexibleServer_updateReplicationRole(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.source(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password"),
		{
			PreConfig: func() { time.Sleep(15 * time.Minute) },
			Config:    r.replica(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_mysql_flexible_server.replica").ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.updateReplicationRole(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_mysql_flexible_server.replica").ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccMySqlFlexibleServer_failover(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.failover(data, "1", "2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.failover(data, "2", "1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password"),
	})
}

func (MySqlFlexibleServerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FlexibleServerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MySQL.FlexibleServerClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving MySql Flexible Server %q (resource group: %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ServerProperties != nil), nil
}

func (MySqlFlexibleServerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mysql-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r MySqlFlexibleServerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  sku_name               = "B_Standard_B1s"
}
`, r.template(data), data.RandomInteger)
}

func (r MySqlFlexibleServerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "import" {
  name                   = azurerm_mysql_flexible_server.test.name
  resource_group_name    = azurerm_mysql_flexible_server.test.resource_group_name
  location               = azurerm_mysql_flexible_server.test.location
  administrator_login    = azurerm_mysql_flexible_server.test.administrator_login
  administrator_password = azurerm_mysql_flexible_server.test.administrator_password
  sku_name               = azurerm_mysql_flexible_server.test.sku_name
}
`, r.basic(data))
}

func (r MySqlFlexibleServerResource) complete(data acceptance.TestData) string {
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
      name = "Microsoft.DBforMySQL/flexibleServers"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acc%[2]d.mysql.database.azure.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                  = "acctestVnetZone%[2]d.com"
  private_dns_zone_name = azurerm_private_dns_zone.test.name
  virtual_network_id    = azurerm_virtual_network.test.id
  resource_group_name   = azurerm_resource_group.test.name
}

resource "azurerm_mysql_flexible_server" "test" {
  name                         = "acctest-fs-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_password       = "QAZwsx123"
  zone                         = "1"
  version                      = "8.0.21"
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false

  storage {
    size_gb = 20
    iops    = 360
  }

  delegated_subnet_id = azurerm_subnet.test.id
  private_dns_zone_id = azurerm_private_dns_zone.test.id
  sku_name            = "GP_Standard_D2ds_v4"

  high_availability {
    mode                      = "ZoneRedundant"
    standby_availability_zone = "2"
  }

  maintenance_window {
    day_of_week  = 0
    start_hour   = 8
    start_minute = 0
  }

  tags = {
    ENV = "Test"
  }

  depends_on = [azurerm_private_dns_zone_virtual_network_link.test]
}
`, r.template(data), data.RandomInteger)
}

func (r MySqlFlexibleServerResource) completeUpdate(data acceptance.TestData) string {
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
      name = "Microsoft.DBforMySQL/flexibleServers"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acc%[2]d.mysql.database.azure.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                  = "acctestVnetZone%[2]d.com"
  private_dns_zone_name = azurerm_private_dns_zone.test.name
  virtual_network_id    = azurerm_virtual_network.test.id
  resource_group_name   = azurerm_resource_group.test.name
}

resource "azurerm_mysql_flexible_server" "test" {
  name                         = "acctest-fs-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_password       = "123wsxQAZ"
  zone                         = "1"
  version                      = "8.0.21"
  backup_retention_days        = 10
  geo_redundant_backup_enabled = false

  storage {
    size_gb           = 32
    iops              = 400
    auto_grow_enabled = false
  }

  delegated_subnet_id = azurerm_subnet.test.id
  private_dns_zone_id = azurerm_private_dns_zone.test.id
  sku_name            = "GP_Standard_D4ds_v4"

  maintenance_window {
    day_of_week  = 0
    start_hour   = 8
    start_minute = 0
  }

  tags = {
    ENV = "Stage"
  }

  depends_on = [azurerm_private_dns_zone_virtual_network_link.test]
}
`, r.template(data), data.RandomInteger)
}

func (r MySqlFlexibleServerResource) updateMaintenanceWindow(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  sku_name               = "B_Standard_B1s"

  maintenance_window {
    day_of_week  = 0
    start_hour   = 8
    start_minute = 0
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MySqlFlexibleServerResource) updateMaintenanceWindowUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  sku_name               = "B_Standard_B1s"

  maintenance_window {
    day_of_week  = 3
    start_hour   = 7
    start_minute = 15
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MySqlFlexibleServerResource) updateSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  sku_name               = "MO_Standard_E2ds_v4"
}
`, r.template(data), data.RandomInteger)
}

func (r MySqlFlexibleServerResource) updateHADisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  sku_name               = "GP_Standard_D2ds_v4"
}
`, r.template(data), data.RandomInteger)
}

func (r MySqlFlexibleServerResource) updateHASameZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"

  high_availability {
    mode = "SameZone"
  }

  sku_name = "GP_Standard_D2ds_v4"
  zone     = "1"
}
`, r.template(data), data.RandomInteger)
}

func (r MySqlFlexibleServerResource) updateHAZoneRedundant(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"

  high_availability {
    mode                      = "ZoneRedundant"
    standby_availability_zone = "2"
  }

  sku_name = "GP_Standard_D2ds_v4"
  zone     = "1"
}
`, r.template(data), data.RandomInteger)
}

func (r MySqlFlexibleServerResource) updateHAZoneRedundantUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"

  high_availability {
    mode                      = "ZoneRedundant"
    standby_availability_zone = "3"
  }

  sku_name = "GP_Standard_D2ds_v4"
  zone     = "1"
}
`, r.template(data), data.RandomInteger)
}

func (r MySqlFlexibleServerResource) pitr(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "pitr" {
  name                              = "acctest-fs-pitr-%d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = azurerm_resource_group.test.location
  create_mode                       = "PointInTimeRestore"
  source_server_id                  = azurerm_mysql_flexible_server.test.id
  point_in_time_restore_time_in_utc = "%s"
}
`, r.basic(data), data.RandomInteger, time.Now().Add(time.Duration(15)*time.Minute).UTC().Format(time.RFC3339))
}

func (r MySqlFlexibleServerResource) source(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  sku_name               = "GP_Standard_D4ds_v4"
}
`, r.template(data), data.RandomInteger)
}

func (r MySqlFlexibleServerResource) replica(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "replica" {
  name                = "acctest-fs-replica-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  create_mode         = "Replica"
  source_server_id    = azurerm_mysql_flexible_server.test.id
}
`, r.source(data), data.RandomInteger)
}

func (r MySqlFlexibleServerResource) updateReplicationRole(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "replica" {
  name                = "acctest-fs-replica-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  create_mode         = "Replica"
  source_server_id    = azurerm_mysql_flexible_server.test.id
  replication_role    = "None"
}
`, r.source(data), data.RandomInteger)
}

func (r MySqlFlexibleServerResource) geoRestoreSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                         = "acctest-fs-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_password       = "QAZwsx123"
  geo_redundant_backup_enabled = true
  sku_name                     = "B_Standard_B1s"
}
`, r.template(data), data.RandomInteger)
}

func (r MySqlFlexibleServerResource) geoRestore(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "geo_restore" {
  name                = "acctest-fs-restore-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  create_mode         = "GeoRestore"
  source_server_id    = azurerm_mysql_flexible_server.test.id
}
`, r.geoRestoreSource(data), data.RandomInteger, os.Getenv("ARM_GEO_RESTORE_LOCATION"))
}

func (r MySqlFlexibleServerResource) updateStorage(data acceptance.TestData, sizeGB int, iops int, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                         = "acctest-fs-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_password       = "QAZwsx123"
  sku_name                     = "GP_Standard_D4ds_v4"
  geo_redundant_backup_enabled = true

  storage {
    size_gb           = %d
    iops              = %d
    auto_grow_enabled = %t
  }
}
`, r.template(data), data.RandomInteger, sizeGB, iops, enabled)
}

func (r MySqlFlexibleServerResource) failover(data acceptance.TestData, primaryZone string, standbyZone string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  sku_name               = "GP_Standard_D2ds_v4"
  zone                   = "%s"

  high_availability {
    mode                      = "ZoneRedundant"
    standby_availability_zone = "%s"
  }
}
`, r.template(data), data.RandomInteger, primaryZone, standbyZone)
}
