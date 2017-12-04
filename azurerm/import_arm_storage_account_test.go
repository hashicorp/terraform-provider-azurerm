package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"

	"github.com/terraform-providers/terraform-provider-azurerm/utils/aztesting"
)

func TestAccAzureRMStorageAccount_importBasic(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"

	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	config := testAccAzureRMStorageAccount_basic(ri, rs, testLocation())

	aztesting.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
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

func TestAccAzureRMStorageAccount_importPremium(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"

	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	config := testAccAzureRMStorageAccount_premium(ri, rs, testLocation())

	aztesting.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
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

func TestAccAzureRMStorageAccount_importNonStandardCasing(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"

	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	config := testAccAzureRMStorageAccount_nonStandardCasing(ri, rs, testLocation())

	aztesting.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
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

func TestAccAzureRMStorageAccount_importBlobEncryption(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"

	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	config := testAccAzureRMStorageAccount_blobEncryption(ri, rs, testLocation())

	aztesting.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
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

func TestAccAzureRMStorageAccount_importFileEncryption(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"

	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	config := testAccAzureRMStorageAccount_fileEncryption(ri, rs, testLocation())

	aztesting.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
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

func TestAccAzureRMStorageAccount_importEnableHttpsTrafficOnly(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"

	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	config := testAccAzureRMStorageAccount_enableHttpsTrafficOnly(ri, rs, testLocation())

	aztesting.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
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
