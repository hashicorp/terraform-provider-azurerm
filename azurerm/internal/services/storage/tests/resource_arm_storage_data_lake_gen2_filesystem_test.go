package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parsers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMStorageDataLakeGen2FileSystem_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_filesystem", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageDataLakeGen2FileSystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageDataLakeGen2FileSystem_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageDataLakeGen2FileSystemExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageDataLakeGen2FileSystem_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_filesystem", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageDataLakeGen2FileSystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageDataLakeGen2FileSystem_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageDataLakeGen2FileSystemExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMStorageDataLakeGen2FileSystem_requiresImport),
		},
	})
}

func TestAccAzureRMStorageDataLakeGen2FileSystem_properties(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_filesystem", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageDataLakeGen2FileSystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageDataLakeGen2FileSystem_properties(data, "aGVsbG8="),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageDataLakeGen2FileSystemExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageDataLakeGen2FileSystem_properties(data, "ZXll"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageDataLakeGen2FileSystemExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMStorageDataLakeGen2FileSystemExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Storage.FileSystemsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		fileSystemName := rs.Primary.Attributes["name"]
		storageID, err := parsers.ParseAccountID(rs.Primary.Attributes["storage_account_id"])
		if err != nil {
			return err
		}

		resp, err := client.GetProperties(ctx, storageID.Name, fileSystemName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: File System %q (Account %q) does not exist", fileSystemName, storageID.Name)
			}

			return fmt.Errorf("Bad: Get on FileSystemsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMStorageDataLakeGen2FileSystemDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Storage.FileSystemsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_data_lake_gen2_filesystem" {
			continue
		}

		fileSystemName := rs.Primary.Attributes["name"]
		storageID, err := parsers.ParseAccountID(rs.Primary.Attributes["storage_account_id"])
		if err != nil {
			return err
		}

		props, err := client.GetProperties(ctx, storageID.Name, fileSystemName)
		if err != nil {
			return nil
		}

		return fmt.Errorf("File System still exists: %+v", props)
	}

	return nil
}

func testAccAzureRMStorageDataLakeGen2FileSystem_basic(data acceptance.TestData) string {
	template := testAccAzureRMStorageDataLakeGen2FileSystem_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%d"
  storage_account_id = azurerm_storage_account.test.id
}
`, template, data.RandomInteger)
}

func testAccAzureRMStorageDataLakeGen2FileSystem_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStorageDataLakeGen2FileSystem_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_data_lake_gen2_filesystem" "import" {
  name               = azurerm_storage_data_lake_gen2_filesystem.test.name
  storage_account_id = azurerm_storage_data_lake_gen2_filesystem.test.storage_account_id
}
`, template)
}

func testAccAzureRMStorageDataLakeGen2FileSystem_properties(data acceptance.TestData, value string) string {
	template := testAccAzureRMStorageDataLakeGen2FileSystem_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%d"
  storage_account_id = azurerm_storage_account.test.id

  properties = {
    key = "%s"
  }
}
`, template, data.RandomInteger, value)
}

func testAccAzureRMStorageDataLakeGen2FileSystem_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
