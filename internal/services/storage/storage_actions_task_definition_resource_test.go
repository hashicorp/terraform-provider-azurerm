// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storageactions/2023-01-01/storagetasks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StorageActionsTaskDefinitionResource struct{}

func TestAccStorageActionsTaskDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task_definition", "test")
	r := StorageActionsTaskDefinitionResource{}
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

func TestAccStorageActionsTaskDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task_definition", "test")
	r := StorageActionsTaskDefinitionResource{}
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

func TestAccStorageActionsTaskDefinition_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task_definition", "test")
	r := StorageActionsTaskDefinitionResource{}
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

func TestAccStorageActionsTaskDefinition_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task_definition", "test")
	r := StorageActionsTaskDefinitionResource{}
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

func TestAccStorageActionsTaskDefinition_deleteWithOtherOperationsInIf(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task_definition", "test")
	r := StorageActionsTaskDefinitionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.deleteWithOtherOperationsInIf(data),
			ExpectError: regexp.MustCompile("`DeleteBlob` operation cannot be combined with other operations"),
		},
	})
}

func TestAccStorageActionsTaskDefinition_deleteWithOtherOperationsInElse(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task_definition", "test")
	r := StorageActionsTaskDefinitionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.deleteWithOtherOperationsInElse(data),
			ExpectError: regexp.MustCompile("`DeleteBlob` operation cannot be combined with other operations"),
		},
	})
}

func (r StorageActionsTaskDefinitionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r StorageActionsTaskDefinitionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_actions_task_definition" "test" {
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
      condition = "[[endsWith(Name, '.docx')]]"

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

func (r StorageActionsTaskDefinitionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_storage_actions_task_definition" "import" {
  name                = azurerm_storage_actions_task_definition.test.name
  resource_group_name = azurerm_storage_actions_task_definition.test.resource_group_name
  location            = azurerm_storage_actions_task_definition.test.location
  description         = azurerm_storage_actions_task_definition.test.description
  enabled             = azurerm_storage_actions_task_definition.test.enabled

  identity {
    type = "SystemAssigned"
  }

  action {
    if {
      condition = "[[endsWith(Name, '.docx')]]"

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

func (r StorageActionsTaskDefinitionResource) complete(data acceptance.TestData) string {
	config := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_actions_task_definition" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  description         = "complete test"
  enabled             = false

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  action {
    if {
      condition = "[[and(equals(AccessTier, 'Cool'), greater(Content-Length, '100'))]]"

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

func (r StorageActionsTaskDefinitionResource) template(data acceptance.TestData) string {
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

// deleteWithOtherOperationsInIf produces an invalid configuration where a
// `DeleteBlob` operation is combined with another operation in the same `if`
// block. The Azure API rejects this with a `ValidationFailed` error.
func (r StorageActionsTaskDefinitionResource) deleteWithOtherOperationsInIf(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_actions_task_definition" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  description         = "delete with other operations should fail"
  enabled             = true

  identity {
    type = "SystemAssigned"
  }

  action {
    if {
      condition = "[[not(equals(BlobType, 'PageBlob'))]]"

      operation {
        name       = "DeleteBlob"
        on_failure = "break"
        on_success = "continue"
      }

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

// deleteWithOtherOperationsInElse produces an invalid configuration where a
// `DeleteBlob` operation is combined with another operation in the same `else`
// block. The Azure API rejects this with a `ValidationFailed` error.
func (r StorageActionsTaskDefinitionResource) deleteWithOtherOperationsInElse(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_actions_task_definition" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  description         = "delete with other operations should fail"
  enabled             = true

  identity {
    type = "SystemAssigned"
  }

  action {
    if {
      condition = "[[less(LastAccessTime, '2024-01-01T00:00:00Z')]]"

      operation {
        name       = "SetBlobTier"
        on_failure = "break"
        on_success = "continue"

        parameters = {
          tier = "Hot"
        }
      }
    }

    else {
      operation {
        name       = "DeleteBlob"
        on_failure = "break"
        on_success = "continue"
      }

      operation {
        name       = "SetBlobTags"
        on_failure = "break"
        on_success = "continue"

        parameters = {
          archived = "true"
        }
      }
    }
  }
}
`, template, data.RandomString)
}
