// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mysql_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MySqlFlexibleServerResource struct{}

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
				check.That(data.ResourceName).Key("replica_capacity").Exists(),
			),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
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
			Config: r.updateStorage(data, 20, 360, true, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.updateStorage(data, 34, 402, false, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.updateStorageNoIOPS(data, 34, false, true),
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
			Config: r.failover(data, "2", "3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.failover(data, "3", "2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccMySqlFlexibleServer_createWithCustomerManagedKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withCustomerManagedKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccMySqlFlexibleServer_updateToCustomerManagedKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withoutCustomerManagedKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.withCustomerManagedKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password"),
	})
}

// this test can fail with an uninformative error, tracked here https://github.com/Azure/azure-rest-api-specs/issues/22980
func TestAccMySqlFlexibleServer_enableGeoRedundantBackup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enableGeoRedundantBackup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccMySqlFlexibleServer_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.identityNone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccMySqlFlexibleServer_writeOnlyPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.11.0"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.writeOnlyPassword(data, "a-secret-from-kv", 1),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("administrator_password_wo_version"),
			{
				Config: r.writeOnlyPassword(data, "a-secret-from-kv-updated", 2),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("administrator_password_wo_version"),
		},
	})
}

func TestAccMySqlFlexibleServer_updateToWriteOnlyPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.11.0"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("administrator_password"),
			{
				Config: r.writeOnlyPassword(data, "a-secret-from-kv", 1),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("administrator_password", "administrator_password_wo_version"),
			{
				Config: r.basic(data),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("administrator_password"),
		},
	})
}

func TestAccMySqlFlexibleServer_publicNetworkAccess(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_flexible_server", "test")
	r := MySqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.publicNetworkAccess(data, servers.EnableStatusEnumDisabled),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.publicNetworkAccess(data, servers.EnableStatusEnumEnabled),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (MySqlFlexibleServerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := servers.ParseFlexibleServerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MySQL.FlexibleServers.Servers.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id.ID(), err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Properties != nil), nil
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
`, data.RandomInteger, data.Locations.Ternary)
}

func (r MySqlFlexibleServerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "_admin_Terraform_892123456789312"
  administrator_password = "QAZwsx123"
  sku_name               = "B_Standard_B1ms"
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
  zone                   = azurerm_mysql_flexible_server.test.zone
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

  depends_on = [azurerm_subnet.test]
}

resource "azurerm_mysql_flexible_server" "test" {
  name                         = "acctest-fs-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_password       = "QAZwsx123"
  zone                         = "3"
  version                      = "8.0.21"
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false

  storage {
    size_gb             = 20
    iops                = 360
    log_on_disk_enabled = true
  }

  delegated_subnet_id = azurerm_subnet.test.id
  private_dns_zone_id = azurerm_private_dns_zone.test.id
  sku_name            = "MO_Standard_E2ds_v4"

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

  depends_on = [azurerm_subnet.test]
}

resource "azurerm_mysql_flexible_server" "test" {
  name                         = "acctest-fs-%[2]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_password       = "123wsxQAZ"
  zone                         = "3"
  version                      = "8.0.21"
  backup_retention_days        = 10
  geo_redundant_backup_enabled = false

  storage {
    size_gb             = 32
    iops                = 400
    auto_grow_enabled   = false
    io_scaling_enabled  = false
    log_on_disk_enabled = false
  }

  delegated_subnet_id = azurerm_subnet.test.id
  private_dns_zone_id = azurerm_private_dns_zone.test.id
  sku_name            = "MO_Standard_E4ds_v4"

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
  administrator_login    = "_admin_Terraform_892123456789312"
  administrator_password = "QAZwsx123"
  sku_name               = "GP_Standard_D2ds_v4"

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
  administrator_login    = "_admin_Terraform_892123456789312"
  administrator_password = "QAZwsx123"
  sku_name               = "GP_Standard_D2ds_v4"

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
  administrator_login    = "_admin_Terraform_892123456789312"
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
  zone                   = "2"
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
  zone     = "2"

  lifecycle {
    ignore_changes = [
      high_availability.0.standby_availability_zone
    ]
  }
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
  version                = "8.0.21"

  high_availability {
    mode                      = "ZoneRedundant"
    standby_availability_zone = "2"
  }

  sku_name = "GP_Standard_D2ds_v4"
  zone     = "3"
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
  version                = "8.0.21"

  high_availability {
    mode                      = "ZoneRedundant"
    standby_availability_zone = "3"
  }

  sku_name = "GP_Standard_D2ds_v4"
  zone     = "2"
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
  zone                              = "2"

  lifecycle {
    ignore_changes = [source_server_id]
  }
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
  version                = "8.0.21"
  zone                   = "2"
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
  version             = "8.0.21"
  zone                = "2"
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
  version             = "8.0.21"
  zone                = "2"

  lifecycle {
    ignore_changes = [source_server_id]
  }
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
  sku_name                     = "B_Standard_B1ms"
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

func (r MySqlFlexibleServerResource) updateStorage(data acceptance.TestData, sizeGB int, iops int, autoGrowEnabled bool, ioScalingEnabled bool) string {
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
  version                      = "8.0.21"
  zone                         = "2"

  storage {
    size_gb            = %d
    iops               = %d
    auto_grow_enabled  = %t
    io_scaling_enabled = %t
  }
}
`, r.template(data), data.RandomInteger, sizeGB, iops, autoGrowEnabled, ioScalingEnabled)
}

func (r MySqlFlexibleServerResource) updateStorageNoIOPS(data acceptance.TestData, sizeGB int, autoGrowEnabled bool, ioScalingEnabled bool) string {
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
  version                      = "8.0.21"
  zone                         = "2"

  storage {
    size_gb            = %d
    auto_grow_enabled  = %t
    io_scaling_enabled = %t
  }
}
`, r.template(data), data.RandomInteger, sizeGB, autoGrowEnabled, ioScalingEnabled)
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
  version                = "8.0.21"

  high_availability {
    mode                      = "ZoneRedundant"
    standby_availability_zone = "%s"
  }
}
`, r.template(data), data.RandomInteger, primaryZone, standbyZone)
}

func (r MySqlFlexibleServerResource) cmkTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestmi%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "server" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id

  key_permissions = ["Get", "List", "WrapKey", "UnwrapKey", "GetRotationPolicy", "SetRotationPolicy"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy", "SetRotationPolicy"]
}

resource "azurerm_key_vault_key" "test" {
  name         = "test"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.server,
  ]
}
`, data.RandomInteger, data.Locations.Ternary, data.RandomString, data.RandomString)
}

func (r MySqlFlexibleServerResource) withoutCustomerManagedKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "_admin_Terraform_892123456789312"
  administrator_password = "QAZwsx123"
  sku_name               = "B_Standard_B1ms"
  zone                   = "2"
}
`, r.cmkTemplate(data), data.RandomInteger)
}

func (r MySqlFlexibleServerResource) withCustomerManagedKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "_admin_Terraform_892123456789312"
  administrator_password = "QAZwsx123"
  sku_name               = "B_Standard_B1ms"
  zone                   = "2"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  customer_managed_key {
    key_vault_key_id                  = azurerm_key_vault_key.test.id
    primary_user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }
}
`, r.cmkTemplate(data), data.RandomInteger)
}

func (r MySqlFlexibleServerResource) enableGeoRedundantBackup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-mysql2-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctestmi%s"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
}

resource "azurerm_key_vault" "test2" {
  name                     = "acctestkv2%s"
  location                 = azurerm_resource_group.test2.location
  resource_group_name      = azurerm_resource_group.test2.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "server2" {
  key_vault_id = azurerm_key_vault.test2.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_user_assigned_identity.test2.principal_id

  key_permissions = ["Get", "List", "WrapKey", "UnwrapKey"]
}

resource "azurerm_key_vault_access_policy" "client2" {
  key_vault_id = azurerm_key_vault.test2.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
}

resource "azurerm_key_vault_key" "test2" {
  name         = "test2"
  key_vault_id = azurerm_key_vault.test2.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client2,
    azurerm_key_vault_access_policy.server2,
  ]
}

resource "azurerm_mysql_flexible_server" "test" {
  name                         = "acctest-fs-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "_admin_Terraform_892123456789312"
  administrator_password       = "QAZwsx123"
  sku_name                     = "B_Standard_B1ms"
  zone                         = "2"
  geo_redundant_backup_enabled = true

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id, azurerm_user_assigned_identity.test2.id]
  }

  customer_managed_key {
    key_vault_key_id                     = azurerm_key_vault_key.test.id
    primary_user_assigned_identity_id    = azurerm_user_assigned_identity.test.id
    geo_backup_key_vault_key_id          = azurerm_key_vault_key.test2.id
    geo_backup_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  }
}
`, r.cmkTemplate(data), data.RandomInteger, data.Locations.Ternary, data.RandomString, data.RandomString, data.RandomInteger)
}

func (r MySqlFlexibleServerResource) identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestmi%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "_admin_Terraform_892123456789312"
  administrator_password = "QAZwsx123"
  sku_name               = "B_Standard_B1ms"
  zone                   = "2"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.template(data), data.RandomString, data.RandomInteger)
}

func (r MySqlFlexibleServerResource) identityNone(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestmi%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "_admin_Terraform_892123456789312"
  administrator_password = "QAZwsx123"
  sku_name               = "B_Standard_B1ms"
  zone                   = "2"
}
`, r.template(data), data.RandomString, data.RandomInteger)
}

func (r MySqlFlexibleServerResource) writeOnlyPassword(data acceptance.TestData, secret string, version int) string {
	return fmt.Sprintf(`
%s

%s

resource "azurerm_mysql_flexible_server" "test" {
  name                              = "acctest-fs-%[3]d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = azurerm_resource_group.test.location
  administrator_login               = "_admin_Terraform_892123456789312"
  administrator_password_wo         = ephemeral.azurerm_key_vault_secret.test.value
  administrator_password_wo_version = %[4]d
  sku_name                          = "B_Standard_B1ms"
}
`, r.template(data), acceptance.WriteOnlyKeyVaultSecretTemplate(data, secret), data.RandomInteger, version)
}

func (r MySqlFlexibleServerResource) publicNetworkAccess(data acceptance.TestData, pna servers.EnableStatusEnum) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "_admin_Terraform_892123456789312"
  administrator_password = "QAZwsx123"
  sku_name               = "B_Standard_B1ms"
  public_network_access  = "%s"
}
`, r.template(data), data.RandomInteger, pna)
}
