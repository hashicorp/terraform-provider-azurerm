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

func TestAccDataSourceMachineLearningRegistry_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_machine_learning_registry", "test")
	r := MachineLearningRegistryDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("main_region.#").HasValue("1"),
				check.That(data.ResourceName).Key("main_region.0.location").Exists(),
				check.That(data.ResourceName).Key("main_region.0.storage_account_type").HasValue("Standard_LRS"),
				check.That(data.ResourceName).Key("main_region.0.hns_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("main_region.0.system_created_storage_account_id").Exists(),
				check.That(data.ResourceName).Key("main_region.0.system_created_container_registry_id").Exists(),
				check.That(data.ResourceName).Key("replication_region.#").HasValue("0"),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("discovery_url").Exists(),
				check.That(data.ResourceName).Key("ml_flow_registry_uri").Exists(),
				check.That(data.ResourceName).Key("managed_resource_group").Exists(),
			),
		},
	})
}

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
				check.That(data.ResourceName).Key("main_region.#").HasValue("1"),
				check.That(data.ResourceName).Key("replication_region.#").HasValue("2"),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("test"),
			),
		},
	})
}

func (r MachineLearningRegistryDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_machine_learning_registry" "test" {
  name                = azurerm_machine_learning_registry.test.name
  resource_group_name = azurerm_machine_learning_registry.test.resource_group_name
}
`, MachineLearningRegistry{}.basic(data))
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
