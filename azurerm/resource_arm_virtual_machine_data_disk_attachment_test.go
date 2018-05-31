package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualMachineDataDiskAttachment_basic(t *testing.T) {
	resourceName := "azurerm_virtual_machine_data_disk_attachment.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMVirtualMachineDataDiskAttachment_basic(ri, location)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_machine_id"),
					resource.TestCheckResourceAttrSet(resourceName, "managed_disk_id"),
					resource.TestCheckResourceAttr(resourceName, "lun", "0"),
					resource.TestCheckResourceAttr(resourceName, "caching", "None"),
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
	config := testAccAzureRMVirtualMachineDataDiskAttachment_multipleDisks(ri, location)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(firstResourceName),
					resource.TestCheckResourceAttrSet(firstResourceName, "name"),
					resource.TestCheckResourceAttrSet(firstResourceName, "virtual_machine_id"),
					resource.TestCheckResourceAttrSet(firstResourceName, "managed_disk_id"),
					resource.TestCheckResourceAttr(firstResourceName, "lun", "10"),
					resource.TestCheckResourceAttr(firstResourceName, "caching", "None"),

					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(secondResourceName),
					resource.TestCheckResourceAttrSet(secondResourceName, "name"),
					resource.TestCheckResourceAttrSet(secondResourceName, "virtual_machine_id"),
					resource.TestCheckResourceAttrSet(secondResourceName, "managed_disk_id"),
					resource.TestCheckResourceAttr(secondResourceName, "lun", "20"),
					resource.TestCheckResourceAttr(secondResourceName, "caching", "ReadOnly"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineDataDiskAttachment_updatingCaching(t *testing.T) {
	resourceName := "azurerm_virtual_machine_data_disk_attachment.test"
	ri := acctest.RandInt()
	location := testLocation()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "caching", "None"),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_readOnly(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "caching", "ReadOnly"),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_readWrite(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "caching", "ReadWrite"),
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

		diskId, err := parseAzureResourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		diskName := diskId.Path["dataDisks"]

		// deliberately not using strings.Equals as this is case sensitive
		for _, disk := range *resp.StorageProfile.DataDisks {
			if *disk.Name == diskName {
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

		diskId, err := parseAzureResourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		diskName := diskId.Path["dataDisks"]

		for _, disk := range *resp.StorageProfile.DataDisks {
			// deliberately not using strings.Equals as this is case sensitive
			if *disk.Name == diskName {
				return fmt.Errorf("Disk %q is still attached to Virtual Machine %q (Resource Group %q)", diskName, virtualMachineName, resourceGroup)
			}
		}
	}

	return nil
}

func testAccAzureRMVirtualMachineDataDiskAttachment_basic(rInt int, location string) string {
	template := testAccAzureRMVirtualMachineDataDiskAttachment_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = "${azurerm_managed_disk.test.id}"
  virtual_machine_id = "${azurerm_virtual_machine.test.id}"
  lun                = "0"
  caching            = "None"
}
`, template)
}

func testAccAzureRMVirtualMachineDataDiskAttachment_multipleDisks(rInt int, location string) string {
	template := testAccAzureRMVirtualMachineDataDiskAttachment_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_data_disk_attachment" "first" {
  managed_disk_id    = "${azurerm_managed_disk.test.id}"
  virtual_machine_id = "${azurerm_virtual_machine.test.id}"
  lun                = "10"
  caching            = "None"
}

resource "azurerm_managed_disk" "second" {
  name                 = "%d-disk2"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_virtual_machine_data_disk_attachment" "second" {
  managed_disk_id    = "${azurerm_managed_disk.second.id}"
  virtual_machine_id = "${azurerm_virtual_machine.test.id}"
  lun                = "20"
  caching            = "ReadOnly"
}
`, template, rInt)
}

func testAccAzureRMVirtualMachineDataDiskAttachment_readOnly(rInt int, location string) string {
	template := testAccAzureRMVirtualMachineDataDiskAttachment_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = "${azurerm_managed_disk.test.id}"
  virtual_machine_id = "${azurerm_virtual_machine.test.id}"
  lun                = "0"
  caching            = "ReadOnly"
}
`, template)
}

func testAccAzureRMVirtualMachineDataDiskAttachment_readWrite(rInt int, location string) string {
	template := testAccAzureRMVirtualMachineDataDiskAttachment_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = "${azurerm_managed_disk.test.id}"
  virtual_machine_id = "${azurerm_virtual_machine.test.id}"
  lun                = "0"
  caching            = "ReadWrite"
}
`, template)
}

func testAccAzureRMVirtualMachineDataDiskAttachment_template(rInt int, location string) string {
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
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
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
  name                 = "%d-disk1"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}
