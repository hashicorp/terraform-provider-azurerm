package tests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMStorageShareFile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_file", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShareFile_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareFileExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageShareFile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_file", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShareFile_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareFileExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMStorageShareFile_requiresImport),
		},
	})
}

func TestAccAzureRMStorageShareFile_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_file", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShareFile_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareFileExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageShareFile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_file", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShareFile_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareFileExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageShareFile_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareFileExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageShareFile_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareFileExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageShareFile_withFile(t *testing.T) {
	sourceBlob, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Failed to create local source blob file")
	}

	if err := testAccAzureRMStorageBlob_populateTempFile(sourceBlob); err != nil {
		t.Fatalf("Error populating temp file: %s", err)
	}
	data := acceptance.BuildTestData(t, "azurerm_storage_share_file", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShareFile_withFile(data, sourceBlob.Name()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareFileExists(data.ResourceName),
				),
			},
			data.ImportStep("source"),
		},
	})
}

func testCheckAzureRMStorageShareFileExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		storageShareID, err := parse.StorageShareDataPlaneID(rs.Primary.Attributes["storage_share_id"])
		if err != nil {
			return err
		}
		path := rs.Primary.Attributes["path"]

		account, err := storageClient.FindAccount(ctx, storageShareID.AccountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for File %q (Share %q): %s", storageShareID.AccountName, name, storageShareID.Name, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", storageShareID.AccountName)
		}

		client, err := storageClient.FileShareFilesClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building FileShare File Client: %s", err)
		}

		resp, err := client.GetProperties(ctx, storageShareID.AccountName, storageShareID.Name, path, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on FileShareFilesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: File %q (File Share %q / Account %q / Resource Group %q) does not exist", name, storageShareID.Name, storageShareID.AccountName, account.ResourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMStorageShareFileDestroy(s *terraform.State) error {
	storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_share_file" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		storageShareID, err := parse.StorageShareDataPlaneID(rs.Primary.Attributes["storage_share_id"])
		if err != nil {
			return err
		}
		path := rs.Primary.Attributes["path"]

		account, err := storageClient.FindAccount(ctx, storageShareID.AccountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for File %q (Share %q): %s", storageShareID.AccountName, name, storageShareID.Name, err)
		}

		// not found, the account's gone
		if account == nil {
			return nil
		}

		if err != nil {
			return fmt.Errorf("Error locating Resource Group for Storage Share File %q (Share %s, Account %s): %s", name, storageShareID.Name, storageShareID.AccountName, err)
		}

		client, err := storageClient.FileShareFilesClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building FileShare File Client: %s", err)
		}

		resp, err := client.GetProperties(ctx, storageShareID.AccountName, storageShareID.Name, path, name)
		if err != nil {
			return nil
		}

		return fmt.Errorf("File Share still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMStorageShareFile_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "test" {
  name                 = "fileshare"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 50
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageShareFile_basic(data acceptance.TestData) string {
	template := testAccAzureRMStorageShareFile_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_file" "test" {
  name                 = "dir"
  storage_share_id     = azurerm_storage_share.test.id

  metadata = {
    hello = "world"
  }
}
`, template)
}

func testAccAzureRMStorageShareFile_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStorageShareFile_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_file" "import" {
  name                 = azurerm_storage_share_file.test.name
  storage_share_id     = azurerm_storage_share.test.id
}
`, template)
}

func testAccAzureRMStorageShareFile_complete(data acceptance.TestData) string {
	template := testAccAzureRMStorageShareFile_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_file" "test" {
  name                 = "dir"
  storage_share_id     = azurerm_storage_share.test.id

  
  content_type        = "test_content_type"
  content_encoding    = "test_encoding"
  content_disposition = "test_content_disposition"

  metadata = {
    hello = "world"
  }
}
`, template)
}

func testAccAzureRMStorageShareFile_withFile(data acceptance.TestData, fileName string) string {
	template := testAccAzureRMStorageShareFile_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_file" "test" {
  name                 = "dir"
  storage_share_id     = azurerm_storage_share.test.id

  source = "%s"

  metadata = {
    hello = "world"
  }
}
`, template, fileName)
}
