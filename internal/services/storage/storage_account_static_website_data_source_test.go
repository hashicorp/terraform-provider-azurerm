// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type StorageAccountStaticWebsiteDataSource struct{}

func TestAccStorageAccountStaticWebsiteDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account_static_website", "test")
	r := StorageAccountStaticWebsiteDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
check.That(data.ResourceName).Key("id").MatchesOtherKey(check.That("azurerm_storage_account.test").Key("id")),
				check.That(data.ResourceName).Key("error_404_document").HasValue("sadpanda_2.html"),
				check.That(data.ResourceName).Key("index_document").HasValue("index_2.html"),
			),
		},
	})
}

func (d StorageAccountStaticWebsiteDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_storage_account_static_website" "test" {
  storage_account_id = azurerm_storage_account_static_website.test.storage_account_id
}
`, StorageAccountStaticWebsiteResource{}.complete(data))
}
