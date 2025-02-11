// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type LogicAppStandardDataSource struct{}

func TestAccLogicAppStandardDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_logic_app_standard", "test")
	r := LogicAppStandardDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,workflowapp"),
				check.That(data.ResourceName).Key("version").HasValue("~4"),
				check.That(data.ResourceName).Key("outbound_ip_addresses").Exists(),
				check.That(data.ResourceName).Key("possible_outbound_ip_addresses").Exists(),
				check.That(data.ResourceName).Key("custom_domain_verification_id").Exists(),
			),
		},
	})
}

func (r LogicAppStandardDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_logic_app_standard" "test" {
  name                = azurerm_logic_app_standard.test.name
  resource_group_name = azurerm_logic_app_standard.test.resource_group_name
}
`, LogicAppStandardResource{}.basic(data))
}
