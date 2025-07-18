// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connections_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ApiConnectionDataSource struct{}

func TestAccApiConnectionDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_connection", "test")
	r := ApiConnectionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("managed_api_id").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Environment").HasValue("Test"),
			),
		},
	})
}

// Note: Data sources don't need Exists functions, so helper functions can be placed directly after test functions
func (ApiConnectionDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_connection" "test" {
  name                = "acctestconn-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  managed_api_id      = data.azurerm_managed_api.test.id
  display_name        = "Test Connection"

  tags = {
    Environment = "Test"
  }
}

data "azurerm_api_connection" "test" {
  name                = azurerm_api_connection.test.name
  resource_group_name = azurerm_api_connection.test.resource_group_name
}
`, ApiConnectionTestResource{}.template(data), data.RandomInteger)
}
