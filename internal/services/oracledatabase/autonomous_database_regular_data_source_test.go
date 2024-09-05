// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracledatabase_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracledatabase"
	"testing"
)

type AutonomousDatabaseRegularDataSource struct{}

func TestAdbsRegularDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracledatabase.AutonomousDatabaseRegularDataSource{}.ResourceType(), "test")
	r := AutonomousDatabaseRegularDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("data_storage_size_in_gbs").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("license_model").Exists(),
			),
		},
	})
}

func (d AutonomousDatabaseRegularDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracledatabase_autonomous_database_regular" "test" {
  name = azurerm_oracledatabase_autonomous_database_regular.test.name
  resource_group_name = azurerm_oracledatabase_autonomous_database_regular.test.resource_group_name
}
`, AdbsRegularResource{}.basic(data))
}
