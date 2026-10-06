// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/data-plane/batch/2022-01-01-15-0/jobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccount"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type BatchJobResource struct{}

func TestAccBatchJob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_job", "test")
	r := BatchJobResource{}

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

func TestAccBatchJob_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_job", "test")
	r := BatchJobResource{}

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

func TestAccBatchJob_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_job", "test")
	r := BatchJobResource{}

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

func TestAccBatchJob_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_job", "test")
	r := BatchJobResource{}

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

func (r BatchJobResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.JobID(state.ID)
	if err != nil {
		return nil, err
	}

	endpoint, err := batch.BatchJobResource{}.GetEndpoint(ctx, clients.Batch.AccountClient, batchaccount.NewBatchAccountID(id.SubscriptionId, id.ResourceGroup, id.BatchAccountName))
	if err != nil {
		return nil, err
	}
	client := clients.Batch.JobsClient.Clone(endpoint)

	idSDK := jobs.NewJobID(endpoint, id.Name)

	resp, err := client.JobGet(ctx, idSDK, jobs.DefaultJobGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r BatchJobResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_batch_job" "test" {
  name          = "acctest-job-%[2]d"
  batch_pool_id = azurerm_batch_pool.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r BatchJobResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_batch_job" "test" {
  name          = "acctest-job-%[2]d"
  batch_pool_id = azurerm_batch_pool.test.id
  display_name  = "acctest-job-display-%[2]d"
  common_environment_properties = {
    env       = "Test"
    terraform = "true"
  }
  priority           = 1
  task_retry_maximum = 1
}
`, r.template(data), data.RandomInteger)
}

func (r BatchJobResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_batch_job" "test" {
  name               = "acctest-job-%[2]d"
  batch_pool_id      = azurerm_batch_pool.test.id
  priority           = 2
  task_retry_maximum = -1
}
`, r.template(data), data.RandomInteger)
}

func (r BatchJobResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_batch_job" "import" {
  name          = azurerm_batch_job.test.name
  batch_pool_id = azurerm_batch_job.test.batch_pool_id
}
`, r.basic(data))
}

func (r BatchJobResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-batchjob-%[1]d"
  location = "west europe"
}

resource "azurerm_batch_account" "test" {
  name                = "acctestaccount%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_batch_pool" "test" {
  name                = "acctest-pool-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 22.04"
  vm_size             = "STANDARD_A1_V2"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.RandomString)
}
