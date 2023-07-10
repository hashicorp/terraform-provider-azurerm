// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ApiManagementApiVersionSetDataSource struct{}

func TestAccDataSourceApiManagementApiVersionSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management_api_version_set", "test")
	r := ApiManagementApiVersionSetDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("api_management_name").Exists(),
			),
		},
	})
}

func (ApiManagementApiVersionSetDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_api_management_api_version_set" "test" {
  name                = azurerm_api_management_api_version_set.test.name
  resource_group_name = azurerm_api_management_api_version_set.test.resource_group_name
  api_management_name = azurerm_api_management_api_version_set.test.api_management_name
}
`, ApiManagementApiVersionSetResource{}.basic(data))
}
