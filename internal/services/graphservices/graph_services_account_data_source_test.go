// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package graphservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AccountDataSourceTestResource struct{}

func TestAccGraphServicesAccountDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_graph_services_account", "test")
	r := AccountDataSourceTestResource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("application_id").Exists(),
				check.That(data.ResourceName).Key("billing_plan_id").Exists(),
			),
		},
	})
}

func (r AccountDataSourceTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_graph_services_account" "test" {
  name                = azurerm_graph_services_account.test.name
  resource_group_name = azurerm_graph_services_account.test.resource_group_name
}
`, AccountTestResource{}.basic(data))
}