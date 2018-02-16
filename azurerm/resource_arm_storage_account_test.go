package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestValidateArmStorageAccountType(t *testing.T) {
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{"standard_lrs", false},
		{"invalid", true},
	}

	for _, test := range testCases {
		_, es := validateArmStorageAccountType(test.input, "account_type")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating account_type %q to fail", test.input)
		}
	}
}

func TestValidateArmStorageAccountName(t *testing.T) {
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{"ab", true},
		{"ABC", true},
		{"abc", false},
		{"123456789012345678901234", false},
		{"1234567890123456789012345", true},
		{"abc12345", false},
	}

	for _, test := range testCases {
		_, es := validateArmStorageAccountName(test.input, "name")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating name %q to fail", test.input)
		}
	}
}

func TestAccAzureRMStorageAccount_basic(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_basic(ri, rs, location)
	postConfig := testAccAzureRMStorageAccount_update(ri, rs, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "account_tier", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "account_replication_type", "LRS"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "production"),
				),
			},

			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "account_tier", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "account_replication_type", "GRS"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "staging"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_premium(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_premium(ri, rs, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "account_tier", "Premium"),
					resource.TestCheckResourceAttr(resourceName, "account_replication_type", "LRS"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "production"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_disappears(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	preConfig := testAccAzureRMStorageAccount_basic(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "account_tier", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "account_replication_type", "LRS"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "production"),
					testCheckAzureRMStorageAccountDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMStorageAccount_blobConnectionString(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	preConfig := testAccAzureRMStorageAccount_basic(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttrSet("azurerm_storage_account.testsa", "primary_blob_connection_string"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_blobEncryption(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_blobEncryption(ri, rs, location)
	postConfig := testAccAzureRMStorageAccount_blobEncryptionDisabled(ri, rs, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "enable_blob_encryption", "true"),
				),
			},

			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "enable_blob_encryption", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_fileEncryption(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_fileEncryption(ri, rs, location)
	postConfig := testAccAzureRMStorageAccount_fileEncryptionDisabled(ri, rs, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "enable_file_encryption", "true"),
				),
			},

			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "enable_file_encryption", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_enableHttpsTrafficOnly(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_enableHttpsTrafficOnly(ri, rs, location)
	postConfig := testAccAzureRMStorageAccount_enableHttpsTrafficOnlyDisabled(ri, rs, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "enable_https_traffic_only", "true"),
				),
			},

			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "enable_https_traffic_only", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_blobStorageWithUpdate(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_blobStorage(ri, rs, location)
	postConfig := testAccAzureRMStorageAccount_blobStorageUpdate(ri, rs, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "account_kind", "BlobStorage"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "access_tier", "Hot"),
				),
			},

			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "access_tier", "Cool"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_storageV2WithUpdate(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_storageV2(ri, rs, location)
	postConfig := testAccAzureRMStorageAccount_storageV2Update(ri, rs, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "account_kind", "StorageV2"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "access_tier", "Hot"),
				),
			},

			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "access_tier", "Cool"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_NonStandardCasing(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	preConfig := testAccAzureRMStorageAccount_nonStandardCasing(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
				),
			},

			{
				Config:             preConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testCheckAzureRMStorageAccountExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		storageAccount := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		// Ensure resource group exists in API
		conn := testAccProvider.Meta().(*ArmClient).storageServiceClient

		resp, err := conn.GetProperties(resourceGroup, storageAccount)
		if err != nil {
			return fmt.Errorf("Bad: Get on storageServiceClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: StorageAccount %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMStorageAccountDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		storageAccount := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		// Ensure resource group exists in API
		conn := testAccProvider.Meta().(*ArmClient).storageServiceClient

		_, err := conn.Delete(resourceGroup, storageAccount)
		if err != nil {
			return fmt.Errorf("Bad: Delete on storageServiceClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMStorageAccountDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).storageServiceClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_account" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.GetProperties(resourceGroup, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Storage Account still exists:\n%#v", resp.AccountProperties)
		}
	}

	return nil
}

func testAccAzureRMStorageAccount_basic(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "%s"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "${azurerm_resource_group.testrg.location}"
    account_tier = "Standard"
    account_replication_type = "LRS"

    tags {
        environment = "production"
    }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_premium(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "%s"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "${azurerm_resource_group.testrg.location}"
    account_tier = "Premium"
    account_replication_type = "LRS"

    tags {
        environment = "production"
    }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_update(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "%s"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "${azurerm_resource_group.testrg.location}"
    account_tier = "Standard"
    account_replication_type = "GRS"

    tags {
        environment = "staging"
    }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_blobEncryption(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "%s"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "${azurerm_resource_group.testrg.location}"
    account_tier = "Standard"
    account_replication_type = "LRS"
    enable_blob_encryption = true

    tags {
        environment = "production"
    }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_blobEncryptionDisabled(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "%s"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "${azurerm_resource_group.testrg.location}"
    account_tier = "Standard"
    account_replication_type = "LRS"
    enable_blob_encryption = false

    tags {
        environment = "production"
    }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_fileEncryption(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "%s"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "${azurerm_resource_group.testrg.location}"
    account_tier = "Standard"
    account_replication_type = "LRS"
    enable_file_encryption = true

    tags {
        environment = "production"
    }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_fileEncryptionDisabled(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "%s"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "${azurerm_resource_group.testrg.location}"
    account_tier = "Standard"
    account_replication_type = "LRS"
    enable_file_encryption = false

    tags {
        environment = "production"
    }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_enableHttpsTrafficOnly(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "%s"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "${azurerm_resource_group.testrg.location}"
    account_tier = "Standard"
    account_replication_type = "LRS"
    enable_https_traffic_only = true

    tags {
        environment = "production"
    }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_enableHttpsTrafficOnlyDisabled(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "%s"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "${azurerm_resource_group.testrg.location}"
    account_tier = "Standard"
    account_replication_type = "LRS"
    enable_https_traffic_only = false

    tags {
        environment = "production"
    }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_blobStorage(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "%s"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "${azurerm_resource_group.testrg.location}"
    account_kind = "BlobStorage"
    account_tier = "Standard"
    account_replication_type = "LRS"

    tags {
        environment = "production"
    }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_blobStorageUpdate(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "%s"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "${azurerm_resource_group.testrg.location}"
    account_kind = "BlobStorage"
    account_tier = "Standard"
    account_replication_type = "LRS"
    access_tier = "Cool"

    tags {
        environment = "production"
    }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_storageV2(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "%s"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "${azurerm_resource_group.testrg.location}"
    account_kind = "StorageV2"
    account_tier = "Standard"
    account_replication_type = "LRS"

    tags {
        environment = "production"
    }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_storageV2Update(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "%s"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "${azurerm_resource_group.testrg.location}"
    account_kind = "StorageV2"
    account_tier = "Standard"
    account_replication_type = "LRS"
    access_tier = "Cool"

    tags {
        environment = "production"
    }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_nonStandardCasing(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "%s"
}
resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"
    location = "${azurerm_resource_group.testrg.location}"
    account_tier = "standard"
    account_replication_type = "lrs"
    tags {
        environment = "production"
    }
}
`, rInt, location, rString)
}
