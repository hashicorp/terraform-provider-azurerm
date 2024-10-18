// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AdbsNCharSetsDataSource struct{}

func TestAdbsNCharSetsDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_adbs_national_character_sets", "test")
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

provider "azurerm" {
  features {}
}

data "azurerm_oracle_adbs_national_character_sets" "test" {
  location = "eastus"
}
`, d.template())
}

func (d AdbsNCharSetsDataSource) template() string {
	return fmt.Sprint(`

data "azurerm_client_config" "current" {}

`)
}
