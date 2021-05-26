package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type StorageManagementPolicyResource struct{}

func TestAccStorageManagementPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy", "test")
	r := StorageManagementPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.name").HasValue("rule1"),
				check.That(data.ResourceName).Key("rule.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("rule.0.filters.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.prefix_match.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.blob_types.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than").HasValue("10"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than").HasValue("50"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.delete_after_days_since_modification_greater_than").HasValue("100"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than").HasValue("30"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageManagementPolicy_singleAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy", "test")
	r := StorageManagementPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.singleAction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.name").HasValue("singleActionRule"),
				check.That(data.ResourceName).Key("rule.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("rule.0.filters.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.prefix_match.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.blob_types.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than").HasValue("10"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than").HasValue("30"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageManagementPolicy_singleActionUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy", "test")
	r := StorageManagementPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.singleAction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.name").HasValue("singleActionRule"),
				check.That(data.ResourceName).Key("rule.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("rule.0.filters.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.prefix_match.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.blob_types.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than").HasValue("10"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than").HasValue("30"),
			),
		},
		data.ImportStep(),
		{
			Config: r.singleActionUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.name").HasValue("singleActionRule"),
				check.That(data.ResourceName).Key("rule.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("rule.0.filters.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.prefix_match.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.blob_types.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.delete_after_days_since_modification_greater_than").HasValue("30"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than").HasValue("30"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageManagementPolicy_multipleRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy", "test")
	r := StorageManagementPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rule.#").HasValue("2"),

				// Rule1
				check.That(data.ResourceName).Key("rule.0.name").HasValue("rule1"),
				check.That(data.ResourceName).Key("rule.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("rule.0.filters.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.prefix_match.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.blob_types.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than").HasValue("10"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than").HasValue("50"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.delete_after_days_since_modification_greater_than").HasValue("100"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than").HasValue("30"),

				// Rule2
				check.That(data.ResourceName).Key("rule.1.name").HasValue("rule2"),
				check.That(data.ResourceName).Key("rule.1.enabled").HasValue("false"),
				check.That(data.ResourceName).Key("rule.1.filters.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.1.filters.0.prefix_match.#").HasValue("2"),
				check.That(data.ResourceName).Key("rule.1.filters.0.blob_types.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.1.actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.1.actions.0.base_blob.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.1.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than").HasValue("11"),
				check.That(data.ResourceName).Key("rule.1.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than").HasValue("51"),
				check.That(data.ResourceName).Key("rule.1.actions.0.base_blob.0.delete_after_days_since_modification_greater_than").HasValue("101"),
				check.That(data.ResourceName).Key("rule.1.actions.0.snapshot.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.1.actions.0.snapshot.0.delete_after_days_since_creation_greater_than").HasValue("31"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageManagementPolicy_updateMultipleRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy", "test")
	r := StorageManagementPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rule.#").HasValue("2"),

				// Rule1
				check.That(data.ResourceName).Key("rule.0.name").HasValue("rule1"),
				check.That(data.ResourceName).Key("rule.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("rule.0.filters.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.prefix_match.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.blob_types.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than").HasValue("10"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than").HasValue("50"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.delete_after_days_since_modification_greater_than").HasValue("100"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than").HasValue("30"),

				// Rule2
				check.That(data.ResourceName).Key("rule.1.name").HasValue("rule2"),
				check.That(data.ResourceName).Key("rule.1.enabled").HasValue("false"),
				check.That(data.ResourceName).Key("rule.1.filters.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.1.filters.0.prefix_match.#").HasValue("2"),
				check.That(data.ResourceName).Key("rule.1.filters.0.blob_types.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.1.actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.1.actions.0.base_blob.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.1.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than").HasValue("11"),
				check.That(data.ResourceName).Key("rule.1.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than").HasValue("51"),
				check.That(data.ResourceName).Key("rule.1.actions.0.base_blob.0.delete_after_days_since_modification_greater_than").HasValue("101"),
				check.That(data.ResourceName).Key("rule.1.actions.0.snapshot.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.1.actions.0.snapshot.0.delete_after_days_since_creation_greater_than").HasValue("31"),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleRuleUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rule.#").HasValue("2"),

				// Rule1
				check.That(data.ResourceName).Key("rule.0.name").HasValue("rule1"),
				check.That(data.ResourceName).Key("rule.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("rule.0.filters.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.prefix_match.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.blob_types.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than").HasValue("10"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than").HasValue("50"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.delete_after_days_since_modification_greater_than").HasValue("100"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than").HasValue("30"),

				// Rule2
				check.That(data.ResourceName).Key("rule.1.name").HasValue("rule2"),
				check.That(data.ResourceName).Key("rule.1.enabled").HasValue("true"), // check updated
				check.That(data.ResourceName).Key("rule.1.filters.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.1.filters.0.prefix_match.#").HasValue("2"),
				check.That(data.ResourceName).Key("rule.1.filters.0.blob_types.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.1.actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.1.actions.0.base_blob.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.1.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than").HasValue("12"),    // check updated
				check.That(data.ResourceName).Key("rule.1.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than").HasValue("52"), // check updated
				check.That(data.ResourceName).Key("rule.1.actions.0.base_blob.0.delete_after_days_since_modification_greater_than").HasValue("102"),         // check updated
				check.That(data.ResourceName).Key("rule.1.actions.0.snapshot.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.1.actions.0.snapshot.0.delete_after_days_since_creation_greater_than").HasValue("32"), // check updated
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageManagementPolicy_blobTypes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy", "test")
	r := StorageManagementPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blobTypes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.name").HasValue("rule1"),
				check.That(data.ResourceName).Key("rule.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("rule.0.filters.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.prefix_match.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.blob_types.#").HasValue("2"),
				check.That(data.ResourceName).Key("rule.0.actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.delete_after_days_since_modification_greater_than").HasValue("100"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than").HasValue("30"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageManagementPolicy_blobIndexMatch(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy", "test")
	r := StorageManagementPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blobIndexMatchDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.blobIndexMatch(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.blobIndexMatchDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageManagementPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy", "test")
	r := StorageManagementPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageManagementPolicy_zero(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy", "test")
	r := StorageManagementPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zero(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageManagementPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy", "test")
	r := StorageManagementPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageManagementPolicyResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	storageAccountId := state.Attributes["storage_account_id"]
	id, err := parse.StorageAccountID(storageAccountId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Storage.ManagementPoliciesClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Management Policy (Account %q / Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r StorageManagementPolicyResource) basic(data acceptance.TestData) string {
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
    name    = "rule-1"
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

func (r StorageManagementPolicyResource) singleAction(data acceptance.TestData) string {
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

func (r StorageManagementPolicyResource) singleActionUpdate(data acceptance.TestData) string {
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

func (r StorageManagementPolicyResource) multipleRule(data acceptance.TestData) string {
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

func (r StorageManagementPolicyResource) multipleRuleUpdated(data acceptance.TestData) string {
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

func (r StorageManagementPolicyResource) blobTypes(data acceptance.TestData) string {
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

func (r StorageManagementPolicyResource) blobIndexMatchTemplate(data acceptance.TestData) string {
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
}

data "azurerm_client_config" "current" {}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Owner"
  principal_id         = data.azurerm_client_config.current.object_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageManagementPolicyResource) blobIndexMatch(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_management_policy" "test" {
  storage_account_id = azurerm_storage_account.test.id

  rule {
    name    = "rule1"
    enabled = true
    filters {
      prefix_match = ["container1/prefix1"]
      blob_types   = ["blockBlob"]

      match_blob_index_tag {
        name  = "tag1"
        value = "val1"
      }

      match_blob_index_tag {
        name      = "tag2"
        operation = "=="
        value     = "val2"
      }
    }
    actions {
      base_blob {
        tier_to_cool_after_days_since_modification_greater_than    = 10
        tier_to_archive_after_days_since_modification_greater_than = 50
        delete_after_days_since_modification_greater_than          = 100
      }
    }
  }

  depends_on = [azurerm_role_assignment.test]
}
`, r.blobIndexMatchTemplate(data))
}

func (r StorageManagementPolicyResource) blobIndexMatchDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

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
    }
  }

  depends_on = [azurerm_role_assignment.test]
}
`, r.blobIndexMatchTemplate(data))
}

func (r StorageManagementPolicyResource) complete(data acceptance.TestData) string {
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
        change_tier_to_archive_after_days_since_creation = 90
        change_tier_to_cool_after_days_since_creation    = 23
        delete_after_days_since_creation_greater_than    = 30
      }
      version {
        change_tier_to_archive_after_days_since_creation = 9
        change_tier_to_cool_after_days_since_creation    = 90
        delete_after_days_since_creation                 = 3
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageManagementPolicyResource) completeUpdate(data acceptance.TestData) string {
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
    name    = "rule2"
    enabled = true
    filters {
      prefix_match = ["container2/prefix2"]
      blob_types   = ["blockBlob"]
    }
    actions {
      base_blob {
        tier_to_cool_after_days_since_modification_greater_than    = 11
        tier_to_archive_after_days_since_modification_greater_than = 51
        delete_after_days_since_modification_greater_than          = 101
      }
      snapshot {
        change_tier_to_archive_after_days_since_creation = 91
        change_tier_to_cool_after_days_since_creation    = 24
        delete_after_days_since_creation_greater_than    = 31
      }
      version {
        change_tier_to_archive_after_days_since_creation = 10
        change_tier_to_cool_after_days_since_creation    = 91
        delete_after_days_since_creation                 = 4
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageManagementPolicyResource) zero(data acceptance.TestData) string {
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
        change_tier_to_archive_after_days_since_creation = 0
        change_tier_to_cool_after_days_since_creation    = 0
        delete_after_days_since_creation_greater_than    = 30
      }
      version {
        change_tier_to_archive_after_days_since_creation = 0
        change_tier_to_cool_after_days_since_creation    = 0
        delete_after_days_since_creation                 = 0
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
