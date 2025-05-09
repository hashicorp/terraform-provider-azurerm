// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

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
