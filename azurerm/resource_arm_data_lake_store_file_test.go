package azurerm

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestValidateAzureDataLakeStoreRemoteFilePath(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			Value:  "bad",
			Errors: 1,
		},
		{
			Value:  "/good/file/path",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		_, errors := validateDataLakeStoreRemoteFilePath()(tc.Value, "unittest")

		if len(errors) != tc.Errors {
			t.Fatalf("Expected validateDataLakeStoreRemoteFilePath to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}

func TestAccAzureRMDataLakeStoreFile_basic(t *testing.T) {
	resourceName := "azurerm_data_lake_store_file.test"

	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeStoreFile_basic(ri, rs, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreFileExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"local_file_path"},
			},
		},
	})
}

func TestAccAzureRMDataLakeStoreFile_largefiles(t *testing.T) {
	resourceName := "azurerm_data_lake_store_file.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	//"large" in this context is anything greater than 4 megabytes
	largeSize := 12 * 1024 * 1024 //12 mb
	data := make([]byte, largeSize)
	rand.Read(data) //fill with random data

	tmpfile, err := ioutil.TempFile("", "azurerm-acc-datalake-file-large")
	if err != nil {
		t.Errorf("Unable to open a temporary file.")
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(data); err != nil {
		t.Errorf("Unable to write to temporary file %q: %v", tmpfile.Name(), err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Errorf("Unable to close temporary file %q: %v", tmpfile.Name(), err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeStoreFile_largefiles(ri, rs, acceptance.Location(), tmpfile.Name()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreFileExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"local_file_path"},
			},
		},
	})
}

func TestAccAzureRMDataLakeStoreFile_requiresimport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_data_lake_store_file.test"

	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeStoreFile_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreFileExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMDataLakeStoreFile_requiresImport(ri, rs, location),
				ExpectError: acceptance.RequiresImportError("azurerm_data_lake_store_file"),
			},
		},
	})
}

func testCheckAzureRMDataLakeStoreFileExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		remoteFilePath := rs.Primary.Attributes["remote_file_path"]
		accountName := rs.Primary.Attributes["account_name"]

		conn := acceptance.AzureProvider.Meta().(*clients.Client).Datalake.StoreFilesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := conn.GetFileStatus(ctx, accountName, remoteFilePath, utils.Bool(true))
		if err != nil {
			return fmt.Errorf("Bad: Get on dataLakeStoreFileClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Date Lake Store File Rule %q (Account %q) does not exist", remoteFilePath, accountName)
		}

		return nil
	}
}

func testCheckAzureRMDataLakeStoreFileDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Datalake.StoreFilesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_lake_store_file" {
			continue
		}

		remoteFilePath := rs.Primary.Attributes["remote_file_path"]
		accountName := rs.Primary.Attributes["account_name"]

		resp, err := conn.GetFileStatus(ctx, accountName, remoteFilePath, utils.Bool(true))
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("Data Lake Store File still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMDataLakeStoreFile_basic(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_lake_store" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "%s"
  firewall_state      = "Disabled"
}

resource "azurerm_data_lake_store_file" "test" {
  remote_file_path = "/test/application_gateway_test.cer"
  account_name     = "${azurerm_data_lake_store.test.name}"
  local_file_path  = "./testdata/application_gateway_test.cer"
}
`, rInt, location, rString, location)
}

func testAccAzureRMDataLakeStoreFile_largefiles(rInt int, rString, location, file string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_lake_store" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "%s"
  firewall_state      = "Disabled"
}

resource "azurerm_data_lake_store_file" "test" {
  remote_file_path = "/test/testAccAzureRMDataLakeStoreFile_largefiles.bin"
  account_name     = "${azurerm_data_lake_store.test.name}"
  local_file_path  = "%s"
}
`, rInt, location, rString, location, file)
}

func testAccAzureRMDataLakeStoreFile_requiresImport(rInt int, rString, location string) string {
	template := testAccAzureRMDataLakeStoreFile_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_store_file" "import" {
  remote_file_path = "${azurerm_data_lake_store_file.test.remote_file_path}"
  account_name     = "${azurerm_data_lake_store_file.test.name}"
  local_file_path  = "./testdata/application_gateway_test.cer"
}
`, template)
}
