// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type BastionHostDataSource struct{}

func TestAccBastionHostDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_bastion_host", "test")
	r := BastionHostDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("sku").Exists(),
				check.That(data.ResourceName).Key("dns_name").Exists(),
				check.That(data.ResourceName).Key("scale_units").Exists(),
				check.That(data.ResourceName).Key("file_copy_enabled").Exists(),
				check.That(data.ResourceName).Key("ip_connect_enabled").Exists(),
				check.That(data.ResourceName).Key("shareable_link_enabled").Exists(),
				check.That(data.ResourceName).Key("session_recording_enabled").Exists(),
				check.That(data.ResourceName).Key("ip_configuration.0.name").Exists(),
				check.That(data.ResourceName).Key("ip_configuration.0.subnet_id").Exists(),
				check.That(data.ResourceName).Key("ip_configuration.0.public_ip_address_id").Exists(),
			),
		},
	})
}

func (BastionHostDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_bastion_host" "test" {
  name                = azurerm_bastion_host.test.name
  resource_group_name = azurerm_bastion_host.test.resource_group_name
}
`, BastionHostResource{}.basic(data))
}
