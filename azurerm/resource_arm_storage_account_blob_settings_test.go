package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"testing"
)

func TestAccAzureRMStorageAccountBlobSettings_basic(t *testing.T) {
	resourceName := "azurerm_storage_account_blob_settings.testsabs"
	storageAccountResourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	createConfig := testAccAzureRMStorageAccountBlobSettings_basic(ri, rs, location)
	updateConfig := testAccAzureRMStorageAccountBlobSettings_update(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountBlobSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(storageAccountResourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_soft_delete", "true"),
					resource.TestCheckResourceAttr(resourceName, "soft_delete_retention_days", "123"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: updateConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(storageAccountResourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_soft_delete", "true"),
					resource.TestCheckResourceAttr(resourceName, "soft_delete_retention_days", "99"),
				),
			},
		},
	})
}

func testAccAzureRMStorageAccountBlobSettings_basic(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSABS-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account_blob_settings" "testsabs" {
  resource_group_name        = "${azurerm_resource_group.testrg.name}"
  storage_account_name       = "${azurerm_storage_account.testsa.name}"
  enable_soft_delete         = true
  soft_delete_retention_days = 123
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccountBlobSettings_update(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSABS-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account_blob_settings" "testsabs" {
  resource_group_name        = "${azurerm_resource_group.testrg.name}"
  storage_account_name       = "${azurerm_storage_account.testsa.name}"
  enable_soft_delete         = true
  soft_delete_retention_days = 99
}
`, rInt, location, rString)
}

func testCheckAzureRMStorageAccountBlobSettingsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).storageBlobServicesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_account_blob_settings" {
			continue
		}

		storageAccountName := rs.Primary.Attributes["storage_account_name"]
		resourceGroupName := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetServiceProperties(ctx, resourceGroupName, storageAccountName)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Storage Account Blob Settings still exist:\n%#v", resp.BlobServicePropertiesProperties)
			}
			return nil
		}
	}

	return nil
}
