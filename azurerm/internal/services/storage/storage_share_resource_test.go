package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMStorageShare_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShare_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageShare_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShare_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMStorageShare_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_storage_share"),
			},
		},
	})
}

func TestAccAzureRMStorageShare_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShare_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(data.ResourceName),
					testCheckAzureRMStorageShareDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMStorageShare_deleteAndRecreate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShare_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageShare_template(data),
			},
			{
				Config: testAccAzureRMStorageShare_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageShare_metaData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShare_metaData(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageShare_metaDataUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageShare_acl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShare_acl(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageShare_aclUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageShare_aclGhostedRecall(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShare_aclGhostedRecall(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageShare_updateQuota(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShare_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMStorageShare_updateQuota(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "quota", "5"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageShare_largeQuota(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShare_largeQuota(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageShare_largeQuotaUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMStorageShareExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		shareName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for Share %q: %s", accountName, shareName, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
		}

		client, err := storageClient.FileSharesClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building FileShare Client: %s", err)
		}

		exists, err := client.Exists(ctx, account.ResourceGroup, accountName, shareName)
		if err != nil {
			return fmt.Errorf("Bad: checking for presence of Share %q (Storage Account: %q): %+v", shareName, accountName, err)
		}
		if exists == nil || !*exists {
			return fmt.Errorf("Bad: Share %q (Storage Account: %q) does not exist", shareName, accountName)
		}

		return nil
	}
}

func testCheckAzureRMStorageShareDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		shareName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for Share %q: %s", accountName, shareName, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
		}

		client, err := storageClient.FileSharesClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building FileShare Client: %s", err)
		}

		if err := client.Delete(ctx, account.ResourceGroup, accountName, shareName); err != nil {
			return fmt.Errorf("Error deleting Share %q (Account %q): %v", shareName, accountName, err)
		}

		return nil
	}
}

func testCheckAzureRMStorageShareDestroy(s *terraform.State) error {
	storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_share" {
			continue
		}

		shareName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for Share %q: %s", accountName, shareName, err)
		}

		// expected since it's been deleted
		if account == nil {
			return nil
		}

		client, err := storageClient.FileSharesClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building FileShare Client: %s", err)
		}

		exists, err := client.Exists(ctx, account.ResourceGroup, accountName, shareName)
		if err != nil {
			return nil
		}

		if exists != nil && *exists {
			return fmt.Errorf("Share still exists!")
		}

		return nil
	}

	return nil
}

func testAccAzureRMStorageShare_basic(data acceptance.TestData) string {
	template := testAccAzureRMStorageShare_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
}
`, template, data.RandomString)
}

func testAccAzureRMStorageShare_metaData(data acceptance.TestData) string {
	template := testAccAzureRMStorageShare_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name

  metadata = {
    hello = "world"
  }
}
`, template, data.RandomString)
}

func testAccAzureRMStorageShare_metaDataUpdated(data acceptance.TestData) string {
	template := testAccAzureRMStorageShare_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name

  metadata = {
    hello = "world"
    happy = "birthday"
  }
}
`, template, data.RandomString)
}

func testAccAzureRMStorageShare_acl(data acceptance.TestData) string {
	template := testAccAzureRMStorageShare_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name

  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "rwd"
      start       = "2019-07-02T09:38:21.0000000Z"
      expiry      = "2019-07-02T10:38:21.0000000Z"
    }
  }
}
`, template, data.RandomString)
}

func testAccAzureRMStorageShare_aclGhostedRecall(data acceptance.TestData) string {
	template := testAccAzureRMStorageShare_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name

  acl {
    id = "GhostedRecall"
    access_policy {
      permissions = "r"
    }
  }
}
`, template, data.RandomString)
}

func testAccAzureRMStorageShare_aclUpdated(data acceptance.TestData) string {
	template := testAccAzureRMStorageShare_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name

  acl {
    id = "AAAANDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "rwd"
      start       = "2019-07-02T09:38:21.0000000Z"
      expiry      = "2019-07-02T10:38:21.0000000Z"
    }
  }
  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "rwd"
      start       = "2019-07-02T09:38:21.0000000Z"
      expiry      = "2019-07-02T10:38:21.0000000Z"
    }
  }
}
`, template, data.RandomString)
}

func testAccAzureRMStorageShare_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStorageShare_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "import" {
  name                 = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_share.test.storage_account_name
}
`, template)
}

func testAccAzureRMStorageShare_updateQuota(data acceptance.TestData) string {
	template := testAccAzureRMStorageShare_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 5
}
`, template, data.RandomString)
}

func testAccAzureRMStorageShare_largeQuota(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storageshare-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestshare%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Premium"
  account_replication_type = "LRS"
  account_kind             = "FileStorage"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 6000
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMStorageShare_largeQuotaUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storageshare-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestshare%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Premium"
  account_replication_type = "LRS"
  account_kind             = "FileStorage"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 10000
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMStorageShare_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
