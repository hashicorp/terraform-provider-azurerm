// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DbSystemShapesDataSource struct{}

func TestDbSystemShapesDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_db_system_shapes", "test")
	r := DbSystemShapesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("db_system_shapes.0.available_core_count").Exists(),
			),
		},
	})
}

func TestDbSystemShapesDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_db_system_shapes", "test")
	r := DbSystemShapesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("db_system_shapes.0.available_core_count").Exists(),
			),
		},
	})
}

func (d DbSystemShapesDataSource) basic() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_oracle_db_system_shapes" "test" {
  location = "eastus"
}
`
}

func (d DbSystemShapesDataSource) complete() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_oracle_db_system_shapes" "test" {
  location = "eastus"
  zone     = "2"
}
`
}
