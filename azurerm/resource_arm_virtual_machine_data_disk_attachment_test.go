package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualMachineDataDiskAttachment_basic(t *testing.T) {
	resourceName := "azurerm_virtual_machine_data_disk_attachment.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMVirtualMachineDataDiskAttachment_basic(ri, location)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_machine_id"),
					resource.TestCheckResourceAttrSet(resourceName, "managed_disk_id"),
					resource.TestCheckResourceAttr(resourceName, "lun", "0"),
					resource.TestCheckResourceAttr(resourceName, "caching", "None"),
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

func TestAccAzureRMVirtualMachineDataDiskAttachment_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_virtual_machine_data_disk_attachment.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualMachineDataDiskAttachment_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_virtual_machine_data_disk_attachment"),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineDataDiskAttachment_multipleDisks(t *testing.T) {
	firstResourceName := "azurerm_virtual_machine_data_disk_attachment.first"
	secondResourceName := "azurerm_virtual_machine_data_disk_attachment.second"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMVirtualMachineDataDiskAttachment_multipleDisks(ri, location)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(firstResourceName),
					resource.TestCheckResourceAttrSet(firstResourceName, "virtual_machine_id"),
					resource.TestCheckResourceAttrSet(firstResourceName, "managed_disk_id"),
					resource.TestCheckResourceAttr(firstResourceName, "lun", "10"),
					resource.TestCheckResourceAttr(firstResourceName, "caching", "None"),

					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(secondResourceName),
					resource.TestCheckResourceAttrSet(secondResourceName, "virtual_machine_id"),
					resource.TestCheckResourceAttrSet(secondResourceName, "managed_disk_id"),
					resource.TestCheckResourceAttr(secondResourceName, "lun", "20"),
					resource.TestCheckResourceAttr(secondResourceName, "caching", "ReadOnly"),
				),
			},
			{
				ResourceName:      firstResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      secondResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMVirtualMachineDataDiskAttachment_updatingCaching(t *testing.T) {
	resourceName := "azurerm_virtual_machine_data_disk_attachment.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	resource.ParallelTest(t, resource.TestCase{
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

func TestAccAzureRMVirtualMachineDataDiskAttachment_updatingWriteAccelerator(t *testing.T) {
	resourceName := "azurerm_virtual_machine_data_disk_attachment.test"
	ri := tf.AccRandTimeInt()
	location := testAltLocation()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_writeAccelerator(ri, location, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "write_accelerator_enabled", "false"),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_writeAccelerator(ri, location, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "write_accelerator_enabled", "true"),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_writeAccelerator(ri, location, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "write_accelerator_enabled", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineDataDiskAttachment_virtualMachineExtension(t *testing.T) {
	resourceName := "azurerm_virtual_machine_data_disk_attachment.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_virtualMachineExtensionPrep(ri, location),
			},
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_virtualMachineExtensionComplete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_machine_id"),
					resource.TestCheckResourceAttrSet(resourceName, "managed_disk_id"),
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

		id, err := azure.ParseAzureResourceID(virtualMachineId)
		if err != nil {
			return err
		}

		virtualMachineName := id.Path["virtualMachines"]
		resourceGroup := id.ResourceGroup

		client := testAccProvider.Meta().(*ArmClient).compute.VMClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, virtualMachineName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on vmClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: VirtualMachine %q (resource group: %q) does not exist", virtualMachineName, resourceGroup)
		}

		diskId, err := azure.ParseAzureResourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		diskName := diskId.Path["dataDisks"]

		// deliberately not using strings.EqualFold as this is case sensitive
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

		id, err := azure.ParseAzureResourceID(virtualMachineId)
		if err != nil {
			return err
		}

		virtualMachineName := id.Path["virtualMachines"]
		resourceGroup := id.ResourceGroup

		client := testAccProvider.Meta().(*ArmClient).compute.VMClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, virtualMachineName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("Bad: Get on vmClient: %+v", err)
		}

		diskId, err := azure.ParseAzureResourceID(rs.Primary.ID)
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

func testAccAzureRMVirtualMachineDataDiskAttachment_requiresImport(rInt int, location string) string {
	template := testAccAzureRMVirtualMachineDataDiskAttachment_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_data_disk_attachment" "import" {
  managed_disk_id    = "${azurerm_virtual_machine_data_disk_attachment.test.managed_disk_id}"
  virtual_machine_id = "${azurerm_virtual_machine_data_disk_attachment.test.virtual_machine_id}"
  lun                = "${azurerm_virtual_machine_data_disk_attachment.test.lun}"
  caching            = "${azurerm_virtual_machine_data_disk_attachment.test.caching}"
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

func testAccAzureRMVirtualMachineDataDiskAttachment_writeAccelerator(rInt int, location string, enabled bool) string {
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

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_M64s"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Premium_LRS"
  }

  os_profile {
    computer_name  = "hn%d"
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
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id           = "${azurerm_managed_disk.test.id}"
  virtual_machine_id        = "${azurerm_virtual_machine.test.id}"
  lun                       = "0"
  caching                   = "None"
  write_accelerator_enabled = %t
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, enabled)
}

func testAccAzureRMVirtualMachineDataDiskAttachment_template(rInt int, location string) string {
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

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_F2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "hn%d"
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

func testAccAzureRMVirtualMachineDataDiskAttachment_virtualMachineExtensionPrep(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestni%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctestvm%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_F4"

  delete_os_disk_on_termination    = true
  delete_data_disks_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_profile {
    computer_name  = "testvm"
    admin_username = "tfuser123"
    admin_password = "Password1234!"
  }

  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "staging"
  }
}

resource "azurerm_virtual_machine_extension" "test" {
  name                 = "random-script"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_machine_name = "${azurerm_virtual_machine.test.name}"
  publisher            = "Microsoft.Azure.Extensions"
  type                 = "CustomScript"
  type_handler_version = "2.0"

  settings = <<SETTINGS
	{
		"commandToExecute": "hostname"
	}
SETTINGS

  tags = {
    environment = "Production"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineDataDiskAttachment_virtualMachineExtensionComplete(rInt int, location string) string {
	template := testAccAzureRMVirtualMachineDataDiskAttachment_virtualMachineExtensionPrep(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_managed_disk" "test" {
  name                 = "acctest%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = "${azurerm_managed_disk.test.id}"
  virtual_machine_id = "${azurerm_virtual_machine.test.id}"
  lun                = "11"
  caching            = "ReadWrite"
}
`, template, rInt)
}
