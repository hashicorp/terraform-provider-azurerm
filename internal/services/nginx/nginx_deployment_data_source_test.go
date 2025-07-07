// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nginx_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type NginxDeploymentDataSource struct{}

func TestAccNginxDeploymentDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_nginx_deployment", "test")
	r := NginxDeploymentDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("nginx_version").Exists(),
				check.That(data.ResourceName).Key("sku").Exists(),
				check.That(data.ResourceName).Key("capacity").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("automatic_upgrade_channel").Exists(),
				check.That(data.ResourceName).Key("dataplane_api_endpoint").Exists(),
			),
		},
	})
}

func (d NginxDeploymentDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_nginx_deployment" "test" {
  name                = azurerm_nginx_deployment.test.name
  resource_group_name = azurerm_nginx_deployment.test.resource_group_name
}
`, DeploymentResource{}.basic(data))
}

func (d NginxDeploymentDataSource) basicAutoscaling(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_nginx_deployment" "test" {
  name                = azurerm_nginx_deployment.test.name
  resource_group_name = azurerm_nginx_deployment.test.resource_group_name
}
`, DeploymentResource{}.basicAutoscaling(data))
}

func TestAccNginxDeploymentDataSource_autoscaling(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_nginx_deployment", "test")
	r := NginxDeploymentDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicAutoscaling(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("auto_scale_profile.0.name").HasValue("test"),
				check.That(data.ResourceName).Key("auto_scale_profile.0.min_capacity").HasValue("10"),
				check.That(data.ResourceName).Key("auto_scale_profile.0.max_capacity").HasValue("30"),
			),
		},
	})
}

func (d NginxDeploymentDataSource) basicNginxAppProtect(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
data "azurerm_nginx_deployment" "test" {
  name                = azurerm_nginx_deployment.test.name
  resource_group_name = azurerm_nginx_deployment.test.resource_group_name
}
`, DeploymentResource{}.basicNginxAppProtect(data))
}

func TestAccNginxDeploymentDataSource_nginxappprotect(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_nginx_deployment", "test")
	r := NginxDeploymentDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicNginxAppProtect(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("web_application_firewall.0.activation_state_enabled").HasValue("true"),
			),
		},
	})
}
