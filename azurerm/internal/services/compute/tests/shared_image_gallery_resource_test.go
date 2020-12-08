package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSharedImageGallery_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_gallery", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageGalleryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSharedImageGallery_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageGalleryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "description", ""),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSharedImageGallery_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_gallery", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageGalleryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSharedImageGallery_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageGalleryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "description", ""),
				),
			},
			{
				Config:      testAccAzureRMSharedImageGallery_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_shared_image_gallery"),
			},
		},
	})
}

func TestAccAzureRMSharedImageGallery_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image_gallery", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageGalleryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSharedImageGallery_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSharedImageGalleryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Shared images and things."),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Hello", "There"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.World", "Example"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMSharedImageGalleryDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.GalleriesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_shared_image_gallery" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Shared Image Gallery still exists:\n%+v", resp)
		}
	}

	return nil
}

func testCheckAzureRMSharedImageGalleryExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.GalleriesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		galleryName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Shared Image Gallery: %s", galleryName)
		}

		resp, err := client.Get(ctx, resourceGroup, galleryName)
		if err != nil {
			return fmt.Errorf("Bad: Get on galleriesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Shared Image Gallery %q (resource group: %q) does not exist", galleryName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMSharedImageGallery_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMSharedImageGallery_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image_gallery" "import" {
  name                = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_shared_image_gallery.test.resource_group_name
  location            = azurerm_shared_image_gallery.test.location
}
`, testAccAzureRMSharedImageGallery_basic(data))
}

func testAccAzureRMSharedImageGallery_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  description         = "Shared images and things."

  tags = {
    Hello = "There"
    World = "Example"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
