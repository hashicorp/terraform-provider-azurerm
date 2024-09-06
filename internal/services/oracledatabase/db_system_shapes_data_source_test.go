// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracledatabase_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"testing"
)

type DbSystemShapesDataSource struct{}

func TestDbSystemShapesDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracledatabase_db_system_shapes", "test")
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
	return fmt.Sprintf(`

%s

data "azurerm_oracledatabase_db_system_shapes" "test" {
  location_name = "eastus"
}
`, d.template())
}

func (d DbSystemShapesDataSource) template() string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

`)
}
