package datashare_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMDataShareDataSetDataLakeGen2File_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_data_lake_gen2", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDataSetDestroy("azurerm_data_share_dataset_data_lake_gen2"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShareDataSetDataLakeGen2File_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareDataSetExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataShareDataSetDataLakeGen2Folder_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_data_lake_gen2", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDataSetDestroy("azurerm_data_share_dataset_data_lake_gen2"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShareDataSetDataLakeGen2Folder_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareDataSetExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataShareDataSetDataLakeGen2FileSystem_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_data_lake_gen2", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDataSetDestroy("azurerm_data_share_dataset_data_lake_gen2"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShareDataSetDataLakeGen2FileSystem_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareDataSetExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataShareDataLakeGen2DataSet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_data_lake_gen2", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDataSetDestroy("azurerm_data_share_dataset_data_lake_gen2"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShareDataSetDataLakeGen2File_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareDataSetExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMDataShareDataLakeGen2DataSet_requiresImport),
		},
	})
}

func testAccAzureRMDataShareDataLakeGen2DataSet_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-datashare-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_share_account" "test" {
  name                = "acctest-dsa-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_data_share" "test" {
  name       = "acctest_ds_%[1]d"
  account_id = azurerm_data_share_account.test.id
  kind       = "CopyBased"
}

resource "azurerm_storage_account" "test" {
  name                     = "accteststr%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%[1]d"
  storage_account_id = azurerm_storage_account.test.id
}

data "azuread_service_principal" "test" {
  display_name = azurerm_data_share_account.test.name
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Reader"
  principal_id         = data.azuread_service_principal.test.object_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(12))
}

func testAccAzureRMDataShareDataSetDataLakeGen2File_basic(data acceptance.TestData) string {
	config := testAccAzureRMDataShareDataLakeGen2DataSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_data_lake_gen2" "test" {
  name               = "acctest-dlds-%d"
  share_id           = azurerm_data_share.test.id
  storage_account_id = azurerm_storage_account.test.id
  file_system_name   = azurerm_storage_data_lake_gen2_filesystem.test.name
  file_path          = "myfile.txt"
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, config, data.RandomInteger)
}

func testAccAzureRMDataShareDataSetDataLakeGen2Folder_basic(data acceptance.TestData) string {
	config := testAccAzureRMDataShareDataLakeGen2DataSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_data_lake_gen2" "test" {
  name               = "acctest-dlds-%d"
  share_id           = azurerm_data_share.test.id
  storage_account_id = azurerm_storage_account.test.id
  file_system_name   = azurerm_storage_data_lake_gen2_filesystem.test.name
  folder_path        = "test"
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, config, data.RandomInteger)
}

func testAccAzureRMDataShareDataSetDataLakeGen2FileSystem_basic(data acceptance.TestData) string {
	config := testAccAzureRMDataShareDataLakeGen2DataSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_data_lake_gen2" "test" {
  name               = "acctest-dlds-%d"
  share_id           = azurerm_data_share.test.id
  storage_account_id = azurerm_storage_account.test.id
  file_system_name   = azurerm_storage_data_lake_gen2_filesystem.test.name
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, config, data.RandomInteger)
}

func testAccAzureRMDataShareDataLakeGen2DataSet_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMDataShareDataSetDataLakeGen2File_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_data_lake_gen2" "import" {
  name               = azurerm_data_share_dataset_data_lake_gen2.test.name
  share_id           = azurerm_data_share.test.id
  storage_account_id = azurerm_data_share_dataset_data_lake_gen2.test.storage_account_id
  file_system_name   = azurerm_data_share_dataset_data_lake_gen2.test.file_system_name
  file_path          = azurerm_data_share_dataset_data_lake_gen2.test.file_path
}
`, config)
}
