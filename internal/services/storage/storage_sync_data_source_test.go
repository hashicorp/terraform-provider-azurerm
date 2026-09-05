// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type StorageSyncDataSource struct{}

func TestAccDataSourceStorageSync_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_sync", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageSyncDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("incoming_traffic_policy").Exists(),
				check.That(data.ResourceName).Key("tags.%").Exists(),
			),
		},
	})
}

func TestAccDataSourceStorageSync_registeredServers(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_sync", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageSyncDataSource{}.registeredServers(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("registered_servers.#").HasValue("1"),
			),
		},
	})
}

func (d StorageSyncDataSource) basic(data acceptance.TestData) string {
	basic := StorageSyncResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_sync" "test" {
  name                = azurerm_storage_sync.test.name
  resource_group_name = azurerm_storage_sync.test.resource_group_name
}
`, basic)
}

func (d StorageSyncDataSource) registeredServers(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_storage_sync" "test" {
  name                = azurerm_storage_sync.test.name
  resource_group_name = azurerm_storage_sync.test.resource_group_name

  depends_on = [terraform_data.afs_server_id]
}
`, StorageSyncServerEndpointResource{}.template(data))
}
