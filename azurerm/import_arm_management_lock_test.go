package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMManagementLock_importResourceGroupReadOnlyBasic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMManagementLock_resourceGroupReadOnlyBasic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_management_lock.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMManagementLock_importResourceGroupReadOnlyComplete(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMManagementLock_resourceGroupReadOnlyComplete(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_management_lock.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMManagementLock_importResourceGroupCanNotDeleteBasic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMManagementLock_resourceGroupCanNotDeleteBasic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_management_lock.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMManagementLock_importResourceGroupCanNotDeleteComplete(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMManagementLock_resourceGroupCanNotDeleteComplete(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_management_lock.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMManagementLock_importPublicIPCanNotDeleteBasic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMManagementLock_publicIPCanNotDeleteBasic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_management_lock.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMManagementLock_importPublicIPReadOnlyBasic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMManagementLock_publicIPReadOnlyBasic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_management_lock.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
