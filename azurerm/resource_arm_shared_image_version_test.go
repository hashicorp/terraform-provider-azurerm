package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSharedImageVersion_basic(t *testing.T) {
	resourceName := "azurerm_shared_image_version.test"

	ri := acctest.RandInt()
	resourceGroup := fmt.Sprintf("acctestRG-%d", ri)
	userName := "testadmin"
	password := "Password1234!"
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", ri)
	sshPort := "22"
	location := testLocation()
	altLocation := testAltLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSharedImageVersionDestroy,
		Steps: []resource.TestStep{
			{
				// need to create a vm and then reference it in the image creation
				Config:  testAccAzureRMSharedImageVersion_setup(ri, location, userName, password, hostName),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
					testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, location),
				),
			},
			{
				Config: testAccAzureRMSharedImageVersion_imageVersion(ri, location, userName, password, hostName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageVersionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "managed_image_id"),
					resource.TestCheckResourceAttr(resourceName, "regions.#", "1"),
				),
			},
			{
				Config: testAccAzureRMSharedImageVersion_imageVersionUpdated(ri, location, altLocation, userName, password, hostName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageVersionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "managed_image_id"),
					resource.TestCheckResourceAttr(resourceName, "regions.#", "2"),
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

func testCheckAzureRMSharedImageVersionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).galleryImageVersionsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_shared_image_version" {
			continue
		}

		imageVersion := rs.Primary.Attributes["name"]
		imageName := rs.Primary.Attributes["image_name"]
		galleryName := rs.Primary.Attributes["gallery_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, galleryName, imageName, imageVersion, "")
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("Shared Image Version still exists:\n%+v", resp)
	}

	return nil
}

func testCheckAzureRMSharedImageVersionExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		imageVersion := rs.Primary.Attributes["name"]
		imageName := rs.Primary.Attributes["image_name"]
		galleryName := rs.Primary.Attributes["gallery_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Shared Image Version: %s", imageName)
		}

		client := testAccProvider.Meta().(*ArmClient).galleryImageVersionsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, galleryName, imageName, imageVersion, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on galleryImageVersionsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Shared Image Version %q (Image %q / Gallery %q / Resource Group: %q) does not exist", imageVersion, imageName, galleryName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMSharedImageVersion_setup(rInt int, location, username, password, hostname string) string {
	return testAccAzureRMImage_standaloneImage_setup(rInt, username, password, hostname, location)
}

func testAccAzureRMSharedImageVersion_provision(rInt int, location, username, password, hostname string) string {
	template := testAccAzureRMImage_standaloneImage_provision(rInt, username, password, hostname, location)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%d"
  gallery_name        = "${azurerm_shared_image_gallery.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  os_type             = "Linux"

  identifier {
    publisher = "AccTesPublisher%d"
    offer     = "AccTesOffer%d"
    sku       = "AccTesSku%d"
  }
}
`, template, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMSharedImageVersion_imageVersion(rInt int, location, username, password, hostname string) string {
	template := testAccAzureRMSharedImageVersion_provision(rInt, location, username, password, hostname)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = "${azurerm_shared_image_gallery.test.name}"
  image_name          = "${azurerm_shared_image.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  managed_image_id    = "${azurerm_image.test.id}"
  regions             = ["${azurerm_resource_group.test.location}"]
}
`, template)
}

func testAccAzureRMSharedImageVersion_imageVersionUpdated(rInt int, location, altLocation, username, password, hostname string) string {
	template := testAccAzureRMSharedImageVersion_provision(rInt, location, username, password, hostname)
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = "${azurerm_shared_image_gallery.test.name}"
  image_name          = "${azurerm_shared_image.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  managed_image_id    = "${azurerm_image.test.id}"
  regions             = ["${azurerm_resource_group.test.location}", "%s"]
}
`, template, altLocation)
}
