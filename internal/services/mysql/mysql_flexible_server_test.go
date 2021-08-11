package mysql_test

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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MysqlFlexibleServerResource struct {
}

func TestAccMysqlFlexibleServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MysqlFlexibleServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("ha_enabled").Exists(),
				check.That(data.ResourceName).Key("storage_profile.#").Exists(),
				check.That(data.ResourceName).Key("byok_enforcement_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("create_mode", "administrator_login_password"),
	})
}

func TestAccMysqlFlexibleServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MysqlFlexibleServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("ha_enabled").Exists(),
				check.That(data.ResourceName).Key("storage_profile.#").Exists(),
				check.That(data.ResourceName).Key("byok_enforcement_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMysqlFlexibleServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MysqlFlexibleServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("byok_enforcement_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("create_mode", "administrator_login_password"),
	})
}

func TestAccMysqlFlexibleServer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MysqlFlexibleServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("byok_enforcement_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("create_mode", "administrator_login_password"),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("byok_enforcement_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("create_mode", "administrator_login_password"),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("byok_enforcement_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("create_mode", "administrator_login_password"),
	})
}

func TestAccMysqlFlexibleServer_updateMaintenanceWindow(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MysqlFlexibleServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("ha_enabled").Exists(),
				check.That(data.ResourceName).Key("storage_profile.#").Exists(),
				check.That(data.ResourceName).Key("byok_enforcement_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("create_mode", "administrator_login_password"),
		{
			Config: r.updateMaintenanceWindow(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("ha_enabled").Exists(),
				check.That(data.ResourceName).Key("storage_profile.#").Exists(),
				check.That(data.ResourceName).Key("byok_enforcement_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("create_mode", "administrator_login_password"),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("ha_enabled").Exists(),
				check.That(data.ResourceName).Key("storage_profile.#").Exists(),
				check.That(data.ResourceName).Key("byok_enforcement_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("create_mode", "administrator_login_password"),
	})
}

func TestAccMysqlFlexibleServer_updateStorageProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MysqlFlexibleServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("ha_enabled").Exists(),
				check.That(data.ResourceName).Key("storage_profile.#").Exists(),
				check.That(data.ResourceName).Key("byok_enforcement_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("create_mode", "administrator_login_password"),
		{
			Config: r.updateStorageProfile(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("ha_enabled").Exists(),
				check.That(data.ResourceName).Key("storage_profile.#").Exists(),
				check.That(data.ResourceName).Key("byok_enforcement_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("create_mode", "administrator_login_password"),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("ha_enabled").Exists(),
				check.That(data.ResourceName).Key("storage_profile.#").Exists(),
				check.That(data.ResourceName).Key("byok_enforcement_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("create_mode", "administrator_login_password"),
	})
}

func TestAccMysqlFlexibleServer_delegatedSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MysqlFlexibleServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("create_mode", "administrator_login_password"),
		{
			Config: r.delegatedSubnet(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("create_mode", "administrator_login_password", "delegated_subnet_id"),
		{
			Config: r.delegatedSubnetRestore(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("create_mode", "administrator_login_password", "delegated_subnet_id"),
	})
}

func TestAccMysqlFlexibleServer_replica(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "replica")
	r := MysqlFlexibleServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.createModeReplica(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("administrator_login").Exists(),
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("ha_enabled").Exists(),
				check.That(data.ResourceName).Key("replication_role").HasValue("Replica"),
				check.That(data.ResourceName).Key("storage_profile.#").Exists(),
				check.That(data.ResourceName).Key("byok_enforcement_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("create_mode", "administrator_login_password"),
	})
}

func TestAccMysqlFlexibleServer_pitr(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MysqlFlexibleServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("ha_enabled").Exists(),
				check.That(data.ResourceName).Key("storage_profile.#").Exists(),
				check.That(data.ResourceName).Key("byok_enforcement_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("create_mode", "administrator_login_password"),
		{
			PreConfig: func() { time.Sleep(10 * time.Minute) },
			Config:    r.createModePitr(data),
			Check: resource.ComposeTestCheckFunc(
				check.That("azurerm_mysql_flexible_server.pitr").ExistsInAzure(r),
				check.That("azurerm_mysql_flexible_server.pitr").Key("version").Exists(),
				check.That("azurerm_mysql_flexible_server.pitr").Key("ha_enabled").Exists(),
				check.That("azurerm_mysql_flexible_server.pitr").Key("storage_profile.#").Exists(),
				check.That("azurerm_mysql_flexible_server.pitr").Key("byok_enforcement_enabled").Exists(),
				check.That("azurerm_mysql_flexible_server.pitr").Key("fqdn").Exists(),
				check.That("azurerm_mysql_flexible_server.pitr").Key("public_network_access_enabled").Exists(),
				check.That("azurerm_mysql_flexible_server.pitr").Key("replica_capacity").Exists(),
			),
		},

		data.ImportStep("create_mode", "creation_source_database_id", "restore_point_in_time"),
	})
}

func (MysqlFlexibleServerResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.FlexibleServerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MySQL.FlexibleServersClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Mysql Flexible Server %q (resource group: %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ServerProperties != nil), nil
}

func (MysqlFlexibleServerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-Mysql-FlexibleServer-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r MysqlFlexibleServerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                         = "acctest-fs-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "cloudsa"
  administrator_login_password = "<administratorLoginPassword>"
  sku {
    name = "Standard_D2ds_v4"
    tier = "GeneralPurpose"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MysqlFlexibleServerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "import" {
  name                         = azurerm_mysql_flexible_server.test.name
  resource_group_name          = azurerm_mysql_flexible_server.test.resource_group_name
  location                     = azurerm_mysql_flexible_server.test.location
  administrator_login          = azurerm_mysql_flexible_server.test.administrator_login
  administrator_login_password = azurerm_mysql_flexible_server.test.administrator_login_password
  sku {
    name = azurerm_mysql_flexible_server.test.sku.0.name
    tier = azurerm_mysql_flexible_server.test.sku.0.tier
  }
}
`, r.basic(data))
}

func (r MysqlFlexibleServerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                         = "acctest-fs-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "cloudsa"
  administrator_login_password = "<administratorLoginPassword>"
  create_mode                  = "Default"
  version                      = "5.7"
  ha_enabled                   = false

  maintenance_window {
    day_of_week  = 2
    start_hour   = 14
    start_minute = 0
  }

  storage_profile {
    backup_retention_days    = 14
    storage_autogrow_enabled = true
    storage_iops             = 2
    storage_mb               = 5120
  }

  identity {
    type = "SystemAssigned"
  }

  sku {
    name = "Standard_D2ds_v4"
    tier = "GeneralPurpose"
  }

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MysqlFlexibleServerResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                         = "acctest-fs-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "cloudsa"
  administrator_login_password = "<administratorLoginPasswordUpdate>"
  create_mode                  = "Default"
  version                      = "5.7"
  ha_enabled                   = true

  maintenance_window {
    day_of_week  = 1
    start_hour   = 10
    start_minute = 30
  }

  storage_profile {
    backup_retention_days    = 7
    storage_autogrow_enabled = false
    storage_iops             = 1
    storage_mb               = 8192
  }

  identity {
    type = "SystemAssigned"
  }

  sku {
    name = "Standard_B2s"
    tier = "Burstable"
  }

  tags = {
    ENV = "Test"
    PRO = "Stage"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MysqlFlexibleServerResource) updateMaintenanceWindow(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                         = "acctest-fs-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "cloudsa"
  administrator_login_password = "<administratorLoginPassword>"
  maintenance_window {
    day_of_week  = 1
    start_hour   = 10
    start_minute = 30
  }
  sku {
    name = "Standard_D2ds_v4"
    tier = "GeneralPurpose"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MysqlFlexibleServerResource) updateStorageProfile(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                         = "acctest-fs-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "cloudsa"
  administrator_login_password = "<administratorLoginPassword>"
  storage_profile {
    backup_retention_days    = 14
    storage_autogrow_enabled = true
    storage_iops             = 2
    storage_mb               = 5120
  }
  sku {
    name = "Standard_D2ds_v4"
    tier = "GeneralPurpose"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MysqlFlexibleServerResource) delegatedSubnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

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
      name = "Microsoft.DBforMySQL/flexibleServers"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_mysql_flexible_server" "test" {
  name                         = "acctest-fs-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "cloudsa"
  administrator_login_password = "<administratorLoginPassword>"
  delegated_subnet_id          = azurerm_subnet.test.id
  sku {
    name = "Standard_D2ds_v4"
    tier = "GeneralPurpose"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MysqlFlexibleServerResource) delegatedSubnetRestore(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

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
      name = "Microsoft.DBForMySql/flexibleServers"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_mysql_flexible_server" "test" {
  name                         = "acctest-fs-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "cloudsa"
  administrator_login_password = "<administratorLoginPassword>"
  sku {
    name = "Standard_D2ds_v4"
    tier = "GeneralPurpose"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MysqlFlexibleServerResource) createModeReplica(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "replica" {
  name                      = "acctest-fs-replica-%d"
  resource_group_name       = azurerm_resource_group.test.name
  location                  = azurerm_resource_group.test.location
  create_mode               = "Replica"
  source_flexible_server_id = azurerm_mysql_flexible_server.test.id
}
`, r.basic(data), data.RandomInteger)
}

func (r MysqlFlexibleServerResource) createModePitr(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "pitr" {
  name                      = "acctest-fs-pitr-%d"
  resource_group_name       = azurerm_resource_group.test.name
  location                  = azurerm_resource_group.test.location
  create_mode               = "PointInTimeRestore"
  source_flexible_server_id = azurerm_mysql_flexible_server.test.id
  restore_point_in_time     = "%s"
}
`, r.basic(data), data.RandomInteger, time.Now().Add(time.Duration(10)*time.Minute).UTC().Format(time.RFC3339))
}
