// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracledatabase_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"testing"
)

type AdbsCharSetsDataSource struct{}

func TestAdbsCharSetsDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracledatabase_adbs_character_sets", "test")
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
	return fmt.Sprintf(`

%s

data "azurerm_oracledatabase_adbs_character_sets" "test" {
  location_name = "eastus"
}
`, d.template())
}

func (d AdbsCharSetsDataSource) template() string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

`)
}
