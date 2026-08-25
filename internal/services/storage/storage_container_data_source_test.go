// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type StorageContainerDataSource struct{}

func TestAccStorageContainerDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_container", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageContainerDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("container_access_type").HasValue("private"),
				check.That(data.ResourceName).Key("has_immutability_policy").HasValue("false"),
				check.That(data.ResourceName).Key("default_encryption_scope").HasValue(fmt.Sprintf("acctestEScontainer%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("encryption_scope_override_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("metadata.%").HasValue("2"),
				check.That(data.ResourceName).Key("metadata.k1").HasValue("v1"),
				check.That(data.ResourceName).Key("metadata.k2").HasValue("v2"),
				check.That(data.ResourceName).Key("url").HasValue(fmt.Sprintf("https://acctestacc%[1]s.blob.core.windows.net/acctest-container-%[1]s", data.RandomString)),
			),
		},
	})
}

func (d StorageContainerDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

data "azurerm_storage_container" "test" {
  name               = azurerm_storage_container.test.name
  storage_account_id = azurerm_storage_account.test.id
}
`, StorageContainerResource{}.complete(data))
}
