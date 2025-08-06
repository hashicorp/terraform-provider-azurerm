// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ContainerRegistryCacheRuleDataSource struct{}

func TestAccDataSourceAzureRMContainerRegistryCacheRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_registry_cache_rule", "test")
	r := ContainerRegistryCacheRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("container_registry_id").Exists(),
				check.That(data.ResourceName).Key("source_repo").Exists(),
				check.That(data.ResourceName).Key("target_repo").Exists(),
			),
		},
	})
}

func (ContainerRegistryCacheRuleDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_registry_cache_rule" "test" {
  name                  = azurerm_container_registry_cache_rule.test.name
  container_registry_id = azurerm_container_registry_cache_rule.test.container_registry_id
}
`, ContainerRegistryCacheRuleResource{}.basic(data))
}
