package azurerm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parsers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMStorageDataLakeGen2FileSystem_basic(t *testing.T) {
	resourceName := "azurerm_storage_data_lake_gen2_filesystem.test"

	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageDataLakeGen2FileSystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageDataLakeGen2FileSystem_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageDataLakeGen2FileSystemExists(resourceName),
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

func TestAccAzureRMStorageDataLakeGen2FileSystem_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_storage_data_lake_gen2_filesystem.test"

	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageDataLakeGen2FileSystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageDataLakeGen2FileSystem_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageDataLakeGen2FileSystemExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMStorageDataLakeGen2FileSystem_requiresImport(ri, rs, location),
				ExpectError: testRequiresImportError("azurerm_storage_data_lake_gen2_filesystem"),
			},
		},
	})
}

func TestAccAzureRMStorageDataLakeGen2FileSystem_properties(t *testing.T) {
	resourceName := "azurerm_storage_data_lake_gen2_filesystem.test"

	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageDataLakeGen2FileSystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageDataLakeGen2FileSystem_properties(ri, rs, location, "aGVsbG8="),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageDataLakeGen2FileSystemExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMStorageDataLakeGen2FileSystem_properties(ri, rs, location, "ZXll"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageDataLakeGen2FileSystemExists(resourceName),
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

func testCheckAzureRMStorageDataLakeGen2FileSystemExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*ArmClient).Storage.FileSystemsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_data_lake_gen2_filesystem" {
			continue
		}

		client := testAccProvider.Meta().(*ArmClient).Storage.FileSystemsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMStorageDataLakeGen2FileSystem_basic(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageDataLakeGen2FileSystem_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%d"
  storage_account_id = azurerm_storage_account.test.id
}
`, template, rInt)
}

func testAccAzureRMStorageDataLakeGen2FileSystem_requiresImport(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageDataLakeGen2FileSystem_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_data_lake_gen2_filesystem" "import" {
  name               = azurerm_storage_data_lake_gen2_filesystem.test.name
  storage_account_id = azurerm_storage_data_lake_gen2_filesystem.storage_account_id
}
`, template)
}

func testAccAzureRMStorageDataLakeGen2FileSystem_properties(rInt int, rString, location, value string) string {
	template := testAccAzureRMStorageDataLakeGen2FileSystem_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%d"
  storage_account_id = azurerm_storage_account.test.id

  properties = {
    key = "%s"
  }
}
`, template, rInt, value)
}

func testAccAzureRMStorageDataLakeGen2FileSystem_template(rInt int, rString, location string) string {
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
`, rInt, location, rString)
}
