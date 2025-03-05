// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nginx_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DataplaneAPIKeyDataSource struct{}

func TestAccDataplaneAPIKeyDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_nginx_dataplane_apikey", "test")
	r := DataplaneAPIKeyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("end_date_time").IsNotEmpty(),
				check.That(data.ResourceName).Key("hint").IsNotEmpty(),
			),
		},
	})
}

func (DataplaneAPIKeyDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_nginx_dataplane_apikey" "test" {
  name                = azurerm_nginx_dataplane_apikey.test.name
  nginx_deployment_id = azurerm_nginx_dataplane_apikey.test.nginx_deployment_id
}
`, DataplaneAPIKeyResource{}.basic(data))
}
