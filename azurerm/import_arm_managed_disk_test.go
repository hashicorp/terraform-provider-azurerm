package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMManagedDisk_importEmpty(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMManagedDisk_empty(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_managed_disk.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMManagedDisk_importEmpty_withZone(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMManagedDisk_empty_withZone(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagedDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_managed_disk.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
