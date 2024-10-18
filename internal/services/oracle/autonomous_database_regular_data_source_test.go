// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
)

type AutonomousDatabaseRegularDataSource struct{}

func TestAdbsRegularDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.AutonomousDatabaseRegularDataSource{}.ResourceType(), "test")
	r := AutonomousDatabaseRegularDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("data_storage_size_in_tbs").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("license_model").Exists(),
			),
		},
	})
}

func (d AutonomousDatabaseRegularDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_autonomous_database" "test" {
  name                = azurerm_oracle_autonomous_database.test.name
  resource_group_name = azurerm_oracle_autonomous_database.test.resource_group_name
}
`, AdbsRegularResource{}.basic(data))
}
