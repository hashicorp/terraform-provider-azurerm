// Copyright © 2025, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"testing"
)

type AdbsCrossRegionDisasterRecoveryResource struct{}

func (a AdbsCrossRegionDisasterRecoveryResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := autonomousdatabases.ParseAutonomousDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Oracle.OracleClient.AutonomousDatabases.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving adbs_crdr %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
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
		data.ImportStep("admin_password"),
	})
}

func TestAdbsCrossRegionDisasterRecoveryResource_update(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.AutonomousDatabaseCrossRegionDisasterRecoveryResource{}.ResourceType(), "adbs_secondary_crdr")
	r := AdbsCrossRegionDisasterRecoveryResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAdbsCrossRegionDisasterRecoveryResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.AutonomousDatabaseCrossRegionDisasterRecoveryResource{}.ResourceType(), "import")
	r := AdbsCrossRegionDisasterRecoveryResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.requiresImport(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
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
  name = "OFakeO%[2]d"

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
  customer_contacts = ["test@test.com"]
}

resource "azurerm_oracle_autonomous_database_cross_region_disaster_recovery" "adbs_secondary_crdr" {
  name = local.crdr_dbname

  display_name                     = local.crdr_dbname
  location                         = "%[4]s"
  database_type                    = "CrossRegionDisasterRecovery"
  source                           = "CrossRegionDisasterRecovery"
  source_id      				   = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.id
  source_ocid                      = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.ocid
  remote_disaster_recovery_type    = "Adg"
  replicate_automatic_backups	   = true
  subnet_id                        = azurerm_subnet.fra_vnet_subnet_test.id
  virtual_network_id               = azurerm_virtual_network.fra_vnet_test.id

  resource_group_name              = azurerm_resource_group.crdr_rg.name
  source_location				   = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.location
  license_model                    = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.license_model
  backup_retention_period_in_days  = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.backup_retention_period_in_days
  auto_scaling_enabled             = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.auto_scaling_enabled
  auto_scaling_for_storage_enabled = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.auto_scaling_for_storage_enabled
  mtls_connection_required         = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.mtls_connection_required
  data_storage_size_in_tbs         = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.data_storage_size_in_tbs
  compute_model                    = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.compute_model
  compute_count                    = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.compute_count
  db_workload                      = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.db_workload
  db_version                       = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.db_version
  admin_password                   = "TestPass#2024#"
  character_set                    = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.character_set
  national_character_set           = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.national_character_set
}
`, a.template(data), data.RandomInteger, "eastus", "westus")
}

func (a AdbsCrossRegionDisasterRecoveryResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_cross_region_disaster_recovery" "import" {
  name = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.name

  display_name                     = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.name
  location                         = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.location
  database_type                    = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.database_type
  source                           = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.source
  source_id      				   = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.source_id
  source_ocid                      = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.source_ocid
  remote_disaster_recovery_type    = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.remote_disaster_recovery_type
  replicate_automatic_backups	   = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.replicate_automatic_backups
  subnet_id                        = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.subnet_id
  virtual_network_id               = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.virtual_network_id
  resource_group_name              = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.resource_group_name
  source_location				   = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.source_location
  license_model                    = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.license_model
  backup_retention_period_in_days  = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.backup_retention_period_in_days
  auto_scaling_enabled             = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.auto_scaling_enabled
  auto_scaling_for_storage_enabled = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.auto_scaling_for_storage_enabled
  mtls_connection_required         = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.mtls_connection_required
  data_storage_size_in_tbs         = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.data_storage_size_in_tbs
  compute_model                    = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.compute_model
  compute_count                    = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.compute_count
  db_workload                      = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.db_workload
  db_version                       = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.db_version
  admin_password                   = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.admin_password
  character_set                    = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.character_set
  national_character_set           = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.national_character_set
}
`, a.complete(data))
}

func (a AdbsCrossRegionDisasterRecoveryResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_oracle_autonomous_database_cross_region_disaster_recovery" "adbs_secondary_crdr" {
    name =  "OFakeC%[2]d"

  display_name                     = "OFakeC%[2]d"
  remote_disaster_recovery_type    = "Adg"
  replicate_automatic_backups	   = false
  subnet_id                        = azurerm_subnet.fra_vnet_subnet_test.id
  virtual_network_id               = azurerm_virtual_network.fra_vnet_test.id
  location                         = "%[4]s"
  database_type                    = "CrossRegionDisasterRecovery"
  source                           = "CrossRegionDisasterRecovery"
  source_id      				   = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.id
  source_ocid                      = azurerm_oracle_autonomous_database.adbs_primary_for_crdr.ocid
  resource_group_name              = azurerm_resource_group.crdr_rg.name
  source_location				   = "%[3]s"
  license_model                    = "LicenseIncluded"
  backup_retention_period_in_days  = 12
  auto_scaling_enabled             = false
  auto_scaling_for_storage_enabled = false
  mtls_connection_required         = false
  data_storage_size_in_tbs         = 1
  compute_model                    = "ECPU"
  compute_count                    = 2
  db_workload                      = "DW"
  admin_password                   = "TestPass#2024#"
  db_version                       = "19c"
  character_set                    = "AL32UTF8"
  national_character_set           = "AL16UTF16"
}
`, a.complete(data), data.RandomInteger, "eastus", "westus")
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
