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
)

func TestAccAzureRMStorageShareDirectory_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(5))
	location := testLocation()
	resourceName := "azurerm_storage_share_directory.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageShareDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShareDirectory_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareDirectoryExists(resourceName),
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

func TestAccAzureRMStorageShareDirectory_uppercase(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(5))
	location := testLocation()
	resourceName := "azurerm_storage_share_directory.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageShareDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShareDirectory_uppercase(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareDirectoryExists(resourceName),
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

func TestAccAzureRMStorageShareDirectory_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(5))
	location := testLocation()
	resourceName := "azurerm_storage_share_directory.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageShareDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShareDirectory_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareDirectoryExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMStorageShareDirectory_requiresImport(ri, rs, location),
				ExpectError: testRequiresImportError("azurerm_storage_share_directory"),
			},
		},
	})
}

func TestAccAzureRMStorageShareDirectory_complete(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(5))
	location := testLocation()
	resourceName := "azurerm_storage_share_directory.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageShareDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShareDirectory_complete(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareDirectoryExists(resourceName),
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

func TestAccAzureRMStorageShareDirectory_update(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(5))
	location := testLocation()
	resourceName := "azurerm_storage_share_directory.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageShareDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShareDirectory_complete(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareDirectoryExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMStorageShareDirectory_updated(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareDirectoryExists(resourceName),
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
func TestAccAzureRMStorageShareDirectory_nested(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(5))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageShareDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShareDirectory_nested(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareDirectoryExists("azurerm_storage_share_directory.parent"),
					testCheckAzureRMStorageShareDirectoryExists("azurerm_storage_share_directory.child"),
				),
			},
		},
	})
}

func testCheckAzureRMStorageShareDirectoryExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		shareName := rs.Primary.Attributes["share_name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		storageClient := testAccProvider.Meta().(*ArmClient).Storage
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error locating Resource Group for Storage Share Directory %q (Share %s, Account %s): %s", name, shareName, accountName, err)
		}
		if resourceGroup == nil {
			return fmt.Errorf("Unable to locate Resource Group for Storage Share Directory %q (Share %s, Account %s) ", name, shareName, accountName)
		}

		client, err := storageClient.FileShareDirectoriesClient(ctx, *resourceGroup, accountName)
		if err != nil {
			return fmt.Errorf("Error building FileShare Client: %s", err)
		}

		resp, err := client.Get(ctx, accountName, shareName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on FileShareDirectoriesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Directory %q (File Share %q / Account %q / Resource Group %q) does not exist", name, shareName, accountName, *resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMStorageShareDirectoryDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_share_directory" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		shareName := rs.Primary.Attributes["share_name"]
		accountName := rs.Primary.Attributes["storage_account_name"]

		storageClient := testAccProvider.Meta().(*ArmClient).Storage
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
		if err != nil {
			return fmt.Errorf("Error locating Resource Group for Storage Share Directory %q (Share %s, Account %s): %s", name, shareName, accountName, err)
		}

		// not found, the account's gone
		if resourceGroup == nil {
			return nil
		}

		client, err := storageClient.FileShareDirectoriesClient(ctx, *resourceGroup, accountName)
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

func testAccAzureRMStorageShareDirectory_basic(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageShareDirectory_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "test" {
  name                 = "dir"
  share_name           = "${azurerm_storage_share.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}
`, template)
}

func testAccAzureRMStorageShareDirectory_uppercase(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageShareDirectory_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "test" {
  name                 = "UpperCaseCharacterS"
  share_name           = "${azurerm_storage_share.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}
`, template)
}

func testAccAzureRMStorageShareDirectory_requiresImport(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageShareDirectory_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "import" {
  name                 = "${azurerm_storage_share_directory.test.name}"
  share_name           = "${azurerm_storage_share_directory.test.share_name}"
  storage_account_name = "${azurerm_storage_share_directory.test.storage_account_name}"
}
`, template)
}

func testAccAzureRMStorageShareDirectory_complete(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageShareDirectory_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "test" {
  name                 = "dir"
  share_name           = "${azurerm_storage_share.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"

  metadata = {
    hello = "world"
  }
}
`, template)
}

func testAccAzureRMStorageShareDirectory_updated(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageShareDirectory_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "test" {
  name                 = "dir"
  share_name           = "${azurerm_storage_share.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"

  metadata = {
    hello    = "world"
    sunshine = "at dawn"
  }
}
`, template)
}

func testAccAzureRMStorageShareDirectory_nested(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageShareDirectory_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "parent" {
  name                 = "parent"
  share_name           = "${azurerm_storage_share.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}

resource "azurerm_storage_share_directory" "child" {
  name                 = "${azurerm_storage_share_directory.parent.name}/child"
  share_name           = "${azurerm_storage_share.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}
`, template)
}

func testAccAzureRMStorageShareDirectory_template(rInt int, rString string, location string) string {
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

resource "azurerm_storage_share" "test" {
  name                 = "fileshare"
  storage_account_name = "${azurerm_storage_account.test.name}"
  quota                = 50
}
`, rInt, location, rString)
}
