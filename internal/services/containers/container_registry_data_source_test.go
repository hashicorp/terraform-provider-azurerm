// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ContainerRegistryDataSource struct{}

func TestAccDataSourceAzureRMContainerRegistry_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_registry", "test")
	r := ContainerRegistryDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("admin_enabled").Exists(),
				check.That(data.ResourceName).Key("login_server").Exists(),
			),
		},
	})
}

func TestAccDataSourceAzureRMContainerRegistry_dataEndpointPremium(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_registry", "test")
	r := ContainerRegistryDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.dataEndpointPremium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("data_endpoint_enabled").HasValue("true"),
			),
		},
	})
}

func (ContainerRegistryDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_registry" "test" {
  name                = azurerm_container_registry.test.name
  resource_group_name = azurerm_container_registry.test.resource_group_name
}
`, ContainerRegistryResource{}.basicManaged(data, "Basic"))
}

func (ContainerRegistryDataSource) dataEndpointPremium(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_registry" "test" {
  name                = azurerm_container_registry.test.name
  resource_group_name = azurerm_container_registry.test.resource_group_name
}
`, ContainerRegistryResource{}.dataEndpointPremium(data, true))
}
