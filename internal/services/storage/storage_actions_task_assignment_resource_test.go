// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storagetaskassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StorageActionsTaskAssignmentResource struct{}

func TestAccStorageActionsTaskAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task_assignment", "test")
	r := StorageActionsTaskAssignmentResource{}
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

func TestAccStorageActionsTaskAssignment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task_assignment", "test")
	r := StorageActionsTaskAssignmentResource{}
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

func TestAccStorageActionsTaskAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task_assignment", "test")
	r := StorageActionsTaskAssignmentResource{}
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

func TestAccStorageActionsTaskAssignment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task_assignment", "test")
	r := StorageActionsTaskAssignmentResource{}
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

func TestAccStorageActionsTaskAssignment_runOnce(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task_assignment", "test")
	r := StorageActionsTaskAssignmentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.runOnce(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageActionsTaskAssignment_mockRun(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task_assignment", "test")
	r := StorageActionsTaskAssignmentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.mockRun(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageActionsTaskAssignmentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := storagetaskassignments.ParseStorageTaskAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Storage.ResourceManager.StorageTaskAssignments
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r StorageActionsTaskAssignmentResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_actions_task_assignment" "test" {
  name               = "acctest%s"
  storage_account_id = azurerm_storage_account.test.id
  task_id            = azurerm_storage_actions_task.test.id
  description        = "basic test"
  enabled            = true

  execution_context {
    trigger {
      type       = "OnSchedule"
      interval   = 1
      start_from = "2030-01-01T00:00:00Z"
      end_by     = "2031-01-01T00:00:00Z"
    }
  }

  report_prefix = "report"
}
`, template, data.RandomString)
}

func (r StorageActionsTaskAssignmentResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_actions_task_assignment" "test" {
  name               = "acctest%s"
  storage_account_id = azurerm_storage_account.test.id
  task_id            = azurerm_storage_actions_task.test.id
  description        = "complete test"
  enabled            = false

  execution_context {
    trigger {
      type       = "OnSchedule"
      interval   = 7
      start_from = "2030-01-01T00:00:00Z"
      end_by     = "2031-01-01T00:00:00Z"
    }

    target {
      prefix         = ["container1/", "container2/prefix/"]
      exclude_prefix = ["container1/skip/"]
    }
  }

  report_prefix = "complete-report"
}
`, template, data.RandomString)
}

func (r StorageActionsTaskAssignmentResource) runOnce(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_actions_task_assignment" "test" {
  name               = "acctest%s"
  storage_account_id = azurerm_storage_account.test.id
  task_id            = azurerm_storage_actions_task.test.id
  description        = "run once test"
  enabled            = true

  execution_context {
    trigger {
      type     = "RunOnce"
      start_on = "2030-01-01T00:00:00Z"
    }
  }

  report_prefix = "report"
}
`, template, data.RandomString)
}

func (r StorageActionsTaskAssignmentResource) mockRun(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_actions_task_assignment" "test" {
  name               = "acctest%s"
  storage_account_id = azurerm_storage_account.test.id
  task_id            = azurerm_storage_actions_task.test.id
  description        = "mock run test"
  enabled            = true

  execution_context {
    trigger {
      type     = "MockRun"
      start_on = "2030-01-01T00:00:00Z"
    }
  }

  report_prefix = "report"
}
`, template, data.RandomString)
}

func (r StorageActionsTaskAssignmentResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_actions_task_assignment" "import" {
  name               = azurerm_storage_actions_task_assignment.test.name
  storage_account_id = azurerm_storage_actions_task_assignment.test.storage_account_id
  task_id            = azurerm_storage_actions_task_assignment.test.task_id
  description        = azurerm_storage_actions_task_assignment.test.description
  enabled            = azurerm_storage_actions_task_assignment.test.enabled

  execution_context {
    trigger {
      type       = "OnSchedule"
      interval   = 1
      start_from = "2030-01-01T00:00:00Z"
      end_by     = "2031-01-01T00:00:00Z"
    }
  }

  report_prefix = "report"
}
`, config)
}

func (r StorageActionsTaskAssignmentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_actions_task" "test" {
  name                = "acctest%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  description         = "for assignment test"
  enabled             = true

  identity {
    type = "SystemAssigned"
  }

  action {
    if {
      condition = "[[equals(AccessTier, 'Cool')]]"

      operation {
        name       = "SetBlobTier"
        on_failure = "break"
        on_success = "continue"

        parameters = {
          tier = "Hot"
        }
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
