package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type NetworkPacketCaptureResource struct {
}

func testAccNetworkPacketCapture_localDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_packet_capture", "test")
	r := NetworkPacketCaptureResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.localDiskConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkPacketCapture_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_packet_capture", "test")
	r := NetworkPacketCaptureResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.localDiskConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.localDiskConfig_requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_network_packet_capture"),
		},
	})
}

func testAccNetworkPacketCapture_storageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_packet_capture", "test")
	r := NetworkPacketCaptureResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.storageAccountConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkPacketCapture_storageAccountAndLocalDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_packet_capture", "test")
	r := NetworkPacketCaptureResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.storageAccountAndLocalDiskConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkPacketCapture_withFilters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_packet_capture", "test")
	r := NetworkPacketCaptureResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.localDiskConfigWithFilters(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t NetworkPacketCaptureResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.PacketCaptureID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.PacketCapturesClient.Get(ctx, id.ResourceGroup, id.NetworkWatcherName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Network Packet Capture (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (NetworkPacketCaptureResource) base(data acceptance.TestData) string {
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
  address_prefix       = "10.0.2.0/24"
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
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osdisk"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

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

func (r NetworkPacketCaptureResource) localDiskConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_packet_capture" "test" {
  name                 = "acctestpc-%d"
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  target_resource_id   = azurerm_virtual_machine.test.id

  storage_location {
    file_path = "/var/captures/packet.cap"
  }

  depends_on = [azurerm_virtual_machine_extension.test]
}
`, r.base(data), data.RandomInteger)
}

func (r NetworkPacketCaptureResource) localDiskConfig_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_packet_capture" "import" {
  name                 = azurerm_network_packet_capture.test.name
  network_watcher_name = azurerm_network_packet_capture.test.network_watcher_name
  resource_group_name  = azurerm_network_packet_capture.test.resource_group_name
  target_resource_id   = azurerm_network_packet_capture.test.target_resource_id

  storage_location {
    file_path = "/var/captures/packet.cap"
  }

  depends_on = [azurerm_virtual_machine_extension.test]
}
`, r.localDiskConfig(data))
}

func (r NetworkPacketCaptureResource) localDiskConfigWithFilters(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_packet_capture" "test" {
  name                 = "acctestpc-%d"
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  target_resource_id   = azurerm_virtual_machine.test.id

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
`, r.base(data), data.RandomInteger)
}

func (r NetworkPacketCaptureResource) storageAccountConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_network_packet_capture" "test" {
  name                 = "acctestpc-%d"
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  target_resource_id   = azurerm_virtual_machine.test.id

  storage_location {
    storage_account_id = azurerm_storage_account.test.id
  }

  depends_on = [azurerm_virtual_machine_extension.test]
}
`, r.base(data), data.RandomString, data.RandomInteger)
}

func (r NetworkPacketCaptureResource) storageAccountAndLocalDiskConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_network_packet_capture" "test" {
  name                 = "acctestpc-%d"
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  target_resource_id   = azurerm_virtual_machine.test.id

  storage_location {
    file_path          = "/var/captures/packet.cap"
    storage_account_id = azurerm_storage_account.test.id
  }

  depends_on = [azurerm_virtual_machine_extension.test]
}
`, r.base(data), data.RandomString, data.RandomInteger)
}
