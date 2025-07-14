// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package paloalto_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type LocalRulestackDataSource struct{}

func TestAccLocalRulestackDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack", "test")

	d := LocalRulestackDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("description").HasValue(fmt.Sprintf("Acceptance Test Desc - %d", data.RandomInteger)),
			),
		},
	})
}

func (d LocalRulestackDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_palo_alto_local_rulestack" "test" {
  name                = azurerm_palo_alto_local_rulestack.test.name
  resource_group_name = azurerm_palo_alto_local_rulestack.test.resource_group_name
}

`, LocalRulestackResource{}.complete(data))
}
