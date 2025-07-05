// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DevCenterGalleryDataSource struct{}

func TestAccDevCenterGalleryDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dev_center_gallery", "test")
	r := DevCenterGalleryDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("dev_center_id").Exists(),
				check.That(data.ResourceName).Key("shared_gallery_id").Exists(),
			),
		},
	})
}

func (d DevCenterGalleryDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_dev_center_gallery" "test" {
  name          = azurerm_dev_center_gallery.test.name
  dev_center_id = azurerm_dev_center_gallery.test.dev_center_id
}
`, DevCenterGalleryTestResource{}.basic(data))
}
