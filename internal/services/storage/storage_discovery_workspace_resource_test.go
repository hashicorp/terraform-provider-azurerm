// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagediscovery/2025-09-01/storagediscoveryworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StorageDiscoveryWorkspaceResource struct{}

func TestAccStorageDiscoveryWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_discovery_workspace", "test")
	r := StorageDiscoveryWorkspaceResource{}

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

func TestAccStorageDiscoveryWorkspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_discovery_workspace", "test")
	r := StorageDiscoveryWorkspaceResource{}

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

func TestAccStorageDiscoveryWorkspace_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_discovery_workspace", "test")
	r := StorageDiscoveryWorkspaceResource{}

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

func (r StorageDiscoveryWorkspaceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := storagediscoveryworkspaces.ParseProviderStorageDiscoveryWorkspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Storage.StorageDiscoveryWorkspacesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r StorageDiscoveryWorkspaceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accteststorage%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

data "azurerm_subscription" "current" {}

resource "azurerm_storage_discovery_workspace" "test" {
  name                = "acctestsdw-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  workspace_roots = [data.azurerm_subscription.current.id]

  scopes {
    display_name = "TestScope"
    resource_types = [
      azurerm_storage_account.test.id,
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageDiscoveryWorkspaceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accteststorage%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

data "azurerm_subscription" "current" {}

resource "azurerm_storage_discovery_workspace" "test" {
  name                = "acctestsdw-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  description         = "Test Storage Discovery Workspace"
  sku                 = "Standard"

  workspace_roots = [data.azurerm_subscription.current.id]

  scopes {
    display_name = "TestScope1"
    resource_types = [
      azurerm_storage_account.test.id,
    ]
    tag_keys_only = ["tag1", "tag2"]
    tags = {
      tag3 = "value3"
    }
  }

  scopes {
    display_name = "TestScope2"
    resource_types = [
      azurerm_storage_account.test.id,
    ]
  }

  tags = {
    environment = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
