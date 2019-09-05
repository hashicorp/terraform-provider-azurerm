package azurerm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func testAccAzureRMNetworkConnectionMonitor_addressBasic(t *testing.T) {
	resourceName := "azurerm_network_connection_monitor.test"

	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_basicAddressConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(resourceName, "location", azure.NormalizeLocation(location)),
					resource.TestCheckResourceAttr(resourceName, "auto_start", "true"),
					resource.TestCheckResourceAttr(resourceName, "interval_in_seconds", "60"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_network_connection_monitor.test"

	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_basicAddressConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMNetworkConnectionMonitor_requiresImportConfig(ri, location),
				ExpectError: testRequiresImportError("azurerm_network_connection_monitor"),
			},
		},
	})

}

func testAccAzureRMNetworkConnectionMonitor_addressComplete(t *testing.T) {
	resourceName := "azurerm_network_connection_monitor.test"

	ri := tf.AccRandTimeInt()
	location := testLocation()
	autoStart := "false"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_completeAddressConfig(ri, location, autoStart),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_start", "false"),
					resource.TestCheckResourceAttr(resourceName, "interval_in_seconds", "30"),
					resource.TestCheckResourceAttr(resourceName, "source.0.port", "20020"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_addressUpdate(t *testing.T) {
	resourceName := "azurerm_network_connection_monitor.test"

	ri := tf.AccRandTimeInt()
	location := testLocation()
	autoStart := "true"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_basicAddressConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMNetworkConnectionMonitor_completeAddressConfig(ri, location, autoStart),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_start", "true"),
					resource.TestCheckResourceAttr(resourceName, "interval_in_seconds", "30"),
					resource.TestCheckResourceAttr(resourceName, "source.0.port", "20020"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_vmBasic(t *testing.T) {
	resourceName := "azurerm_network_connection_monitor.test"

	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_basicVmConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(resourceName, "location", azure.NormalizeLocation(location)),
					resource.TestCheckResourceAttr(resourceName, "auto_start", "true"),
					resource.TestCheckResourceAttr(resourceName, "interval_in_seconds", "60"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_vmComplete(t *testing.T) {
	resourceName := "azurerm_network_connection_monitor.test"

	ri := tf.AccRandTimeInt()
	location := testLocation()
	autoStart := "false"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_completeVmConfig(ri, location, autoStart),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_start", "false"),
					resource.TestCheckResourceAttr(resourceName, "interval_in_seconds", "30"),
					resource.TestCheckResourceAttr(resourceName, "source.0.port", "20020"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_vmUpdate(t *testing.T) {
	resourceName := "azurerm_network_connection_monitor.test"

	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_basicVmConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMNetworkConnectionMonitor_completeVmConfig(ri, location, "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_start", "true"),
					resource.TestCheckResourceAttr(resourceName, "interval_in_seconds", "30"),
					resource.TestCheckResourceAttr(resourceName, "source.0.port", "20020"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_destinationUpdate(t *testing.T) {
	resourceName := "azurerm_network_connection_monitor.test"

	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkConnectionMonitor_basicAddressConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "destination.0.address"),
				),
			},
			{
				Config: testAccAzureRMNetworkConnectionMonitor_basicVmConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "destination.0.virtual_machine_id"),
				),
			},
			{
				Config: testAccAzureRMNetworkConnectionMonitor_basicAddressConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkConnectionMonitorExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "destination.0.address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_missingDestination(t *testing.T) {
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMNetworkConnectionMonitor_missingDestinationConfig(ri, location),
				ExpectError: regexp.MustCompile("Error: either `destination.virtual_machine_id` or `destination.address` must be specified"),
			},
		},
	})
}

func testAccAzureRMNetworkConnectionMonitor_conflictingDestinations(t *testing.T) {
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkConnectionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMNetworkConnectionMonitor_conflictingDestinationsConfig(ri, location),
				ExpectError: regexp.MustCompile("conflicts with destination.0.address"),
			},
		},
	})
}

func testCheckAzureRMNetworkConnectionMonitorExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ArmClient).network.ConnectionMonitorsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
	client := testAccProvider.Meta().(*ArmClient).network.ConnectionMonitorsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMNetworkConnectionMonitor_baseConfig(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_watcher" "test" {
  name                = "acctnw-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "src" {
  name                = "acctni-src%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "src" {
  name                  = "acctvm-src%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.src.id}"]
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
  location                   = "${azurerm_resource_group.test.location}"
  resource_group_name        = "${azurerm_resource_group.test.name}"
  virtual_machine_name       = "${azurerm_virtual_machine.src.name}"
  publisher                  = "Microsoft.Azure.NetworkWatcher"
  type                       = "NetworkWatcherAgentLinux"
  type_handler_version       = "1.4"
  auto_upgrade_minor_version = true
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMNetworkConnectionMonitor_baseWithDestConfig(rInt int, location string) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseConfig(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "dest" {
  name                = "acctni-dest%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "dest" {
  name                  = "acctvm-dest%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.dest.id}"]
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
`, config, rInt, rInt, rInt, rInt)
}

func testAccAzureRMNetworkConnectionMonitor_basicAddressConfig(rInt int, location string) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseConfig(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name                 = "acctestcm-%d"
  network_watcher_name = "${azurerm_network_watcher.test.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_network_watcher.test.location}"

  source {
    virtual_machine_id = "${azurerm_virtual_machine.src.id}"
  }

  destination {
    address = "terraform.io"
    port    = 80
  }

  depends_on = ["azurerm_virtual_machine_extension.src"]
}
`, config, rInt)
}

func testAccAzureRMNetworkConnectionMonitor_completeAddressConfig(rInt int, location, autoStart string) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseConfig(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name                 = "acctestcm-%d"
  network_watcher_name = "${azurerm_network_watcher.test.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_network_watcher.test.location}"

  auto_start          = %s
  interval_in_seconds = 30

  source {
    virtual_machine_id = "${azurerm_virtual_machine.src.id}"
    port               = 20020
  }

  destination {
    address = "terraform.io"
    port    = 443
  }

  tags = {
    env = "test"
  }

  depends_on = ["azurerm_virtual_machine_extension.src"]
}
`, config, rInt, autoStart)
}

func testAccAzureRMNetworkConnectionMonitor_basicVmConfig(rInt int, location string) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseWithDestConfig(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name                 = "acctestcm-%d"
  network_watcher_name = "${azurerm_network_watcher.test.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_network_watcher.test.location}"

  source {
    virtual_machine_id = "${azurerm_virtual_machine.src.id}"
  }

  destination {
    virtual_machine_id = "${azurerm_virtual_machine.dest.id}"
    port               = 80
  }

  depends_on = ["azurerm_virtual_machine_extension.src"]
}
`, config, rInt)
}

func testAccAzureRMNetworkConnectionMonitor_completeVmConfig(rInt int, location, autoStart string) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseWithDestConfig(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name                 = "acctestcm-%d"
  network_watcher_name = "${azurerm_network_watcher.test.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_network_watcher.test.location}"

  auto_start          = %s
  interval_in_seconds = 30

  source {
    virtual_machine_id = "${azurerm_virtual_machine.src.id}"
    port               = 20020
  }

  destination {
    virtual_machine_id = "${azurerm_virtual_machine.dest.id}"
    port               = 443
  }

  tags = {
    env = "test"
  }

  depends_on = ["azurerm_virtual_machine_extension.src"]
}
`, config, rInt, autoStart)
}

func testAccAzureRMNetworkConnectionMonitor_missingDestinationConfig(rInt int, location string) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseConfig(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name                 = "acctestcm-%d"
  network_watcher_name = "${azurerm_network_watcher.test.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_network_watcher.test.location}"

  source {
    virtual_machine_id = "${azurerm_virtual_machine.src.id}"
  }

  destination {
    port = 80
  }

  depends_on = ["azurerm_virtual_machine_extension.src"]
}
`, config, rInt)
}

func testAccAzureRMNetworkConnectionMonitor_conflictingDestinationsConfig(rInt int, location string) string {
	config := testAccAzureRMNetworkConnectionMonitor_baseConfig(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "test" {
  name                 = "acctestcm-%d"
  network_watcher_name = "${azurerm_network_watcher.test.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_network_watcher.test.location}"

  source {
    virtual_machine_id = "${azurerm_virtual_machine.src.id}"
  }

  destination {
    address            = "terraform.io"
    virtual_machine_id = "${azurerm_virtual_machine.src.id}"
    port               = 80
  }

  depends_on = ["azurerm_virtual_machine_extension.src"]
}
`, config, rInt)
}

func testAccAzureRMNetworkConnectionMonitor_requiresImportConfig(rInt int, location string) string {
	config := testAccAzureRMNetworkConnectionMonitor_basicAddressConfig(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_network_connection_monitor" "import" {
  name                 = "${azurerm_network_connection_monitor.test.name}"
  network_watcher_name = "${azurerm_network_connection_monitor.test.network_watcher_name}"
  resource_group_name  = "${azurerm_network_connection_monitor.test.resource_group_name}"
  location             = "${azurerm_network_connection_monitor.test.location}"

  source {
    virtual_machine_id = "${azurerm_virtual_machine.src.id}"
  }

  destination {
    address = "terraform.io"
    port    = 80
  }

  depends_on = ["azurerm_virtual_machine_extension.src"]
}
`, config)
}
