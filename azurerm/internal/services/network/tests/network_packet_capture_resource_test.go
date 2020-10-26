package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func testAccAzureRMNetworkPacketCapture_localDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_packet_capture", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkPacketCaptureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNetworkPacketCapture_localDiskConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkPacketCaptureExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkPacketCapture_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_packet_capture", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkPacketCaptureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNetworkPacketCapture_localDiskConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkPacketCaptureExists(data.ResourceName),
				),
			},
			{
				Config:      testAzureRMNetworkPacketCapture_localDiskConfig_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_network_packet_capture"),
			},
		},
	})
}
func testAccAzureRMNetworkPacketCapture_storageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_packet_capture", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkPacketCaptureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNetworkPacketCapture_storageAccountConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkPacketCaptureExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkPacketCapture_storageAccountAndLocalDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_packet_capture", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkPacketCaptureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNetworkPacketCapture_storageAccountAndLocalDiskConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkPacketCaptureExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkPacketCapture_withFilters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_packet_capture", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkPacketCaptureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMNetworkPacketCapture_localDiskConfigWithFilters(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkPacketCaptureExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMNetworkPacketCaptureExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.PacketCapturesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		watcherName := rs.Primary.Attributes["network_watcher_name"]
		NetworkPacketCaptureName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, watcherName, NetworkPacketCaptureName)
		if err != nil {
			return fmt.Errorf("Bad: Get on NetworkPacketCapturesClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Packet Capture does not exist: %s", NetworkPacketCaptureName)
		}

		return nil
	}
}

func testCheckAzureRMNetworkPacketCaptureDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.PacketCapturesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_network_packet_capture" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		watcherName := rs.Primary.Attributes["network_watcher_name"]
		NetworkPacketCaptureName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, watcherName, NetworkPacketCaptureName)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Packet Capture still exists:%s", *resp.Name)
		}
	}

	return nil
}

func testAzureRMNetworkPacketCapture_base(data acceptance.TestData) string {
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

func testAzureRMNetworkPacketCapture_localDiskConfig(data acceptance.TestData) string {
	config := testAzureRMNetworkPacketCapture_base(data)
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
`, config, data.RandomInteger)
}

func testAzureRMNetworkPacketCapture_localDiskConfig_requiresImport(data acceptance.TestData) string {
	config := testAzureRMNetworkPacketCapture_localDiskConfig(data)
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
`, config)
}

func testAzureRMNetworkPacketCapture_localDiskConfigWithFilters(data acceptance.TestData) string {
	config := testAzureRMNetworkPacketCapture_base(data)
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
`, config, data.RandomInteger)
}

func testAzureRMNetworkPacketCapture_storageAccountConfig(data acceptance.TestData) string {
	config := testAzureRMNetworkPacketCapture_base(data)
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
`, config, data.RandomString, data.RandomInteger)
}

func testAzureRMNetworkPacketCapture_storageAccountAndLocalDiskConfig(data acceptance.TestData) string {
	config := testAzureRMNetworkPacketCapture_base(data)
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
`, config, data.RandomString, data.RandomInteger)
}
