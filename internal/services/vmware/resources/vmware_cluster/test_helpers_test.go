// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vmware_cluster_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type VmwarePrivateCloudResource struct{}

func (r VmwarePrivateCloudResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vmware_private_cloud" "test" {
  name                = "acctest-PC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "av36"

  management_cluster {
    size = 3
  }
  network_subnet_cidr = "192.168.48.0/22"
}
`, r.template(data), data.RandomInteger)
}

func (VmwarePrivateCloudResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  # In Vmware acctest, please disable correlation request id, else the continuous operations like update or delete will not be triggered
  # issue https://github.com/Azure/azure-rest-api-specs/issues/14086
  disable_correlation_request_id = true
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-Vmware-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
