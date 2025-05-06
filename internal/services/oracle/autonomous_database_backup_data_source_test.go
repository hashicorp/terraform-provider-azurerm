// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AutonomousDatabaseBackupDataSourceTest struct{}

func TestAccAutonomousDatabaseBackupDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_autonomous_database_backup", "test")
	r := AutonomousDatabaseBackupDataSourceTest{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("autonomous_database_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("backup_type").Exists(),
				check.That(data.ResourceName).Key("retention_period_in_days").Exists(),
				check.That(data.ResourceName).Key("autonomous_database_ocid").Exists(),
				check.That(data.ResourceName).Key("autonomous_database_backup_ocid").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
			),
		},
	})
}

func (r AutonomousDatabaseBackupDataSourceTest) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_backup" "test" {
  name                      = "backup%[2]d"
  display_name              = "backup%[2]d"
  resource_group_name       = azurerm_resource_group.test.name
  location                  = "%[3]s"
  autonomous_database_name  = azurerm_oracle_autonomous_database.test.name
  retention_period_in_days  = 30
}

data "azurerm_oracle_autonomous_database_backup" "test" {
  name                      = azurerm_oracle_autonomous_database_backup.test.name
  resource_group_name       = azurerm_resource_group.test.name
  autonomous_database_name  = azurerm_oracle_autonomous_database.test.name
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AutonomousDatabaseBackupDataSourceTest) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest%[1]d_vnet"
  address_space       = ["10.0.0.0/16"]
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "eacctest%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "delegation"

    service_delegation {
      actions = [
        "Microsoft.Network/networkinterfaces/*",
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
      name = "Oracle.Database/networkAttachments"
    }
  }
}

resource "azurerm_oracle_autonomous_database" "test" {
  name                             = "OFake%[1]d"
  display_name                     = "OFake%[1]d"
  resource_group_name              = azurerm_resource_group.test.name
  location                         = "%[2]s"
  compute_model                    = "ECPU"
  compute_count                    = 2
  license_model                    = "BringYourOwnLicense"
  backup_retention_period_in_days  = 12
  auto_scaling_enabled             = false
  auto_scaling_for_storage_enabled = false
  mtls_connection_required         = false
  data_storage_size_in_tbs         = 1
  db_workload                      = "OLTP"
  admin_password                   = "TestPass#2024#"
  db_version                       = "19c"
  character_set                    = "AL32UTF8"
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_subnet.test.id
  virtual_network_id               = azurerm_virtual_network.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}
