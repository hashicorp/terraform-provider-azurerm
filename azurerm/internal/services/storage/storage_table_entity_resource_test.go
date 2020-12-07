package storage_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/table/entities"
)

func TestAccAzureRMTableEntity_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table_entity", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTableEntityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTableEntity_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTableEntityExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMTableEntity_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table_entity", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTableEntityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTableEntity_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTableEntityExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMTableEntity_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_storage_table_entity"),
			},
		},
	})
}

func TestAccAzureRMTableEntity_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table_entity", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMTableEntityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTableEntity_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTableEntityExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMTableEntity_updated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTableEntityExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMTableEntityExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		tableName := rs.Primary.Attributes["table_name"]
		accountName := rs.Primary.Attributes["storage_account_name"]
		partitionKey := rs.Primary.Attributes["partition_key"]
		rowKey := rs.Primary.Attributes["row_key"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error locating Resource Group for Storage Table Entity (Partition Key %q / Row Key %q) (Table %q / Account %q): %v", partitionKey, rowKey, tableName, accountName, err)
		}
		if account == nil {
			return fmt.Errorf("Storage Account %q was not found!", accountName)
		}

		client, err := storageClient.TableEntityClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building Table Entity Client: %s", err)
		}

		input := entities.GetEntityInput{
			PartitionKey:  partitionKey,
			RowKey:        rowKey,
			MetaDataLevel: entities.NoMetaData,
		}
		resp, err := client.Get(ctx, accountName, tableName, input)
		if err != nil {
			return fmt.Errorf("Bad: Get on Table EntityClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Entity (Partition Key %q / Row Key %q) (Table %q / Account %q / Resource Group %q) does not exist", partitionKey, rowKey, tableName, accountName, account.ResourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMTableEntityDestroy(s *terraform.State) error {
	storageClient := acceptance.AzureProvider.Meta().(*clients.Client).Storage
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_table_entity" {
			continue
		}

		tableName := rs.Primary.Attributes["table_name"]
		accountName := rs.Primary.Attributes["storage_account_name"]
		partitionKey := rs.Primary.Attributes["parititon_key"]
		rowKey := rs.Primary.Attributes["row_key"]

		account, err := storageClient.FindAccount(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error locating Resource Group for Storage Table Entity (Partition Key %q / Row Key %q) (Table %q / Account %q): %v", partitionKey, rowKey, tableName, accountName, err)
		}

		// not found, the account's gone
		if account == nil {
			return nil
		}

		client, err := storageClient.TableEntityClient(ctx, *account)
		if err != nil {
			return fmt.Errorf("Error building TableEntity Client: %s", err)
		}

		input := entities.GetEntityInput{
			PartitionKey: partitionKey,
			RowKey:       rowKey,
		}
		resp, err := client.Get(ctx, accountName, tableName, input)
		if err != nil {
			return fmt.Errorf("Bad: Get on Table Entity: %+v", err)
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Table Entity still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMTableEntity_basic(data acceptance.TestData) string {
	template := testAccAzureRMTableEntity_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_table_entity" "test" {
  storage_account_name = azurerm_storage_account.test.name
  table_name           = azurerm_storage_table.test.name

  partition_key = "test_partition%d"
  row_key       = "test_row%d"
  entity = {
    Foo = "Bar"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTableEntity_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMTableEntity_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_table_entity" "import" {
  storage_account_name = azurerm_storage_account.test.name
  table_name           = azurerm_storage_table.test.name

  partition_key = "test_partition%d"
  row_key       = "test_row%d"
  entity = {
    Foo = "Bar"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTableEntity_updated(data acceptance.TestData) string {
	template := testAccAzureRMTableEntity_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_table_entity" "test" {
  storage_account_name = azurerm_storage_account.test.name
  table_name           = azurerm_storage_table.test.name

  partition_key = "test_partition%d"
  row_key       = "test_row%d"
  entity = {
    Foo  = "Bar"
    Test = "Updated"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMTableEntity_template(data acceptance.TestData) string {
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

resource "azurerm_storage_table" "test" {
  name                 = "acctestst%d"
  storage_account_name = azurerm_storage_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}
