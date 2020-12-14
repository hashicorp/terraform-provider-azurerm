package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMStorageManagementPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountManagementPolicyDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageManagementPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountManagementPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.name", "rule1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.prefix_match.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.prefix_match.3439697764", "container1/prefix1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.1068358194", "blockBlob"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than", "50"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.delete_after_days_since_modification_greater_than", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.snapshot.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than", "30"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageManagementPolicy_singleAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountManagementPolicyDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageManagementPolicy_singleAction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountManagementPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.name", "singleActionRule"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.prefix_match.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.prefix_match.3439697764", "container1/prefix1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.1068358194", "blockBlob"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.snapshot.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than", "30"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageManagementPolicy_singleActionUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountManagementPolicyDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageManagementPolicy_singleAction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountManagementPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.name", "singleActionRule"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.prefix_match.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.prefix_match.3439697764", "container1/prefix1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.1068358194", "blockBlob"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.snapshot.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than", "30"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageManagementPolicy_singleActionUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountManagementPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.name", "singleActionRule"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.prefix_match.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.prefix_match.3439697764", "container1/prefix1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.1068358194", "blockBlob"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.delete_after_days_since_modification_greater_than", "30"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.snapshot.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than", "30"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageManagementPolicy_multipleRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountManagementPolicyDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageManagementPolicy_multipleRule(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountManagementPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "2"),

					// Rule1
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.name", "rule1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.prefix_match.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.prefix_match.3439697764", "container1/prefix1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.1068358194", "blockBlob"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than", "50"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.delete_after_days_since_modification_greater_than", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.snapshot.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than", "30"),

					// Rule2
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.name", "rule2"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.filters.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.filters.0.prefix_match.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.filters.0.prefix_match.4102595489", "container2/prefix1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.filters.0.prefix_match.1837232667", "container2/prefix2"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.filters.0.blob_types.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.filters.0.blob_types.1068358194", "blockBlob"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.0.base_blob.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than", "11"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than", "51"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.0.base_blob.0.delete_after_days_since_modification_greater_than", "101"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.0.snapshot.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.0.snapshot.0.delete_after_days_since_creation_greater_than", "31"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageManagementPolicy_updateMultipleRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountManagementPolicyDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageManagementPolicy_multipleRule(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountManagementPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "2"),

					// Rule1
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.name", "rule1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.prefix_match.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.prefix_match.3439697764", "container1/prefix1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.1068358194", "blockBlob"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than", "50"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.delete_after_days_since_modification_greater_than", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.snapshot.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than", "30"),

					// Rule2
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.name", "rule2"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.filters.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.filters.0.prefix_match.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.filters.0.prefix_match.4102595489", "container2/prefix1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.filters.0.prefix_match.1837232667", "container2/prefix2"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.filters.0.blob_types.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.filters.0.blob_types.1068358194", "blockBlob"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.0.base_blob.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than", "11"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than", "51"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.0.base_blob.0.delete_after_days_since_modification_greater_than", "101"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.0.snapshot.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.0.snapshot.0.delete_after_days_since_creation_greater_than", "31"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageManagementPolicy_multipleRuleUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountManagementPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "2"),

					// Rule1
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.name", "rule1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.prefix_match.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.prefix_match.3439697764", "container1/prefix1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.1068358194", "blockBlob"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than", "50"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.delete_after_days_since_modification_greater_than", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.snapshot.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than", "30"),

					// Rule2
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.name", "rule2"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.enabled", "true"), // check updated
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.filters.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.filters.0.prefix_match.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.filters.0.prefix_match.4102595489", "container2/prefix1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.filters.0.prefix_match.1837232667", "container2/prefix2"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.filters.0.blob_types.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.filters.0.blob_types.1068358194", "blockBlob"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.0.base_blob.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than", "12"),    // check updated
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than", "52"), // check updated
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.0.base_blob.0.delete_after_days_since_modification_greater_than", "102"),         // check updated
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.0.snapshot.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.1.actions.0.snapshot.0.delete_after_days_since_creation_greater_than", "32"), // check updated
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageManagementPolicy_blobTypes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountManagementPolicyDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageManagementPolicy_blobTypes(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountManagementPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.name", "rule1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.prefix_match.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.prefix_match.3439697764", "container1/prefix1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.1068358194", "blockBlob"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.932666486", "appendBlob"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.delete_after_days_since_modification_greater_than", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.snapshot.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than", "30"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMStorageAccountManagementPolicyDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "azurerm_storage_management_policy" {
				continue
			}
			storageAccountID := rs.Primary.Attributes["storage_account_id"]

			exists, err := testCheckAzureRMStorageAccountManagementPolicyExistsInternal(storageAccountID)
			if err != nil {
				return fmt.Errorf("Error checking if item has been destroyed: %s", err)
			}
			if exists {
				return fmt.Errorf("Bad: Storage Account Management Policy '%q' still exists", storageAccountID)
			}
		}

		return nil
	}
}

func testCheckAzureRMStorageAccountManagementPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		storageAccountID := rs.Primary.Attributes["storage_account_id"]

		exists, err := testCheckAzureRMStorageAccountManagementPolicyExistsInternal(storageAccountID)
		if err != nil {
			return fmt.Errorf("Error checking if item exists: %s", err)
		}
		if !exists {
			return fmt.Errorf("Bad: Storage Account Management Policy '%q' does not exist", storageAccountID)
		}

		return nil
	}
}

func testCheckAzureRMStorageAccountManagementPolicyExistsInternal(storageAccountID string) (bool, error) {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Storage.ManagementPoliciesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	rid, err := parse.StorageAccountID(storageAccountID)
	if err != nil {
		return false, fmt.Errorf("Bad: Failed to parse ID (id: %s): %+v", storageAccountID, err)
	}

	response, err := conn.Get(ctx, rid.ResourceGroup, rid.Name)
	if err != nil {
		if utils.ResponseWasNotFound(response.Response) {
			return false, nil
		}
		return false, fmt.Errorf("Bad: Get on storageAccount ManagementPolicy client (id: %s): %+v", storageAccountID, err)
	}

	return true, nil
}

func testAccAzureRMStorageManagementPolicy_basic(data acceptance.TestData) string {
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
  account_kind             = "BlobStorage"
}

resource "azurerm_storage_management_policy" "test" {
  storage_account_id = azurerm_storage_account.test.id

  rule {
    name    = "rule1"
    enabled = true
    filters {
      prefix_match = ["container1/prefix1"]
      blob_types   = ["blockBlob"]
    }
    actions {
      base_blob {
        tier_to_cool_after_days_since_modification_greater_than    = 10
        tier_to_archive_after_days_since_modification_greater_than = 50
        delete_after_days_since_modification_greater_than          = 100
      }
      snapshot {
        delete_after_days_since_creation_greater_than = 30
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageManagementPolicy_singleAction(data acceptance.TestData) string {
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
  account_kind             = "BlobStorage"
}

resource "azurerm_storage_management_policy" "test" {
  storage_account_id = azurerm_storage_account.test.id

  rule {
    name    = "singleActionRule"
    enabled = true
    filters {
      prefix_match = ["container1/prefix1"]
      blob_types   = ["blockBlob"]
    }
    actions {
      base_blob {
        tier_to_cool_after_days_since_modification_greater_than = 10
      }
      snapshot {
        delete_after_days_since_creation_greater_than = 30
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageManagementPolicy_singleActionUpdate(data acceptance.TestData) string {
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
  account_kind             = "BlobStorage"
}

resource "azurerm_storage_management_policy" "test" {
  storage_account_id = azurerm_storage_account.test.id

  rule {
    name    = "singleActionRule"
    enabled = true
    filters {
      prefix_match = ["container1/prefix1"]
      blob_types   = ["blockBlob"]
    }
    actions {
      base_blob {
        delete_after_days_since_modification_greater_than = 30
      }
      snapshot {
        delete_after_days_since_creation_greater_than = 30
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageManagementPolicy_multipleRule(data acceptance.TestData) string {
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
  account_kind             = "BlobStorage"
}

resource "azurerm_storage_management_policy" "test" {
  storage_account_id = azurerm_storage_account.test.id

  rule {
    name    = "rule1"
    enabled = true
    filters {
      prefix_match = ["container1/prefix1"]
      blob_types   = ["blockBlob"]
    }
    actions {
      base_blob {
        tier_to_cool_after_days_since_modification_greater_than    = 10
        tier_to_archive_after_days_since_modification_greater_than = 50
        delete_after_days_since_modification_greater_than          = 100
      }
      snapshot {
        delete_after_days_since_creation_greater_than = 30
      }
    }
  }
  rule {
    name    = "rule2"
    enabled = false
    filters {
      prefix_match = ["container2/prefix1", "container2/prefix2"]
      blob_types   = ["blockBlob"]
    }
    actions {
      base_blob {
        tier_to_cool_after_days_since_modification_greater_than    = 11
        tier_to_archive_after_days_since_modification_greater_than = 51
        delete_after_days_since_modification_greater_than          = 101
      }
      snapshot {
        delete_after_days_since_creation_greater_than = 31
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageManagementPolicy_multipleRuleUpdated(data acceptance.TestData) string {
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
  account_kind             = "BlobStorage"
}

resource "azurerm_storage_management_policy" "test" {
  storage_account_id = azurerm_storage_account.test.id

  rule {
    name    = "rule1"
    enabled = true
    filters {
      prefix_match = ["container1/prefix1"]
      blob_types   = ["blockBlob"]
    }
    actions {
      base_blob {
        tier_to_cool_after_days_since_modification_greater_than    = 10
        tier_to_archive_after_days_since_modification_greater_than = 50
        delete_after_days_since_modification_greater_than          = 100
      }
      snapshot {
        delete_after_days_since_creation_greater_than = 30
      }
    }
  }
  rule {
    name    = "rule2"
    enabled = true
    filters {
      prefix_match = ["container2/prefix1", "container2/prefix2"]
      blob_types   = ["blockBlob"]
    }
    actions {
      base_blob {
        tier_to_cool_after_days_since_modification_greater_than    = 12
        tier_to_archive_after_days_since_modification_greater_than = 52
        delete_after_days_since_modification_greater_than          = 102
      }
      snapshot {
        delete_after_days_since_creation_greater_than = 32
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMStorageManagementPolicy_blobTypes(data acceptance.TestData) string {
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
  account_kind             = "BlobStorage"
}

resource "azurerm_storage_management_policy" "test" {
  storage_account_id = azurerm_storage_account.test.id

  rule {
    name    = "rule1"
    enabled = true
    filters {
      prefix_match = ["container1/prefix1"]
      blob_types   = ["blockBlob", "appendBlob"]
    }
    actions {
      base_blob {
        delete_after_days_since_modification_greater_than = 100
      }
      snapshot {
        delete_after_days_since_creation_greater_than = 30
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
