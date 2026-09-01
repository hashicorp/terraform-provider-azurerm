// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type dataSourceStorageShare struct{}

func TestAccStorageShareDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_share", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: dataSourceStorageShare{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("quota").HasValue("5"),
				check.That(data.ResourceName).Key("metadata.%").HasValue("2"),
				check.That(data.ResourceName).Key("metadata.hello").HasValue("world"),
				check.That(data.ResourceName).Key("metadata.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("rbac_scope_id").MatchesRegex(regexp.MustCompile(`/fileshares/`)),
			),
		},
	})
}

func (d dataSourceStorageShare) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_storage_share" "test" {
  name               = azurerm_storage_share.test.name
  storage_account_id = azurerm_storage_account.test.id
}
`, StorageShareResource{}.complete(data))
}
