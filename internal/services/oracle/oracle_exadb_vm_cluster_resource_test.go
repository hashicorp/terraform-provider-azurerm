// Copyright © 2025, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/exadbvmclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ExadbVmClusterResource struct{}

func (a ExadbVmClusterResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := exadbvmclusters.ParseExadbVMClusterID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Oracle.OracleClient25.ExadbVMClusters.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestExadbVmClusterResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ExadbVmClusterResource{}.ResourceType(), "test")
	r := ExadbVmClusterResource{}
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

func TestExadbVmClusterResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ExadbVmClusterResource{}.ResourceType(), "test")
	r := ExadbVmClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("data_collection_options", "system_version"),
	})
}

func TestExadbVmClusterResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ExadbVmClusterResource{}.ResourceType(), "test")
	r := ExadbVmClusterResource{}
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

func TestExadbVmClusterResource_update(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ExadbVmClusterResource{}.ResourceType(), "test")
	r := ExadbVmClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("data_collection_options"),
	})
}

func (a ExadbVmClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
  %s
resource "azurerm_oracle_exa_db_vm_cluster" "test" {
  location                        = "%[3]s"
  name                            = "OFakeVmacctest%[2]d"
  zones               			  = ["2"]
  resource_group_name             = azurerm_resource_group.test.name
  exascale_db_storage_vault_id	  = azurerm_oracle_exascale_db_storage_vault.test.id
  display_name                    = "OFakeVmacctest%[2]d"
  enabled_ecpu_count              = 16
  grid_image_ocid				  = "ocid1.dbpatch.oc1.iad.anuwcljtt5t4sqqaoajnkveobp3ryw7zlfrrcf6tb2ndvygp54z2gbds2hxa"
  hostname                        = "host"
  node_count              		  = 2
  shape                      	  = "EXADBXS"
  ssh_public_keys                 = ["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"]
  subnet_id                       = azurerm_subnet.virtual_network_subnet.id
  total_ecpu_count                = 32
  vm_file_system_storage {
    total_size_in_gbs = 440
  }	
  virtual_network_id              = azurerm_virtual_network.virtual_network.id
}`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a ExadbVmClusterResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
  %s
resource "azurerm_oracle_exa_db_vm_cluster" "test" {
  location                        = "%[3]s"
  name                            = "OFakeVmacctest%[2]d"
  zones               			   = ["2"]
  resource_group_name             = azurerm_resource_group.test.name
  exascale_db_storage_vault_id	  = azurerm_oracle_exascale_db_storage_vault.test.id
  display_name                    = "OFakeVmacctest%[2]d"
  enabled_ecpu_count              = 16
  hostname                        = "host"
  node_count              		  = 2
  shape                      	  = "EXADBXS"
  ssh_public_keys                 = ["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"]
  subnet_id                       = azurerm_subnet.virtual_network_subnet.id
  total_ecpu_count                = 32
  vm_file_system_storage {
    total_size_in_gbs = 440
  }	
  virtual_network_id              = azurerm_virtual_network.virtual_network.id
  backup_subnet_cidr              = "10.1.0.0/16"
  cluster_name              	  = "clustername"
  data_collection_options {
    diagnostics_events_enabled = true
    health_monitoring_enabled  = true
    incident_logs_enabled      = true
  }
  domain                      = "ociactsubnet.ociactvnet.oraclevcn.com"
  grid_image_ocid             = "ocid1.dbpatch.oc1.iad.anuwcljtt5t4sqqaoajnkveobp3ryw7zlfrrcf6tb2ndvygp54z2gbds2hxa"
  license_model               = "BringYourOwnLicense"
  private_zone_ocid     	  = "private_zoneocid"
  scan_listener_port_tcp      = 1521
  scan_listener_port_tcp_ssl  = 2484
  system_version              = "23.1.23.0.0.250207"
  time_zone        			  = "UTC"
}`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a ExadbVmClusterResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_oracle_exa_db_vm_cluster" "test" {
  location                        = "%[3]s"
  name                            = "OFakeVmacctest%[2]d"
  zones               			   = ["2"]
  resource_group_name             = azurerm_resource_group.test.name
  exascale_db_storage_vault_id	  = azurerm_oracle_exascale_db_storage_vault.test.id
  display_name                    = "OFakeVmacctest%[2]d"
  enabled_ecpu_count              = 16
  grid_image_ocid                 = "ocid1.dbpatch.oc1.iad.anuwcljtt5t4sqqaoajnkveobp3ryw7zlfrrcf6tb2ndvygp54z2gbds2hxa"
  hostname                        = "host"
  node_count              		  = 4
  shape                      	  = "EXADBXS"
  ssh_public_keys                 = ["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"]
  subnet_id                       = azurerm_subnet.virtual_network_subnet.id
  total_ecpu_count                = 64
  vm_file_system_storage {
    total_size_in_gbs = 440
  }	
  virtual_network_id              = azurerm_virtual_network.virtual_network.id
}`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a ExadbVmClusterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_exa_db_vm_cluster" "import" {
  location                        = azurerm_oracle_exa_db_vm_cluster.test.location
  name                            = azurerm_oracle_exa_db_vm_cluster.test.name
  resource_group_name             = azurerm_oracle_exa_db_vm_cluster.test.resource_group_name
  exascale_db_storage_vault_id    = azurerm_oracle_exa_db_vm_cluster.test.exascale_db_storage_vault_id
  display_name                    = azurerm_oracle_exa_db_vm_cluster.test.display_name
  enabled_ecpu_count      		  = azurerm_oracle_exa_db_vm_cluster.test.enabled_ecpu_count
  grid_image_ocid                 = azurerm_oracle_exa_db_vm_cluster.test.grid_image_ocid
  hostname   					  = azurerm_oracle_exa_db_vm_cluster.test.hostname
  node_count                      = azurerm_oracle_exa_db_vm_cluster.test.node_count
  shape		                      = azurerm_oracle_exa_db_vm_cluster.test.shape
  ssh_public_keys                 = azurerm_oracle_exa_db_vm_cluster.test.ssh_public_keys
  subnet_id                       = azurerm_oracle_exa_db_vm_cluster.test.subnet_id
  total_ecpu_count                = azurerm_oracle_exa_db_vm_cluster.test.total_ecpu_count
  vm_file_system_storage {
    total_size_in_gbs = 440
  }	  virtual_network_id              = azurerm_oracle_exa_db_vm_cluster.test.virtual_network_id
  zones                                = azurerm_oracle_exascale_db_storage_vault.test.zones
}
`, a.basic(data))
}

func (a ExadbVmClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "virtual_network" {
  name                = "actvnet"
  address_space       = ["10.0.0.0/16"]
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "virtual_network_subnet" {
  name                 = "actsubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.virtual_network.name
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

resource "azurerm_oracle_exascale_db_storage_vault" "test" {
  name                = "OFakeacctest%[1]d"
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  display_name        = "OFakeacctest%[1]d"
  description         = "description"
  high_capacity_database_storage_input {
    total_size_in_gbs = 300
  }
  additional_flash_cache_in_percent = 100
  zones               = ["2"]
}

`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
