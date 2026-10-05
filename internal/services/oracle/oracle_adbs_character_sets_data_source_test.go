// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AdbsCharSetsDataSource struct{}

func TestAdbsCharSetsDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_adbs_character_sets", "test")
	r := AdbsCharSetsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("character_sets.0.character_set").Exists(),
			),
		},
	})
}

func (d AdbsCharSetsDataSource) basic() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_oracle_adbs_character_sets" "test" {
  location = "eastus"
}
`
}
