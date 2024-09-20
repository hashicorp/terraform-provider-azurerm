// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package subscription_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ExtendedLocationsDataSource struct{}

func TestAccDataSourceExtendedLocations_westEurope(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_extended_locations", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config:      ExtendedLocationsDataSource{}.basic("westeurope"),
			ExpectError: regexp.MustCompile("no extended locations were found for the location \"westeurope\""),
		},
	})
}

func TestAccDataSourceExtendedLocations_westUS(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_extended_locations", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: ExtendedLocationsDataSource{}.basic("westus"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("extended_locations.#").HasValue("1"),
				check.That(data.ResourceName).Key("extended_locations.0").HasValue("losangeles"),
			),
		},
	})
}

func (d ExtendedLocationsDataSource) basic(location string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_extended_locations" "test" {
  location = "%s"
}
`, location)
}
