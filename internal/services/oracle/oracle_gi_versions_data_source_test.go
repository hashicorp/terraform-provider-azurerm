// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type GiVersionsDataSource struct{}

func TestGiVersionsDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_gi_versions", "test")
	r := GiVersionsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("versions.#").HasValue("10"),
			),
		},
	})
}

func TestGiVersionsDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_gi_versions", "test")
	r := GiVersionsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("versions.#").HasValue("2"),
			),
		},
	})
}

func (d GiVersionsDataSource) basic() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_oracle_gi_versions" "test" {
  location = "eastus"
}
`
}

func (d GiVersionsDataSource) complete() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_oracle_gi_versions" "test" {
  location = "eastus"
  shape    = "Exadata.X9M"
  zone     = "2"
}
`
}
