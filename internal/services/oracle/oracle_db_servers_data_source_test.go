// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DBServersDataSource struct{}

func TestDBServersDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_db_servers", "test")
	r := DBServersDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("db_servers.0.compartment_id").Exists(),
			),
		},
	})
}

func (d DBServersDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_db_servers" "test" {
  resource_group_name               = azurerm_resource_group.test.name
  cloud_exadata_infrastructure_name = azurerm_oracle_exadata_infrastructure.test.name
}
`, d.template(data))
}

func (d DBServersDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

`, ExadataInfraResource{}.basic(data))
}
