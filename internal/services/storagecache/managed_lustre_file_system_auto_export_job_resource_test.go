package storagecache_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/autoexportjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagedLustreFileSystemAutoExportJobResource struct{}

func (r ManagedLustreFileSystemAutoExportJobResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := autoexportjobs.ParseAutoExportJobID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.StorageCache.AutoExportJobs
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestAccManagedLustreFileSystemExportJob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_lustre_file_system_auto_export_job", "test")
	r := ManagedLustreFileSystemAutoExportJobResource{}

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

func TestAccManagedLustreFileSystemExportJob_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_lustre_file_system_auto_export_job", "test")
	r := ManagedLustreFileSystemAutoExportJobResource{}

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

func TestAccManagedLustreFileSystemExportJob_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_lustre_file_system_auto_export_job", "test")
	r := ManagedLustreFileSystemAutoExportJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedLustreFileSystemExportJob_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_lustre_file_system_auto_export_job", "test")
	r := ManagedLustreFileSystemAutoExportJobResource{}

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

func (r ManagedLustreFileSystemAutoExportJobResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_lustre_file_system_auto_export_job" "test" {
  name                 = "acctest-amlfsaej-%d"
  resource_group_name  = azurerm_resource_group.test.name
  aml_file_system_name = azurerm_managed_lustre_file_system.test.name
  location             = azurerm_resource_group.test.location

  auto_export_prefixes = ["/"]
}
`, ManagedLustreFileSystemResource{}.complete(data), data.RandomInteger)
}

func (r ManagedLustreFileSystemAutoExportJobResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_lustre_file_system_auto_export_job" "test" {
  name                 = "acctest-amlfsaej-%d"
  resource_group_name  = azurerm_resource_group.test.name
  aml_file_system_name = azurerm_managed_lustre_file_system.test.name
  location             = azurerm_resource_group.test.location

  auto_export_prefixes = ["/"]
  admin_status         = "Enable"
}
`, ManagedLustreFileSystemResource{}.complete(data), data.RandomInteger)
}

func (r ManagedLustreFileSystemAutoExportJobResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_lustre_file_system_auto_export_job" "import" {
	name                 = azurerm_managed_lustre_file_system_auto_export_job.test.name
	resource_group_name  = azurerm_managed_lustre_file_system_auto_export_job.test.resource_group_name
	aml_file_system_name = azurerm_managed_lustre_file_system_auto_export_job.test.aml_file_system_name
	location             = azurerm_managed_lustre_file_system_auto_export_job.test.location
	
	auto_export_prefixes = azurerm_managed_lustre_file_system_auto_export_job.test.auto_export_prefixes
	admin_status         = azurerm_managed_lustre_file_system_auto_export_job.test.admin_status
`, r.basic(data))
}
