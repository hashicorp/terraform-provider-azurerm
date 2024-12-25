// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ComputeSkusDataSource struct{}

func TestAccDataSourceComputeSkus_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_compute_skus", "test")
	r := ComputeSkusDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("skus.#").HasValue("1"),
				check.That(data.ResourceName).Key("skus.1").DoesNotExist(),
				check.That(data.ResourceName).Key("skus.0.name").HasValue("Standard_DS2_v2"),
				check.That(data.ResourceName).Key("skus.0.capabilities.#").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceComputeSkus_withCapabilities(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_compute_skus", "test")
	r := ComputeSkusDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.withCapabilities(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("skus.#").HasValue("1"),
				check.That(data.ResourceName).Key("skus.1").DoesNotExist(),
				check.That(data.ResourceName).Key("skus.0.name").HasValue("Standard_DS2_v2"),
				check.That(data.ResourceName).Key("skus.0.capabilities.%").Exists(),
			),
		},
	})
}

func TestAccDataSourceComputeSkus_allSkus(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_compute_skus", "test")
	r := ComputeSkusDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.allSkus(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("skus.0.name").Exists(),
				check.That(data.ResourceName).Key("skus.1.name").Exists(),
				check.That(data.ResourceName).Key("skus.2.name").Exists(),
			),
		},
	})
}

func (ComputeSkusDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_compute_skus" "test" {
  name     = "Standard_DS2_v2"
  location = "%s"
}
`, data.Locations.Primary)
}

func (ComputeSkusDataSource) withCapabilities(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_compute_skus" "test" {
  name                 = "Standard_DS2_v2"
  location             = "%s"
  include_capabilities = true
}
`, data.Locations.Primary)
}

func (ComputeSkusDataSource) allSkus(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_compute_skus" "test" {
  location = "%s"
}
`, data.Locations.Primary)
}
