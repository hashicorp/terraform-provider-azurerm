// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SharedImageGalleryDataSource struct{}

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

func TestAccDataSourceSharedImageGallery_imageNames(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image_gallery", "test")
	r := SharedImageGalleryDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.imageNames(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("image_names.#").HasValue("1"),
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

func (SharedImageGalleryDataSource) imageNames(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image_gallery" "test" {
  name                = azurerm_shared_image.test.gallery_name
  resource_group_name = azurerm_shared_image.test.resource_group_name
}
`, SharedImageResource{}.basic(data))
}
