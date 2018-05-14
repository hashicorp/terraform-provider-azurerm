package azurerm

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualMachineDataDiskAttachment_singleVHD(t *testing.T) {
	resourceName := "azurerm_virtual_machine_data_disk_attachment.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMVirtualMachineDataDiskAttachment_singleVHD(ri, location)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineDataDiskAttachment_singleManagedDisk(t *testing.T) {
	resourceName := "azurerm_virtual_machine_data_disk_attachment.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMVirtualMachineDataDiskAttachment_singleManagedDisk(ri, location)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineDataDiskAttachment_existingManagedDisk(t *testing.T) {
	resourceName := "azurerm_virtual_machine_data_disk_attachment.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMVirtualMachineDataDiskAttachment_existingManagedDisk(ri, location)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineDataDiskAttachment_multipleDisks(t *testing.T) {
	firstResourceName := "azurerm_virtual_machine_data_disk_attachment.first"
	secondResourceName := "azurerm_virtual_machine_data_disk_attachment.second"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMVirtualMachineDataDiskAttachment_multipleDisks(ri, location, 1, 2)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(firstResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "disk_size_gb", "10"),
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(secondResourceName),
					resource.TestCheckResourceAttr(secondResourceName, "disk_size_gb", "20"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineDataDiskAttachment_multipleDisksUpdate(t *testing.T) {
	firstResourceName := "azurerm_virtual_machine_data_disk_attachment.first"
	secondResourceName := "azurerm_virtual_machine_data_disk_attachment.second"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMVirtualMachineDataDiskAttachment_multipleDisks(ri, location, 1, 2)
	updatedConfig := testAccAzureRMVirtualMachineDataDiskAttachment_multipleDisks(ri, location, 3, 2)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(firstResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "disk_size_gb", "10"),
					resource.TestCheckResourceAttr(firstResourceName, "lun", "1"),
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(secondResourceName),
					resource.TestCheckResourceAttr(secondResourceName, "disk_size_gb", "20"),
					resource.TestCheckResourceAttr(secondResourceName, "lun", "2"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(firstResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "disk_size_gb", "10"),
					resource.TestCheckResourceAttr(firstResourceName, "lun", "3"),
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(secondResourceName),
					resource.TestCheckResourceAttr(secondResourceName, "disk_size_gb", "20"),
					resource.TestCheckResourceAttr(secondResourceName, "lun", "2"),
				),
			},
		},
	})
}

func testCheckAzureRMVirtualMachineDataDiskAttachmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		diskName := rs.Primary.Attributes["name"]
		virtualMachineId := rs.Primary.Attributes["virtual_machine_id"]

		id, err := parseAzureResourceID(virtualMachineId)
		if err != nil {
			return err
		}

		virtualMachineName := id.Path["virtualMachines"]
		resourceGroup := id.ResourceGroup

		client := testAccProvider.Meta().(*ArmClient).vmClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, virtualMachineName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on vmClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: VirtualMachine %q (resource group: %q) does not exist", virtualMachineName, resourceGroup)
		}

		// does the disk exist?
		for _, disk := range *resp.StorageProfile.DataDisks {
			if strings.EqualFold(*disk.Name, diskName) {
				return nil
			}
		}

		return fmt.Errorf("Disk %q was not found on Virtual Machine %q (Resource Group %q)", diskName, virtualMachineName, resourceName)
	}
}

func testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_machine_data_disk_attachment" {
			continue
		}

		diskName := rs.Primary.Attributes["name"]
		virtualMachineId := rs.Primary.Attributes["virtual_machine_id"]

		id, err := parseAzureResourceID(virtualMachineId)
		if err != nil {
			return err
		}

		virtualMachineName := id.Path["virtualMachines"]
		resourceGroup := id.ResourceGroup

		client := testAccProvider.Meta().(*ArmClient).vmClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, virtualMachineName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("Bad: Get on vmClient: %+v", err)
		}

		// does the disk exist?
		for _, disk := range *resp.StorageProfile.DataDisks {
			if strings.EqualFold(*disk.Name, diskName) {
				return fmt.Errorf("Disk %q is still attached to Virtual Machine %q (Resource Group %q)", diskName, virtualMachineName, resourceGroup)
			}
		}
	}

	return nil
}

func testAccAzureRMVirtualMachineDataDiskAttachment_singleVHD(rInt int, location string) string {
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
  name = "accsa%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location = "${azurerm_resource_group.test.location}"
  account_tier = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name = "vhds"
  resource_group_name = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_F2"

  storage_image_reference {
    publisher = "Canonical"
    offer = "UbuntuServer"
    sku = "16.04-LTS"
    version = "latest"
  }

  storage_os_disk {
    name = "osd-%d"
    vhd_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb = "45"
  }

  os_profile {
    computer_name = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  name               = "disk1-%d"
  virtual_machine_id = "${azurerm_virtual_machine.test.id}"
  vhd_uri            = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/mydatadisk1.vhd"
  disk_size_gb       = 10
  lun                = 1
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineDataDiskAttachment_singleManagedDisk(rInt int, location string) string {
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

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_F2"

  storage_image_reference {
    publisher = "Canonical"
    offer = "UbuntuServer"
    sku = "16.04-LTS"
    version = "latest"
  }

  storage_os_disk {
    name = "osd-%d"
    caching = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb = "50"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  name               = "disk1-%d"
  virtual_machine_id = "${azurerm_virtual_machine.test.id}"
  managed_disk_type  = "Standard_LRS"
  disk_size_gb       = 10
  lun                = 1
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineDataDiskAttachment_existingManagedDisk(rInt int, location string) string {
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

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_F2"

  storage_image_reference {
    publisher = "Canonical"
    offer = "UbuntuServer"
    sku = "16.04-LTS"
    version = "latest"
  }

  storage_os_disk {
    name = "osd-%d"
    caching = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb = "50"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_managed_disk" "test" {
  name = "acctestd-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  storage_account_type = "Standard_LRS"
  create_option = "Empty"
  disk_size_gb = 10
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  name               = "${azurerm_managed_disk.test.name}"
  virtual_machine_id = "${azurerm_virtual_machine.test.id}"
  managed_disk_id    = "${azurerm_managed_disk.test.id}"
  managed_disk_type  = "Standard_LRS"
  lun                = 1
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineDataDiskAttachment_multipleDisks(rInt int, location string, firstDiskLun int, secondDiskLun int) string {
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

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_F2"

  storage_image_reference {
    publisher = "Canonical"
    offer = "UbuntuServer"
    sku = "16.04-LTS"
    version = "latest"
  }

  storage_os_disk {
    name = "osd-%d"
    caching = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb = "50"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_virtual_machine_data_disk_attachment" "first" {
  name               = "disk1-%d"
  virtual_machine_id = "${azurerm_virtual_machine.test.id}"
  managed_disk_type  = "Standard_LRS"
  disk_size_gb       = 10
  lun                = %d
}

resource "azurerm_virtual_machine_data_disk_attachment" "second" {
  name               = "disk2-%d"
  virtual_machine_id = "${azurerm_virtual_machine.test.id}"
  managed_disk_type  = "Standard_LRS"
  disk_size_gb       = 20
  lun                = %d
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt, firstDiskLun, rInt, secondDiskLun)
}
