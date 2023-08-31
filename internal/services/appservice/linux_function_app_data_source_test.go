// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type LinuxFunctionAppDataSource struct{}

func TestAccLinuxFunctionAppDataSource_standardComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_linux_function_app", "test")
	d := LinuxFunctionAppDataSource{}

	ipListRegex := regexp.MustCompile(`(([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})(,){0,1})+`)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.standardComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("outbound_ip_addresses").MatchesRegex(ipListRegex),
				check.That(data.ResourceName).Key("outbound_ip_address_list.#").Exists(),
				check.That(data.ResourceName).Key("possible_outbound_ip_addresses").MatchesRegex(ipListRegex),
				check.That(data.ResourceName).Key("possible_outbound_ip_address_list.#").Exists(),
				check.That(data.ResourceName).Key("default_hostname").HasValue(fmt.Sprintf("acctest-lfa-%d.azurewebsites.net", data.RandomInteger)),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
			),
		},
	})
}

func (LinuxFunctionAppDataSource) standardComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_linux_function_app" "test" {
  name                = azurerm_linux_function_app.test.name
  resource_group_name = azurerm_linux_function_app.test.resource_group_name
}
`, LinuxFunctionAppResource{}.standardComplete(data))
}
