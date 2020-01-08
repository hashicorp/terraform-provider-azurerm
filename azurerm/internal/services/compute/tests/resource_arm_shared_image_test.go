package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSharedImage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSharedImage_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "description", ""),
				),
			},
			data.ImportStep(),
		},
	})
}
func TestAccAzureRMSharedImage_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_shared_image", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSharedImage_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "description", ""),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSharedImage_requiresImport),
		},
	})
}

func TestAccAzureRMSharedImage_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSharedImage_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "os_type", "Linux"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Wubba lubba dub dub"),
					resource.TestCheckResourceAttr(data.ResourceName, "eula", "Do you agree there's infinite Rick's and Infinite Morty's?"),
					resource.TestCheckResourceAttr(data.ResourceName, "privacy_statement_uri", "https://council.of.ricks/privacy-statement"),
					resource.TestCheckResourceAttr(data.ResourceName, "release_note_uri", "https://council.of.ricks/changelog.md"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMSharedImageDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.GalleryImagesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.GalleryImagesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMSharedImage_basic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMSharedImage_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMSharedImage_basic(data)
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
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMSharedImage_complete(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
