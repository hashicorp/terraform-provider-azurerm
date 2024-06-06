// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nginx_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type NginxConfigurationDataSource struct{}

func TestAccNginxConfigurationDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_nginx_configuration", "test")
	r := NginxConfigurationDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("root_file").Exists(),
				check.That(data.ResourceName).Key("config_file.0.content").Exists(),
			),
		},
	})
}

func (d NginxConfigurationDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_nginx_configuration" "test" {
  nginx_deployment_id = azurerm_nginx_deployment.test.id

  depends_on = [azurerm_nginx_configuration.test]
}
`, ConfigurationResource{}.basic(data))
}
