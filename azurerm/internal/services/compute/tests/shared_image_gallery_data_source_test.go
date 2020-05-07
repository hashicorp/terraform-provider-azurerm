package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMSharedImageGallery_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_gallery", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageGalleryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSharedImageGallery_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMSharedImageGallery_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_gallery", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSharedImageGalleryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSharedImageGallery_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Shared images and things."),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Hello", "There"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.World", "Example"),
				),
			},
		},
	})
}

func testAccDataSourceSharedImageGallery_basic(data acceptance.TestData) string {
	template := testAccAzureRMSharedImageGallery_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_shared_image_gallery" "test" {
  name                = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_shared_image_gallery.test.resource_group_name
}
`, template)
}

func testAccDataSourceSharedImageGallery_complete(data acceptance.TestData) string {
	template := testAccAzureRMSharedImageGallery_complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_shared_image_gallery" "test" {
  name                = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_shared_image_gallery.test.resource_group_name
}
`, template)
}
