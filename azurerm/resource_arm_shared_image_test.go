package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSharedImage_basic(t *testing.T) {
	resourceName := "azurerm_shared_image.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSharedImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSharedImage_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
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
func TestAccAzureRMSharedImage_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_shared_image.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSharedImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSharedImage_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				Config:      testAccAzureRMSharedImage_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_shared_image"),
			},
		},
	})
}

func TestAccAzureRMSharedImage_complete(t *testing.T) {
	resourceName := "azurerm_shared_image.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSharedImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSharedImage_complete(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Linux"),
					resource.TestCheckResourceAttr(resourceName, "description", "Wubba lubba dub dub"),
					resource.TestCheckResourceAttr(resourceName, "eula", "Do you agree there's infinite Rick's and Infinite Morty's?"),
					resource.TestCheckResourceAttr(resourceName, "privacy_statement_uri", "https://council.of.ricks/privacy-statement"),
					resource.TestCheckResourceAttr(resourceName, "release_note_uri", "https://council.of.ricks/changelog.md"),
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

func testCheckAzureRMSharedImageDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).compute.GalleryImagesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_shared_image" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		galleryName := rs.Primary.Attributes["gallery_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, galleryName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Shared Image still exists:\n%+v", resp)
		}
	}

	return nil
}

func testCheckAzureRMSharedImageExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		imageName := rs.Primary.Attributes["name"]
		galleryName := rs.Primary.Attributes["gallery_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Shared Image: %s", imageName)
		}

		client := testAccProvider.Meta().(*ArmClient).compute.GalleryImagesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, galleryName, imageName)
		if err != nil {
			return fmt.Errorf("Bad: Get on galleryImagesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Shared Image %q (Gallery %q / Resource Group: %q) does not exist", imageName, galleryName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMSharedImage_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

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
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMSharedImage_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image" "import" {
  name                = "${azurerm_shared_image.test.name}"
  gallery_name        = "${azurerm_shared_image.test.gallery_name}"
  resource_group_name = "${azurerm_shared_image.test.resource_group_name}"
  location            = "${azurerm_shared_image.test.location}"
  os_type             = "${azurerm_shared_image.test.os_type}"

  identifier {
    publisher = "AccTesPublisher%d"
    offer     = "AccTesOffer%d"
    sku       = "AccTesSku%d"
  }
}
`, testAccAzureRMSharedImage_basic(rInt, location), rInt, rInt, rInt)
}

func testAccAzureRMSharedImage_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_shared_image" "test" {
  name                  = "acctestimg%d"
  gallery_name          = "${azurerm_shared_image_gallery.test.name}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  location              = "${azurerm_resource_group.test.location}"
  os_type               = "Linux"
  description           = "Wubba lubba dub dub"
  eula                  = "Do you agree there's infinite Rick's and Infinite Morty's?"
  privacy_statement_uri = "https://council.of.ricks/privacy-statement"
  release_note_uri      = "https://council.of.ricks/changelog.md"

  identifier {
    publisher = "AccTesPublisher%d"
    offer     = "AccTesOffer%d"
    sku       = "AccTesSku%d"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}
