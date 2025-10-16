package oracle_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/dbsystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DbSystemResource struct{}

func (a DbSystemResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dbsystems.ParseDbSystemID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Oracle.OracleClient.DbSystems.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestDbSystemResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.DbSystemResource{}.ResourceType(), "test")
	r := DbSystemResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "initial_data_storage_size_in_gb", "pluggable_database_name", "resource_anchor_id"),
	})
}

func TestDbSystemResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.DbSystemResource{}.ResourceType(), "test")
	r := DbSystemResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "initial_data_storage_size_in_gb", "pluggable_database_name", "resource_anchor_id"),
	})
}

func TestDbSystemResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.DbSystemResource{}.ResourceType(), "test")
	r := DbSystemResource{}
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

func (a DbSystemResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
  %s
resource "azurerm_oracle_db_system" "test" {
  location                        = "%[3]s"
  zones               			      = ["2"]
  name                            = "acctest%[2]d"
  resource_group_name             = azurerm_resource_group.test.name
  admin_password			            = "testAdminPassword123##"
  compute_count                   = 4
  compute_model                   = "ECPU"
  database_edition      		      = "EnterpriseEdition"
  database_system_options {
    storage_management = "LVM"
  }
  database_version				        = "19.27.0.0"
  hostname                        = "dbhostname"
  license_model                   = "LicenseIncluded"
  network_anchor_id               = "/subscriptions/049e5678-fbb1-4861-93f3-7528bd0779fd/resourceGroups/white-glove/providers/Oracle.Database/networkAnchors/terraform-na"
  resource_anchor_id              = "/subscriptions/049e5678-fbb1-4861-93f3-7528bd0779fd/resourceGroups/white-glove/providers/Oracle.Database/resourceAnchors/ra-white-glove"
  shape                        	  = "VM.Standard.x86"
  source                  		    = "None"
  ssh_public_keys                 = ["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"]
}`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a DbSystemResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
  %s
resource "azurerm_oracle_db_system" "test" {
  location                        = "%[3]s"
  zones              			        = ["2"]
  name                            = "acctest%[2]d"
  resource_group_name             = azurerm_resource_group.test.name
  admin_password			            = "testAdminPassword123##"
  compute_count                   = 4
  compute_model                   = "ECPU"
  database_edition          		  = "EnterpriseEdition"
  database_system_options {
    storage_management = "LVM"
  }
  database_version				        = "19.27.0.0"
  disk_redundancy                 = "Normal"
  hostname                        = "dbhostname"
  license_model                   = "LicenseIncluded"
  network_anchor_id               = "/subscriptions/049e5678-fbb1-4861-93f3-7528bd0779fd/resourceGroups/white-glove/providers/Oracle.Database/networkAnchors/terraform-na"
  resource_anchor_id              = "/subscriptions/049e5678-fbb1-4861-93f3-7528bd0779fd/resourceGroups/white-glove/providers/Oracle.Database/resourceAnchors/ra-white-glove"
  source                     		  = "None"
  shape                        	  = "VM.Standard.x86"
  ssh_public_keys                 = ["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"]
  pluggable_database_name         = "testPdbName"
  storage_volume_performance_mode = "HighPerformance"
  display_name                	  = "acctest%[2]d"
  initial_data_storage_size_in_gb = 256
  node_count			          		  = 1
  tags = {
    test = "testTag1"
  }
  time_zone                        = "UTC"

}`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a DbSystemResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_db_system" "import" {
  location                        = azurerm_oracle_db_system.test.location
  name                            = azurerm_oracle_db_system.test.name
  zones               			      = azurerm_oracle_db_system.test.zones
  resource_group_name             = azurerm_oracle_db_system.test.resource_group_name
  admin_password                  = azurerm_oracle_db_system.test.admin_password
  compute_count                   = azurerm_oracle_db_system.test.compute_count
  compute_model                   = azurerm_oracle_db_system.test.compute_model
  database_edition      		      = azurerm_oracle_db_system.test.database_edition
  database_version				        = azurerm_oracle_db_system.test.database_version
  hostname                        = azurerm_oracle_db_system.test.hostname
  license_model                   = azurerm_oracle_db_system.test.license_model
  network_anchor_id               = azurerm_oracle_db_system.test.network_anchor_id
  resource_anchor_id              = azurerm_oracle_db_system.test.resource_anchor_id
  shape                        	  = azurerm_oracle_db_system.test.shape
  source                  		    = azurerm_oracle_db_system.test.source
  ssh_public_keys                 = azurerm_oracle_db_system.test.ssh_public_keys

}`, a.basic(data))
}

func (a DbSystemResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
