package datashare_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataShareDataSetDataLakeGen1File_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_data_lake_gen1", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataShareDataSetDestroy("azurerm_data_share_dataset_data_lake_gen1"),
		Steps: []resource.TestStep{
			{
				Config: testAccDataShareDataSetDataLakeGen1File_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataShareDataSetExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccDataShareDataSetDataLakeGen1Folder_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_data_lake_gen1", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataShareDataSetDestroy("azurerm_data_share_dataset_data_lake_gen1"),
		Steps: []resource.TestStep{
			{
				Config: testAccDataShareDataSetDataLakeGen1Folder_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataShareDataSetExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccDataShareDataSetDataLakeGen1_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_data_lake_gen1", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataShareDataSetDestroy("azurerm_data_share_dataset_data_lake_gen1"),
		Steps: []resource.TestStep{
			{
				Config: testAccDataShareDataSetDataLakeGen1File_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataShareDataSetExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccDataShareDataSetDataLakeGen1_requiresImport),
		},
	})
}

func testAccDataShareDataSetDataLakeGen1_template(data acceptance.TestData) string {
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
  name                = "acctest-DSA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  identity {
    type = "SystemAssigned"
  }

  tags = {
    env = "Test"
  }
}

resource "azurerm_data_share" "test" {
  name       = "acctest_DS_%[1]d"
  account_id = azurerm_data_share_account.test.id
  kind       = "CopyBased"
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctestdls%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  firewall_state      = "Disabled"
}

resource "azurerm_data_lake_store_file" "test" {
  account_name     = azurerm_data_lake_store.test.name
  local_file_path  = "./testdata/application_gateway_test.cer"
  remote_file_path = "/test/application_gateway_test.cer"
}

data "azuread_service_principal" "test" {
  display_name = azurerm_data_share_account.test.name
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_data_lake_store.test.id
  role_definition_name = "Owner"
  principal_id         = data.azuread_service_principal.test.object_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(12))
}

func testAccDataShareDataSetDataLakeGen1File_basic(data acceptance.TestData) string {
	config := testAccDataShareDataSetDataLakeGen1_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_data_lake_gen1" "test" {
  name               = "acctest-DSDL1-%d"
  data_share_id      = azurerm_data_share.test.id
  data_lake_store_id = azurerm_data_lake_store.test.id
  file_name          = "application_gateway_test.cer"
  folder_path        = "test"
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, config, data.RandomInteger)
}

func testAccDataShareDataSetDataLakeGen1Folder_basic(data acceptance.TestData) string {
	config := testAccDataShareDataSetDataLakeGen1_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_data_lake_gen1" "test" {
  name               = "acctest-DSDL1-%d"
  data_share_id      = azurerm_data_share.test.id
  data_lake_store_id = azurerm_data_lake_store.test.id
  folder_path        = "test"
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, config, data.RandomInteger)
}

func testAccDataShareDataSetDataLakeGen1_requiresImport(data acceptance.TestData) string {
	config := testAccDataShareDataSetDataLakeGen1File_basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_data_share_dataset_data_lake_gen1" "import" {
  name               = azurerm_data_share_dataset_data_lake_gen1.test.name
  data_share_id      = azurerm_data_share.test.id
  data_lake_store_id = azurerm_data_share_dataset_data_lake_gen1.test.data_lake_store_id
  folder_path        = azurerm_data_share_dataset_data_lake_gen1.test.folder_path
}
`, config)
}
