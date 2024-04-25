// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/connectionmonitors"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetworkConnectionMonitorResource struct{}

const LetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const NumberBytes = "1234567890"
const SpecialBytes = "!@#$%^()"

func generateRandomPassword(n int) string {
	b := make([]byte, n)
	for i := range b {
		r := rand.Int()
		switch r % 3 {
		case 0:
			b[i] = LetterBytes[rand.Intn(len(LetterBytes))]
		case 1:
			b[i] = SpecialBytes[rand.Intn(len(SpecialBytes))]
		case 2:
			b[i] = NumberBytes[rand.Intn(len(NumberBytes))]
		}
	}
	return string(b)
}

func testAccNetworkConnectionMonitor_addressBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAddressConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAddressConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImportConfig(data),
			ExpectError: acceptance.RequiresImportError("azurerm_network_connection_monitor"),
		},
	})
}

func testAccNetworkConnectionMonitor_addressComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeAddressConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_addressUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAddressConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.completeAddressConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_vmBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicVmConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_vmComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeVmConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_vmUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicVmConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.completeVmConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_destinationUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAddressConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basicVmConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basicAddressConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_missingDestination(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config:      r.missingDestinationConfig(data),
			ExpectError: regexp.MustCompile("should have at least one destination"),
		},
	})
}

func testAccNetworkConnectionMonitor_conflictingDestinations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config:      r.conflictingDestinationsConfig(data),
			ExpectError: regexp.MustCompile("don't allow creating different endpoints for the same VM"),
		},
	})
}

func testAccNetworkConnectionMonitor_withAddressAndVirtualMachineId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.withAddressAndVirtualMachineIdConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_httpConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.httpConfigurationConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_icmpConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.icmpConfigurationConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_updateEndpointIPAddressAndCoverageLevel(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.endpointIPAddressAndCoverageLevel(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateEndpointIPAddressAndCoverageLevel(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_azureArcVM(t *testing.T) {
	if os.Getenv("CLIENT_SECRET") == "" {
		t.Skip("Skipping as CLIENT_SECRET is not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}
	randomUUID, _ := uuid.GenerateUUID()
	password := generateRandomPassword(10)

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureArcVM(data, randomUUID, password),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t NetworkConnectionMonitorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := connectionmonitors.ParseConnectionMonitorID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.ConnectionMonitors.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Network Connection Monitor (%s): %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (NetworkConnectionMonitorResource) baseConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-Watcher-%d"
  location = "%s"
}

resource "azurerm_network_watcher" "test" {
  name                = "acctest-Watcher-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-Vnet-%d"
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

resource "azurerm_network_interface" "src" {
  name                = "acctest-SrcNIC-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "src" {
  name                  = "acctest-SrcVM-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.src.id]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osdisk-src%d"
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

resource "azurerm_virtual_machine_extension" "src" {
  name                       = "acctest-VMExtension"
  virtual_machine_id         = azurerm_virtual_machine.src.id
  publisher                  = "Microsoft.Azure.NetworkWatcher"
  type                       = "NetworkWatcherAgentLinux"
  type_handler_version       = "1.4"
  auto_upgrade_minor_version = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r NetworkConnectionMonitorResource) baseWithDestConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "dest" {
  name                = "acctest-DestNic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "dest" {
  name                  = "acctest-DestVM-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.dest.id]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osdisk-dest%d"
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
`, r.baseConfig(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r NetworkConnectionMonitorResource) basicAddressConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name               = "source"
    target_resource_id = azurerm_virtual_machine.src.id
  }

  endpoint {
    name    = "destination"
    address = "pluginsdk.io"
  }

  test_configuration {
    name     = "tcp"
    protocol = "Tcp"

    tcp_configuration {
      port                      = 80
      destination_port_behavior = "None"
    }
  }

  test_group {
    name                     = "testtg"
    destination_endpoints    = ["destination"]
    source_endpoints         = ["source"]
    test_configuration_names = ["tcp"]
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, r.baseConfig(data), data.RandomInteger)
}

func (r NetworkConnectionMonitorResource) completeAddressConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-LAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name               = "source"
    target_resource_id = azurerm_virtual_machine.src.id

    filter {
      item {
        address = azurerm_virtual_machine.src.id
        type    = "AgentAddress"
      }

      type = "Include"
    }
  }

  endpoint {
    name    = "destination"
    address = "pluginsdk.io"
  }

  test_configuration {
    name                      = "tcp"
    protocol                  = "Tcp"
    test_frequency_in_seconds = 40
    preferred_ip_version      = "IPv4"

    tcp_configuration {
      port                      = 80
      destination_port_behavior = "ListenIfAvailable"
    }

    success_threshold {
      checks_failed_percent = 50
      round_trip_time_ms    = 40
    }
  }

  test_group {
    name                     = "testtg"
    destination_endpoints    = ["destination"]
    source_endpoints         = ["source"]
    test_configuration_names = ["tcp"]
    enabled                  = true
  }

  notes = "testNote"

  output_workspace_resource_ids = [azurerm_log_analytics_workspace.test.id]

  tags = {
    ENv = "Test"
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, r.baseConfig(data), data.RandomInteger, data.RandomInteger)
}

func (r NetworkConnectionMonitorResource) basicVmConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name               = "source"
    target_resource_id = azurerm_virtual_machine.src.id
  }

  endpoint {
    name               = "destination"
    target_resource_id = azurerm_virtual_machine.dest.id
  }

  test_configuration {
    name     = "tcp"
    protocol = "Tcp"

    tcp_configuration {
      port = 80
    }
  }

  test_group {
    name                     = "testtg"
    destination_endpoints    = ["destination"]
    source_endpoints         = ["source"]
    test_configuration_names = ["tcp"]
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, r.baseWithDestConfig(data), data.RandomInteger)
}

func (r NetworkConnectionMonitorResource) withAddressAndVirtualMachineIdConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name               = "source"
    target_resource_id = azurerm_virtual_machine.src.id
  }

  endpoint {
    name               = "destination"
    target_resource_id = azurerm_virtual_machine.dest.id
    address            = azurerm_network_interface.dest.private_ip_address
  }

  test_configuration {
    name     = "tcp"
    protocol = "Tcp"

    tcp_configuration {
      port = 80
    }
  }

  test_group {
    name                     = "testtg"
    destination_endpoints    = ["destination"]
    source_endpoints         = ["source"]
    test_configuration_names = ["tcp"]
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, r.baseWithDestConfig(data), data.RandomInteger)
}

func (r NetworkConnectionMonitorResource) completeVmConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name               = "source"
    target_resource_id = azurerm_virtual_machine.src.id

    filter {
      item {
        address = azurerm_virtual_machine.src.id
        type    = "AgentAddress"
      }

      type = "Include"
    }
  }

  endpoint {
    name               = "destination"
    target_resource_id = azurerm_virtual_machine.dest.id
  }

  test_configuration {
    name                      = "tcp"
    protocol                  = "Tcp"
    test_frequency_in_seconds = 40

    tcp_configuration {
      port = 80
    }
  }

  test_group {
    name                     = "testtg"
    destination_endpoints    = ["destination"]
    source_endpoints         = ["source"]
    test_configuration_names = ["tcp"]
    enabled                  = true
  }

  tags = {
    ENv = "Test"
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, r.baseWithDestConfig(data), data.RandomInteger)
}

func (r NetworkConnectionMonitorResource) missingDestinationConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name               = "source"
    target_resource_id = azurerm_virtual_machine.src.id
  }

  test_configuration {
    name     = "tcp"
    protocol = "Tcp"

    tcp_configuration {
      port = 80
    }
  }

  test_group {
    name                     = "testtg"
    destination_endpoints    = []
    source_endpoints         = ["source"]
    test_configuration_names = ["tcp"]
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, r.baseConfig(data), data.RandomInteger)
}

func (r NetworkConnectionMonitorResource) conflictingDestinationsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name               = "source"
    target_resource_id = azurerm_virtual_machine.src.id
  }

  endpoint {
    name               = "destination"
    target_resource_id = azurerm_virtual_machine.src.id
  }

  test_configuration {
    name     = "tcp"
    protocol = "Tcp"

    tcp_configuration {
      port = 80
    }
  }

  test_group {
    name                     = "testtg"
    destination_endpoints    = ["destination"]
    source_endpoints         = ["source"]
    test_configuration_names = ["tcp"]
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, r.baseConfig(data), data.RandomInteger)
}

func (r NetworkConnectionMonitorResource) requiresImportConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "import" {
  name               = azurerm_network_connection_monitor.test.name
  network_watcher_id = azurerm_network_connection_monitor.test.network_watcher_id
  location           = azurerm_network_connection_monitor.test.location

  endpoint {
    name               = "source"
    target_resource_id = azurerm_virtual_machine.src.id
  }

  endpoint {
    name    = "destination"
    address = "pluginsdk.io"
  }

  test_configuration {
    name     = "tcp"
    protocol = "Tcp"

    tcp_configuration {
      port = 80
    }
  }

  test_group {
    name                     = "testtg"
    destination_endpoints    = ["destination"]
    source_endpoints         = ["source"]
    test_configuration_names = ["tcp"]
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, r.basicAddressConfig(data))
}

func (r NetworkConnectionMonitorResource) httpConfigurationConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name               = "source"
    target_resource_id = azurerm_virtual_machine.src.id
  }

  endpoint {
    name               = "destination"
    target_resource_id = azurerm_virtual_machine.dest.id
  }

  test_configuration {
    name     = "tcp"
    protocol = "Http"

    http_configuration {
      method                   = "Get"
      port                     = 80
      path                     = "/a/b"
      prefer_https             = false
      valid_status_code_ranges = ["200"]

      request_header {
        name  = "testHeader"
        value = "testVal"
      }
    }
  }

  test_group {
    name                     = "testtg"
    destination_endpoints    = ["destination"]
    source_endpoints         = ["source"]
    test_configuration_names = ["tcp"]
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, r.baseWithDestConfig(data), data.RandomInteger)
}

func (r NetworkConnectionMonitorResource) icmpConfigurationConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name               = "source"
    target_resource_id = azurerm_virtual_machine.src.id
  }

  endpoint {
    name    = "destination"
    address = "pluginsdk.io"
  }

  test_configuration {
    name     = "tcp"
    protocol = "Icmp"

    icmp_configuration {
      trace_route_enabled = true
    }
  }

  test_group {
    name                     = "testtg"
    destination_endpoints    = ["destination"]
    source_endpoints         = ["source"]
    test_configuration_names = ["tcp"]
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, r.baseConfig(data), data.RandomInteger)
}

func (r NetworkConnectionMonitorResource) endpointIPAddressAndCoverageLevel(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name                  = "source"
    target_resource_type  = "AzureVNet"
    target_resource_id    = azurerm_virtual_network.test.id
    included_ip_addresses = azurerm_subnet.test.address_prefixes
    excluded_ip_addresses = ["10.0.2.2", "10.0.2.3"]
    coverage_level        = "Default"
  }

  endpoint {
    address = "pluginsdk.io"
    name    = "destination"
  }

  test_configuration {
    name     = "tcp"
    protocol = "Tcp"

    tcp_configuration {
      port = 80
    }
  }

  test_group {
    name                     = "testtg"
    destination_endpoints    = ["destination"]
    source_endpoints         = ["source"]
    test_configuration_names = ["tcp"]
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, r.baseConfig(data), data.RandomInteger)
}

func (r NetworkConnectionMonitorResource) updateEndpointIPAddressAndCoverageLevel(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test2" {
  name                 = "accttest-Subnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.3.0/24"]
}

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name                  = "source"
    target_resource_type  = "AzureVNet"
    target_resource_id    = azurerm_virtual_network.test.id
    included_ip_addresses = azurerm_subnet.test2.address_prefixes
    excluded_ip_addresses = ["10.0.3.2"]
    coverage_level        = "Average"
  }

  endpoint {
    address = "pluginsdk.io"
    name    = "destination"
  }

  test_configuration {
    name     = "tcp"
    protocol = "Tcp"

    tcp_configuration {
      port = 80
    }
  }

  test_group {
    name                     = "testtg"
    destination_endpoints    = ["destination"]
    source_endpoints         = ["source"]
    test_configuration_names = ["tcp"]
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, r.baseConfig(data), data.RandomInteger, data.RandomInteger)
}

func (r NetworkConnectionMonitorResource) azureArcVM(data acceptance.TestData, randomUUID string, password string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

data "azurerm_client_config" "current" {}

# note: real-life usage prefer random_uuid resource in registry.terraform.io/hashicorp/random
locals {
  random_uuid = "%s"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-monitor-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
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
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_network_security_group" "my_terraform_nsg" {
  name                = "myNetworkSecurityGroup"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule {
    name                       = "SSH"
    priority                   = 1001
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

resource "azurerm_network_interface_security_group_association" "example" {
  network_interface_id      = azurerm_network_interface.test.id
  network_security_group_id = azurerm_network_security_group.my_terraform_nsg.id
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctestVM-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  admin_password                  = "%s"
  provision_vm_agent              = false
  allow_extension_operations      = false
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  connection {
    type     = "ssh"
    host     = azurerm_public_ip.test.ip_address
    user     = "adminuser"
    password = "%s"
  }

  provisioner "file" {
    content = templatefile("scripts/install_arc.sh.tftpl", {
      resource_group_name = azurerm_resource_group.test.name
      uuid                = local.random_uuid
      location            = azurerm_resource_group.test.location
      tenant_id           = data.azurerm_client_config.current.tenant_id
      client_id           = data.azurerm_client_config.current.client_id
      client_secret       = "%s"
      subscription_id     = data.azurerm_client_config.current.subscription_id
    })
    destination = "/home/adminuser/install_arc_agent.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo apt-get install -y python-ctypes",
      "sudo sed -i 's/\r$//' /home/adminuser/install_arc_agent.sh",
      "sudo chmod +x /home/adminuser/install_arc_agent.sh",
      "bash /home/adminuser/install_arc_agent.sh",
    ]
  }
}

data "azurerm_arc_machine" "test" {
  name                = azurerm_linux_virtual_machine.test.name
  resource_group_name = azurerm_resource_group.test.name
  depends_on = [
    azurerm_linux_virtual_machine.test
  ]
}

resource "azurerm_arc_machine_extension" "test" {
  name                      = "acctest-hcme-%d"
  arc_machine_id            = data.azurerm_arc_machine.test.id
  publisher                 = "Microsoft.Azure.NetworkWatcher"
  type                      = "NetworkWatcherAgentLinux"
  type_handler_version      = "1.4"
  automatic_upgrade_enabled = true
  location                  = data.azurerm_arc_machine.test.location
}

resource "azurerm_network_watcher" "test" {
  name                = "acctest-Watcher-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name                 = "source"
    target_resource_type = "AzureArcVM"
    target_resource_id   = data.azurerm_arc_machine.test.id
  }

  endpoint {
    address = "pluginsdk.io"
    name    = "destination"
  }

  test_configuration {
    name                      = "IcmpConnection"
    protocol                  = "Icmp"
    test_frequency_in_seconds = 60

    icmp_configuration {
      trace_route_enabled = true
    }
  }

  test_group {
    name                     = "testtg"
    destination_endpoints    = ["destination"]
    source_endpoints         = ["source"]
    test_configuration_names = ["IcmpConnection"]
  }

  depends_on = [azurerm_arc_machine_extension.test]
}
`, randomUUID, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, password, password, os.Getenv("CLIENT_SECRET"), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
