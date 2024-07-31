// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/packetcaptures"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VirtualMachinePacketCaptureResource struct{}

func testAccVirtualMachinePacketCapture_localDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_packet_capture", "test")
	r := VirtualMachinePacketCaptureResource{}

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

func testAccVirtualMachinePacketCapture_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_packet_capture", "test")
	r := VirtualMachinePacketCaptureResource{}

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

func testAccVirtualMachinePacketCapture_storageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_packet_capture", "test")
	r := VirtualMachinePacketCaptureResource{}

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

func testAccVirtualMachinePacketCapture_storageAccountAndLocalDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_packet_capture", "test")
	r := VirtualMachinePacketCaptureResource{}

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

func testAccVirtualMachinePacketCapture_withFilters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_packet_capture", "test")
	r := VirtualMachinePacketCaptureResource{}

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

func (t VirtualMachinePacketCaptureResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := packetcaptures.ParsePacketCaptureID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.PacketCaptures.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (VirtualMachinePacketCaptureResource) template(data acceptance.TestData) string {
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

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_F2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osdisk"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  delete_os_disk_on_termination = true

  os_profile {
    computer_name  = "hostname%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_virtual_machine_extension" "test" {
  name                       = "network-watcher"
  virtual_machine_id         = azurerm_virtual_machine.test.id
  publisher                  = "Microsoft.Azure.NetworkWatcher"
  type                       = "NetworkWatcherAgentLinux"
  type_handler_version       = "1.4"
  auto_upgrade_minor_version = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r VirtualMachinePacketCaptureResource) localDiskConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_packet_capture" "test" {
  name               = "acctestpc-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  virtual_machine_id = azurerm_virtual_machine.test.id

  storage_location {
    file_path = "/var/captures/packet.cap"
  }

  depends_on = [azurerm_virtual_machine_extension.test]
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualMachinePacketCaptureResource) localDiskConfig_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_packet_capture" "import" {
  name               = azurerm_virtual_machine_packet_capture.test.name
  network_watcher_id = azurerm_virtual_machine_packet_capture.test.network_watcher_id
  virtual_machine_id = azurerm_virtual_machine_packet_capture.test.virtual_machine_id

  storage_location {
    file_path = "/var/captures/packet.cap"
  }

  depends_on = [azurerm_virtual_machine_extension.test]
}
`, r.localDiskConfig(data))
}

func (r VirtualMachinePacketCaptureResource) localDiskConfigWithFilters(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_packet_capture" "test" {
  name               = "acctestpc-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  virtual_machine_id = azurerm_virtual_machine.test.id

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

  depends_on = [azurerm_virtual_machine_extension.test]
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualMachinePacketCaptureResource) storageAccountConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_virtual_machine_packet_capture" "test" {
  name               = "acctestpc-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  virtual_machine_id = azurerm_virtual_machine.test.id

  storage_location {
    storage_account_id = azurerm_storage_account.test.id
  }

  depends_on = [azurerm_virtual_machine_extension.test]
}
`, r.template(data), data.RandomString, data.RandomInteger)
}

func (r VirtualMachinePacketCaptureResource) storageAccountAndLocalDiskConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_virtual_machine_packet_capture" "test" {
  name               = "acctestpc-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  virtual_machine_id = azurerm_virtual_machine.test.id

  storage_location {
    file_path          = "/var/captures/packet.cap"
    storage_account_id = azurerm_storage_account.test.id
  }

  depends_on = [azurerm_virtual_machine_extension.test]
}
`, r.template(data), data.RandomString, data.RandomInteger)
}
