// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudvmclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CloudVmClusterResource struct{}

func (a CloudVmClusterResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := cloudvmclusters.ParseCloudVMClusterID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Oracle.OracleClient.CloudVMClusters.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestCloudVmClusterResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.CloudVmClusterResource{}.ResourceType(), "test")
	r := CloudVmClusterResource{}
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

func TestCloudVmClusterResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.CloudVmClusterResource{}.ResourceType(), "test")
	r := CloudVmClusterResource{}
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

func TestCloudVmClusterResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.CloudVmClusterResource{}.ResourceType(), "test")
	r := CloudVmClusterResource{}
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

func TestCloudVmClusterResource_update(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.CloudVmClusterResource{}.ResourceType(), "test")
	r := CloudVmClusterResource{}
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
		data.ImportStep(),
	})
}

func (a CloudVmClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
  %s
resource "azurerm_oracle_cloud_vm_cluster" "test" {
  location                        = "%[3]s"
  name                            = "OFakeVmacctest%[2]d"
  resource_group_name             = azurerm_resource_group.test.name
  cloud_exadata_infrastructure_id = azurerm_oracle_exadata_infrastructure.test.id
  cpu_core_count                  = 4
  data_storage_size_in_tbs        = 2
  db_node_storage_size_in_gbs     = 120
  db_servers                      = [for obj in data.azurerm_oracle_db_servers.test.db_servers : obj.ocid]
  display_name                    = "OFakeVmacctest%[2]d"
  gi_version                      = "23.0.0.0"
  license_model                   = "BringYourOwnLicense"
  memory_size_in_gbs              = 60
  hostname                        = "hostname"
  ssh_public_keys                 = ["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"]
  subnet_id                       = azurerm_subnet.virtual_network_subnet.id
  virtual_network_id              = azurerm_virtual_network.virtual_network.id
}`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a CloudVmClusterResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
  %s
resource "azurerm_oracle_cloud_vm_cluster" "test" {
  location                        = "%[3]s"
  name                            = "OFakeVmacctest%[2]d"
  resource_group_name             = azurerm_resource_group.test.name
  cloud_exadata_infrastructure_id = azurerm_oracle_exadata_infrastructure.test.id
  cpu_core_count                  = 4
  data_collection_options {
    diagnostics_events_enabled = true
    health_monitoring_enabled  = true
    incident_logs_enabled      = true
  }
  data_storage_size_in_tbs    = 2
  db_node_storage_size_in_gbs = 120
  db_servers                  = [for obj in data.azurerm_oracle_db_servers.test.db_servers : obj.ocid]
  display_name                = "OFakeVmacctest%[2]d"
  domain                      = "ociofakeacctes.com"
  gi_version                  = "23.0.0.0"
  local_backup_enabled        = true
  sparse_diskgroup_enabled    = true
  license_model               = "BringYourOwnLicense"
  memory_size_in_gbs          = 60
  hostname                    = "hostname"
  ssh_public_keys             = ["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"]
  subnet_id                   = azurerm_subnet.virtual_network_subnet.id
  scan_listener_port_tcp      = 1521
  scan_listener_port_tcp_ssl  = 2484
  system_version              = "23.1.23.0.0.250207"
  tags = {
    test = "testTag1"
  }
  time_zone          = "UTC"
  zone_id            = "ocid1.dns-zone.oc1.iad.aaaaaaaac7lyw74bnybmlek7nrsd5h3v5kjfv3aiw62menpuuwoder7yhmpa"
  virtual_network_id = azurerm_virtual_network.virtual_network.id
}`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a CloudVmClusterResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_oracle_cloud_vm_cluster" "test" {
  location                        = "%[3]s"
  name                            = "OFakeVmacctest%[2]d"
  resource_group_name             = azurerm_resource_group.test.name
  cloud_exadata_infrastructure_id = azurerm_oracle_exadata_infrastructure.test.id
  cpu_core_count                  = 4
  data_storage_size_in_tbs        = 2
  db_node_storage_size_in_gbs     = 120
  db_servers                      = [for obj in data.azurerm_oracle_db_servers.test.db_servers : obj.ocid]
  display_name                    = "OFakeVmacctest%[2]d"
  gi_version                      = "23.0.0.0"
  license_model                   = "BringYourOwnLicense"
  memory_size_in_gbs              = 60
  hostname                        = "hostname"
  ssh_public_keys                 = ["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"]
  subnet_id                       = azurerm_subnet.virtual_network_subnet.id
  tags = {
    test = "testTag1"
  }
  virtual_network_id = azurerm_virtual_network.virtual_network.id
}`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a CloudVmClusterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_cloud_vm_cluster" "import" {
  location                        = azurerm_oracle_cloud_vm_cluster.test.location
  name                            = azurerm_oracle_cloud_vm_cluster.test.name
  resource_group_name             = azurerm_oracle_cloud_vm_cluster.test.resource_group_name
  cloud_exadata_infrastructure_id = azurerm_oracle_cloud_vm_cluster.test.cloud_exadata_infrastructure_id
  cpu_core_count                  = azurerm_oracle_cloud_vm_cluster.test.cpu_core_count
  data_storage_size_in_tbs        = azurerm_oracle_cloud_vm_cluster.test.data_storage_size_in_tbs
  db_node_storage_size_in_gbs     = azurerm_oracle_cloud_vm_cluster.test.db_node_storage_size_in_gbs
  db_servers                      = azurerm_oracle_cloud_vm_cluster.test.db_servers
  display_name                    = azurerm_oracle_cloud_vm_cluster.test.display_name
  gi_version                      = azurerm_oracle_cloud_vm_cluster.test.gi_version
  license_model                   = azurerm_oracle_cloud_vm_cluster.test.license_model
  memory_size_in_gbs              = azurerm_oracle_cloud_vm_cluster.test.memory_size_in_gbs
  hostname                        = azurerm_oracle_cloud_vm_cluster.test.hostname
  ssh_public_keys                 = azurerm_oracle_cloud_vm_cluster.test.ssh_public_keys
  subnet_id                       = azurerm_oracle_cloud_vm_cluster.test.subnet_id
  virtual_network_id              = azurerm_oracle_cloud_vm_cluster.test.virtual_network_id
}
`, a.basic(data))
}

func (a CloudVmClusterResource) template(data acceptance.TestData) string {
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
  name                = "OFakeacctest%[1]d_vnet"
  address_space       = ["10.0.0.0/16"]
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "virtual_network_subnet" {
  name                 = "OFakeacctest%[1]d"
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

resource "azurerm_oracle_exadata_infrastructure" "test" {
  name                = "OFakeacctest%[1]d"
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  compute_count       = "2"
  display_name        = "OFakeacctest%[1]d"
  shape               = "Exadata.X9M"
  storage_count       = "3"
  zones               = ["3"]
}

data "azurerm_oracle_db_servers" "test" {
  resource_group_name               = azurerm_resource_group.test.name
  cloud_exadata_infrastructure_name = azurerm_oracle_exadata_infrastructure.test.name
}


`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
