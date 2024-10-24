// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

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
