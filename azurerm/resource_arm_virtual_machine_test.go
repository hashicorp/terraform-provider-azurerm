package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMVirtualMachine_winTimeZone(t *testing.T) {
	resourceName := "azurerm_virtual_machine.test"
	var vm compute.VirtualMachine
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMVirtualMachine_winTimeZone(ri, testLocation())
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineExists("azurerm_virtual_machine.test", &vm),
					resource.TestCheckResourceAttr(resourceName, "os_profile_windows_config.59207889.timezone", "Pacific Standard Time"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachine_SystemAssignedIdentity(t *testing.T) {
	var vm compute.VirtualMachine
	resourceName := "azurerm_virtual_machine.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMVirtualMachineSystemAssignedIdentity(ri, testLocation())
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineExists(resourceName, &vm),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "SystemAssigned"),
					resource.TestCheckResourceAttr(resourceName, "identity.0.identity_ids.#", "0"),
					resource.TestMatchResourceAttr(resourceName, "identity.0.principal_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachine_UserAssignedIdentity(t *testing.T) {
	var vm compute.VirtualMachine
	resourceName := "azurerm_virtual_machine.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(14)
	config := testAccAzureRMVirtualMachineUserAssignedIdentity(ri, testLocation(), rs)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineExists(resourceName, &vm),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "UserAssigned"),
					resource.TestCheckResourceAttr(resourceName, "identity.0.identity_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "identity.0.principal_id", ""),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachine_multipleAssignedIdentity(t *testing.T) {
	var vm compute.VirtualMachine
	resourceName := "azurerm_virtual_machine.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(14)
	config := testAccAzureRMVirtualMachineMultipleAssignedIdentity(ri, testLocation(), rs)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineExists(resourceName, &vm),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "SystemAssigned, UserAssigned"),
					resource.TestCheckResourceAttr(resourceName, "identity.0.identity_ids.#", "1"),
					resource.TestMatchResourceAttr(resourceName, "identity.0.principal_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func testCheckAzureRMVirtualMachineExists(resourceName string, vm *compute.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		vmName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for virtual machine: %s", vmName)
		}

		client := testAccProvider.Meta().(*ArmClient).compute.VMClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, vmName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on vmClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: VirtualMachine %q (resource group: %q) does not exist", vmName, resourceGroup)
		}

		*vm = resp

		return nil
	}
}

func testCheckAzureRMVirtualMachineDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).compute.VMClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_machine" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name, "")

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("Virtual Machine still exists:\n%#v", resp.VirtualMachineProperties)
	}

	return nil
}

func testAccAzureRMVirtualMachine_winTimeZone(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2012-Datacenter"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  os_profile {
    computer_name  = "winhost01"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_windows_config {
    timezone = "Pacific Standard Time"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineSystemAssignedIdentity(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb  = "45"
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineUserAssignedIdentity(rInt int, location string, rString string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  name = "acctest%s"
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb  = "45"
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }

  identity {
    type         = "UserAssigned"
    identity_ids = ["${azurerm_user_assigned_identity.test.id}"]
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rString, rInt, rInt)
}

func testAccAzureRMVirtualMachineMultipleAssignedIdentity(rInt int, location string, rString string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb  = "45"
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = ["${azurerm_user_assigned_identity.test.id}"]
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rString, rInt, rInt)
}
