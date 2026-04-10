// Copyright IBM Corp. 2014, 2025
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

type StorageMoverNfsFileShareTargetEndpointTestResource struct{}

func TestAccStorageMoverNfsFileShareTargetEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_nfs_file_share_target_endpoint", "test")
	r := StorageMoverNfsFileShareTargetEndpointTestResource{}
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

func TestAccStorageMoverNfsFileShareTargetEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_nfs_file_share_target_endpoint", "test")
	r := StorageMoverNfsFileShareTargetEndpointTestResource{}
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

func TestAccStorageMoverNfsFileShareTargetEndpoint_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_nfs_file_share_target_endpoint", "test")
	r := StorageMoverNfsFileShareTargetEndpointTestResource{}
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

func TestAccStorageMoverNfsFileShareTargetEndpoint_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_nfs_file_share_target_endpoint", "test")
	r := StorageMoverNfsFileShareTargetEndpointTestResource{}
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

func (r StorageMoverNfsFileShareTargetEndpointTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r StorageMoverNfsFileShareTargetEndpointTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Premium"
  account_replication_type = "LRS"
  account_kind             = "FileStorage"
}

resource "azurerm_storage_share" "test" {
  name               = "acctestshare%s"
  storage_account_id = azurerm_storage_account.test.id
  enabled_protocol   = "NFS"
  quota              = 100
}

resource "azurerm_storage_mover" "test" {
  name                = "acctest-ssm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger)
}

func (r StorageMoverNfsFileShareTargetEndpointTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_mover_nfs_file_share_target_endpoint" "test" {
  name               = "acctest-smnfse-%d"
  storage_mover_id   = azurerm_storage_mover.test.id
  storage_account_id = azurerm_storage_account.test.id
  file_share_name    = azurerm_storage_share.test.name
}
`, r.template(data), data.RandomInteger)
}

func (r StorageMoverNfsFileShareTargetEndpointTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_mover_nfs_file_share_target_endpoint" "import" {
  name               = azurerm_storage_mover_nfs_file_share_target_endpoint.test.name
  storage_mover_id   = azurerm_storage_mover.test.id
  storage_account_id = azurerm_storage_account.test.id
  file_share_name    = azurerm_storage_share.test.name
}
`, r.basic(data))
}

func (r StorageMoverNfsFileShareTargetEndpointTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_mover_nfs_file_share_target_endpoint" "test" {
  name               = "acctest-smnfse-%d"
  storage_mover_id   = azurerm_storage_mover.test.id
  storage_account_id = azurerm_storage_account.test.id
  file_share_name    = azurerm_storage_share.test.name
  description        = "Example NFS File Share Endpoint Description"
}
`, r.template(data), data.RandomInteger)
}

func (r StorageMoverNfsFileShareTargetEndpointTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_mover_nfs_file_share_target_endpoint" "test" {
  name               = "acctest-smnfse-%d"
  storage_mover_id   = azurerm_storage_mover.test.id
  storage_account_id = azurerm_storage_account.test.id
  file_share_name    = azurerm_storage_share.test.name
  description        = "Updated NFS File Share Endpoint Description"
}
`, r.template(data), data.RandomInteger)
}
