package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMStorageShareDirectory_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShareDirectory_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareDirectoryExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageShareDirectory_uppercase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShareDirectory_uppercase(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareDirectoryExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageShareDirectory_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShareDirectory_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareDirectoryExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMStorageShareDirectory_requiresImport),
		},
	})
}

func TestAccAzureRMStorageShareDirectory_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShareDirectory_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareDirectoryExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageShareDirectory_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShareDirectory_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareDirectoryExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageShareDirectory_updated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareDirectoryExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}
func TestAccAzureRMStorageShareDirectory_nested(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "parent")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShareDirectory_nested(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareDirectoryExists(data.ResourceName),
					testCheckAzureRMStorageShareDirectoryExists("azurerm_storage_share_directory.child_one"),
					testCheckAzureRMStorageShareDirectoryExists("azurerm_storage_share_directory.child_two"),
					testCheckAzureRMStorageShareDirectoryExists("azurerm_storage_share_directory.multiple_child_one"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMStorageShareDirectoryExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		shareName := rs.Primary.Attributes["share_name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for Directory %q (Share %q): %s", accountName, name, shareName, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
		}

		client, err := storageClient.FileShareDirectoriesClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building FileShare Client: %s", err)
		}

		resp, err := client.Get(ctx, accountName, shareName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on FileShareDirectoriesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Directory %q (File Share %q / Account %q / Resource Group %q) does not exist", name, shareName, accountName, account.ResourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMStorageShareDirectoryDestroy(s *terraform.State) error {
	storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_share_directory" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		shareName := rs.Primary.Attributes["share_name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for Directory %q (Share %q): %s", accountName, name, shareName, err)
		}

		// not found, the account's gone
		if account == nil {
			return nil
		}

		if err != nil {
			return fmt.Errorf("Error locating Resource Group for Storage Share Directory %q (Share %s, Account %s): %s", name, shareName, accountName, err)
		}

		client, err := storageClient.FileShareDirectoriesClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building FileShare Client: %s", err)
		}

		resp, err := client.Get(ctx, accountName, shareName, name)
		if err != nil {
			return nil
		}

		return fmt.Errorf("File Share still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMStorageShareDirectory_basic(data acceptance.TestData) string {
	template := testAccAzureRMStorageShareDirectory_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "test" {
  name                 = "dir"
  share_name           = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_account.test.name
}
`, template)
}

func testAccAzureRMStorageShareDirectory_uppercase(data acceptance.TestData) string {
	template := testAccAzureRMStorageShareDirectory_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "test" {
  name                 = "UpperCaseCharacterS"
  share_name           = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_account.test.name
}
`, template)
}

func testAccAzureRMStorageShareDirectory_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStorageShareDirectory_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "import" {
  name                 = azurerm_storage_share_directory.test.name
  share_name           = azurerm_storage_share_directory.test.share_name
  storage_account_name = azurerm_storage_share_directory.test.storage_account_name
}
`, template)
}

func testAccAzureRMStorageShareDirectory_complete(data acceptance.TestData) string {
	template := testAccAzureRMStorageShareDirectory_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "test" {
  name                 = "dir"
  share_name           = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_account.test.name

  metadata = {
    hello = "world"
  }
}
`, template)
}

func testAccAzureRMStorageShareDirectory_updated(data acceptance.TestData) string {
	template := testAccAzureRMStorageShareDirectory_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "test" {
  name                 = "dir"
  share_name           = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_account.test.name

  metadata = {
    hello    = "world"
    sunshine = "at dawn"
  }
}
`, template)
}

func testAccAzureRMStorageShareDirectory_nested(data acceptance.TestData) string {
	template := testAccAzureRMStorageShareDirectory_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "parent" {
  name                 = "123--parent-dir"
  share_name           = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_storage_share_directory" "child_one" {
  name                 = "${azurerm_storage_share_directory.parent.name}/child1"
  share_name           = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_storage_share_directory" "child_two" {
  name                 = "${azurerm_storage_share_directory.child_one.name}/childtwo--123"
  share_name           = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_storage_share_directory" "multiple_child_one" {
  name                 = "${azurerm_storage_share_directory.parent.name}/c"
  share_name           = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_account.test.name
}
`, template)
}

func testAccAzureRMStorageShareDirectory_template(data acceptance.TestData) string {
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
