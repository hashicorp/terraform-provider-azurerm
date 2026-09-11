package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

// TestAccStorageQueue_V1ToV2_4810 tests the regular migration path that uses `resource_manager_id`
// It uses v4.81.0 as the setup version because it is the last release where the resource was fully functional.
func TestAccStorageQueue_V1ToV2_4810(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_queue", "test")
	r := StorageQueueResource{}

	data.ResourceRegressionTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicV1(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("https://acctestacc%s.queue.core.windows.net/acctestmysamplequeue-%d", data.RandomString, data.RandomInteger)),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.Storage/storageAccounts/acctestacc%[3]s/queueServices/default/queues/acctestmysamplequeue-%[2]d", data.Subscriptions.Primary, data.RandomInteger, data.RandomString)),
			),
		},
	}, "4.81.0")
}

// TestAccStorageQueue_V1ToV2_3380 tests the fallback `FindAccount` path of the state migration
// It uses v3.38.0 as the setup version because `resource_manager_id` was introduced in v3.39.0.
func TestAccStorageQueue_V1ToV2_3380(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_queue", "test")
	r := StorageQueueResource{}

	data.ResourceRegressionTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicV1(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("https://acctestacc%s.queue.core.windows.net/acctestmysamplequeue-%d", data.RandomString, data.RandomInteger)),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.Storage/storageAccounts/acctestacc%[3]s/queueServices/default/queues/acctestmysamplequeue-%[2]d", data.Subscriptions.Primary, data.RandomInteger, data.RandomString)),
			),
		},
	}, "3.38.0")
}

// TestAccStorageQueue_V1ToV2_4810_AzureADAuth covers https://github.com/hashicorp/terraform-provider-azurerm/issues/32950
// A queue created in v4 using `storage_account_name` with `storage_use_azuread = true` persisted a
// `resource_manager_id` containing an empty `resourceGroups` segment, because that Azure AD path never
// resolved the Resource Group. This broke the v1->v2 state migration when upgrading to 5.x. It uses
// v4.81.0 as the setup version because it is the last release where the resource was fully functional.
func TestAccStorageQueue_V1ToV2_4810_AzureADAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_queue", "test")
	r := StorageQueueResource{}

	data.ResourceRegressionTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureADAuthV1(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("https://acctestacc%s.queue.core.windows.net/acctestmysamplequeue-%d", data.RandomString, data.RandomInteger)),
			),
		},
		{
			Config: r.basicAzureADAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.Storage/storageAccounts/acctestacc%[3]s/queueServices/default/queues/acctestmysamplequeue-%[2]d", data.Subscriptions.Primary, data.RandomInteger, data.RandomString)),
			),
		},
	}, "4.81.0")
}

func (r StorageQueueResource) basicV1(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "test" {
  name                 = "acctestmysamplequeue-%d"
  storage_account_name = azurerm_storage_account.test.name
}
`, r.template(data), data.RandomInteger)
}

func (r StorageQueueResource) azureADAuthV1(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  storage_use_azuread = true
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_queue" "test" {
  name                 = "acctestmysamplequeue-%[1]d"
  storage_account_name = azurerm_storage_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
