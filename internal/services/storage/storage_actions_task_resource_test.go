// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storageactions/2023-01-01/storagetasks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StorageActionsTaskResource struct{}

func TestAccStorageActionsTask_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task", "test")
	r := StorageActionsTaskResource{}
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

func (r StorageActionsTaskResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_actions_task" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  description         = "basic test"
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
`, template, data.RandomString)
}

func TestAccStorageActionsTask_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task", "test")
	r := StorageActionsTaskResource{}
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

func (r StorageActionsTaskResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_storage_actions_task" "import" {
  name                = azurerm_storage_actions_task.test.name
  resource_group_name = azurerm_storage_actions_task.test.resource_group_name
  location            = azurerm_storage_actions_task.test.location
  description         = azurerm_storage_actions_task.test.description
  enabled             = azurerm_storage_actions_task.test.enabled

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
`, config)
}

func TestAccStorageActionsTask_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task", "test")
	r := StorageActionsTaskResource{}
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

func (r StorageActionsTaskResource) complete(data acceptance.TestData) string {
	config := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_actions_task" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  description         = "complete test"
  enabled             = true

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
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

      operation {
        name       = "SetBlobTags"
        on_failure = "break"
        on_success = "continue"

        parameters = {
          processed = "true"
        }
      }
    }

    else {
      operation {
        name       = "DeleteBlob"
        on_failure = "break"
        on_success = "continue"
      }
    }
  }

  tags = {
    environment = "test"
  }
}
`, config, data.RandomString)
}

func TestAccStorageActionsTask_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task", "test")
	r := StorageActionsTaskResource{}
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

func (r StorageActionsTaskResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := storagetasks.ParseStorageTaskID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Storage.StorageTasksClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r StorageActionsTaskResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
