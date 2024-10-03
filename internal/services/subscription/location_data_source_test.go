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

type LocationsDataSource struct{}

func TestAccLocationDataSource_NonExistingRegion(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_location", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config:      LocationsDataSource{}.basic("not-existing-region"),
			ExpectError: regexp.MustCompile("\"not-existing-region\" was not found in the list of supported Azure Locations"),
		},
	})
}

func TestAccLocationDataSource_eastUS(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_location", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: LocationsDataSource{}.basic("eastus"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue("East US"),
				check.That(data.ResourceName).Key("zone_mappings.0.logical_zone").HasValue("1"),
				check.That(data.ResourceName).Key("zone_mappings.1.logical_zone").HasValue("2"),
				check.That(data.ResourceName).Key("zone_mappings.2.logical_zone").HasValue("3"),
				check.That(data.ResourceName).Key("zone_mappings.0.physical_zone").IsNotEmpty(),
				check.That(data.ResourceName).Key("zone_mappings.1.physical_zone").IsNotEmpty(),
				check.That(data.ResourceName).Key("zone_mappings.2.physical_zone").IsNotEmpty(),
			),
		},
	})
}

func (d LocationsDataSource) basic(location string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_location" "test" {
  location = "%s"
}
`, location)
}
