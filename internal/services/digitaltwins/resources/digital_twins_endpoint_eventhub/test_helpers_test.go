// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package digital_twins_endpoint_eventhub_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type DigitalTwinsInstanceResource struct{}

func (r DigitalTwinsInstanceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_instance" "test" {
  name                = "acctest-DT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, r.template(data), data.RandomInteger)
}

func (DigitalTwinsInstanceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-dtwin-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
