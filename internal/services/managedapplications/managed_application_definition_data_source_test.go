// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedapplications_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ManagedApplicationDefinitionDataSource struct{}

func TestAccManagedApplicationDefinitionDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_application_definition", "test")
	r := ManagedApplicationDefinitionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func (ManagedApplicationDefinitionDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_managed_application_definition" "test" {
  name                = azurerm_managed_application_definition.test.name
  resource_group_name = azurerm_managed_application_definition.test.resource_group_name
}
`, ManagedApplicationDefinitionResource{}.basic(data))
}
