// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SchedulerDataSource struct{}

func TestAccDurableTaskSchedulerDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_durable_task_scheduler", "test")
	r := SchedulerDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("sku_name").HasValue("Consumption"),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("endpoint").Exists(),
				check.That(data.ResourceName).Key("ip_allow_list.#").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccDurableTaskSchedulerDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_durable_task_scheduler", "test")
	r := SchedulerDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("sku_name").HasValue("Dedicated"),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("endpoint").Exists(),
				check.That(data.ResourceName).Key("ip_allow_list.#").HasValue("2"),
				check.That(data.ResourceName).Key("capacity").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("test"),
				check.That(data.ResourceName).Key("tags.purpose").HasValue("acceptance-testing"),
			),
		},
	})
}

func (r SchedulerDataSource) basic(data acceptance.TestData) string {
	template := SchedulerResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_durable_task_scheduler" "test" {
  name                = azurerm_durable_task_scheduler.test.name
  resource_group_name = azurerm_durable_task_scheduler.test.resource_group_name
}
`, template)
}

func (r SchedulerDataSource) complete(data acceptance.TestData) string {
	template := SchedulerResource{}.complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_durable_task_scheduler" "test" {
  name                = azurerm_durable_task_scheduler.test.name
  resource_group_name = azurerm_durable_task_scheduler.test.resource_group_name
}
`, template)
}
