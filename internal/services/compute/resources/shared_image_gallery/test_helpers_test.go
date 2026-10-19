// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package shared_image_gallery_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type SharedImageResource struct{}

func (SharedImageResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}
resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_shared_image" "test" {
  name                = "acctestimg%[2]d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"
  identifier {
    publisher = "AccTesPublisher%[2]d"
    offer     = "AccTesOffer%[2]d"
    sku       = "AccTesSku%[2]d"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}
