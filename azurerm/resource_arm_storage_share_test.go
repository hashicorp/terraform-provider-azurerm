package azurerm

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMStorageShare_basic(t *testing.T) {
	var sS storage.Share

	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMStorageShare_basic(ri, rs, testLocation())
	resourceName := "azurerm_storage_share.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(resourceName, &sS),
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

func TestAccAzureRMStorageShare_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	var sS storage.Share

	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	resourceName := "azurerm_storage_share.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageShare_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(resourceName, &sS),
				),
			},
			{
				Config:      testAccAzureRMStorageShare_requiresImport(ri, rs, location),
				ExpectError: testRequiresImportError("azurerm_storage_share"),
			},
		},
	})
}

func TestAccAzureRMStorageShare_disappears(t *testing.T) {
	var sS storage.Share

	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMStorageShare_basic(ri, rs, testLocation())
	resourceName := "azurerm_storage_share.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(resourceName, &sS),
					testAccARMStorageShareDisappears(resourceName, &sS),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMStorageShare_updateQuota(t *testing.T) {
	var sS storage.Share

	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMStorageShare_basic(ri, rs, testLocation())
	config2 := testAccAzureRMStorageShare_updateQuota(ri, rs, testLocation())
	resourceName := "azurerm_storage_share.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(resourceName, &sS),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageShareExists(resourceName, &sS),
					resource.TestCheckResourceAttr(resourceName, "quota", "5"),
				),
			},
		},
	})
}

func testCheckAzureRMStorageShareExists(resourceName string, sS *storage.Share) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		storageAccountName := rs.Primary.Attributes["storage_account_name"]
		resourceGroupName, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for share: %s", name)
		}

		armClient := testAccProvider.Meta().(*ArmClient)
		ctx := armClient.StopContext
		fileClient, accountExists, err := armClient.getFileServiceClientForStorageAccount(ctx, resourceGroupName, storageAccountName)
		if err != nil {
			return err
		}
		if !accountExists {
			return fmt.Errorf("Bad: Storage Account %q does not exist", storageAccountName)
		}

		shares, err := fileClient.ListShares(storage.ListSharesParameters{
			Prefix:  name,
			Timeout: 90,
		})
		if err != nil {
			return fmt.Errorf("Error listing Storage Share %q shares (storage account: %q) : %+v", name, storageAccountName, err)
		}

		if len(shares.Shares) == 0 {
			return fmt.Errorf("Bad: Share %q (storage account: %q) does not exist", name, storageAccountName)
		}

		var found bool
		for _, share := range shares.Shares {
			if share.Name == name {
				found = true
				*sS = share
			}
		}

		if !found {
			return fmt.Errorf("Bad: Share %q (storage account: %q) does not exist", name, storageAccountName)
		}

		return nil
	}
}

func testAccARMStorageShareDisappears(resourceName string, sS *storage.Share) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		armClient := testAccProvider.Meta().(*ArmClient)
		ctx := armClient.StopContext

		storageAccountName := rs.Primary.Attributes["storage_account_name"]
		resourceGroupName, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for storage share: %s", sS.Name)
		}

		fileClient, accountExists, err := armClient.getFileServiceClientForStorageAccount(ctx, resourceGroupName, storageAccountName)
		if err != nil {
			return err
		}
		if !accountExists {
			log.Printf("[INFO]Storage Account %q doesn't exist so the share won't exist", storageAccountName)
			return nil
		}

		reference := fileClient.GetShareReference(sS.Name)
		options := &storage.FileRequestOptions{}
		err = reference.Create(options)
		if err != nil {
			return fmt.Errorf("Error creating Storage Share %q reference (storage account: %q) : %+v", sS.Name, storageAccountName, err)
		}

		if _, err = reference.DeleteIfExists(options); err != nil {
			return fmt.Errorf("Error deleting storage Share %q: %s", sS.Name, err)
		}

		return nil
	}
}

func testCheckAzureRMStorageShareDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_share" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		storageAccountName := rs.Primary.Attributes["storage_account_name"]
		resourceGroupName, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for share: %s", name)
		}

		armClient := testAccProvider.Meta().(*ArmClient)
		ctx := armClient.StopContext
		fileClient, accountExists, err := armClient.getFileServiceClientForStorageAccount(ctx, resourceGroupName, storageAccountName)
		if err != nil {
			//If we can't get keys then the blob can't exist
			return nil
		}
		if !accountExists {
			return nil
		}

		shares, err := fileClient.ListShares(storage.ListSharesParameters{
			Prefix:  name,
			Timeout: 90,
		})

		if err != nil {
			return nil
		}

		var found bool
		for _, share := range shares.Shares {
			if share.Name == name {
				found = true
			}
		}

		if found {
			return fmt.Errorf("Bad: Share %q (storage account: %q) still exists", name, storageAccountName)
		}
	}

	return nil
}

func testAccAzureRMStorageShare_basic(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_share" "test" {
  name                 = "testshare%[3]s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}
`, rInt, location, rString)
}

func testAccAzureRMStorageShare_requiresImport(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageShare_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "import" {
  name                 = "${azurerm_storage_share.test.name}"
  resource_group_name  = "${azurerm_storage_share.test.resource_group_name}"
  storage_account_name = "${azurerm_storage_share.test.storage_account_name}"
}
`, template)
}

func testAccAzureRMStorageShare_updateQuota(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_share" "test" {
  name                 = "testshare%[3]s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
  quota                = 5
}
`, rInt, location, rString)
}

func TestValidateArmStorageShareName(t *testing.T) {
	validNames := []string{
		"valid-name",
		"valid02-name",
	}
	for _, v := range validNames {
		_, errors := validateArmStorageShareName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Share Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"InvalidName1",
		"-invalidname1",
		"invalid_name",
		"invalid!",
		"double-hyphen--invalid",
		"ww",
		strings.Repeat("w", 65),
	}
	for _, v := range invalidNames {
		_, errors := validateArmStorageShareName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Share Name", v)
		}
	}
}
