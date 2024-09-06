// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracledatabase_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"testing"
)

type GiVersionsDataSource struct{}

func TestGiVersionsDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracledatabase_gi_versions", "test")
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
	return fmt.Sprintf(`

%s

data "azurerm_oracledatabase_gi_versions" "test" {
  location_name = "eastus"
}
`, d.template())
}

func (d GiVersionsDataSource) template() string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

`)
}
