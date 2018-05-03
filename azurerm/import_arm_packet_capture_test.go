package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func testAccAzureRMPacketCapture_importBasic(t *testing.T) {
	rInt := acctest.RandInt()
	rString := acctest.RandString(6)
	location := testLocation()

	resourceName := "azurerm_packet_capture.test"
	config := testAzureRMPacketCapture_storageAccountAndLocalDiskConfig(rInt, rString, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPacketCaptureDestroy,
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
