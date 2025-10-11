// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
)

type ResourceAnchorDataSource struct{}

func TestResourceAnchorDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ResourceAnchorDataSource{}.ResourceType(), "test")
	r := ResourceAnchorDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("linked_compartment_id").Exists(),
				check.That(data.ResourceName).Key("provisioning_state").Exists(),
			),
		},
	})
}

func (d ResourceAnchorDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_database_resource_anchor" "test" {
  name                = azurerm_oracle_database_resource_anchor.test.name
  resource_group_name = azurerm_oracle_database_resource_anchor.test.resource_group_name
}
`, ResourceAnchorResource{}.basic(data))
}
