// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DataFactoryTriggerScheduleDataSource struct{}

func TestAccDataFactoryTriggerScheduleDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_factory_trigger_schedule", "test")
	r := DataFactoryTriggerScheduleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func (DataFactoryTriggerScheduleDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_data_factory_trigger_schedule" "test" {
  name            = azurerm_data_factory_trigger_schedule.test.name
  data_factory_id = azurerm_data_factory.test.id
}
`, TriggerScheduleResource{}.basic(data))
}
