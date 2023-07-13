// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SpringCloudServiceDataSource struct{}

func TestAccDataSourceSpringCloudService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_spring_cloud_service", "test")
	r := SpringCloudServiceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("outbound_public_ip_addresses.0").Exists(),
			),
		},
	})
}

func (SpringCloudServiceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_spring_cloud_service" "test" {
  name                = azurerm_spring_cloud_service.test.name
  resource_group_name = azurerm_spring_cloud_service.test.resource_group_name
}
`, SpringCloudServiceResource{}.basic(data))
}
