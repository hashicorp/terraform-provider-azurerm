package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDataLakeStoreFile_basic(t *testing.T) {
	resourceName := "azurerm_data_lake_store_file.test"

	ri := acctest.RandInt()
	rs := acctest.RandString(4)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeStoreFile_basic(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreFileExists(resourceName),
				),
			},
		},
	})
}

func testCheckAzureRMDataLakeStoreFileExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		remoteFilePath := rs.Primary.Attributes["remote_file_path"]
		accountName := rs.Primary.Attributes["account_name"]

		conn := testAccProvider.Meta().(*ArmClient).dataLakeStoreFilesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
	conn := testAccProvider.Meta().(*ArmClient).dataLakeStoreFilesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMDataLakeStoreFile_basic(rInt int, rs, location string) string {
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
  remote_file_path    = "/test/application_gateway_test.cer"
  account_name        = "${azurerm_data_lake_store.test.name}"
  local_file_path     = "./testdata/application_gateway_test.cer"
}
`, rInt, location, rs, location)
}
