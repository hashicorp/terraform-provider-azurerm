// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AutonomousDatabaseCloneResource struct{}

func TestAccAutonomousDatabaseClone_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_clone", "test")
	r := AutonomousDatabaseCloneResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("clone_type").HasValue("Full"),
				check.That(data.ResourceName).Key("source").HasValue("Database"),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctestADB-%d-clone", data.RandomInteger)),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccAutonomousDatabaseClone_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_clone", "test")
	r := AutonomousDatabaseCloneResource{}

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

func TestAccAutonomousDatabaseClone_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_clone", "test")
	r := AutonomousDatabaseCloneResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("clone_type").HasValue("Full"),
				check.That(data.ResourceName).Key("source").HasValue("Database"),
				check.That(data.ResourceName).Key("is_refreshable_clone").HasValue("true"),
				check.That(data.ResourceName).Key("refreshable_model").HasValue("Manual"),
				check.That(data.ResourceName).Key("auto_scaling_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auto_scaling_for_storage_enabled").HasValue("true"),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccAutonomousDatabaseClone_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_clone", "test")
	r := AutonomousDatabaseCloneResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("compute_count").HasValue("2"),
				check.That(data.ResourceName).Key("data_storage_size_in_tbs").HasValue("1"),
			),
		},
		data.ImportStep("admin_password"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("compute_count").HasValue("4"),
				check.That(data.ResourceName).Key("data_storage_size_in_tbs").HasValue("2"),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccAutonomousDatabaseClone_metadataClone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_clone", "test")
	r := AutonomousDatabaseCloneResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.metadataClone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("clone_type").HasValue("Metadata"),
				check.That(data.ResourceName).Key("source").HasValue("Database"),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func (r AutonomousDatabaseCloneResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := autonomousdatabases.ParseAutonomousDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Oracle.OracleClient.AutonomousDatabases
	resp, err := client.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r AutonomousDatabaseCloneResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_clone" "test" {
  name                = "ADB%dclone"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  source_id  = azurerm_oracle_autonomous_database.source.id
  clone_type = "Full"
  source     = "Database"

  admin_password                    = "BEstrO0ng_#11"
  backup_retention_period_in_days   = 7
  character_set                     = "AL32UTF8"
  compute_count                     = 2.0
  compute_model                     = "ECPU"
  data_storage_size_in_tbs         = 1
  db_version                       = "19c"
  db_workload                      = "OLTP"
  display_name                     = "ADB%dclone"
  license_model                    = "LicenseIncluded"
  auto_scaling_enabled             = false
  auto_scaling_for_storage_enabled = false
  mtls_connection_required         = false
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_subnet.test.id
  virtual_network_id               = azurerm_virtual_network.test.id

  tags = {
    Environment = "Test"
    Purpose     = "Clone"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r AutonomousDatabaseCloneResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_clone" "import" {
  name                = azurerm_oracle_autonomous_database_clone.test.name
  resource_group_name = azurerm_oracle_autonomous_database_clone.test.resource_group_name
  location            = azurerm_oracle_autonomous_database_clone.test.location

  source_id  = azurerm_oracle_autonomous_database_clone.test.source_id
  clone_type = azurerm_oracle_autonomous_database_clone.test.clone_type
  source     = azurerm_oracle_autonomous_database_clone.test.source

  admin_password                    = azurerm_oracle_autonomous_database_clone.test.admin_password
  backup_retention_period_in_days   = azurerm_oracle_autonomous_database_clone.test.backup_retention_period_in_days
  character_set                     = azurerm_oracle_autonomous_database_clone.test.character_set
  compute_count                     = azurerm_oracle_autonomous_database_clone.test.compute_count
  compute_model                     = azurerm_oracle_autonomous_database_clone.test.compute_model
  data_storage_size_in_tbs         = azurerm_oracle_autonomous_database_clone.test.data_storage_size_in_tbs
  db_version                       = azurerm_oracle_autonomous_database_clone.test.db_version
  db_workload                      = azurerm_oracle_autonomous_database_clone.test.db_workload
  display_name                     = azurerm_oracle_autonomous_database_clone.test.display_name
  license_model                    = azurerm_oracle_autonomous_database_clone.test.license_model
  auto_scaling_enabled             = azurerm_oracle_autonomous_database_clone.test.auto_scaling_enabled
  auto_scaling_for_storage_enabled = azurerm_oracle_autonomous_database_clone.test.auto_scaling_for_storage_enabled
  mtls_connection_required         = azurerm_oracle_autonomous_database_clone.test.mtls_connection_required
  national_character_set           = azurerm_oracle_autonomous_database_clone.test.national_character_set
  subnet_id                        = azurerm_oracle_autonomous_database_clone.test.subnet_id
  virtual_network_id               = azurerm_oracle_autonomous_database_clone.test.virtual_network_id

  tags = azurerm_oracle_autonomous_database_clone.test.tags
}
`, r.basic(data))
}

func (r AutonomousDatabaseCloneResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_clone" "test" {
  name                = "ADB%dclone"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  source_id  = azurerm_oracle_autonomous_database.source.id
  clone_type = "Full"
  source     = "Database"

  admin_password                    = "BEstrO0ng_#11"
  backup_retention_period_in_days   = 15
  character_set                     = "AL32UTF8"
  compute_count                     = 4.0
  compute_model                     = "ECPU"
  data_storage_size_in_tbs         = 2
  db_version                       = "19c"
  db_workload                      = "OLTP"
  display_name                     = "ADB%dclone"
  license_model                    = "LicenseIncluded"
  auto_scaling_enabled             = true
  auto_scaling_for_storage_enabled = true
  mtls_connection_required         = false
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_subnet.test.id
  virtual_network_id               = azurerm_virtual_network.test.id

  # Clone-specific optional fields
  is_refreshable_clone = true
  refreshable_model    = "Manual"

  customer_contacts = ["test@example.com"]

  tags = {
    Environment = "Test"
    Purpose     = "CompleteClone"
    Type        = "Full"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r AutonomousDatabaseCloneResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_clone" "test" {
  name                = "ADB%dclone"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  source_id  = azurerm_oracle_autonomous_database.source.id
  clone_type = "Full"
  source     = "Database"

  admin_password                    = "BEstrO0ng_#11"
  backup_retention_period_in_days   = 7
  character_set                     = "AL32UTF8"
  compute_count                     = 4.0  # Updated from 2.0
  compute_model                     = "ECPU"
  data_storage_size_in_tbs         = 2    # Updated from 1
  db_version                       = "19c"
  db_workload                      = "OLTP"
  display_name                     = "ADB%dclone"
  license_model                    = "LicenseIncluded"
  auto_scaling_enabled             = true  # Updated from false
  auto_scaling_for_storage_enabled = true  # Updated from false
  mtls_connection_required         = false
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_subnet.test.id
  virtual_network_id               = azurerm_virtual_network.test.id

  tags = {
    Environment = "Test"
    Purpose     = "Clone"
    Updated     = "true"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r AutonomousDatabaseCloneResource) metadataClone(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_clone" "test" {
  name                = "ADB%dclone"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  source_id  = azurerm_oracle_autonomous_database.source.id
  clone_type = "Metadata"  
  source     = "Database"

  admin_password                    = "BEstrO0ng_#11"
  backup_retention_period_in_days   = 7
  character_set                     = "AL32UTF8"
  compute_count                     = 2.0
  compute_model                     = "ECPU"
  data_storage_size_in_tbs         = 1
  db_version                       = "19c"
  db_workload                      = "OLTP"
  display_name                     = "ADB%dclone"
  license_model                    = "LicenseIncluded"
  auto_scaling_enabled             = false
  auto_scaling_for_storage_enabled = false
  mtls_connection_required         = false
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_subnet.test.id
  virtual_network_id               = azurerm_virtual_network.test.id

  tags = {
    Environment = "Test"
    Purpose     = "MetadataClone"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r AutonomousDatabaseCloneResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-oracle-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "Oracle.Database.networkAttachments"

    service_delegation {
      name    = "Oracle.Database/networkAttachments"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

# Source autonomous database to clone from
resource "azurerm_oracle_autonomous_database" "source" {
  name                = "ADB%dsource"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  admin_password                    = "BEstrO0ng_#11"
  backup_retention_period_in_days   = 7
  character_set                     = "AL32UTF8"
  compute_count                     = 2.0
  compute_model                     = "ECPU"
  data_storage_size_in_tbs         = 1
  db_version                       = "19c"
  db_workload                      = "OLTP"
  display_name                     = "ADB%dsource"
  license_model                    = "LicenseIncluded"
  auto_scaling_enabled             = false
  auto_scaling_for_storage_enabled = false
  mtls_connection_required         = false
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_subnet.test.id
  virtual_network_id               = azurerm_virtual_network.test.id

  tags = {
    Environment = "Test"
    Purpose     = "Source"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
