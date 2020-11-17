package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMStorageTable_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageTable_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageTableExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageTable_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageTable_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageTableExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMStorageTable_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_storage_table"),
			},
		},
	})
}

func TestAccAzureRMStorageTable_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageTable_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageTableExists("azurerm_storage_table.test"),
					testAccARMStorageTableDisappears("azurerm_storage_table.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMStorageTable_acl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageTable_acl(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageTableExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageTable_aclUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageTableExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMStorageTableExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		tableName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for Table %q: %s", accountName, tableName, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
		}

		client, err := storageClient.TablesClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building Table Client: %s", err)
		}

		exists, err := client.Exists(ctx, account.ResourceGroup, accountName, tableName)
		if err != nil {
			return fmt.Errorf("retrieving Table %q: %s", tableName, accountName)
		}
		if exists == nil || !*exists {
			return fmt.Errorf("Table %q doesn't exist in Account %q!", tableName, accountName)
		}

		return nil
	}
}

func testAccARMStorageTableDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		tableName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for Table %q: %s", accountName, tableName, err)
		}
		if account == nil {
			return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
		}

		client, err := storageClient.TablesClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building Table Client: %s", err)
		}

		exists, err := client.Exists(ctx, account.ResourceGroup, accountName, tableName)
		if err != nil {
			return fmt.Errorf("Error retrieving Table %q: %s", tableName, accountName)
		}
		if exists == nil || !*exists {
			return fmt.Errorf("Table %q doesn't exist in Account %q so it can't be deleted!", tableName, accountName)
		}

		if err := client.Delete(ctx, account.ResourceGroup, accountName, tableName); err != nil {
			return fmt.Errorf("deleting Table %q: %s", tableName, err)
		}

		return nil
	}
}

func testCheckAzureRMStorageTableDestroy(s *terraform.State) error {
	storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_table" {
			continue
		}

		tableName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %q for Table %q: %s", accountName, tableName, err)
		}
		if account == nil {
			return nil
		}

		client, err := storageClient.TablesClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building Table Client: %s", err)
		}

		exists, err := client.Exists(ctx, account.ResourceGroup, accountName, tableName)
		if err != nil {
			return nil
		}
		if exists != nil && *exists {
			return fmt.Errorf("Table still exists")
		}

		return nil
	}

	return nil
}

func testAccAzureRMStorageTable_basic(data acceptance.TestData) string {
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

resource "azurerm_storage_table" "test" {
  name                 = "acctestst%d"
  storage_account_name = azurerm_storage_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func testAccAzureRMStorageTable_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStorageTable_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_table" "import" {
  name                 = azurerm_storage_table.test.name
  storage_account_name = azurerm_storage_table.test.storage_account_name
}
`, template)
}

func testAccAzureRMStorageTable_acl(data acceptance.TestData) string {
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

resource "azurerm_storage_table" "test" {
  name                 = "acctestst%d"
  storage_account_name = azurerm_storage_account.test.name
  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "raud"
      start       = "2020-11-26T08:49:37.0000000Z"
      expiry      = "2020-11-27T08:49:37.0000000Z"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func testAccAzureRMStorageTable_aclUpdated(data acceptance.TestData) string {
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

resource "azurerm_storage_table" "test" {
  name                 = "acctestst%d"
  storage_account_name = azurerm_storage_account.test.name

  acl {
    id = "AAAANDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "raud"
      start       = "2020-11-26T08:49:37.0000000Z"
      expiry      = "2020-11-27T08:49:37.0000000Z"
    }
  }
  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "raud"
      start       = "2019-07-02T09:38:21.0000000Z"
      expiry      = "2019-07-02T10:38:21.0000000Z"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}
