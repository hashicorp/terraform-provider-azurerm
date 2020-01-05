package tests

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
		_, es := storage.ValidateArmStorageAccountType(test.input, "account_type")

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
		_, es := storage.ValidateArmStorageAccountName(test.input, "name")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating name %q to fail", test.input)
		}
	}
}

func TestAccAzureRMStorageAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	preConfig := testAccAzureRMStorageAccount_basic(data)
	postConfig := testAccAzureRMStorageAccount_update(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "account_tier", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_replication_type", "LRS"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "production"),
				),
			},
			data.ImportStep(),
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "account_tier", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_replication_type", "GRS"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "staging"),
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
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMStorageAccount_requiresImport),
		},
	})
}

func TestAccAzureRMStorageAccount_writeLock(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_writeLock(data),
			},
			{
				// works around a bug in the test suite where the Storage Account won't be re-read after the Lock's provisioned
				Config: testAccAzureRMStorageAccount_writeLock(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "account_tier", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_replication_type", "LRS"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "production"),
					resource.TestCheckResourceAttr(data.ResourceName, "primary_connection_string", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "secondary_connection_string", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "primary_blob_connection_string", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "secondary_blob_connection_string", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "primary_access_key", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "secondary_access_key", ""),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageAccount_premium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	preConfig := testAccAzureRMStorageAccount_premium(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "account_tier", "Premium"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_replication_type", "LRS"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "production"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageAccount_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	preConfig := testAccAzureRMStorageAccount_basic(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "account_tier", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_replication_type", "LRS"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "production"),
					testCheckAzureRMStorageAccountDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMStorageAccount_blobConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	preConfig := testAccAzureRMStorageAccount_basic(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_blob_connection_string"),
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
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	preConfig := testAccAzureRMStorageAccount_blobEncryption(data)
	postConfig := testAccAzureRMStorageAccount_blobEncryptionDisabled(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_blob_encryption", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_blob_encryption", "false"),
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
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	preConfig := testAccAzureRMStorageAccount_fileEncryption(data)
	postConfig := testAccAzureRMStorageAccount_fileEncryptionDisabled(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_file_encryption", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_file_encryption", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_enableHttpsTrafficOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	preConfig := testAccAzureRMStorageAccount_enableHttpsTrafficOnly(data)
	postConfig := testAccAzureRMStorageAccount_enableHttpsTrafficOnlyDisabled(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_https_traffic_only", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_https_traffic_only", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_isHnsEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	preConfig := testAccAzureRMStorageAccount_isHnsEnabledTrue(data)
	postConfig := testAccAzureRMStorageAccount_isHnsEnabledFalse(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "is_hns_enabled", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "is_hns_enabled", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_blobStorageWithUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	preConfig := testAccAzureRMStorageAccount_blobStorage(data)
	postConfig := testAccAzureRMStorageAccount_blobStorageUpdate(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "account_kind", "BlobStorage"),
					resource.TestCheckResourceAttr(data.ResourceName, "access_tier", "Hot"),
				),
			},
			data.ImportStep(),
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "access_tier", "Cool"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_blockBlobStorage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_blockBlobStorage(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "account_kind", "BlockBlobStorage"),
					resource.TestCheckResourceAttr(data.ResourceName, "access_tier", ""),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageAccount_fileStorageWithUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	preConfig := testAccAzureRMStorageAccount_fileStorage(data)
	postConfig := testAccAzureRMStorageAccount_fileStorageUpdate(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "account_kind", "FileStorage"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_tier", "Premium"),
					resource.TestCheckResourceAttr(data.ResourceName, "access_tier", "Hot"),
				),
			},
			data.ImportStep(),
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "account_tier", "Premium"),
					resource.TestCheckResourceAttr(data.ResourceName, "access_tier", "Cool"),
				),
			},
		},
	})
}
func TestAccAzureRMStorageAccount_storageV2WithUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	preConfig := testAccAzureRMStorageAccount_storageV2(data)
	postConfig := testAccAzureRMStorageAccount_storageV2Update(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "account_kind", "StorageV2"),
					resource.TestCheckResourceAttr(data.ResourceName, "access_tier", "Hot"),
				),
			},
			data.ImportStep(),
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "access_tier", "Cool"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_NonStandardCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	preConfig := testAccAzureRMStorageAccount_nonStandardCasing(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config:             preConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAzureRMStorageAccount_enableIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	config := testAccAzureRMStorageAccount_identity(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.tenant_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_updateResourceByEnablingIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	basicResourceNoManagedIdentity := testAccAzureRMStorageAccount_basic(data)
	managedIdentityEnabled := testAccAzureRMStorageAccount_identity(data)

	uuidMatch := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: basicResourceNoManagedIdentity,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.#", "0"),
				),
			},
			{
				Config: managedIdentityEnabled,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", uuidMatch),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.tenant_id", uuidMatch),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_networkRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	preConfig := testAccAzureRMStorageAccount_networkRules(data)
	postConfig := testAccAzureRMStorageAccount_networkRulesUpdate(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rules.0.default_action", "Deny"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rules.0.ip_rules.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rules.0.virtual_network_subnet_ids.#", "1"),
				),
			},
			data.ImportStep(),
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rules.0.default_action", "Deny"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rules.0.ip_rules.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rules.0.virtual_network_subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rules.0.bypass.#", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_networkRulesDeleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	preConfig := testAccAzureRMStorageAccount_networkRules(data)
	postConfig := testAccAzureRMStorageAccount_networkRulesReverted(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rules.0.default_action", "Deny"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rules.0.ip_rules.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rules.0.virtual_network_subnet_ids.#", "1"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rules.0.default_action", "Allow"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_enableAdvancedThreatProtection(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	preConfig := testAccAzureRMStorageAccount_enableAdvancedThreatProtection(data)
	postConfig := testAccAzureRMStorageAccount_enableAdvancedThreatProtectionDisabled(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_advanced_threat_protection", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_advanced_threat_protection", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageAccount_blobProperties(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_blobProperties(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "blob_properties.0.delete_retention_policy.0.days", "300"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_blobPropertiesUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "blob_properties.0.delete_retention_policy.0.days", "7"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageAccount_queueProperties(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_queueProperties(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_queuePropertiesUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
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
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Storage.AccountsClient

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
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Storage.AccountsClient

		if _, err := conn.Delete(ctx, resourceGroup, storageAccount); err != nil {
			return fmt.Errorf("Bad: Delete on storageServiceClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMStorageAccountDestroy(s *terraform.State) error {
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Storage.AccountsClient

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

func testAccAzureRMStorageAccount_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStorageAccount_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "import" {
  name                     = "${azurerm_storage_account.test.name}"
  resource_group_name      = "${azurerm_storage_account.test.resource_group_name}"
  location                 = "${azurerm_storage_account.test.location}"
  account_tier             = "${azurerm_storage_account.test.account_tier}"
  account_replication_type = "${azurerm_storage_account.test.account_replication_type}"
}
`, template)
}

func testAccAzureRMStorageAccount_writeLock(data acceptance.TestData) string {
	template := testAccAzureRMStorageAccount_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = "${azurerm_storage_account.test.id}"
  lock_level = "ReadOnly"
}
`, template, data.RandomInteger)
}

func testAccAzureRMStorageAccount_premium(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Premium"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_blobEncryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  enable_blob_encryption   = true

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_blobEncryptionDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  enable_blob_encryption   = false

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_fileEncryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  enable_file_encryption   = true

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_fileEncryptionDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  enable_file_encryption   = false

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_enableHttpsTrafficOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                  = "${azurerm_resource_group.test.location}"
  account_tier              = "Standard"
  account_replication_type  = "LRS"
  enable_https_traffic_only = true

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_enableHttpsTrafficOnlyDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                  = "${azurerm_resource_group.test.location}"
  account_tier              = "Standard"
  account_replication_type  = "LRS"
  enable_https_traffic_only = false

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_isHnsEnabledTrue(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  is_hns_enabled           = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_isHnsEnabledFalse(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  is_hns_enabled           = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_blobStorage(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_blobStorageUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  access_tier              = "Cool"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_blockBlobStorage(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_kind             = "BlockBlobStorage"
  account_tier             = "Premium"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_fileStorage(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_kind             = "FileStorage"
  account_tier             = "Premium"
  account_replication_type = "LRS"
  access_tier              = "Hot"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_fileStorageUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_kind             = "FileStorage"
  account_tier             = "Premium"
  account_replication_type = "LRS"
  access_tier              = "Cool"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_storageV2(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_storageV2Update(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  access_tier              = "Cool"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_nonStandardCasing(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "standard"
  account_replication_type = "lrs"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_networkRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString)
}

func testAccAzureRMStorageAccount_networkRulesUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString)
}

func testAccAzureRMStorageAccount_networkRulesReverted(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString)
}

func testAccAzureRMStorageAccount_enableAdvancedThreatProtection(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                              = "unlikely23exst2acct%s"
  resource_group_name               = "${azurerm_resource_group.test.name}"
  location                          = "${azurerm_resource_group.test.location}"
  account_tier                      = "Standard"
  account_replication_type          = "LRS"
  enable_advanced_threat_protection = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_enableAdvancedThreatProtectionDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                              = "unlikely23exst2acct%s"
  resource_group_name               = "${azurerm_resource_group.test.name}"
  location                          = "${azurerm_resource_group.test.location}"
  account_tier                      = "Standard"
  account_replication_type          = "LRS"
  enable_advanced_threat_protection = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_blobProperties(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  blob_properties {
    delete_retention_policy {
      days = 300
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_blobPropertiesUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  blob_properties {
    delete_retention_policy {
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_queueProperties(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_queuePropertiesUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
