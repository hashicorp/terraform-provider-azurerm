package mysql_test

import (
	"context"
	"fmt"
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

func TestAccMySqlFlexibleServer_basic_5_7(t *testing.T) {
	testAccMySqlFlexibleServer_basic(t, "5.7")
}

func TestAccMySqlFlexibleServer_basic_8_0_21(t *testing.T) {
	testAccMySqlFlexibleServer_basic(t, "8.0.21")
}

func testAccMySqlFlexibleServer_basic(t *testing.T, version string) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, version),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccMySqlFlexibleServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "5.7"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, "5.7"),
			ExpectError: acceptance.RequiresImportError("azurerm_mysql_flexible_server"),
		},
	})
}

func TestAccMySqlFlexibleServer_complete_5_7(t *testing.T) {
	testAccMySqlFlexibleServer_complete(t, "5.7")
}

func TestAccMySqlFlexibleServer_complete_8_0_21(t *testing.T) {
	testAccMySqlFlexibleServer_complete(t, "8.0.21")
}

func testAccMySqlFlexibleServer_complete(t *testing.T, version string) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, version),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccMySqlFlexibleServer_completeUpdate_5_7(t *testing.T) {
	testAccMySqlFlexibleServer_completeUpdate(t, "5.7")
}

func TestAccMySqlFlexibleServer_completeUpdate_8_0_21(t *testing.T) {
	testAccMySqlFlexibleServer_completeUpdate(t, "8.0.21")
}

func testAccMySqlFlexibleServer_completeUpdate(t *testing.T, version string) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, version),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.completeUpdate(data, version),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccMySqlFlexibleServer_updateMaintenanceWindow_5_7(t *testing.T) {
	testAccMySqlFlexibleServer_updateMaintenanceWindow(t, "5.7")
}

func TestAccMySqlFlexibleServer_updateMaintenanceWindow_8_0_21(t *testing.T) {
	testAccMySqlFlexibleServer_updateMaintenanceWindow(t, "8.0.21")
}

func testAccMySqlFlexibleServer_updateMaintenanceWindow(t *testing.T, version string) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, version),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.updateMaintenanceWindow(data, version),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.updateMaintenanceWindowUpdated(data, version),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.basic(data, version),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccMySqlFlexibleServer_updateSku_5_7(t *testing.T) {
	testAccMySqlFlexibleServer_updateSku(t, "5.7")
}

func TestAccMySqlFlexibleServer_updateSku_8_0_21(t *testing.T) {
	testAccMySqlFlexibleServer_updateSku(t, "8.0.21")
}

func testAccMySqlFlexibleServer_updateSku(t *testing.T, version string) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, version),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.updateSku(data, version),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.basic(data, version),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccMySqlFlexibleServer_pitr_5_7(t *testing.T) {
	testAccMySqlFlexibleServer_pitr(t, "5.7")
}

func TestAccMySqlFlexibleServer_pitr_8_0_21(t *testing.T) {
	testAccMySqlFlexibleServer_pitr(t, "8.0.21")
}

func testAccMySqlFlexibleServer_pitr(t *testing.T, version string) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, version),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			PreConfig: func() { time.Sleep(15 * time.Minute) },
			Config:    r.pitr(data, version),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_mysql_flexible_server.pitr").ExistsInAzure(r),
				check.That("azurerm_mysql_flexible_server.pitr").Key("zone").Exists(),
				check.That("azurerm_mysql_flexible_server.pitr").Key("fqdn").Exists(),
				check.That("azurerm_mysql_flexible_server.pitr").Key("public_network_access_enabled").Exists(),
				check.That("azurerm_mysql_flexible_server.pitr").Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode", "point_in_time_restore_time_in_utc"),
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

func (r MySqlFlexibleServerResource) basic(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  storage {
    size_gb = 32
  }
  version  = "%s"
  sku_name = "B_Standard_B1s"
}
`, r.template(data), data.RandomInteger, version)
}

func (r MySqlFlexibleServerResource) requiresImport(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "import" {
  name                   = azurerm_mysql_flexible_server.test.name
  resource_group_name    = azurerm_mysql_flexible_server.test.resource_group_name
  location               = azurerm_mysql_flexible_server.test.location
  administrator_login    = azurerm_mysql_flexible_server.test.administrator_login
  administrator_password = azurerm_mysql_flexible_server.test.administrator_password
  version                = azurerm_mysql_flexible_server.test.version
  storage {
    size_gb = azurerm_mysql_flexible_server.test.storage.0.size_gb
  }
  sku_name = azurerm_mysql_flexible_server.test.sku_name
}
`, r.basic(data, version))
}

func (r MySqlFlexibleServerResource) complete(data acceptance.TestData, version string) string {
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
  name                   = "acctest-fs-%[2]d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  zone                   = "1"
  version                = "%[3]s"
  backup_retention_days  = 7
  storage {
    size_gb = 32
    iops    = 400
auto_grow_enabled = true
  }
  delegated_subnet_id = azurerm_subnet.test.id
  private_dns_zone_id = azurerm_private_dns_zone.test.id
  sku_name            = "GP_Standard_D2ds_v4"

  high_availability {
    mode                      = "ZoneRedundant"
    standby_availability_zone = "1"
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
`, r.template(data), data.RandomInteger, version)
}

func (r MySqlFlexibleServerResource) completeUpdate(data acceptance.TestData, version string) string {
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
  version                      = "%[3]s"
  backup_retention_days        = 10
  geo_redundant_backup_enabled = true
  storage {
    size_gb           = 64
    iops              = 1280
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
`, r.template(data), data.RandomInteger, version)
}

func (r MySqlFlexibleServerResource) updateMaintenanceWindow(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  storage {
    size_gb = 32
  }
  version  = "%s"
  sku_name = "B_Standard_B1s"

  maintenance_window {
    day_of_week  = 0
    start_hour   = 8
    start_minute = 0
  }
}
`, r.template(data), data.RandomInteger, version)
}

func (r MySqlFlexibleServerResource) updateMaintenanceWindowUpdated(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  storage {
    size_gb = 32
  }
  version  = "%s"
  sku_name = "B_Standard_B1s"

  maintenance_window {
    day_of_week  = 3
    start_hour   = 7
    start_minute = 15
  }
}
`, r.template(data), data.RandomInteger, version)
}

func (r MySqlFlexibleServerResource) updateSku(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  storage {
    size_gb = 32
  }
  version  = "%s"
  sku_name = "MO_Standard_E2ds_v4"
}
`, r.template(data), data.RandomInteger, version)
}

func (r MySqlFlexibleServerResource) pitr(data acceptance.TestData, version string) string {
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
`, r.basic(data, version), data.RandomInteger, time.Now().Add(time.Duration(15)*time.Minute).UTC().Format(time.RFC3339))
}
