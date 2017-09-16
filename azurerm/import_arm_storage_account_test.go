package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMStorageAccount_importBasic(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"

	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	config := testAccAzureRMStorageAccount_basic(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
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
