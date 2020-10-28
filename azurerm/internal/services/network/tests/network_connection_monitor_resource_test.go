package tests

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
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
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "location", azure.NormalizeLocation(data.Locations.Primary)),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_start", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval_in_seconds", "60"),
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

	autoStart := "false"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_completeAddressConfig(data, autoStart),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_start", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval_in_seconds", "30"),
					resource.TestCheckResourceAttr(data.ResourceName, "source.0.port", "20020"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.env", "test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_addressUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")

	autoStart := "true"

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
				Config: testAccAzureRMNetworkConnectionMonitor_completeAddressConfig(data, autoStart),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_start", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval_in_seconds", "30"),
					resource.TestCheckResourceAttr(data.ResourceName, "source.0.port", "20020"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.env", "test"),
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
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "location", azure.NormalizeLocation(data.Locations.Primary)),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_start", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval_in_seconds", "60"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_vmComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_connection_monitor", "test")

	autoStart := "false"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_completeVmConfig(data, autoStart),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_start", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval_in_seconds", "30"),
					resource.TestCheckResourceAttr(data.ResourceName, "source.0.port", "20020"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.env", "test"),
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
				Config: testAccAzureRMNetworkConnectionMonitor_completeVmConfig(data, "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_start", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "interval_in_seconds", "30"),
					resource.TestCheckResourceAttr(data.ResourceName, "source.0.port", "20020"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.env", "test"),
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
					resource.TestCheckResourceAttrSet(data.ResourceName, "destination.0.address"),
				),
			},
			{
				Config: testAccAzureRMNetworkConnectionMonitor_basicVmConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "destination.0.virtual_machine_id"),
				),
			},
			{
				Config: testAccAzureRMNetworkConnectionMonitor_basicAddressConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "destination.0.address"),
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
				ExpectError: regexp.MustCompile("Error: either `destination.virtual_machine_id` or `destination.address` must be specified"),
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
				ExpectError: regexp.MustCompile("conflicts with destination.0.address"),
			},
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

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		watcherName := rs.Primary.Attributes["network_watcher_name"]
		NetworkConnectionMonitorName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, watcherName, NetworkConnectionMonitorName)
		if err != nil {
			return fmt.Errorf("Bad: Get on NetworkConnectionMonitorsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Connection Monitor does not exist: %s", NetworkConnectionMonitorName)
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

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		watcherName := rs.Primary.Attributes["network_watcher_name"]
		NetworkConnectionMonitorName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, watcherName, NetworkConnectionMonitorName)

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
  name     = "acctestRG-watcher-%d"
  location = "%s"
}

resource "azurerm_network_watcher" "test" {
  name                = "acctnw-%d"
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
  name               = "internal"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.0.2.0/24"
}

resource "azurerm_network_interface" "src" {
  name                = "acctni-src%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "src" {
  name                  = "acctvm-src%d"
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
  name                       = "network-watcher"
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
  name                = "acctni-dest%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "dest" {
  name                  = "acctvm-dest%d"
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
  name                 = "acctestcm-%d"
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_network_watcher.test.location

  source {
    virtual_machine_id = azurerm_virtual_machine.src.id
  }

  destination {
    address = "terraform.io"
    port    = 80
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, config, data.RandomInteger)
}

func testAccAzureRMNetworkConnectionMonitor_completeAddressConfig(data acceptance.TestData, autoStart string) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name                 = "acctestcm-%d"
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_network_watcher.test.location

  auto_start          = %s
  interval_in_seconds = 30

  source {
    virtual_machine_id = azurerm_virtual_machine.src.id
    port               = 20020
  }

  destination {
    address = "terraform.io"
    port    = 443
  }

  tags = {
    env = "test"
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, config, data.RandomInteger, autoStart)
}

func testAccAzureRMNetworkConnectionMonitor_basicVmConfig(data acceptance.TestData) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseWithDestConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name                 = "acctestcm-%d"
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_network_watcher.test.location

  source {
    virtual_machine_id = azurerm_virtual_machine.src.id
  }

  destination {
    virtual_machine_id = azurerm_virtual_machine.dest.id
    port               = 80
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, config, data.RandomInteger)
}

func testAccAzureRMNetworkConnectionMonitor_completeVmConfig(data acceptance.TestData, autoStart string) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseWithDestConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name                 = "acctestcm-%d"
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_network_watcher.test.location

  auto_start          = %s
  interval_in_seconds = 30

  source {
    virtual_machine_id = azurerm_virtual_machine.src.id
    port               = 20020
  }

  destination {
    virtual_machine_id = azurerm_virtual_machine.dest.id
    port               = 443
  }

  tags = {
    env = "test"
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, config, data.RandomInteger, autoStart)
}

func testAccAzureRMNetworkConnectionMonitor_missingDestinationConfig(data acceptance.TestData) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name                 = "acctestcm-%d"
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_network_watcher.test.location

  source {
    virtual_machine_id = azurerm_virtual_machine.src.id
  }

  destination {
    port = 80
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
  name                 = "acctestcm-%d"
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_network_watcher.test.location

  source {
    virtual_machine_id = azurerm_virtual_machine.src.id
  }

  destination {
    address            = "terraform.io"
    virtual_machine_id = azurerm_virtual_machine.src.id
    port               = 80
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
  name                 = azurerm_network_connection_monitor.test.name
  network_watcher_name = azurerm_network_connection_monitor.test.network_watcher_name
  resource_group_name  = azurerm_network_connection_monitor.test.resource_group_name
  location             = azurerm_network_connection_monitor.test.location

  source {
    virtual_machine_id = azurerm_virtual_machine.src.id
  }

  destination {
    address = "terraform.io"
    port    = 80
  }

  depends_on = [azurerm_virtual_machine_extension.src]
}
`, config)
}
