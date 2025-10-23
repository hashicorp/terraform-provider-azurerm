// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AdbsCrossRegionDisasterRecoveryResource struct{}

func (a AdbsCrossRegionDisasterRecoveryResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := autonomousdatabases.ParseAutonomousDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Oracle.OracleClient.AutonomousDatabases.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestAdbsCrossRegionDisasterRecoveryResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.AutonomousDatabaseCrossRegionDisasterRecoveryResource{}.ResourceType(), "adbs_secondary_crdr")
	r := AdbsCrossRegionDisasterRecoveryResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAdbsCrossRegionDisasterRecoveryResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.AutonomousDatabaseCrossRegionDisasterRecoveryResource{}.ResourceType(), "adbs_secondary_crdr")
	r := AdbsCrossRegionDisasterRecoveryResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAdbsCrossRegionDisasterRecoveryResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.AutonomousDatabaseCrossRegionDisasterRecoveryResource{}.ResourceType(), "adbs_secondary_crdr")
	r := AdbsCrossRegionDisasterRecoveryResource{}
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

func (a AdbsCrossRegionDisasterRecoveryResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
provider "azurerm" {
  features {}
}

locals {
  crdr_dbname = "OFakeN%[2]d"
}

resource "azurerm_oracle_autonomous_database" "adbs_primary_for_crdr" {
  name                             = "OFakeO%[2]d"
  display_name                     = "OFakeO%[2]d"
  resource_group_name              = azurerm_resource_group.crdr_rg.name
  location                         = "%[3]s"
  compute_model                    = "ECPU"
  compute_count                    = 2
  license_model                    = "LicenseIncluded"
  backup_retention_period_in_days  = 12
  auto_scaling_enabled             = false
  auto_scaling_for_storage_enabled = false
  mtls_connection_required         = false
  data_storage_size_in_tbs         = 1
  db_workload                      = "DW"
  admin_password                   = "TestPass#2024#"
  db_version                       = "19c"
  character_set                    = "AL32UTF8"
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_subnet.iad_vnet_subnet_test.id
  virtual_network_id               = azurerm_virtual_network.iad_vnet_test.id
  customer_contacts                = ["test@test.com"]
}

resource "azurerm_oracle_autonomous_database_cross_region_disaster_recovery" "adbs_secondary_crdr" {
  name                          = local.crdr_dbname
  display_name                  = local.crdr_dbname
  location                      = "%[4]s"
  resource_group_name           = azurerm_resource_group.crdr_rg.name
  subnet_id                     = azurerm_subnet.fra_vnet_subnet_test.id
  virtual_network_id            = azurerm_virtual_network.fra_vnet_test.id
  source_autonomous_database_id = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.id

}
`, a.template(data), data.RandomInteger, "eastus", "westus")
}

func (a AdbsCrossRegionDisasterRecoveryResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

locals {
  crdr_dbname = "OFakeC%[2]d"
}

resource "azurerm_oracle_autonomous_database" "adbs_primary_for_crdr" {
  name                             = "OFakeO%[2]d"
  display_name                     = "OFakeO%[2]d"
  resource_group_name              = azurerm_resource_group.crdr_rg.name
  location                         = "%[3]s"
  compute_model                    = "ECPU"
  compute_count                    = 2
  license_model                    = "LicenseIncluded"
  backup_retention_period_in_days  = 12
  auto_scaling_enabled             = false
  auto_scaling_for_storage_enabled = false
  mtls_connection_required         = false
  data_storage_size_in_tbs         = 1
  db_workload                      = "DW"
  admin_password                   = "TestPass#2024#"
  db_version                       = "19c"
  character_set                    = "AL32UTF8"
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_subnet.iad_vnet_subnet_test.id
  virtual_network_id               = azurerm_virtual_network.iad_vnet_test.id
  customer_contacts                = ["test@test.com"]
}

resource "azurerm_oracle_autonomous_database_cross_region_disaster_recovery" "adbs_secondary_crdr" {
  name                          = local.crdr_dbname
  resource_group_name           = azurerm_resource_group.crdr_rg.name
  display_name                  = local.crdr_dbname
  location                      = "%[4]s"
  source_autonomous_database_id = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.id
  subnet_id                     = azurerm_subnet.fra_vnet_subnet_test.id
  virtual_network_id            = azurerm_virtual_network.fra_vnet_test.id

  replicate_automatic_backups_enabled = true

  tags = {
    Purpose = "basic-acceptance"
  }

}
`, a.template(data), data.RandomInteger, "eastus", "westus")
}

func (a AdbsCrossRegionDisasterRecoveryResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_oracle_autonomous_database_cross_region_disaster_recovery" "import" {
  name                          = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.name
  display_name                  = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.display_name
  location                      = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.location
  resource_group_name           = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.resource_group_name
  subnet_id                     = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.subnet_id
  virtual_network_id            = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.virtual_network_id
  source_autonomous_database_id = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.source_autonomous_database_id
}
`, a.basic(data))
}

func (a AdbsCrossRegionDisasterRecoveryResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "crdr_rg" {
  name     = "CRDRacctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "iad_vnet_test" {
  name                = "CRDR_iad_acctest%[1]d_vnet"
  address_space       = ["10.0.0.0/16"]
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.crdr_rg.name
}

resource "azurerm_subnet" "iad_vnet_subnet_test" {
  name                 = "CRDR_iad_subnet_acctest%[1]d"
  resource_group_name  = azurerm_resource_group.crdr_rg.name
  virtual_network_name = azurerm_virtual_network.iad_vnet_test.name
  address_prefixes     = ["10.0.0.0/26"]

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

resource "azurerm_virtual_network" "fra_vnet_test" {
  name                = "CRDR_fra_acctest%[1]d_vnet"
  address_space       = ["10.0.0.0/16"]
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.crdr_rg.name
}

resource "azurerm_subnet" "fra_vnet_subnet_test" {
  name                 = "CRDR_fra_subnet_acctest%[1]d"
  resource_group_name  = azurerm_resource_group.crdr_rg.name
  virtual_network_name = azurerm_virtual_network.fra_vnet_test.name
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

`, data.RandomInteger, "eastus", "westus")
}
