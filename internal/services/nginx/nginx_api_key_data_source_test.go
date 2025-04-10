// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nginx_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type APIKeyDataSource struct{}

func TestAccAPIKeyDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_nginx_api_key", "test")
	r := APIKeyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("end_date_time").IsNotEmpty(),
				check.That(data.ResourceName).Key("hint").IsNotEmpty(),
			),
		},
	})
}

func (APIKeyDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_nginx_api_key" "test" {
  name                = azurerm_nginx_api_key.test.name
  nginx_deployment_id = azurerm_nginx_api_key.test.nginx_deployment_id
}
`, APIKeyResource{}.complete(data))
}
