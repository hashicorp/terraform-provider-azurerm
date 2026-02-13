// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MachineLearningRegistryDataSource struct{}

func TestAccDataSourceMachineLearningRegistry_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_machine_learning_registry", "test")
	r := MachineLearningRegistryDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("primary_region.#").HasValue("1"),
				check.That(data.ResourceName).Key("replication_region.#").HasValue("2"),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
			),
		},
	})
}

func (r MachineLearningRegistryDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_machine_learning_registry" "test" {
  name                = azurerm_machine_learning_registry.test.name
  resource_group_name = azurerm_machine_learning_registry.test.resource_group_name
}
`, MachineLearningRegistry{}.complete(data))
}
