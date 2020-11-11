package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualMachineDataDiskAttachment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_data_disk_attachment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_machine_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "managed_disk_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "lun", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "caching", "None"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualMachineDataDiskAttachment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_data_disk_attachment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualMachineDataDiskAttachment_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_virtual_machine_data_disk_attachment"),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineDataDiskAttachment_destroy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_data_disk_attachment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(data.ResourceName),
					testCheckAzureRMVirtualMachineDataDiskAttachmentDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMVirtualMachineDataDiskAttachment_multipleDisks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_data_disk_attachment", "first")

	secondResourceName := "azurerm_virtual_machine_data_disk_attachment.second"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_multipleDisks(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_machine_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "managed_disk_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "lun", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "caching", "None"),

					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(secondResourceName),
					resource.TestCheckResourceAttrSet(secondResourceName, "virtual_machine_id"),
					resource.TestCheckResourceAttrSet(secondResourceName, "managed_disk_id"),
					resource.TestCheckResourceAttr(secondResourceName, "lun", "20"),
					resource.TestCheckResourceAttr(secondResourceName, "caching", "ReadOnly"),
				),
			},
			data.ImportStep(),
			{
				ResourceName:      secondResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMVirtualMachineDataDiskAttachment_updatingCaching(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_data_disk_attachment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "caching", "None"),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_readOnly(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "caching", "ReadOnly"),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_readWrite(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "caching", "ReadWrite"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineDataDiskAttachment_updatingWriteAccelerator(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_data_disk_attachment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_writeAccelerator(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "write_accelerator_enabled", "false"),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_writeAccelerator(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "write_accelerator_enabled", "true"),
				),
			},
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_writeAccelerator(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "write_accelerator_enabled", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineDataDiskAttachment_managedServiceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_data_disk_attachment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_managedServiceIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_machine_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "managed_disk_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "lun", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "caching", "None"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualMachineDataDiskAttachment_virtualMachineExtension(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_data_disk_attachment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDataDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_virtualMachineExtensionPrep(data),
			},
			{
				Config: testAccAzureRMVirtualMachineDataDiskAttachment_virtualMachineExtensionComplete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineDataDiskAttachmentExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_machine_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "managed_disk_id"),
				),
			},
		},
	})
}

func testCheckAzureRMVirtualMachineDataDiskAttachmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testCheckAzureRMVirtualMachineDataDiskAttachmentDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.VMClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

		outputDisks := make([]compute.DataDisk, 0)
		for _, disk := range *resp.StorageProfile.DataDisks {
			// deliberately not using strings.Equals as this is case sensitive
			if *disk.Name == diskName {
				continue
			}

			outputDisks = append(outputDisks, disk)
		}
		resp.StorageProfile.DataDisks = &outputDisks

		// fixes #2485
		resp.Identity = nil
		// fixes #1600
		resp.Resources = nil

		future, err := client.CreateOrUpdate(ctx, resourceGroup, virtualMachineName, resp)
		if err != nil {
			return fmt.Errorf("Error updating Virtual Machine %q (Resource Group %q) with Disk %q: %+v", virtualMachineName, resourceGroup, diskName, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for Virtual Machine %q (Resource Group %q) to finish updating Disk %q: %+v", virtualMachineName, resourceGroup, diskName, err)
		}

		return nil
	}
}

func testAccAzureRMVirtualMachineDataDiskAttachment_basic(data acceptance.TestData) string {
	template := testAccAzureRMVirtualMachineDataDiskAttachment_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "0"
  caching            = "None"
}
`, template)
}

func testAccAzureRMVirtualMachineDataDiskAttachment_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMVirtualMachineDataDiskAttachment_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_data_disk_attachment" "import" {
  managed_disk_id    = azurerm_virtual_machine_data_disk_attachment.test.managed_disk_id
  virtual_machine_id = azurerm_virtual_machine_data_disk_attachment.test.virtual_machine_id
  lun                = azurerm_virtual_machine_data_disk_attachment.test.lun
  caching            = azurerm_virtual_machine_data_disk_attachment.test.caching
}
`, template)
}

func testAccAzureRMVirtualMachineDataDiskAttachment_managedServiceIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
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

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_managed_disk" "test" {
  name                 = "%d-disk1"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "0"
  caching            = "None"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualMachineDataDiskAttachment_multipleDisks(data acceptance.TestData) string {
	template := testAccAzureRMVirtualMachineDataDiskAttachment_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_data_disk_attachment" "first" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "10"
  caching            = "None"
}

resource "azurerm_managed_disk" "second" {
  name                 = "%d-disk2"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_virtual_machine_data_disk_attachment" "second" {
  managed_disk_id    = azurerm_managed_disk.second.id
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "20"
  caching            = "ReadOnly"
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualMachineDataDiskAttachment_readOnly(data acceptance.TestData) string {
	template := testAccAzureRMVirtualMachineDataDiskAttachment_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "0"
  caching            = "ReadOnly"
}
`, template)
}

func testAccAzureRMVirtualMachineDataDiskAttachment_readWrite(data acceptance.TestData) string {
	template := testAccAzureRMVirtualMachineDataDiskAttachment_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "0"
  caching            = "ReadWrite"
}
`, template)
}

func testAccAzureRMVirtualMachineDataDiskAttachment_writeAccelerator(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
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
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id           = azurerm_managed_disk.test.id
  virtual_machine_id        = azurerm_virtual_machine.test.id
  lun                       = "0"
  caching                   = "None"
  write_accelerator_enabled = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, enabled)
}

func testAccAzureRMVirtualMachineDataDiskAttachment_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
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
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualMachineDataDiskAttachment_virtualMachineExtensionPrep(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestni%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctestvm%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
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
  virtual_machine_id   = azurerm_virtual_machine.test.id
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualMachineDataDiskAttachment_virtualMachineExtensionComplete(data acceptance.TestData) string {
	template := testAccAzureRMVirtualMachineDataDiskAttachment_virtualMachineExtensionPrep(data)
	return fmt.Sprintf(`
%s

resource "azurerm_managed_disk" "test" {
  name                 = "acctest%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "11"
  caching            = "ReadWrite"
}
`, template, data.RandomInteger)
}
