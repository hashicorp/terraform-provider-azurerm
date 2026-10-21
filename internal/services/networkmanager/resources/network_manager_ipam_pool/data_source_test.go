// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_ipam_pool_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type ManagerIpamPoolDataSource struct{}

func (d ManagerIpamPoolDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_network_manager_ipam_pool" "test" {
  name               = azurerm_network_manager_ipam_pool.test.name
  network_manager_id = azurerm_network_manager_ipam_pool.test.network_manager_id
}
`, ManagerIpamPoolResource{}.complete(data))
}
