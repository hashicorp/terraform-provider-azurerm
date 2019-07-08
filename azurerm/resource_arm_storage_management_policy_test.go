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
		CheckDestroy: testCheckAzureRMBatchCertificateDestroy, // TODO add this
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

// TODO - add multiple rules, filters and (?)actions
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
