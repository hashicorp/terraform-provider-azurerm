// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracledatabase_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"testing"
)

type AdbsNCharSetsDataSource struct{}

func TestAdbsNCharSetsDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracledatabase_adbs_national_character_sets", "test")
	r := AdbsNCharSetsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("character_sets.0.character_set").Exists(),
			),
		},
	})
}

func (d AdbsNCharSetsDataSource) basic() string {
	return fmt.Sprintf(`

%s

data "azurerm_oracledatabase_adbs_national_character_sets" "test" {
  location_name = "eastus"
}
`, d.template())
}

func (d AdbsNCharSetsDataSource) template() string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

`)
}
