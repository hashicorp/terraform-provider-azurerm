package azurerm

import (
	"fmt"
	// "net/http"
	// "os"
	// "regexp"
	"testing"

	// "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	// "github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMStorageManagementPolicy_basic(t *testing.T) {
	resourceName := "azurerm_storage_management_policy.testpolicy"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	config := testAccAzureRMStorageManagementPolicy_basic(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.name", "rule1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.type", "Lifecycle"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.filters.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.filters.0.prefix_match.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.filters.0.prefix_match.3439697764", "container1/prefix1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.filters.0.blob_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.filters.0.blob_types.1068358194", "blockBlob"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.actions.0.base_blob.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than", "10"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than", "50"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.actions.0.base_blob.0.delete_after_days_since_modification_greater_than", "100"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.actions.0.snapshot.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than", "30"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageManagementPolicy_multipleRule(t *testing.T) {
	resourceName := "azurerm_storage_management_policy.testpolicy"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	config := testAccAzureRMStorageManagementPolicy_multipleRule(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rule.#", "2"),
					// Rule1
					resource.TestCheckResourceAttr(resourceName, "rule.0.name", "rule1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.type", "Lifecycle"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.filters.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.filters.0.prefix_match.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.filters.0.prefix_match.3439697764", "container1/prefix1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.filters.0.blob_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.filters.0.blob_types.1068358194", "blockBlob"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.actions.0.base_blob.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than", "10"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than", "50"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.actions.0.base_blob.0.delete_after_days_since_modification_greater_than", "100"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.actions.0.snapshot.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than", "30"),
					// Rule2
					resource.TestCheckResourceAttr(resourceName, "rule.1.name", "rule2"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.type", "Lifecycle"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.filters.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.filters.0.prefix_match.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.filters.0.prefix_match.4102595489", "container2/prefix1"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.filters.0.prefix_match.1837232667", "container2/prefix2"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.filters.0.blob_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.filters.0.blob_types.1068358194", "blockBlob"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.actions.0.base_blob.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than", "11"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than", "51"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.actions.0.base_blob.0.delete_after_days_since_modification_greater_than", "101"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.actions.0.snapshot.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.actions.0.snapshot.0.delete_after_days_since_creation_greater_than", "31"),
				),
			},
		},
	})
}
func testAccAzureRMStorageManagementPolicy_basic(rInt int, rString string, location string) string {
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
	account_kind             = "BlobStorage"
}

resource "azurerm_storage_management_policy" "testpolicy" {
	storage_account_id = "${azurerm_storage_account.testsa.id}"

	rule {
		name    = "rule1"
		enabled = true
		type    = "Lifecycle"
		filters {
			prefix_match = [ "container1/prefix1" ]
			blob_types  = [ "blockBlob" ]
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
`, rInt, location, rString)
}

func testAccAzureRMStorageManagementPolicy_multipleRule(rInt int, rString string, location string) string {
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
	account_kind             = "BlobStorage"
}

resource "azurerm_storage_management_policy" "testpolicy" {
	storage_account_id = "${azurerm_storage_account.testsa.id}"

	rule {
		name    = "rule1"
		enabled = true
		type    = "Lifecycle"
		filters {
			prefix_match = [ "container1/prefix1" ]
			blob_types  = [ "blockBlob" ]
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
		type    = "Lifecycle"
		filters {
		  prefix_match = [ "container2/prefix1", "container2/prefix2" ]
		  blob_types  = [ "blockBlob" ]
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
`, rInt, location, rString)
}
