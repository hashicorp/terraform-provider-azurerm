package azurerm

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/table/entities"
)

func TestAccAzureRMTableEntity_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(5))
	location := testLocation()
	resourceName := "azurerm_storage_table_entity.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTableEntityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTableEntity_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTableEntityExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMTableEntity_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(5))
	location := testLocation()
	resourceName := "azurerm_storage_table_entity.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTableEntityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTableEntity_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTableEntityExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMTableEntity_requiresImport(ri, rs, location),
				ExpectError: testRequiresImportError("azurerm_storage_table_entity"),
			},
		},
	})
}

func TestAccAzureRMTableEntity_update(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(5))
	location := testLocation()
	resourceName := "azurerm_storage_table_entity.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMTableEntityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMTableEntity_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTableEntityExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMTableEntity_updated(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMTableEntityExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMTableEntityExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		tableName := rs.Primary.Attributes["table_name"]
		accountName := rs.Primary.Attributes["storage_account_name"]
		partitionKey := rs.Primary.Attributes["partition_key"]
		rowKey := rs.Primary.Attributes["row_key"]

		storageClient := testAccProvider.Meta().(*ArmClient).storage
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error locating Resource Group for Storage Table Entity (Partition Key %q / Row Key %q) (Table %q / Account %q): %v", partitionKey, rowKey, tableName, accountName, err)
		}
		if resourceGroup == nil {
			return fmt.Errorf("Unable to locate Resource Group for Storage Table Entity (Partition Key %q / Row Key %q) (Table %q / Account %q)", partitionKey, rowKey, tableName, accountName)
		}

		client, err := storageClient.TableEntityClient(ctx, *resourceGroup, accountName)
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
			return fmt.Errorf("Bad: Entity (Partition Key %q / Row Key %q) (Table %q / Account %q / Resource Group %q) does not exist", partitionKey, rowKey, tableName, accountName, *resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMTableEntityDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_table_entity" {
			continue
		}

		tableName := rs.Primary.Attributes["table_name"]
		accountName := rs.Primary.Attributes["storage_account_name"]
		partitionKey := rs.Primary.Attributes["parititon_key"]
		rowKey := rs.Primary.Attributes["row_key"]

		storageClient := testAccProvider.Meta().(*ArmClient).storage
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error locating Resource Group for Storage Table Entity (Partition Key %q / Row Key %q) (Table %q / Account %q / Resource Group %q): %v", partitionKey, rowKey, tableName, accountName, *resourceGroup, err)
		}

		// not found, the account's gone
		if resourceGroup == nil {
			return nil
		}

		client, err := storageClient.TableEntityClient(ctx, *resourceGroup, accountName)
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

func testAccAzureRMTableEntity_basic(rInt int, rString string, location string) string {
	template := testAccAzureRMTableEntity_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_table_entity" "test" {
  storage_account_name = "${azurerm_storage_account.test.name}"
  table_name           = "${azurerm_storage_table.test.name}"
  
  partition_key = "test_partition%d"
  row_key       = "test_row%d"
  entity = {
    Foo = "Bar"
  }
}
`, template, rInt, rInt)
}

func testAccAzureRMTableEntity_requiresImport(rInt int, rString string, location string) string {
	template := testAccAzureRMTableEntity_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_table_entity" "test" {
  storage_account_name = "${azurerm_storage_account.test.name}"
  table_name           = "${azurerm_storage_table.test.name}"
	
  partition_key = "test_partition%d"
  row_key       = "test_row%d"
  entity = {
    Foo = "Bar"
  }
}
`, template, rInt, rInt)
}

func testAccAzureRMTableEntity_updated(rInt int, rString string, location string) string {
	template := testAccAzureRMTableEntity_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_table_entity" "test" {
  storage_account_name = "${azurerm_storage_account.test.name}"
  table_name           = "${azurerm_storage_table.test.name}"
	
  partition_key = "test_partition%d"
  row_key       = "test_row%d"
  entity = {
	Foo = "Bar"
	Test = "Updated"
  }
}
`, template, rInt, rInt)
}

func testAccAzureRMTableEntity_template(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_table" "test" {
  name                 = "acctestst%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}
`, rInt, location, rString, rInt)
}
