package azurerm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMStorageTable_basic(t *testing.T) {
	resourceName := "azurerm_storage_table.test"

	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMStorageTable_basic(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageTableExists(resourceName),
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

func TestAccAzureRMStorageTable_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_storage_table.test"

	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageTable_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageTableExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMStorageTable_requiresImport(ri, rs, location),
				ExpectError: testRequiresImportError("azurerm_storage_table"),
			},
		},
	})
}

func TestAccAzureRMStorageTable_disappears(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMStorageTable_basic(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
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
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	resourceName := "azurerm_storage_table.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageTable_acl(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageTableExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMStorageTable_aclUpdated(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageTableExists(resourceName),
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

func testCheckAzureRMStorageTableExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		tableName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		storageClient := testAccProvider.Meta().(*ArmClient).storage
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error locating Resource Group for Storage Table %q (Account %s): %s", tableName, accountName, err)
		}
		if resourceGroup == nil {
			return fmt.Errorf("Unable to locate Resource Group for Storage Table %q (Account %s)", tableName, accountName)
		}

		client, err := storageClient.TablesClient(ctx, *resourceGroup, accountName)
		if err != nil {
			return fmt.Errorf("Error building Table Client: %s", err)
		}

		props, err := client.Exists(ctx, accountName, tableName)
		if err != nil {
			if utils.ResponseWasNotFound(props) {
				return fmt.Errorf("Table %q doesn't exist in Account %q!", tableName, accountName)
			}

			return fmt.Errorf("Error retrieving Table %q: %s", tableName, accountName)
		}

		return nil
	}
}

func testAccARMStorageTableDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		tableName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		storageClient := testAccProvider.Meta().(*ArmClient).storage
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error locating Resource Group for Storage Table %q (Account %s): %s", tableName, accountName, err)
		}
		if resourceGroup == nil {
			return fmt.Errorf("Unable to locate Resource Group for Storage Table %q (Account %s)", tableName, accountName)
		}

		client, err := storageClient.TablesClient(ctx, *resourceGroup, accountName)
		if err != nil {
			return fmt.Errorf("Error building Table Client: %s", err)
		}

		props, err := client.Exists(ctx, accountName, tableName)
		if err != nil {
			if utils.ResponseWasNotFound(props) {
				return fmt.Errorf("Table %q doesn't exist in Account %q so it can't be deleted!", tableName, accountName)
			}

			return fmt.Errorf("Error retrieving Table %q: %s", tableName, accountName)
		}

		if _, err := client.Delete(ctx, accountName, tableName); err != nil {
			return fmt.Errorf("Error deleting Table %q: %s", tableName, err)
		}

		return nil
	}
}

func testCheckAzureRMStorageTableDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_table" {
			continue
		}

		tableName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		storageClient := testAccProvider.Meta().(*ArmClient).storage
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error locating Resource Group for Storage Table %q (Account %s): %s", tableName, accountName, err)
		}
		if resourceGroup == nil {
			return nil
		}

		client, err := storageClient.TablesClient(ctx, *resourceGroup, accountName)
		if err != nil {
			return fmt.Errorf("Error building Table Client: %s", err)
		}

		props, err := client.Exists(ctx, accountName, tableName)
		if err != nil {
			return nil
		}

		return fmt.Errorf("Table still exists: %+v", props)
	}

	return nil
}

func TestValidateArmStorageTableName(t *testing.T) {
	validNames := []string{
		"mytable01",
		"mytable",
		"myTable",
		"MYTABLE",
		"tbl",
		strings.Repeat("w", 63),
	}
	for _, v := range validNames {
		_, errors := validateArmStorageTableName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Storage Table Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"table",
		"-invalidname1",
		"invalid_name",
		"invalid!",
		"ww",
		strings.Repeat("w", 64),
	}
	for _, v := range invalidNames {
		_, errors := validateArmStorageTableName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Storage Table Name", v)
		}
	}
}

func testAccAzureRMStorageTable_basic(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_table" "test" {
  name                 = "acctestst%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}
`, rInt, location, rString, rInt)
}

func testAccAzureRMStorageTable_requiresImport(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageTable_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_table" "import" {
  name                 = "${azurerm_storage_table.test.name}"
  resource_group_name  = "${azurerm_storage_table.test.resource_group_name}"
  storage_account_name = "${azurerm_storage_table.test.storage_account_name}"
}
`, template)
}

func testAccAzureRMStorageTable_acl(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_table" "test" {
  name                 = "acctestst%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "raud"
      start       = "2020-11-26T08:49:37.0000000Z"
      expiry      = "2020-11-27T08:49:37.0000000Z"
    }
  }
}
`, rInt, location, rString, rInt)
}

func testAccAzureRMStorageTable_aclUpdated(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_table" "test" {
  name                 = "acctestst%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"

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
`, rInt, location, rString, rInt)
}
