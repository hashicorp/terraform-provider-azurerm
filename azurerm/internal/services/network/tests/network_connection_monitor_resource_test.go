package tests

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
)

func testAccAzureRMNetworkConnectionMonitor_addressBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_basicAddressConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_basicAddressConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMNetworkConnectionMonitor_requiresImportConfig(data),
				ExpectError: acceptance.RequiresImportError("azurerm_network_connection_monitor"),
			},
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_addressComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_completeAddressConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_addressUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_basicAddressConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMNetworkConnectionMonitor_completeAddressConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_vmBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_basicVmConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_vmComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_completeVmConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_vmUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_basicVmConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMNetworkConnectionMonitor_completeVmConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_destinationUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_basicAddressConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMNetworkConnectionMonitor_basicVmConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMNetworkConnectionMonitor_basicAddressConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_missingDestination(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMNetworkConnectionMonitor_missingDestinationConfig(data),
				ExpectError: regexp.MustCompile("must have at least 2 endpoints"),
			},
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_conflictingDestinations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMNetworkConnectionMonitor_conflictingDestinationsConfig(data),
				ExpectError: regexp.MustCompile("don't allow creating different endpoints for the same VM"),
			},
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_withAddressAndVirtualMachineId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_withAddressAndVirtualMachineIdConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_httpConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_httpConfigurationConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_icmpConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_icmpConfigurationConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMNetworkConnectionMonitorExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ConnectionMonitorsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.NetworkConnectionMonitorID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.WatcherName, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on NetworkConnectionMonitorsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Connection Monitor does not exist: %s", id.Name)
		}

		return nil
	}
}

func testCheckAzureRMNetworkConnectionMonitorDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ConnectionMonitorsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_network_connection_monitor" {
			continue
		}

		id, err := parse.NetworkConnectionMonitorID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.WatcherName, id.Name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Connection Monitor still exists:%s", *resp.Name)
		}
	}

	return nil
}

func testAccAzureRMNetworkConnectionMonitor_baseConfig(data acceptance.TestData) string {
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
  address_prefix       = "10.0.2.0/24"
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

func testAccAzureRMNetworkConnectionMonitor_baseWithDestConfig(data acceptance.TestData) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseConfig(data)
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
`, config, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkConnectionMonitor_basicAddressConfig(data acceptance.TestData) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name               = "source"
    virtual_machine_id = azurerm_virtual_machine.src.id
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
    name                  = "testtg"
    destination_endpoints = ["destination"]
    sources               = ["source"]
    test_configurations   = ["tcp"]
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, config, data.RandomInteger)
}

func testAccAzureRMNetworkConnectionMonitor_completeAddressConfig(data acceptance.TestData) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseConfig(data)
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
    virtual_machine_id = azurerm_virtual_machine.src.id

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
    name                       = "tcp"
    protocol                   = "Tcp"
    test_frequency_iin_seconds = 40
    preferred_ip_version       = "IPv4"

    tcp_configuration {
      port = 80
    }

    success_threshold {
      checks_failed_percent = 50
      round_trip_time_ms    = 40
    }
  }

  test_group {
    name                  = "testtg"
    destination_endpoints = ["destination"]
    sources               = ["source"]
    test_configurations   = ["tcp"]
    enabled               = true
  }

  notes = "testNote"

  output_workspace_resource_ids = [azurerm_log_analytics_workspace.test.id]

  tags = {
    ENv = "Test"
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, config, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkConnectionMonitor_basicVmConfig(data acceptance.TestData) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseWithDestConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name               = "source"
    virtual_machine_id = azurerm_virtual_machine.src.id
  }

  endpoint {
    name               = "destination"
    virtual_machine_id = azurerm_virtual_machine.dest.id
  }

  test_configuration {
    name     = "tcp"
    protocol = "Tcp"

    tcp_configuration {
      port = 80
    }
  }

  test_group {
    name                  = "testtg"
    destination_endpoints = ["destination"]
    sources               = ["source"]
    test_configurations   = ["tcp"]
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, config, data.RandomInteger)
}

func testAccAzureRMNetworkConnectionMonitor_withAddressAndVirtualMachineIdConfig(data acceptance.TestData) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseWithDestConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name               = "source"
    virtual_machine_id = azurerm_virtual_machine.src.id
  }

  endpoint {
    name               = "destination"
    virtual_machine_id = azurerm_virtual_machine.dest.id
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
    name                  = "testtg"
    destination_endpoints = ["destination"]
    sources               = ["source"]
    test_configurations   = ["tcp"]
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, config, data.RandomInteger)
}

func testAccAzureRMNetworkConnectionMonitor_completeVmConfig(data acceptance.TestData) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseWithDestConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name               = "source"
    virtual_machine_id = azurerm_virtual_machine.src.id

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
    virtual_machine_id = azurerm_virtual_machine.dest.id
  }

  test_configuration {
    name                       = "tcp"
    protocol                   = "Tcp"
    test_frequency_iin_seconds = 40

    tcp_configuration {
      port = 80
    }
  }

  test_group {
    name                  = "testtg"
    destination_endpoints = ["destination"]
    sources               = ["source"]
    test_configurations   = ["tcp"]
    enabled               = true
  }

  tags = {
    ENv = "Test"
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, config, data.RandomInteger)
}

func testAccAzureRMNetworkConnectionMonitor_missingDestinationConfig(data acceptance.TestData) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name               = "source"
    virtual_machine_id = azurerm_virtual_machine.src.id
  }

  test_configuration {
    name     = "tcp"
    protocol = "Tcp"

    tcp_configuration {
      port = 80
    }
  }

  test_group {
    name                  = "testtg"
    destination_endpoints = ["destination"]
    sources               = ["source"]
    test_configurations   = ["tcp"]
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, config, data.RandomInteger)
}

func testAccAzureRMNetworkConnectionMonitor_conflictingDestinationsConfig(data acceptance.TestData) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name               = "source"
    virtual_machine_id = azurerm_virtual_machine.src.id
  }

  endpoint {
    name               = "destination"
    virtual_machine_id = azurerm_virtual_machine.src.id
  }

  test_configuration {
    name     = "tcp"
    protocol = "Tcp"

    tcp_configuration {
      port = 80
    }
  }

  test_group {
    name                  = "testtg"
    destination_endpoints = ["destination"]
    sources               = ["source"]
    test_configurations   = ["tcp"]
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, config, data.RandomInteger)
}

func testAccAzureRMNetworkConnectionMonitor_requiresImportConfig(data acceptance.TestData) string {
	config := testAccAzureRMNetworkConnectionMonitor_basicAddressConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "import" {
  name               = azurerm_network_connection_monitor.test.name
  network_watcher_id = azurerm_network_connection_monitor.test.network_watcher_id
  location           = azurerm_network_connection_monitor.test.location

  endpoint {
    name               = "source"
    virtual_machine_id = azurerm_virtual_machine.src.id
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
    name                  = "testtg"
    destination_endpoints = ["destination"]
    sources               = ["source"]
    test_configurations   = ["tcp"]
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, config)
}

func testAccAzureRMNetworkConnectionMonitor_httpConfigurationConfig(data acceptance.TestData) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name               = "source"
    virtual_machine_id = azurerm_virtual_machine.src.id
  }

  endpoint {
    name    = "destination"
    address = "terraform.io"
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
    name                  = "testtg"
    destination_endpoints = ["destination"]
    sources               = ["source"]
    test_configurations   = ["tcp"]
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, config, data.RandomInteger)
}

func testAccAzureRMNetworkConnectionMonitor_icmpConfigurationConfig(data acceptance.TestData) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name               = "acctest-CM-%d"
  network_watcher_id = azurerm_network_watcher.test.id
  location           = azurerm_network_watcher.test.location

  endpoint {
    name               = "source"
    virtual_machine_id = azurerm_virtual_machine.src.id
  }

  endpoint {
    name    = "destination"
    address = "terraform.io"
  }

  test_configuration {
    name     = "tcp"
    protocol = "Icmp"

    icmp_configuration {
      trace_route_disabled = false
    }
  }

  test_group {
    name                  = "testtg"
    destination_endpoints = ["destination"]
    sources               = ["source"]
    test_configurations   = ["tcp"]
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, config, data.RandomInteger)
}
