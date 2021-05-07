package network_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type NetworkConnectionMonitorResource struct {
}

func testAccNetworkConnectionMonitor_addressBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicAddressConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicAddressConfig(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.completeAddressConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_addressUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicAddressConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.completeAddressConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_vmBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicVmConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_vmComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.completeVmConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_vmUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicVmConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.completeVmConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_destinationUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicAddressConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basicVmConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basicAddressConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_missingDestination(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.missingDestinationConfig(data),
			ExpectError: regexp.MustCompile("should have at least one destination"),
		},
	})
}

func testAccNetworkConnectionMonitor_conflictingDestinations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.conflictingDestinationsConfig(data),
			ExpectError: regexp.MustCompile("don't allow creating different endpoints for the same VM"),
		},
	})
}

func testAccNetworkConnectionMonitor_withAddressAndVirtualMachineId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withAddressAndVirtualMachineIdConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_httpConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.httpConfigurationConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_icmpConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.icmpConfigurationConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_endpointDeprecated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.endpointDeprecated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicAddressConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.endpointDeprecated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkConnectionMonitor_updateEndpointIPAddressAndCoverageLevel(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")
	r := NetworkConnectionMonitorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.endpointIPAddressAndCoverageLevel(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateEndpointIPAddressAndCoverageLevel(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t NetworkConnectionMonitorResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ConnectionMonitorID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.ConnectionMonitorsClient.Get(ctx, id.ResourceGroup, id.NetworkWatcherName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Network Connection Monitor (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
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
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osdisk-src%d"
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
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osdisk-dest%d"
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
    address = "terraform.io"
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

func (r NetworkConnectionMonitorResource) completeAddressConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-LAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "pergb2018"
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
    address = "terraform.io"
  }

  test_configuration {
    name                      = "tcp"
    protocol                  = "Tcp"
    test_frequency_in_seconds = 40
    preferred_ip_version      = "IPv4"

    tcp_configuration {
      port = 80
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
    address = "terraform.io"
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
    address = "terraform.io"
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

func (r NetworkConnectionMonitorResource) endpointDeprecated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name                 = "source"
    target_resource_type = "MMAWorkspaceMachine"
    address              = "test.internal.domain.com"
    target_resource_id   = azurerm_log_analytics_workspace.test.id
  }

  endpoint {
    name    = "destination"
    address = "terraform.io"
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
    address = "terraform.io"
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
    address = "terraform.io"
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
