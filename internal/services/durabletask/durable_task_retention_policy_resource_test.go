// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package durabletask_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/retentionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type RetentionPolicyResource struct{}

func TestAccDurableTaskRetentionPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_retention_policy", "test")
	r := RetentionPolicyResource{}

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

func TestAccDurableTaskRetentionPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_retention_policy", "test")
	r := RetentionPolicyResource{}

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

func TestAccDurableTaskRetentionPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_retention_policy", "test")
	r := RetentionPolicyResource{}

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

func TestAccDurableTaskRetentionPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_retention_policy", "test")
	r := RetentionPolicyResource{}

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

func (r RetentionPolicyResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	// The retention policy uses the scheduler ID for API calls
	// Parse our synthetic retention policy ID to get the scheduler details
	id, err := retentionpolicies.ParseSchedulerID(state.Attributes["scheduler_id"])
	if err != nil {
		return nil, err
	}

	_, err = client.DurableTask.RetentionPoliciesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving retention policy for %s: %v", id, err)
	}

	return pointer.To(true), nil
}

func (r RetentionPolicyResource) basic(data acceptance.TestData) string {
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

resource "azurerm_durable_task_retention_policy" "test" {
  scheduler_id = azurerm_durable_task_scheduler.test.id

  retention_policy {
    retention_period_in_days = 30
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r RetentionPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_retention_policy" "import" {
  scheduler_id = azurerm_durable_task_retention_policy.test.scheduler_id

  retention_policy {
    retention_period_in_days = 30
  }
}
`, r.basic(data))
}

func (r RetentionPolicyResource) complete(data acceptance.TestData) string {
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

resource "azurerm_durable_task_retention_policy" "test" {
  scheduler_id = azurerm_durable_task_scheduler.test.id

  retention_policy {
    retention_period_in_days = 30
    orchestration_state      = "Completed"
  }

  retention_policy {
    retention_period_in_days = 7
    orchestration_state      = "Failed"
  }

  retention_policy {
    retention_period_in_days = 14
    orchestration_state      = "Terminated"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
