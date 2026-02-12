// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagemover_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StorageMoverSmbMountEndpointTestResource struct{}

func TestAccStorageMoverSmbMountEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_smb_mount_endpoint", "test")
	r := StorageMoverSmbMountEndpointTestResource{}
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

func TestAccStorageMoverSmbMountEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_smb_mount_endpoint", "test")
	r := StorageMoverSmbMountEndpointTestResource{}
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

func TestAccStorageMoverSmbMountEndpoint_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_smb_mount_endpoint", "test")
	r := StorageMoverSmbMountEndpointTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password_uri"),
	})
}

func TestAccStorageMoverSmbMountEndpoint_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_smb_mount_endpoint", "test")
	r := StorageMoverSmbMountEndpointTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password_uri"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password_uri"),
	})
}

func (r StorageMoverSmbMountEndpointTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := endpoints.ParseEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.StorageMover.EndpointsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r StorageMoverSmbMountEndpointTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_storage_mover" "test" {
  name                = "acctest-ssm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r StorageMoverSmbMountEndpointTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_mover_smb_mount_endpoint" "test" {
  name             = "acctest-smse-%d"
  storage_mover_id = azurerm_storage_mover.test.id
  host             = "192.168.0.1"
  share_name       = "testshare"
}
`, template, data.RandomInteger)
}

func (r StorageMoverSmbMountEndpointTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_mover_smb_mount_endpoint" "import" {
  name             = azurerm_storage_mover_smb_mount_endpoint.test.name
  storage_mover_id = azurerm_storage_mover.test.id
  host             = azurerm_storage_mover_smb_mount_endpoint.test.host
  share_name       = azurerm_storage_mover_smb_mount_endpoint.test.share_name
}
`, config)
}

func (r StorageMoverSmbMountEndpointTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv%s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Get",
      "Set",
      "Delete",
      "Purge",
    ]
  }
}

resource "azurerm_key_vault_secret" "username" {
  name         = "smb-username"
  value        = "testuser"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_key_vault_secret" "password" {
  name         = "smb-password"
  value        = "testpassword123!"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_storage_mover_smb_mount_endpoint" "test" {
  name             = "acctest-smse-%d"
  storage_mover_id = azurerm_storage_mover.test.id
  host             = "192.168.0.1"
  share_name       = "testshare"
  username_uri     = azurerm_key_vault_secret.username.versionless_id
  password_uri     = azurerm_key_vault_secret.password.versionless_id
  description      = "Example SMB Mount Endpoint Description"
}
`, template, data.RandomString, data.RandomInteger)
}

func (r StorageMoverSmbMountEndpointTestResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv%s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Get",
      "Set",
      "Delete",
      "Purge",
    ]
  }
}

resource "azurerm_key_vault_secret" "username" {
  name         = "smb-username"
  value        = "testuser"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_key_vault_secret" "password" {
  name         = "smb-password"
  value        = "testpassword123!"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_storage_mover_smb_mount_endpoint" "test" {
  name             = "acctest-smse-%d"
  storage_mover_id = azurerm_storage_mover.test.id
  host             = "192.168.0.1"
  share_name       = "testshare"
  username_uri     = azurerm_key_vault_secret.username.versionless_id
  password_uri     = azurerm_key_vault_secret.password.versionless_id
  description      = "Updated SMB Mount Endpoint Description"
}
`, template, data.RandomString, data.RandomInteger)
}
