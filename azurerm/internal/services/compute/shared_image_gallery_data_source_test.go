package compute_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type SharedImageGalleryDataSource struct {
}

func TestAccDataSourceSharedImageGallery_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_gallery", "test")
	r := SharedImageGalleryDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceSharedImageGallery_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_gallery", "test")
	r := SharedImageGalleryDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("description").HasValue("Shared images and things."),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.Hello").HasValue("There"),
				check.That(data.ResourceName).Key("tags.World").HasValue("Example"),
			),
		},
	})
}

func (SharedImageGalleryDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image_gallery" "test" {
  name                = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_shared_image_gallery.test.resource_group_name
}
`, SharedImageGalleryResource{}.basic(data))
}

func (SharedImageGalleryDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image_gallery" "test" {
  name                = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_shared_image_gallery.test.resource_group_name
}
`, SharedImageGalleryResource{}.complete(data))
}
