// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

// TestAccStorageContainer_V1ToV2_4810 tests the regular migration path that uses `resource_manager_id`
// It uses v4.81.0 as the setup version because it is the last release where the resource was fully functional.
func TestAccStorageContainer_V1ToV2_4810(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_container", "test")
	r := StorageContainerResource{}

	data.ResourceRegressionTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicV1(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("https://acctestacc%s.blob.core.windows.net/vhds", data.RandomString)),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.Storage/storageAccounts/acctestacc%[3]s/blobServices/default/containers/vhds", data.Subscriptions.Primary, data.RandomInteger, data.RandomString)),
			),
		},
	}, "4.81.0")
}

func (r StorageContainerResource) basicV1(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}
`, r.template(data))
}
