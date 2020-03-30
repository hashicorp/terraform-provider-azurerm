package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMStorageImportExportJob_importJobBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_import_export_job", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageImportJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStorageImportExportJob_importJobBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageImportJobExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMStorageImportExportJob_exportJobBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_import_export_job", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageExportJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStorageImportExportJob_exportJobBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageExportJobExists(data.ResourceName),
				),
			},
		},
	})
}

func testAccDataSourceStorageImportExportJob_importJobBasic(data acceptance.TestData) string {
	config := testAccAzureRMStorageImportJob_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_import_export_job" "test" {
  name                = azurerm_import_job.test.name
  resource_group_name = azurerm_import_job.test.resource_group_name
}
`, config)
}

func testAccDataSourceStorageImportExportJob_exportJobBasic(data acceptance.TestData) string {
	config := testAccAzureRMStorageExportJob_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_import_export_job" "test" {
  name                = azurerm_export_job.test.name
  resource_group_name = azurerm_export_job.test.resource_group_name
}
`, config)
}
