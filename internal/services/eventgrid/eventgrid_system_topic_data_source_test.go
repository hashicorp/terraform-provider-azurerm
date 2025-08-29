// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type EventGridSystemTopicDataSource struct{}

func TestAccEventGridSystemTopicDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventgrid_system_topic", "test")
	r := EventGridSystemTopicDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Foo").HasValue("Bar"),
				check.That(data.ResourceName).Key("source_arm_resource_id").Exists(),
				check.That(data.ResourceName).Key("topic_type").Exists(),
				check.That(data.ResourceName).Key("metric_arm_resource_id").Exists(),
			),
		},
	})
}

func TestAccEventGridSystemTopicDataSource_basicWithSystemManagedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventgrid_system_topic", "test")
	r := EventGridSystemTopicDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicWithSystemManagedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("0"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
	})
}

func TestAccEventGridSystemTopicDataSource_basicWithUserAssignedManagedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventgrid_system_topic", "test")
	r := EventGridSystemTopicDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicWithUserAssignedManagedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsEmpty(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsEmpty(),
			),
		},
	})
}

func (EventGridSystemTopicDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_eventgrid_system_topic" "test" {
  name                = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, EventGridSystemTopicResource{}.complete(data))
}

func (EventGridSystemTopicDataSource) basicWithSystemManagedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_eventgrid_system_topic" "test" {
  name                = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, EventGridSystemTopicResource{}.basicWithSystemManagedIdentity(data))
}

func (EventGridSystemTopicDataSource) basicWithUserAssignedManagedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_eventgrid_system_topic" "test" {
  name                = azurerm_eventgrid_system_topic.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, EventGridSystemTopicResource{}.basicWithUserAssignedManagedIdentity(data))
}
