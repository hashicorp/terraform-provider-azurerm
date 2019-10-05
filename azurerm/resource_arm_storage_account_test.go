package azurerm

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
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
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_basic(ri, rs, location)
	postConfig := testAccAzureRMStorageAccount_update(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
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
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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

func TestAccAzureRMStorageAccount_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMStorageAccount_requiresImport(ri, rs, location),
				ExpectError: testRequiresImportError("azurerm_storage_account"),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_writeLock(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_writeLock(ri, rs, location),
			},
			{
				// works around a bug in the test suite where the Storage Account won't be re-read after the Lock's provisioned
				Config: testAccAzureRMStorageAccount_writeLock(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "account_tier", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "account_replication_type", "LRS"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "production"),
					resource.TestCheckResourceAttr(resourceName, "primary_connection_string", ""),
					resource.TestCheckResourceAttr(resourceName, "secondary_connection_string", ""),
					resource.TestCheckResourceAttr(resourceName, "primary_blob_connection_string", ""),
					resource.TestCheckResourceAttr(resourceName, "secondary_blob_connection_string", ""),
					resource.TestCheckResourceAttr(resourceName, "primary_access_key", ""),
					resource.TestCheckResourceAttr(resourceName, "secondary_access_key", ""),
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

func TestAccAzureRMStorageAccount_premium(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_premium(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMStorageAccount_disappears(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	preConfig := testAccAzureRMStorageAccount_basic(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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
	resourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	preConfig := testAccAzureRMStorageAccount_basic(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttrSet("azurerm_storage_account.testsa", "primary_blob_connection_string"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_blobEncryption(t *testing.T) {
	_, exists := os.LookupEnv("TF_ACC_STORAGE_ENCRYPTION_DISABLE")
	if !exists {
		t.Skip("`TF_ACC_STORAGE_ENCRYPTION_DISABLE` isn't specified - skipping since disabling encryption is generally disabled")
	}

	resourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_blobEncryption(ri, rs, location)
	postConfig := testAccAzureRMStorageAccount_blobEncryptionDisabled(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "enable_blob_encryption", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "enable_blob_encryption", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_fileEncryption(t *testing.T) {
	_, exists := os.LookupEnv("TF_ACC_STORAGE_ENCRYPTION_DISABLE")
	if !exists {
		t.Skip("`TF_ACC_STORAGE_ENCRYPTION_DISABLE` isn't specified - skipping since disabling encryption is generally disabled")
	}

	resourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_fileEncryption(ri, rs, location)
	postConfig := testAccAzureRMStorageAccount_fileEncryptionDisabled(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "enable_file_encryption", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "enable_file_encryption", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_enableHttpsTrafficOnly(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_enableHttpsTrafficOnly(ri, rs, location)
	postConfig := testAccAzureRMStorageAccount_enableHttpsTrafficOnlyDisabled(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "enable_https_traffic_only", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "enable_https_traffic_only", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_isHnsEnabled(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_isHnsEnabledTrue(ri, rs, location)
	postConfig := testAccAzureRMStorageAccount_isHnsEnabledFalse(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "is_hns_enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "is_hns_enabled", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_blobStorageWithUpdate(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_blobStorage(ri, rs, location)
	postConfig := testAccAzureRMStorageAccount_blobStorageUpdate(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "account_kind", "BlobStorage"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "access_tier", "Hot"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "access_tier", "Cool"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_blockBlobStorage(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_blockBlobStorage(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "account_kind", "BlockBlobStorage"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "access_tier", ""),
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

func TestAccAzureRMStorageAccount_fileStorageWithUpdate(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_fileStorage(ri, rs, location)
	postConfig := testAccAzureRMStorageAccount_fileStorageUpdate(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "account_kind", "FileStorage"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "account_tier", "Premium"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "access_tier", "Hot"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "account_tier", "Premium"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "access_tier", "Cool"),
				),
			},
		},
	})
}
func TestAccAzureRMStorageAccount_storageV2WithUpdate(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_storageV2(ri, rs, location)
	postConfig := testAccAzureRMStorageAccount_storageV2Update(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "account_kind", "StorageV2"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "access_tier", "Hot"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "access_tier", "Cool"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_NonStandardCasing(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	preConfig := testAccAzureRMStorageAccount_nonStandardCasing(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             preConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAzureRMStorageAccount_enableIdentity(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"

	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	config := testAccAzureRMStorageAccount_identity(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(resourceName, "identity.0.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(resourceName, "identity.0.tenant_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_updateResourceByEnablingIdentity(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"

	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)

	basicResourceNoManagedIdentity := testAccAzureRMStorageAccount_basic(ri, rs, testLocation())
	managedIdentityEnabled := testAccAzureRMStorageAccount_identity(ri, rs, testLocation())

	uuidMatch := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: basicResourceNoManagedIdentity,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.#", "0"),
				),
			},
			{
				Config: managedIdentityEnabled,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(resourceName, "identity.0.principal_id", uuidMatch),
					resource.TestMatchResourceAttr(resourceName, "identity.0.tenant_id", uuidMatch),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_networkRules(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_networkRules(ri, rs, location)
	postConfig := testAccAzureRMStorageAccount_networkRulesUpdate(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "network_rules.0.default_action", "Deny"),
					resource.TestCheckResourceAttr(resourceName, "network_rules.0.ip_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_rules.0.virtual_network_subnet_ids.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "network_rules.0.default_action", "Deny"),
					resource.TestCheckResourceAttr(resourceName, "network_rules.0.ip_rules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "network_rules.0.virtual_network_subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_rules.0.bypass.#", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_networkRulesDeleted(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_networkRules(ri, rs, location)
	postConfig := testAccAzureRMStorageAccount_networkRulesReverted(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "network_rules.0.default_action", "Deny"),
					resource.TestCheckResourceAttr(resourceName, "network_rules.0.ip_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_rules.0.virtual_network_subnet_ids.#", "1"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "network_rules.0.default_action", "Allow"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_enableAdvancedThreatProtection(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_enableAdvancedThreatProtection(ri, rs, location)
	postConfig := testAccAzureRMStorageAccount_enableAdvancedThreatProtectionDisabled(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "enable_advanced_threat_protection", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "enable_advanced_threat_protection", "false"),
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

func TestAccAzureRMStorageAccount_queueProperties(t *testing.T) {
	resourceName := "azurerm_storage_account.testsa"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccAzureRMStorageAccount_queueProperties(ri, rs, location)
	postConfig := testAccAzureRMStorageAccount_queuePropertiesUpdated(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(resourceName),
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

func testCheckAzureRMStorageAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		storageAccount := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		// Ensure resource group exists in API
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		conn := testAccProvider.Meta().(*ArmClient).Storage.AccountsClient

		resp, err := conn.GetProperties(ctx, resourceGroup, storageAccount, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on storageServiceClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: StorageAccount %q (resource group: %q) does not exist", storageAccount, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMStorageAccountDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		storageAccount := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		// Ensure resource group exists in API
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		conn := testAccProvider.Meta().(*ArmClient).Storage.AccountsClient

		if _, err := conn.Delete(ctx, resourceGroup, storageAccount); err != nil {
			return fmt.Errorf("Bad: Delete on storageServiceClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMStorageAccountDestroy(s *terraform.State) error {
	ctx := testAccProvider.Meta().(*ArmClient).StopContext
	conn := testAccProvider.Meta().(*ArmClient).Storage.AccountsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_account" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.GetProperties(ctx, resourceGroup, name, "")
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
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_requiresImport(rInt int, rString string, location string) string {
	template := testAccAzureRMStorageAccount_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "import" {
  name                     = "${azurerm_storage_account.testsa.name}"
  resource_group_name      = "${azurerm_storage_account.testrg.resource_group_name}"
  location                 = "${azurerm_storage_account.testrg.location}"
  account_tier             = "${azurerm_storage_account.testrg.account_tier}"
  account_replication_type = "${azurerm_storage_account.testrg.account_replication_type}"
}
`, template)
}

func testAccAzureRMStorageAccount_writeLock(rInt int, rString, location string) string {
	template := testAccAzureRMStorageAccount_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = "${azurerm_storage_account.testsa.id}"
  lock_level = "ReadOnly"
}
`, template, rInt)
}

func testAccAzureRMStorageAccount_premium(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Premium"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_update(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_blobEncryption(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  enable_blob_encryption   = true

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_blobEncryptionDisabled(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  enable_blob_encryption   = false

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_fileEncryption(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  enable_file_encryption   = true

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_fileEncryptionDisabled(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  enable_file_encryption   = false

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_enableHttpsTrafficOnly(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                  = "${azurerm_resource_group.testrg.location}"
  account_tier              = "Standard"
  account_replication_type  = "LRS"
  enable_https_traffic_only = true

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_enableHttpsTrafficOnlyDisabled(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                  = "${azurerm_resource_group.testrg.location}"
  account_tier              = "Standard"
  account_replication_type  = "LRS"
  enable_https_traffic_only = false

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_isHnsEnabledTrue(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  is_hns_enabled           = true
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_isHnsEnabledFalse(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  is_hns_enabled           = false
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_blobStorage(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_blobStorageUpdate(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  access_tier              = "Cool"

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_blockBlobStorage(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_kind             = "BlockBlobStorage"
  account_tier             = "Premium"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_fileStorage(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_kind             = "FileStorage"
  account_tier             = "Premium"
  account_replication_type = "LRS"
  access_tier              = "Hot"

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_fileStorageUpdate(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_kind             = "FileStorage"
  account_tier             = "Premium"
  account_replication_type = "LRS"
  access_tier              = "Cool"

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_storageV2(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_storageV2Update(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  access_tier              = "Cool"

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_nonStandardCasing(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = "${azurerm_resource_group.testrg.name}"
  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "standard"
  account_replication_type = "lrs"

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_identity(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "production"
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_networkRules(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.testrg.location}"
  resource_group_name = "${azurerm_resource_group.testrg.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.testrg.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_storage_account" "testsa" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = "${azurerm_resource_group.testrg.name}"
  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  network_rules {
    default_action             = "Deny"
    ip_rules                   = ["127.0.0.1"]
    virtual_network_subnet_ids = ["${azurerm_subnet.test.id}"]
  }

  tags = {
    environment = "production"
  }
}
`, rInt, location, rInt, rInt, rString)
}

func testAccAzureRMStorageAccount_networkRulesUpdate(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.testrg.location}"
  resource_group_name = "${azurerm_resource_group.testrg.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.testrg.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_storage_account" "testsa" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = "${azurerm_resource_group.testrg.name}"
  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  network_rules {
    default_action = "Deny"
    ip_rules       = ["127.0.0.1", "127.0.0.2"]
    bypass         = ["Logging", "Metrics"]
  }

  tags = {
    environment = "production"
  }
}
`, rInt, location, rInt, rInt, rString)
}

func testAccAzureRMStorageAccount_networkRulesReverted(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.testrg.location}"
  resource_group_name = "${azurerm_resource_group.testrg.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.testrg.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_storage_account" "testsa" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = "${azurerm_resource_group.testrg.name}"
  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  network_rules {
    default_action             = "Allow"
    ip_rules                   = ["127.0.0.1"]
    virtual_network_subnet_ids = ["${azurerm_subnet.test.id}"]
  }

  tags = {
    environment = "production"
  }
}
`, rInt, location, rInt, rInt, rString)
}

func testAccAzureRMStorageAccount_enableAdvancedThreatProtection(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                              = "unlikely23exst2acct%s"
  resource_group_name               = "${azurerm_resource_group.testrg.name}"
  location                          = "${azurerm_resource_group.testrg.location}"
  account_tier                      = "Standard"
  account_replication_type          = "LRS"
  enable_advanced_threat_protection = true
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_enableAdvancedThreatProtectionDisabled(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                              = "unlikely23exst2acct%s"
  resource_group_name               = "${azurerm_resource_group.testrg.name}"
  location                          = "${azurerm_resource_group.testrg.location}"
  account_tier                      = "Standard"
  account_replication_type          = "LRS"
  enable_advanced_threat_protection = false
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_queueProperties(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  queue_properties {
    cors_rule {
      allowed_origins    = ["http://www.example.com"]
      exposed_headers    = ["x-tempo-*"]
      allowed_headers    = ["x-tempo-*"]
      allowed_methods    = ["GET", "PUT"]
      max_age_in_seconds = "500"
    }

    logging {
      version               = "1.0"
      delete                = true
      read                  = true
      write                 = true
      retention_policy_days = 7
    }

    hour_metrics {
      version               = "1.0"
      enabled               = false
      retention_policy_days = 7
    }

    minute_metrics {
      version               = "1.0"
      enabled               = false
      retention_policy_days = 7
    }
  }
}
`, rInt, location, rString)
}

func testAccAzureRMStorageAccount_queuePropertiesUpdated(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "testsa" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  queue_properties {
    cors_rule {
      allowed_origins    = ["http://www.example.com"]
      exposed_headers    = ["x-tempo-*", "x-method-*"]
      allowed_headers    = ["*"]
      allowed_methods    = ["GET"]
      max_age_in_seconds = "2000000000"
		}
		cors_rule {
      allowed_origins    = ["http://www.test.com"]
      exposed_headers    = ["x-tempo-*"]
      allowed_headers    = ["*"]
      allowed_methods    = ["PUT"]
      max_age_in_seconds = "1000"
    }
    logging {
      version               = "1.0"
      delete                = true
      read                  = true
      write                 = true
      retention_policy_days = 7
    }

    hour_metrics {
      version               = "1.0"
      enabled               = true
      retention_policy_days = 7
      include_apis          = true
    }

    minute_metrics {
      version               = "1.0"
      enabled               = true
      include_apis          = false
      retention_policy_days = 7
    }
  }
}
`, rInt, location, rString)
}
