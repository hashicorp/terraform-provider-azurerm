// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/packetcaptures"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VirtualMachineScaleSetPacketCaptureResource struct{}

func testAccVirtualMachineScaleSetPacketCapture_localDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_packet_capture", "test")
	r := VirtualMachineScaleSetPacketCaptureResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.localDiskConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccVirtualMachineScaleSetPacketCapture_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_packet_capture", "test")
	r := VirtualMachineScaleSetPacketCaptureResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.localDiskConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.localDiskConfig_requiresImport),
	})
}

func testAccVirtualMachineScaleSetPacketCapture_storageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_packet_capture", "test")
	r := VirtualMachineScaleSetPacketCaptureResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccountConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccVirtualMachineScaleSetPacketCapture_storageAccountAndLocalDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_packet_capture", "test")
	r := VirtualMachineScaleSetPacketCaptureResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccountAndLocalDiskConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccVirtualMachineScaleSetPacketCapture_withFilters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_packet_capture", "test")
	r := VirtualMachineScaleSetPacketCaptureResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.localDiskConfigWithFilters(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccVirtualMachineScaleSetPacketCapture_machineScope(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_packet_capture", "test")
	r := VirtualMachineScaleSetPacketCaptureResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.machineScope(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t VirtualMachineScaleSetPacketCaptureResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := packetcaptures.ParsePacketCaptureID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.PacketCaptures.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (VirtualMachineScaleSetPacketCaptureResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-watcher-%d"
  location = "%s"
}

resource "azurerm_network_watcher" "test" {
  name                = "acctestnw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                 = "acctestvmss-%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  sku                  = "Standard_F2"
  instances            = 4
  admin_username       = "adminuser"
  admin_password       = "P@ssword1234!"
  computer_name_prefix = "my-linux-computer-name-prefix"
  upgrade_mode         = "Automatic"

  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}

resource "azurerm_virtual_machine_scale_set_extension" "test" {
  name                         = "network-watcher"
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id
  publisher                    = "Microsoft.Azure.NetworkWatcher"
  type                         = "NetworkWatcherAgentLinux"
  type_handler_version         = "1.4"
  auto_upgrade_minor_version   = true
  automatic_upgrade_enabled    = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r VirtualMachineScaleSetPacketCaptureResource) localDiskConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_packet_capture" "test" {
  name                         = "acctestpc-%d"
  network_watcher_id           = azurerm_network_watcher.test.id
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id

  storage_location {
    file_path = "/var/captures/packet.cap"
  }

  depends_on = [azurerm_virtual_machine_scale_set_extension.test]
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualMachineScaleSetPacketCaptureResource) localDiskConfig_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_packet_capture" "import" {
  name                         = azurerm_virtual_machine_scale_set_packet_capture.test.name
  network_watcher_id           = azurerm_virtual_machine_scale_set_packet_capture.test.network_watcher_id
  virtual_machine_scale_set_id = azurerm_virtual_machine_scale_set_packet_capture.test.virtual_machine_scale_set_id

  storage_location {
    file_path = "/var/captures/packet.cap"
  }

  depends_on = [azurerm_virtual_machine_scale_set_extension.test]
}
`, r.localDiskConfig(data))
}

func (r VirtualMachineScaleSetPacketCaptureResource) localDiskConfigWithFilters(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_packet_capture" "test" {
  name                         = "acctestpc-%d"
  network_watcher_id           = azurerm_network_watcher.test.id
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id

  storage_location {
    file_path = "/var/captures/packet.cap"
  }

  filter {
    local_ip_address = "127.0.0.1"
    local_port       = "8080;9020;"
    protocol         = "TCP"
  }

  filter {
    local_ip_address = "127.0.0.1"
    local_port       = "80;443;"
    protocol         = "UDP"
  }

  depends_on = [azurerm_virtual_machine_scale_set_extension.test]
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualMachineScaleSetPacketCaptureResource) storageAccountConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_virtual_machine_scale_set_packet_capture" "test" {
  name                         = "acctestpc-%d"
  network_watcher_id           = azurerm_network_watcher.test.id
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id

  storage_location {
    storage_account_id = azurerm_storage_account.test.id
  }

  depends_on = [azurerm_virtual_machine_scale_set_extension.test]
}
`, r.template(data), data.RandomString, data.RandomInteger)
}

func (r VirtualMachineScaleSetPacketCaptureResource) storageAccountAndLocalDiskConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_virtual_machine_scale_set_packet_capture" "test" {
  name                         = "acctestpc-%d"
  network_watcher_id           = azurerm_network_watcher.test.id
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id

  storage_location {
    file_path          = "/var/captures/packet.cap"
    storage_account_id = azurerm_storage_account.test.id
  }

  depends_on = [azurerm_virtual_machine_scale_set_extension.test]
}
`, r.template(data), data.RandomString, data.RandomInteger)
}

func (r VirtualMachineScaleSetPacketCaptureResource) machineScope(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_scale_set_packet_capture" "test" {
  name                         = "acctestpc-%d"
  network_watcher_id           = azurerm_network_watcher.test.id
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.test.id

  storage_location {
    file_path = "/var/captures/packet.cap"
  }

  machine_scope {
    include_instance_ids = ["0", "1"]
    exclude_instance_ids = ["2", "3"]
  }

  depends_on = [azurerm_virtual_machine_scale_set_extension.test]
}
`, r.template(data), data.RandomInteger)
}
