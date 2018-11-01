package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMSnapshot_import(t *testing.T) {
	resourceName := "azurerm_snapshot.test"
	ri := acctest.RandInt()
	config := testAccAzureRMSnapshot_fromManagedDisk(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_uri"},
			},
		},
	})
}

func TestAccAzureRMSnapshot_importEncryption(t *testing.T) {
	resourceName := "azurerm_snapshot.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	config := testAccAzureRMSnapshot_encryption(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_uri"},
			},
		},
	})
}
