package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMImage_importStandalone(t *testing.T) {
	ri := acctest.RandInt()
	resourceGroup := fmt.Sprintf("acctestRG-%d", ri)
	userName := "testadmin"
	password := "Password1234s!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", ri)
	sshPort := "22"
	location := testLocation()
	preConfig := testAccAzureRMImage_standaloneImage_setup(ri, userName, password, hostName, location)
	postConfig := testAccAzureRMImage_standaloneImage_provision(ri, userName, password, hostName, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMImageDestroy,
		Steps: []resource.TestStep{
			{
				//need to create a vm and then reference it in the image creation
				Config:  preConfig,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
					testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, location),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMImageExists("azurerm_image.test", true),
				),
			},
			{
				ResourceName:      "azurerm_image.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
