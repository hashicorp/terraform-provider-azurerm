package tests

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"
)

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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_basic(data),
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
				Config: testAccAzureRMStorageAccount_update(data),
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

func TestAccAzureRMStorageAccount_tagCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_tagCount(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_premium(data),
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_basic(data),
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_blob_connection_string"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_enableHttpsTrafficOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_enableHttpsTrafficOnly(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_https_traffic_only", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_enableHttpsTrafficOnlyDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_https_traffic_only", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_minTLSVersion(t *testing.T) {
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
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_minTLSVersion(data, "TLS1_0"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_minTLSVersion(data, "TLS1_1"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_minTLSVersion(data, "TLS1_2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_minTLSVersion(data, "TLS1_1"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_allowBlobPublicAccess(t *testing.T) {
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
					resource.TestCheckResourceAttr(data.ResourceName, "allow_blob_public_access", "false"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_allowBlobPublicAccess(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_blob_public_access", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_disAllowBlobPublicAccess(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_blob_public_access", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_isHnsEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_isHnsEnabledTrue(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "is_hns_enabled", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_isHnsEnabledFalse(data),
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_blobStorage(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "account_kind", "BlobStorage"),
					resource.TestCheckResourceAttr(data.ResourceName, "access_tier", "Hot"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_blobStorageUpdate(data),
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_fileStorage(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "account_kind", "FileStorage"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_tier", "Premium"),
					resource.TestCheckResourceAttr(data.ResourceName, "access_tier", "Hot"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_fileStorageUpdate(data),
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_storageV2(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "account_kind", "StorageV2"),
					resource.TestCheckResourceAttr(data.ResourceName, "access_tier", "Hot"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_storageV2Update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "access_tier", "Cool"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_storageV1ToV2Update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_storageToV2Prep(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "account_kind", "Storage"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_storageToV2Update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "account_kind", "StorageV2"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_NonStandardCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_nonStandardCasing(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config:             testAccAzureRMStorageAccount_nonStandardCasing(data),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAzureRMStorageAccount_enableIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_identity(data),
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_networkRules(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rules.0.default_action", "Deny"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rules.0.ip_rules.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rules.0.virtual_network_subnet_ids.#", "1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_networkRulesUpdate(data),
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_networkRules(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rules.0.default_action", "Deny"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rules.0.ip_rules.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rules.0.virtual_network_subnet_ids.#", "1"),
				),
			},
			{
				Config: testAccAzureRMStorageAccount_networkRulesReverted(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rules.0.default_action", "Allow"),
				),
			},
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
				Config: testAccAzureRMStorageAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_blobProperties(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "blob_properties.0.cors_rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "blob_properties.0.delete_retention_policy.0.days", "300"),
				),
			},
			data.ImportStep(),
			{
				PreConfig: func() { time.Sleep(10 * time.Minute) },
				Config:    testAccAzureRMStorageAccount_blobPropertiesUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "blob_properties.0.cors_rule.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "blob_properties.0.delete_retention_policy.0.days", "7"),
				),
			},
			data.ImportStep(),
			{
				PreConfig: func() { time.Sleep(10 * time.Minute) },
				Config:    testAccAzureRMStorageAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageAccount_blobPropertiesWithSoftDeleteContainerEnabled(t *testing.T) {
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
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_blobPropertiesWithSoftDeleteContainerEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				PreConfig: func() { time.Sleep(7 * time.Minute) },
				Config:    testAccAzureRMStorageAccount_blobPropertiesUpdatedWithSoftDeleteContainerEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				PreConfig: func() { time.Sleep(7 * time.Minute) },
				Config:    testAccAzureRMStorageAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
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

func TestAccAzureRMStorageAccount_staticWebsiteEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_staticWebsiteEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				// Disabled
				Config: testAccAzureRMStorageAccount_storageV2(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_staticWebsiteEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageAccount_staticWebsitePropertiesForStorageV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_staticWebsitePropertiesForStorageV2(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_staticWebsitePropertiesUpdatedForStorageV2(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageAccount_staticWebsitePropertiesForBlockBlobStorage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_staticWebsitePropertiesForBlockBlobStorage(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_staticWebsitePropertiesUpdatedForBlockBlobStorage(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageAccount_replicationTypeGZRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageAccount_replicationTypeGZRS(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "account_replication_type", "GZRS"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_replicationTypeRAGZRS(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "account_replication_type", "RAGZRS"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_largeFileShare(t *testing.T) {
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
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_largeFileShareDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageAccount_largeFileShareEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config:      testAccAzureRMStorageAccount_largeFileShareDisabled(data),
				ExpectError: regexp.MustCompile("`large_file_share_enabled` cannot be disabled once it's been enabled"),
			},
		},
	})
}

func testCheckAzureRMStorageAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure resource group exists in API
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Storage.AccountsClient

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		storageAccount := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

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
		// Ensure resource group exists in API
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Storage.AccountsClient

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		storageAccount := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

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
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_tagCount(data acceptance.TestData) string {
	tags := ""
	for i := 0; i < 50; i++ {
		tags += fmt.Sprintf("t%d = \"v%d\"\n", i, i)
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
            %s
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, tags)
}

func testAccAzureRMStorageAccount_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStorageAccount_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "import" {
  name                     = azurerm_storage_account.test.name
  resource_group_name      = azurerm_storage_account.test.resource_group_name
  location                 = azurerm_storage_account.test.location
  account_tier             = azurerm_storage_account.test.account_tier
  account_replication_type = azurerm_storage_account.test.account_replication_type
}
`, template)
}

func testAccAzureRMStorageAccount_writeLock(data acceptance.TestData) string {
	template := testAccAzureRMStorageAccount_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = azurerm_storage_account.test.id
  lock_level = "ReadOnly"
}
`, template, data.RandomInteger)
}

func testAccAzureRMStorageAccount_premium(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
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
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_enableHttpsTrafficOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                  = azurerm_resource_group.test.location
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
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                  = azurerm_resource_group.test.location
  account_tier              = "Standard"
  account_replication_type  = "LRS"
  enable_https_traffic_only = false

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_minTLSVersion(data acceptance.TestData, tlsVersion string) string {
	return fmt.Sprintf(`

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  min_tls_version          = "%s"

  tags = {
    environment = "production"
  }
}

	`, data.RandomInteger, data.Locations.Primary, data.RandomString, tlsVersion)
}

func testAccAzureRMStorageAccount_allowBlobPublicAccess(data acceptance.TestData) string {
	return fmt.Sprintf(`

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  allow_blob_public_access = true

  tags = {
    environment = "production"
  }
}

	`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_disAllowBlobPublicAccess(data acceptance.TestData) string {
	return fmt.Sprintf(`

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  allow_blob_public_access = false

  tags = {
    environment = "production"
  }
}

	`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_isHnsEnabledTrue(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  is_hns_enabled           = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_isHnsEnabledFalse(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  is_hns_enabled           = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_blobStorage(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
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
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
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
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
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
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
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
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
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
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
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
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
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

func testAccAzureRMStorageAccount_storageToV2Prep(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_kind             = "Storage"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_storageToV2Update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_nonStandardCasing(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
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
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
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
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  network_rules {
    default_action             = "Deny"
    ip_rules                   = ["127.0.0.1"]
    virtual_network_subnet_ids = [azurerm_subnet.test.id]
  }

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString)
}

func testAccAzureRMStorageAccount_networkRulesUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
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
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  network_rules {
    default_action             = "Allow"
    ip_rules                   = ["127.0.0.1"]
    virtual_network_subnet_ids = [azurerm_subnet.test.id]
  }

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString)
}

func testAccAzureRMStorageAccount_blobProperties(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  blob_properties {
    cors_rule {
      allowed_origins    = ["http://www.example.com"]
      exposed_headers    = ["x-tempo-*"]
      allowed_headers    = ["x-tempo-*"]
      allowed_methods    = ["GET", "PUT", "PATCH"]
      max_age_in_seconds = "500"
    }

    delete_retention_policy {
      days = 300
    }
    versioning_enabled  = true
    change_feed_enabled = true
    restore_policy {
      days = 5
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_blobPropertiesUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  blob_properties {
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

    delete_retention_policy {
    }

    versioning_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_blobPropertiesWithSoftDeleteContainerEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  blob_properties {
    cors_rule {
      allowed_origins    = ["http://www.example.com"]
      exposed_headers    = ["x-tempo-*"]
      allowed_headers    = ["x-tempo-*"]
      allowed_methods    = ["GET", "PUT", "PATCH"]
      max_age_in_seconds = "500"
    }

    delete_retention_policy {
      days = 300
    }
    versioning_enabled  = true
    change_feed_enabled = true
    restore_policy {
      days = 5
    }
    container_delete_retention_policy {
      days = 7
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_blobPropertiesUpdatedWithSoftDeleteContainerEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestAzureRMSA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  blob_properties {
    cors_rule {
      allowed_origins    = ["http://www.example.com"]
      exposed_headers    = ["x-tempo-*"]
      allowed_headers    = ["x-tempo-*"]
      allowed_methods    = ["GET", "PUT", "PATCH"]
      max_age_in_seconds = "500"
    }

    delete_retention_policy {
      days = 300
    }
    versioning_enabled  = true
    change_feed_enabled = true
    restore_policy {
      days = 5
    }
    container_delete_retention_policy {
      days = 3
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_queueProperties(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
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
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
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

func testAccAzureRMStorageAccount_staticWebsiteEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  static_website {}
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_staticWebsitePropertiesForStorageV2(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  static_website {
    index_document     = "index.html"
    error_404_document = "404.html"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_staticWebsitePropertiesForBlockBlobStorage(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_kind             = "BlockBlobStorage"
  account_tier             = "Premium"
  account_replication_type = "LRS"

  static_website {
    index_document     = "index.html"
    error_404_document = "404.html"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_staticWebsitePropertiesUpdatedForStorageV2(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  static_website {
    index_document     = "index-2.html"
    error_404_document = "404-2.html"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_staticWebsitePropertiesUpdatedForBlockBlobStorage(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_kind             = "BlockBlobStorage"
  account_tier             = "Premium"
  account_replication_type = "LRS"

  static_website {
    index_document     = "index-2.html"
    error_404_document = "404-2.html"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_replicationTypeGZRS(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GZRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_replicationTypeRAGZRS(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "RAGZRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_largeFileShareDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  large_file_share_enabled = false

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageAccount_largeFileShareEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  large_file_share_enabled = true

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
