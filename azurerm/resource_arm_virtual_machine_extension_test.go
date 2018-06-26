package azurerm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMVirtualMachineExtension_basic(t *testing.T) {
	resourceName := "azurerm_virtual_machine_extension.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMVirtualMachineExtension_basic(ri, location)
	postConfig := testAccAzureRMVirtualMachineExtension_basicUpdate(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineExtensionExists(resourceName),
					resource.TestMatchResourceAttr(resourceName, "settings", regexp.MustCompile("hostname")),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineExtensionExists(resourceName),
					resource.TestMatchResourceAttr(resourceName, "settings", regexp.MustCompile("whoami")),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineExtension_concurrent(t *testing.T) {
	firstResourceName := "azurerm_virtual_machine_extension.test"
	secondResourceName := "azurerm_virtual_machine_extension.test2"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineExtension_concurrent(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineExtensionExists(firstResourceName),
					testCheckAzureRMVirtualMachineExtensionExists(secondResourceName),
					resource.TestMatchResourceAttr(firstResourceName, "settings", regexp.MustCompile("hostname")),
					resource.TestMatchResourceAttr(secondResourceName, "settings", regexp.MustCompile("whoami")),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineExtension_linuxDiagnostics(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineExtension_linuxDiagnostics(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineExtensionExists("azurerm_virtual_machine_extension.test"),
				),
			},
		},
	})
}

func testCheckAzureRMVirtualMachineExtensionExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		vmName := rs.Primary.Attributes["virtual_machine_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).vmExtensionClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, vmName, name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on vmExtensionClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: VirtualMachine Extension %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineExtensionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).vmExtensionClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_machine_extension" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		vmName := rs.Primary.Attributes["virtual_machine_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, vmName, name, "")

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Virtual Machine Extension still exists:\n%#v", resp.VirtualMachineExtensionProperties)
		}
	}

	return nil
}

func testAccAzureRMVirtualMachineExtension_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_storage_account" "test" {
    name                     = "accsa%d"
    resource_group_name      = "${azurerm_resource_group.test.name}"
    location                 = "${azurerm_resource_group.test.location}"
    account_tier             = "Standard"
    account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine" "test" {
    name = "acctvm-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    network_interface_ids = ["${azurerm_network_interface.test.id}"]
    vm_size = "Standard_A0"

    storage_image_reference {
		publisher = "Canonical"
		offer = "UbuntuServer"
		sku = "16.04-LTS"
		version = "latest"
    }

    storage_os_disk {
        name = "myosdisk1"
        vhd_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
        caching = "ReadWrite"
        create_option = "FromImage"
    }

    os_profile {
		computer_name = "hostname%d"
		admin_username = "testadmin"
		admin_password = "Password1234!"
    }

    os_profile_linux_config {
	disable_password_authentication = false
   }
}

resource "azurerm_virtual_machine_extension" "test" {
    name = "acctvme-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_machine_name = "${azurerm_virtual_machine.test.name}"
    publisher = "Microsoft.Azure.Extensions"
    type = "CustomScript"
    type_handler_version = "2.0"

    settings = <<SETTINGS
	{
		"commandToExecute": "hostname"
	}
SETTINGS

	tags {
		environment = "Production"
	}
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineExtension_basicUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_storage_account" "test" {
    name                     = "accsa%d"
    resource_group_name      = "${azurerm_resource_group.test.name}"
    location                 = "${azurerm_resource_group.test.location}"
    account_tier             = "Standard"
    account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine" "test" {
    name = "acctvm-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    network_interface_ids = ["${azurerm_network_interface.test.id}"]
    vm_size = "Standard_A0"

    storage_image_reference {
		publisher = "Canonical"
		offer = "UbuntuServer"
		sku = "16.04-LTS"
		version = "latest"
    }

    storage_os_disk {
        name = "myosdisk1"
        vhd_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
        caching = "ReadWrite"
        create_option = "FromImage"
    }

    os_profile {
		computer_name = "hostname%d"
		admin_username = "testadmin"
		admin_password = "Password1234!"
    }

    os_profile_linux_config {
	disable_password_authentication = false
   }
}

resource "azurerm_virtual_machine_extension" "test" {
    name = "acctvme-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_machine_name = "${azurerm_virtual_machine.test.name}"
    publisher = "Microsoft.Azure.Extensions"
    type = "CustomScript"
    type_handler_version = "2.0"

    settings = <<SETTINGS
	{
		"commandToExecute": "whoami"
	}
SETTINGS

	tags {
		environment = "Production"
		cost_center = "MSFT"
	}
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineExtension_concurrent(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_storage_account" "test" {
    name                     = "accsa%d"
    resource_group_name      = "${azurerm_resource_group.test.name}"
    location                 = "${azurerm_resource_group.test.location}"
    account_tier             = "Standard"
    account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine" "test" {
    name = "acctvm-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    network_interface_ids = ["${azurerm_network_interface.test.id}"]
    vm_size = "Standard_A0"

    storage_image_reference {
	publisher = "Canonical"
	offer = "UbuntuServer"
	sku = "16.04-LTS"
	version = "latest"
    }

    storage_os_disk {
        name = "myosdisk1"
        vhd_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
        caching = "ReadWrite"
        create_option = "FromImage"
    }

    os_profile {
	computer_name = "hostname%d"
	admin_username = "testadmin"
	admin_password = "Password1234!"
    }

    os_profile_linux_config {
	disable_password_authentication = false
   }
}

resource "azurerm_virtual_machine_extension" "test" {
    name = "acctvme-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_machine_name = "${azurerm_virtual_machine.test.name}"
    publisher = "Microsoft.Azure.Extensions"
    type = "CustomScript"
    type_handler_version = "2.0"

    settings = <<SETTINGS
	{
		"commandToExecute": "hostname"
	}
SETTINGS
}

resource "azurerm_virtual_machine_extension" "test2" {
    name = "acctvme-%d-2"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_machine_name = "${azurerm_virtual_machine.test.name}"
    publisher = "Microsoft.OSTCExtensions"
    type = "CustomScriptForLinux"
    type_handler_version = "1.5"

    settings = <<SETTINGS
	{
		"commandToExecute": "whoami"
	}
SETTINGS
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineExtension_linuxDiagnostics(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_storage_account" "test" {
    name                     = "accsa%d"
    resource_group_name      = "${azurerm_resource_group.test.name}"
    location                 = "${azurerm_resource_group.test.location}"
    account_tier             = "Standard"
    account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine" "test" {
    name = "acctvm-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    network_interface_ids = ["${azurerm_network_interface.test.id}"]
    vm_size = "Standard_A0"

    storage_image_reference {
	publisher = "Canonical"
	offer = "UbuntuServer"
	sku = "16.04-LTS"
	version = "latest"
    }

    storage_os_disk {
        name = "myosdisk1"
        vhd_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
        caching = "ReadWrite"
        create_option = "FromImage"
    }

    os_profile {
	computer_name = "hostname%d"
	admin_username = "testadmin"
	admin_password = "Password1234!"
    }

    os_profile_linux_config {
	disable_password_authentication = false
   }
}

resource "azurerm_virtual_machine_extension" "test" {
    name = "acctvme-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_machine_name = "${azurerm_virtual_machine.test.name}"
    publisher = "Microsoft.OSTCExtensions"
    type = "LinuxDiagnostic"
    type_handler_version = "2.3"

    protected_settings = <<SETTINGS
	{
		"storageAccountName": "${azurerm_storage_account.test.name}",
        "storageAccountKey": "${azurerm_storage_account.test.primary_access_key}"
	}
SETTINGS

	tags {
		environment = "Production"
	}
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}
