// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/schedulers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DurableTaskSchedulerResource struct{}

func TestAccDurableTaskScheduler_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := DurableTaskSchedulerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDurableTaskScheduler_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := DurableTaskSchedulerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDurableTaskScheduler_skipImportCheckOnCreate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := DurableTaskSchedulerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.template(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.createSchedulerOutsideTerraform(data), "azurerm_resource_group.test"),
			),
		},
		{
			Config: r.withSkipImportCheck(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_allow_list.#").HasValue("1"),
				check.That(data.ResourceName).Key("ip_allow_list.0").HasValue("10.0.0.0/8"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDurableTaskScheduler_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := DurableTaskSchedulerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDurableTaskScheduler_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := DurableTaskSchedulerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
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

func TestAccDurableTaskScheduler_dedicatedWithCapacity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := DurableTaskSchedulerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dedicatedWithCapacity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDurableTaskScheduler_dedicatedWithoutCapacityFails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := DurableTaskSchedulerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
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

func TestAccDurableTaskScheduler_consumptionWithCapacityFails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := DurableTaskSchedulerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.consumptionWithCapacity(data),
			ExpectError: regexp.MustCompile("`capacity` can only be configured when `sku_name` is set to `Dedicated`"),
		},
	})
}

func TestAccDurableTaskScheduler_dedicatedWithTooMuchCapacityFails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := DurableTaskSchedulerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.dedicatedWithTooMuchCapacity(data),
			ExpectError: regexp.MustCompile(`expected capacity to be in the range \(1 - 3\), got 4`),
		},
	})
}

func (r DurableTaskSchedulerResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := schedulers.ParseSchedulerID(state.ID)
	if err != nil {
		return nil, err
	}

	if _, err = client.DurableTask.SchedulersClient.Get(ctx, *id); err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id, err)
	}

	return pointer.To(true), nil
}

func (r DurableTaskSchedulerResource) createSchedulerOutsideTerraform(data acceptance.TestData) acceptance.ClientCheckFunc {
	return func(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
		id, err := schedulers.ParseSchedulerID(fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DurableTask/schedulers/acctestdts%s", data.Subscriptions.Primary, state.Attributes["name"], data.RandomString))
		if err != nil {
			return err
		}

		properties := schedulers.Scheduler{
			Location: location.Normalize(state.Attributes["location"]),
			Properties: &schedulers.SchedulerProperties{
				Sku: schedulers.SchedulerSku{
					Name: schedulers.SchedulerSkuNameConsumption,
				},
				IPAllowlist: []string{"0.0.0.0/0"},
			},
		}

		if err := client.DurableTask.SchedulersClient.CreateOrUpdateThenPoll(ctx, *id, properties); err != nil {
			return fmt.Errorf("creating scheduler outside terraform: %+v", err)
		}

		return nil
	}
}

func (r DurableTaskSchedulerResource) template(data acceptance.TestData) string {
	// Durable Task schedulers are only supported in specific regions.
	data.Locations.Primary = "northeurope"

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
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_scheduler" "test" {
  name                = "acctestdts%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Consumption"
  ip_allow_list       = ["0.0.0.0/0"]
}
`, template, data.RandomString)
}

func (r DurableTaskSchedulerResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_scheduler" "import" {
  name                = azurerm_durable_task_scheduler.test.name
  resource_group_name = azurerm_durable_task_scheduler.test.resource_group_name
  location            = azurerm_durable_task_scheduler.test.location
  sku_name            = azurerm_durable_task_scheduler.test.sku_name
  ip_allow_list       = azurerm_durable_task_scheduler.test.ip_allow_list
}
`, template)
}

func (r DurableTaskSchedulerResource) withSkipImportCheck(data acceptance.TestData) string {
	// Durable Task schedulers are only supported in specific regions.
	data.Locations.Primary = "northeurope"

	return fmt.Sprintf(`
provider "azurerm" {
  features {
    skip_import_check_on_create_and_allow_overwriting_existing_resources = true
  }
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
  ip_allow_list       = ["10.0.0.0/8"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r DurableTaskSchedulerResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_scheduler" "test" {
  name                = "acctestdts%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Consumption"
  ip_allow_list       = ["10.0.0.0/8", "192.168.0.0/16"]

  tags = {
    environment = "staging"
  }
}
`, template, data.RandomString)
}

func (r DurableTaskSchedulerResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_scheduler" "test" {
  name                = "acctestdts%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Dedicated"
  ip_allow_list       = ["10.0.0.0/8", "192.168.0.0/16"]
  capacity            = 2

  tags = {
    environment = "test"
    purpose     = "acceptance-testing"
  }
}
`, template, data.RandomString)
}

func (r DurableTaskSchedulerResource) dedicatedWithCapacity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_scheduler" "test" {
  name                = "acctestdts%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Dedicated"
  ip_allow_list       = ["0.0.0.0/0"]
  capacity            = 2
}
`, template, data.RandomString)
}

func (r DurableTaskSchedulerResource) dedicatedWithoutCapacity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_scheduler" "test" {
  name                = "acctestdts%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Dedicated"
  ip_allow_list       = ["0.0.0.0/0"]
}
`, template, data.RandomString)
}

func (r DurableTaskSchedulerResource) consumptionWithCapacity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_scheduler" "test" {
  name                = "acctestdts%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Consumption"
  ip_allow_list       = ["0.0.0.0/0"]
  capacity            = 1
}
`, template, data.RandomString)
}

func (r DurableTaskSchedulerResource) dedicatedWithTooMuchCapacity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_scheduler" "test" {
  name                = "acctestdts%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Dedicated"
  ip_allow_list       = ["0.0.0.0/0"]
  capacity            = 4
}
`, template, data.RandomString)
}
