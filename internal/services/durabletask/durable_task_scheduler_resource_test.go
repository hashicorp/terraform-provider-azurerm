// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package durabletask_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/schedulers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SchedulerResource struct{}

func TestAccDurableTaskScheduler_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := SchedulerResource{}

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
	r := SchedulerResource{}

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

func TestAccDurableTaskScheduler_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := SchedulerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku_name").HasValue("Dedicated"),
				check.That(data.ResourceName).Key("capacity").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDurableTaskScheduler_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")
	r := SchedulerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
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

func (r SchedulerResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := schedulers.ParseSchedulerID(state.ID)
	if err != nil {
		return nil, err
	}

	_, err = client.DurableTask.SchedulersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id, err)
	}

	return pointer.To(true), nil
}

func (r SchedulerResource) basic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r SchedulerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_scheduler" "import" {
  name                = azurerm_durable_task_scheduler.test.name
  resource_group_name = azurerm_durable_task_scheduler.test.resource_group_name
  location            = azurerm_durable_task_scheduler.test.location
  sku_name            = azurerm_durable_task_scheduler.test.sku_name
  ip_allow_list       = azurerm_durable_task_scheduler.test.ip_allow_list
}
`, r.basic(data))
}

func (r SchedulerResource) complete(data acceptance.TestData) string {
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
  sku_name            = "Dedicated"
  ip_allow_list       = ["10.0.0.0/8", "192.168.0.0/16"]
  capacity            = 1

  tags = {
    environment = "test"
    purpose     = "acceptance-testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
