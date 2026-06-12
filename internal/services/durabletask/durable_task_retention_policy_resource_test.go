// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/retentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/schedulers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/durabletask"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DurableTaskRetentionPolicyResource struct{}

func TestAccDurableTaskRetentionPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_retention_policy", "test")
	r := DurableTaskRetentionPolicyResource{}

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
	r := DurableTaskRetentionPolicyResource{}

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

func TestAccDurableTaskRetentionPolicy_skipImportCheckOnCreate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_retention_policy", "test")
	r := DurableTaskRetentionPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.template(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.createRetentionPolicyOutsideTerraform(), "azurerm_durable_task_scheduler.test"),
			),
		},
		{
			Config: r.withSkipImportCheck(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_retention_period_in_days").HasValue("30"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDurableTaskRetentionPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_retention_policy", "test")
	r := DurableTaskRetentionPolicyResource{}

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
	r := DurableTaskRetentionPolicyResource{}

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

func (r DurableTaskRetentionPolicyResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := durabletask.ParseRetentionPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	schedulerId := retentionpolicies.NewSchedulerID(id.SubscriptionId, id.ResourceGroupName, id.SchedulerName)

	if _, err = client.DurableTask.RetentionPoliciesClient.Get(ctx, schedulerId); err != nil {
		return nil, fmt.Errorf("retrieving retention policy for %s: %v", schedulerId, err)
	}

	return pointer.To(true), nil
}

func (r DurableTaskRetentionPolicyResource) createRetentionPolicyOutsideTerraform() acceptance.ClientCheckFunc {
	return func(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
		schedulerId, err := schedulers.ParseSchedulerID(state.ID)
		if err != nil {
			return err
		}

		properties := retentionpolicies.RetentionPolicy{
			Properties: &retentionpolicies.RetentionPolicyProperties{
				RetentionPolicies: &[]retentionpolicies.RetentionPolicyDetails{
					{
						RetentionPeriodInDays: 14,
					},
				},
			},
		}

		retentionPolicyId := retentionpolicies.NewSchedulerID(schedulerId.SubscriptionId, schedulerId.ResourceGroupName, schedulerId.SchedulerName)
		if err := client.DurableTask.RetentionPoliciesClient.CreateOrReplaceThenPoll(ctx, retentionPolicyId, properties); err != nil {
			return fmt.Errorf("creating retention policy outside terraform: %+v", err)
		}

		return nil
	}
}

func (r DurableTaskRetentionPolicyResource) template(data acceptance.TestData) string {
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

resource "azurerm_durable_task_scheduler" "test" {
  name                = "acctestdts%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Consumption"
  ip_allow_list       = ["0.0.0.0/0"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r DurableTaskRetentionPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_retention_policy" "test" {
  durable_task_scheduler_id        = azurerm_durable_task_scheduler.test.id
  default_retention_period_in_days = 30
}
`, r.template(data))
}

func (r DurableTaskRetentionPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_retention_policy" "import" {
  durable_task_scheduler_id        = azurerm_durable_task_retention_policy.test.durable_task_scheduler_id
  default_retention_period_in_days = 30
}
`, r.basic(data))
}

func (r DurableTaskRetentionPolicyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_durable_task_retention_policy" "test" {
  durable_task_scheduler_id = azurerm_durable_task_scheduler.test.id

  canceled_retention_period_in_days   = 21
  completed_retention_period_in_days  = 30
  failed_retention_period_in_days     = 7
  terminated_retention_period_in_days = 14
}
`, r.template(data))
}

func (r DurableTaskRetentionPolicyResource) withSkipImportCheck(data acceptance.TestData) string {
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
  ip_allow_list       = ["0.0.0.0/0"]
}

resource "azurerm_durable_task_retention_policy" "test" {
  durable_task_scheduler_id        = azurerm_durable_task_scheduler.test.id
  default_retention_period_in_days = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
