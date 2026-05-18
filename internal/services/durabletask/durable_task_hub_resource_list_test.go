// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccDurableTaskHubList_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_hub", "test")
	r := DurableTaskHubResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicListQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_durable_task_hub.list", 2),
				},
			},
		},
	})
}

func (r DurableTaskHubResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-durabletask-%d"
  location = "%s"
}

resource "azurerm_durable_task_scheduler" "test" {
  name                = "acctestdts%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Consumption"
  ip_allow_list       = ["0.0.0.0/0"]
}

resource "azurerm_durable_task_hub" "test" {
  name         = "acctest%s"
  scheduler_id = azurerm_durable_task_scheduler.test.id
}

resource "azurerm_durable_task_hub" "test2" {
  name         = "acctest%s2"
  scheduler_id = azurerm_durable_task_scheduler.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString)
}

func (r DurableTaskHubResource) basicListQuery(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_durable_task_hub" "list" {
  provider = azurerm
  config {
    scheduler_id = "/subscriptions/%s/resourceGroups/acctestRG-durabletask-%d/providers/Microsoft.DurableTask/schedulers/acctestdts%s"
  }
}
`, data.Subscriptions.Primary, data.RandomInteger, data.RandomString)
}
