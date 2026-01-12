// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagemover_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageMoverSmbFileShareEndpointTestResource struct{}

func TestAccStorageMoverSmbFileShareEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_smb_file_share_endpoint", "test")
	r := StorageMoverSmbFileShareEndpointTestResource{}
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

func TestAccStorageMoverSmbFileShareEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_smb_file_share_endpoint", "test")
	r := StorageMoverSmbFileShareEndpointTestResource{}
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

func TestAccStorageMoverSmbFileShareEndpoint_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_smb_file_share_endpoint", "test")
	r := StorageMoverSmbFileShareEndpointTestResource{}
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

func TestAccStorageMoverSmbFileShareEndpoint_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_smb_file_share_endpoint", "test")
	r := StorageMoverSmbFileShareEndpointTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
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
	})
}

func (r StorageMoverSmbFileShareEndpointTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := endpoints.ParseEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.StorageMover.EndpointsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r StorageMoverSmbFileShareEndpointTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "test" {
  name               = "accshare%s"
  storage_account_id = azurerm_storage_account.test.id
  quota              = 50
}

resource "azurerm_storage_mover" "test" {
  name                = "acctest-ssm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger)
}

func (r StorageMoverSmbFileShareEndpointTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_mover_smb_file_share_endpoint" "test" {
  name               = "acctest-smsfse-%d"
  storage_mover_id   = azurerm_storage_mover.test.id
  storage_account_id = azurerm_storage_account.test.id
  file_share_name    = azurerm_storage_share.test.name
}
`, template, data.RandomInteger)
}

func (r StorageMoverSmbFileShareEndpointTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_mover_smb_file_share_endpoint" "import" {
  name               = azurerm_storage_mover_smb_file_share_endpoint.test.name
  storage_mover_id   = azurerm_storage_mover.test.id
  storage_account_id = azurerm_storage_account.test.id
  file_share_name    = azurerm_storage_share.test.name
}
`, config)
}

func (r StorageMoverSmbFileShareEndpointTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_mover_smb_file_share_endpoint" "test" {
  name               = "acctest-smsfse-%d"
  storage_mover_id   = azurerm_storage_mover.test.id
  storage_account_id = azurerm_storage_account.test.id
  file_share_name    = azurerm_storage_share.test.name
  description        = "Example SMB File Share Endpoint Description"
}
`, template, data.RandomInteger)
}

func (r StorageMoverSmbFileShareEndpointTestResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_mover_smb_file_share_endpoint" "test" {
  name               = "acctest-smsfse-%d"
  storage_mover_id   = azurerm_storage_mover.test.id
  storage_account_id = azurerm_storage_account.test.id
  file_share_name    = azurerm_storage_share.test.name
  description        = "Updated SMB File Share Endpoint Description"
}
`, template, data.RandomInteger)
}
