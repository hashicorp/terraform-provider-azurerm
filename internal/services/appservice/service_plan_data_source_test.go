// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ServicePlanDataSource struct{}

func TestAccServicePlanDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_service_plan", "test")
	d := ServicePlanDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				// TODO - rest of the sane properties
			),
		},
	})
}

func (ServicePlanDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data azurerm_service_plan test {
  name                = azurerm_service_plan.test.name
  resource_group_name = azurerm_service_plan.test.resource_group_name
}
`, ServicePlanResource{}.complete(data))
}
