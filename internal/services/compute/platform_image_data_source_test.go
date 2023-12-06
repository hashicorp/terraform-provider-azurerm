// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PlatformImageDataSource struct{}

func TestAccDataSourcePlatformImage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_platform_image", "test")
	r := PlatformImageDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("offer").HasValue("0001-com-ubuntu-server-jammy"),
				check.That(data.ResourceName).Key("sku").HasValue("22_04-lts"),
			),
		},
	})
}

func TestAccDataSourcePlatformImage_withVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_platform_image", "test")
	r := PlatformImageDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.withVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("publisher").HasValue("Canonical"),
				check.That(data.ResourceName).Key("offer").HasValue("0001-com-ubuntu-server-jammy"),
				check.That(data.ResourceName).Key("sku").HasValue("22_04-lts"),
				check.That(data.ResourceName).Key("version").HasValue("22.04.202310040"),
			),
		},
	})
}

func (PlatformImageDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_platform_image" "test" {
  location  = "%s"
  publisher = "Canonical"
  offer     = "0001-com-ubuntu-server-jammy"
  sku       = "22_04-lts"
}
`, data.Locations.Primary)
}

func (PlatformImageDataSource) withVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_platform_image" "test" {
  location  = "%s"
  publisher = "Canonical"
  offer     = "0001-com-ubuntu-server-jammy"
  sku       = "22_04-lts"
  version   = "22.04.202310040"
}
`, data.Locations.Primary)
}
