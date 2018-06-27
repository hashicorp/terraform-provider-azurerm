package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMVirtualMachineDataDiskAttachment_importBasic(t *testing.T) {
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMVirtualMachineDataDiskAttachment_importMultipleDisks(t *testing.T) {
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
			},
			{
				ResourceName:      "azurerm_virtual_machine_data_disk_attachment.first",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "azurerm_virtual_machine_data_disk_attachment.second",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
