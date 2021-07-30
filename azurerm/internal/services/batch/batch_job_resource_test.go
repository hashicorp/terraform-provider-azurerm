package batch_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type BatchJobResource struct{}

func TestAccBatchJob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_job", "test")
	r := BatchJobResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBatchJob_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_job", "test")
	r := BatchJobResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBatchJob_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_job", "test")
	r := BatchJobResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBatchJob_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_job", "test")
	r := BatchJobResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r BatchJobResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.JobID(state.ID)
	if err != nil {
		return nil, err
	}

	client, err := clients.Batch.JobClient(ctx, parse.NewAccountID(id.SubscriptionId, id.ResourceGroup, id.BatchAccountName))
	if err != nil {
		return nil, err
	}

	if resp, err := client.Get(ctx, id.Name, "", "", nil, nil, nil, nil, "", "", nil, nil); err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r BatchJobResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_batch_job" "test" {
  name          = "testaccbj-%d"
  batch_pool_id = azurerm_batch_pool.test.id
}
`, template, data.RandomInteger)
}

func (r BatchJobResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_batch_job" "test" {
  name          = "testaccbj-%[2]d"
  batch_pool_id = azurerm_batch_pool.test.id
  display_name  = "testaccbj-display-%[2]d"
  common_environment_properties = {
    env = "Test"
  }
  priority           = 1
  task_retry_maximum = 1
}
`, template, data.RandomInteger)
}

func (r BatchJobResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_batch_job" "test" {
  name               = "testaccbj-%d"
  batch_pool_id      = azurerm_batch_pool.test.id
  priority           = 2
  task_retry_maximum = -1
}
`, template, data.RandomInteger)
}

func (r BatchJobResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_batch_job" "import" {
  name          = azurerm_batch_job.test.name
  batch_pool_id = azurerm_batch_job.test.batch_pool_id
}
`, template)
}

func (r BatchJobResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-batchjob-%[1]d"
  location = "west europe"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }
}
`, data.RandomInteger, data.RandomString)
}
