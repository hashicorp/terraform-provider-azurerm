// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storagecache_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2025-07-01/autoimportjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagedLustreFileSystemAutoImportJobResource struct{}

func (r ManagedLustreFileSystemAutoImportJobResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := autoimportjobs.ParseAutoImportJobID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.StorageCache.AutoImportJobs
	resp, err := client.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestAccManagedLustreFileSystemAutoImportJob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_lustre_file_system_auto_import_job", "test")
	r := ManagedLustreFileSystemAutoImportJobResource{}

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

func TestAccManagedLustreFileSystemAutoImportJob_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_lustre_file_system_auto_import_job", "test")
	r := ManagedLustreFileSystemAutoImportJobResource{}

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

func TestAccManagedLustreFileSystemAutoImportJob_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_lustre_file_system_auto_import_job", "test")
	r := ManagedLustreFileSystemAutoImportJobResource{}

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
	})
}

func TestAccManagedLustreFileSystemAutoImportJob_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_lustre_file_system_auto_import_job", "test")
	r := ManagedLustreFileSystemAutoImportJobResource{}

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

func (r ManagedLustreFileSystemAutoImportJobResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_lustre_file_system_auto_import_job" "test" {
  name                          = "acctest-amlfsaij-%d"
  managed_lustre_file_system_id = azurerm_managed_lustre_file_system.test.id
  location                      = azurerm_resource_group.test.location

  auto_import_prefixes = ["/"]
}
`, ManagedLustreFileSystemResource{}.complete(data), data.RandomInteger)
}

func (r ManagedLustreFileSystemAutoImportJobResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_lustre_file_system_auto_import_job" "test" {
  name                          = "acctest-amlfsaij-%d"
  managed_lustre_file_system_id = azurerm_managed_lustre_file_system.test.id
  location                      = azurerm_resource_group.test.location

  auto_import_prefixes = ["/"]
  admin_status_enabled = false
}
`, ManagedLustreFileSystemResource{}.complete(data), data.RandomInteger)
}

func (r ManagedLustreFileSystemAutoImportJobResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_lustre_file_system_auto_import_job" "test" {
  name                          = "acctest-amlfsaij-%d"
  managed_lustre_file_system_id = azurerm_managed_lustre_file_system.test.id
  location                      = azurerm_resource_group.test.location

  auto_import_prefixes     = ["/", "/import"]
  conflict_resolution_mode = "Skip"
  deletions_enabled        = true
  maximum_errors           = 10
  admin_status_enabled     = false
}
`, ManagedLustreFileSystemResource{}.complete(data), data.RandomInteger)
}

func (r ManagedLustreFileSystemAutoImportJobResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_lustre_file_system_auto_import_job" "import" {
  name                          = azurerm_managed_lustre_file_system_auto_import_job.test.name
  managed_lustre_file_system_id = azurerm_managed_lustre_file_system_auto_import_job.test.managed_lustre_file_system_id
  location                      = azurerm_managed_lustre_file_system_auto_import_job.test.location

  auto_import_prefixes = azurerm_managed_lustre_file_system_auto_import_job.test.auto_import_prefixes
}
`, r.basic(data))
}
