// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/schedulers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DurableTaskSchedulerResource struct{}

// TestAccDurableTask runs all of the Durable Task acceptance tests sequentially.
//
// NOTE: these tests are combined into a single sequential test rather than run in
// parallel because the number of Durable Task Schedulers that can exist per
// subscription and region is capped by a low quota, so provisioning multiple
// schedulers concurrently exhausts that quota and causes the tests to fail.
func TestAccDurableTask(t *testing.T) {
	testCases := map[string]map[string]func(t *testing.T){
		"scheduler": {
			"basic":                    testAccDurableTaskScheduler_basic,
			"requiresImport":           testAccDurableTaskScheduler_requiresImport,
			"complete":                 testAccDurableTaskScheduler_complete,
			"update":                   testAccDurableTaskScheduler_update,
			"dedicatedWithCapacity":    testAccDurableTaskScheduler_dedicatedWithCapacity,
			"dedicatedWithoutCapacity": testAccDurableTaskScheduler_dedicatedWithoutCapacityFails,
			"consumptionWithCapacity":  testAccDurableTaskScheduler_consumptionWithCapacityFails,
			"resourceIdentity":         testAccDurableTaskScheduler_resourceIdentity,
		},
		"schedulerDataSource": {
			"basic":    testAccDurableTaskSchedulerDataSource_basic,
			"complete": testAccDurableTaskSchedulerDataSource_complete,
		},
		"schedulerList": {
			"basic": testAccDurableTaskSchedulerList_basic,
		},
		"hub": {
			"basic":            testAccDurableTaskHub_basic,
			"requiresImport":   testAccDurableTaskHub_requiresImport,
			"resourceIdentity": testAccDurableTaskHub_resourceIdentity,
		},
		"hubList": {
			"basic": testAccDurableTaskHubList_basic,
		},
		"retentionPolicy": {
			"basic":            testAccDurableTaskRetentionPolicy_basic,
			"requiresImport":   testAccDurableTaskRetentionPolicy_requiresImport,
			"complete":         testAccDurableTaskRetentionPolicy_complete,
			"update":           testAccDurableTaskRetentionPolicy_update,
			"resourceIdentity": testAccDurableTaskRetentionPolicy_resourceIdentity,
		},
	}

	for group, tests := range testCases {
		t.Run(group, func(t *testing.T) {
			for name, tc := range tests {
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccDurableTaskScheduler_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := DurableTaskSchedulerResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccDurableTaskScheduler_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := DurableTaskSchedulerResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccDurableTaskScheduler_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := DurableTaskSchedulerResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccDurableTaskScheduler_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := DurableTaskSchedulerResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccDurableTaskScheduler_dedicatedWithCapacity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := DurableTaskSchedulerResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.dedicatedWithCapacity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccDurableTaskScheduler_dedicatedWithoutCapacityFails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := DurableTaskSchedulerResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.dedicatedWithCapacity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config:      r.dedicatedWithoutCapacity(data),
			ExpectError: regexp.MustCompile("`capacity` must be configured when `sku_name` is set to `Dedicated`"),
		},
	})
}

func testAccDurableTaskScheduler_consumptionWithCapacityFails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := DurableTaskSchedulerResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config:      r.consumptionWithCapacity(data),
			ExpectError: regexp.MustCompile("`capacity` can only be configured when `sku_name` is set to `Dedicated`"),
		},
	})
}

func (r DurableTaskSchedulerResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := schedulers.ParseSchedulerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.DurableTask.SchedulersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r DurableTaskSchedulerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-durabletask-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r DurableTaskSchedulerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_scheduler" "test" {
  name                = "acctestdts%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Consumption"
  ip_allowlist        = ["0.0.0.0/0"]
}
`, r.template(data), data.RandomString)
}

func (r DurableTaskSchedulerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_scheduler" "import" {
  name                = azurerm_durable_task_scheduler.test.name
  resource_group_name = azurerm_durable_task_scheduler.test.resource_group_name
  location            = azurerm_durable_task_scheduler.test.location
  sku_name            = azurerm_durable_task_scheduler.test.sku_name
  ip_allowlist        = azurerm_durable_task_scheduler.test.ip_allowlist
}
`, r.basic(data))
}

func (r DurableTaskSchedulerResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_scheduler" "test" {
  name                = "acctestdts%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Consumption"
  ip_allowlist        = ["10.0.0.0/8", "192.168.0.0/16"]

  tags = {
    environment = "staging"
  }
}
`, r.template(data), data.RandomString)
}

func (r DurableTaskSchedulerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_scheduler" "test" {
  name                = "acctestdts%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Dedicated"
  ip_allowlist        = ["10.0.0.0/8", "192.168.0.0/16"]
  capacity            = 2

  tags = {
    environment = "test"
    purpose     = "acceptance-testing"
  }
}
`, r.template(data), data.RandomString)
}

func (r DurableTaskSchedulerResource) dedicatedWithCapacity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_scheduler" "test" {
  name                = "acctestdts%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Dedicated"
  ip_allowlist        = ["0.0.0.0/0"]
  capacity            = 2
}
`, r.template(data), data.RandomString)
}

func (r DurableTaskSchedulerResource) dedicatedWithoutCapacity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_scheduler" "test" {
  name                = "acctestdts%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Dedicated"
  ip_allowlist        = ["0.0.0.0/0"]
}
`, r.template(data), data.RandomString)
}

func (r DurableTaskSchedulerResource) consumptionWithCapacity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_scheduler" "test" {
  name                = "acctestdts%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Consumption"
  ip_allowlist        = ["0.0.0.0/0"]
  capacity            = 1
}
`, r.template(data), data.RandomString)
}
