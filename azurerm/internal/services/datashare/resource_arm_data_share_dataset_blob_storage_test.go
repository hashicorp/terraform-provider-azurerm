package datashare_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDataShareDataSetBlobStorageFile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_blob_storage", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDataSetDestroy("azurerm_data_share_dataset_blob_storage"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShareDataSetBlobStorageFile_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareDataSetExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataShareDataSetBlobStorageFolder_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_blob_storage", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDataSetDestroy("azurerm_data_share_dataset_blob_storage"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShareDataSetBlobStorageFolder_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareDataSetExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataShareDataSetBlobStorageContainer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_blob_storage", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDataSetDestroy("azurerm_data_share_dataset_blob_storage"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShareDataSetBlobStorageContainer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareDataSetExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataShareDataSetBlobStorage_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_blob_storage", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDataSetDestroy("azurerm_data_share_dataset_blob_storage"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShareDataSetBlobStorageFile_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareDataSetExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMDataShareDataSetBlobStorage_requiresImport),
		},
	})
}

func testCheckAzureRMDataShareDataSetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataShare.DataSetClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("DataShare DataSet not found: %s", resourceName)
		}
		id, err := parse.DataSetID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.AccountName, id.ShareName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: data share data set %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on DataShare DataSet Client: %+v", err)
		}
		return nil
	}
}

// nolint
func testCheckAzureRMDataShareDataSetDestroy(resourceTypeName string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataShare.DataSetClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceTypeName {
				continue
			}
			id, err := parse.DataSetID(rs.Primary.ID)
			if err != nil {
				return err
			}
			if resp, err := client.Get(ctx, id.ResourceGroup, id.AccountName, id.ShareName, id.Name); err != nil {
				if !utils.ResponseWasNotFound(resp.Response) {
					return fmt.Errorf("bad: get on data share data set client: %+v", err)
				}
			}
			return nil
		}
		return nil
	}
}

func testAccAzureRMDataShareDataSetBlobStorage_template(data acceptance.TestData) string {
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
}

resource "azurerm_data_share" "test" {
  name       = "acctest_DS_%[1]d"
  account_id = azurerm_data_share_account.test.id
  kind       = "CopyBased"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "RAGRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctest-sc-%[1]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "container"
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

func testAccAzureRMDataShareDataSetBlobStorageFile_basic(data acceptance.TestData) string {
	config := testAccAzureRMDataShareDataSetBlobStorage_template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_share_dataset_blob_storage" "test" {
  name           = "acctest-DSDSBS-file-%[2]d"
  data_share_id  = azurerm_data_share.test.id
  container_name = azurerm_storage_container.test.name
  storage_account {
    name                = azurerm_storage_account.test.name
    resource_group_name = azurerm_storage_account.test.resource_group_name
    subscription_id     = "%[3]s"
  }
  file_path = "myfile.txt"
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, config, data.RandomInteger, os.Getenv("ARM_SUBSCRIPTION_ID"))
}

func testAccAzureRMDataShareDataSetBlobStorageFolder_basic(data acceptance.TestData) string {
	config := testAccAzureRMDataShareDataSetBlobStorage_template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_share_dataset_blob_storage" "test" {
  name           = "acctest-DSDSBS-folder-%[2]d"
  data_share_id  = azurerm_data_share.test.id
  container_name = azurerm_storage_container.test.name
  storage_account {
    name                = azurerm_storage_account.test.name
    resource_group_name = azurerm_storage_account.test.resource_group_name
    subscription_id     = "%[3]s"
  }
  folder_path = "test"
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, config, data.RandomInteger, os.Getenv("ARM_SUBSCRIPTION_ID"))
}

func testAccAzureRMDataShareDataSetBlobStorageContainer_basic(data acceptance.TestData) string {
	config := testAccAzureRMDataShareDataSetBlobStorage_template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_share_dataset_blob_storage" "test" {
  name           = "acctest-DSDSBS-folder-%[2]d"
  data_share_id  = azurerm_data_share.test.id
  container_name = azurerm_storage_container.test.name
  storage_account {
    name                = azurerm_storage_account.test.name
    resource_group_name = azurerm_storage_account.test.resource_group_name
    subscription_id     = "%[3]s"
  }
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, config, data.RandomInteger, os.Getenv("ARM_SUBSCRIPTION_ID"))
}

func testAccAzureRMDataShareDataSetBlobStorage_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMDataShareDataSetBlobStorageFile_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_blob_storage" "import" {
  name           = azurerm_data_share_dataset_blob_storage.test.name
  data_share_id  = azurerm_data_share.test.id
  container_name = azurerm_data_share_dataset_blob_storage.test.container_name
  storage_account {
    name                = azurerm_data_share_dataset_blob_storage.test.storage_account.0.name
    resource_group_name = azurerm_data_share_dataset_blob_storage.test.storage_account.0.resource_group_name
    subscription_id     = azurerm_data_share_dataset_blob_storage.test.storage_account.0.subscription_id
  }
}
`, config)
}
