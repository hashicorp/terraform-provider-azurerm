package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageManagementPolicyRuleResource struct{}

func TestAccStorageManagementPolicyRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy_rule", "test")
	r := StorageManagementPolicyRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageManagementPolicyRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy_rule", "test")
	r := StorageManagementPolicyRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageManagementPolicyRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy_rule", "test")
	r := StorageManagementPolicyRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageManagementPolicyRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_management_policy_rule", "test")
	r := StorageManagementPolicyRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r StorageManagementPolicyRuleResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := clients.Storage.ManagementPoliciesClient

	id, err := parse.StorageAccountManagementPolicyRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	policyId := parse.NewStorageAccountManagementPolicyID(id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.ManagementPolicyName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.StorageAccountName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", policyId, err)
	}
	res := storage.StorageAccountManagementPolicyRuleResource{}
	return utils.Bool(res.GetRuleInPolicy(resp, id.RuleName) != nil), nil
}

func (r StorageManagementPolicyRuleResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_management_policy_rule" "test" {
  name                 = "rule-%s"
  management_policy_id = azurerm_storage_management_policy.test.id
  filter {
    prefix_match = ["container1/prefix1"]
    blob_types   = ["blockBlob"]
  }
  action {
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
`, template, data.RandomString)
}

func (r StorageManagementPolicyRuleResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_management_policy_rule" "test" {
  name                 = "rule-%s"
  management_policy_id = azurerm_storage_management_policy.test.id
  enabled              = false
  filter {
    prefix_match = ["container2/prefix2"]
    blob_types   = ["blockBlob"]
  }
  action {
    base_blob {
      tier_to_cool_after_days_since_modification_greater_than        = 11
      tier_to_archive_after_days_since_modification_greater_than     = 51
      tier_to_archive_after_days_since_last_tier_change_greater_than = 20
      delete_after_days_since_modification_greater_than              = 101
    }
    snapshot {
      change_tier_to_archive_after_days_since_creation               = 91
      tier_to_archive_after_days_since_last_tier_change_greater_than = 20
      change_tier_to_cool_after_days_since_creation                  = 24
      delete_after_days_since_creation_greater_than                  = 31
    }
    version {
      change_tier_to_archive_after_days_since_creation               = 10
      tier_to_archive_after_days_since_last_tier_change_greater_than = 20
      change_tier_to_cool_after_days_since_creation                  = 91
      delete_after_days_since_creation                               = 4
    }
  }
}
`, template, data.RandomString)
}

func (r StorageManagementPolicyRuleResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_management_policy_rule" "import" {
  name                 = azurerm_storage_management_policy_rule.test.name
  management_policy_id = azurerm_storage_management_policy_rule.test.management_policy_id
  action {}
}
`, template)
}

func (r StorageManagementPolicyRuleResource) template(data acceptance.TestData) string {
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

  lifecycle {
    ignore_changes = [rule]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
