// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type EventGridDomainDataSource struct{}

func TestAccEventGridDomainDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventgrid_domain", "test")
	r := EventGridDomainDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("endpoint").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("input_mapping_fields.0.topic").Exists(),
				check.That(data.ResourceName).Key("input_mapping_fields.0.topic").Exists(),
				check.That(data.ResourceName).Key("input_mapping_default_values.0.data_version").Exists(),
				check.That(data.ResourceName).Key("input_mapping_default_values.0.subject").Exists(),
				check.That(data.ResourceName).Key("tags.%").Exists(),
				check.That(data.ResourceName).Key("tags.foo").Exists(),
				check.That(data.ResourceName).Key("inbound_ip_rule.#").Exists(),
				check.That(data.ResourceName).Key("inbound_ip_rule.0.ip_mask").Exists(),
				check.That(data.ResourceName).Key("inbound_ip_rule.1.ip_mask").Exists(),
				check.That(data.ResourceName).Key("inbound_ip_rule.0.action").Exists(),
				check.That(data.ResourceName).Key("inbound_ip_rule.1.action").Exists(),
				check.That(data.ResourceName).Key("identity.#").HasValue("0"),
			),
		},
	})
}

func TestAccEventGridDomainDataSource_systemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventgrid_domain", "test")
	r := EventGridDomainDataSource{}

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

func TestAccEventGridDomainDataSource_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventgrid_domain", "test")
	r := EventGridDomainDataSource{}

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

func (EventGridDomainDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_eventgrid_domain" "test" {
  name                = azurerm_eventgrid_domain.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, EventGridDomainResource{}.complete(data))
}

func (EventGridDomainDataSource) basicWithSystemManagedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_eventgrid_domain" "test" {
  name                = azurerm_eventgrid_domain.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, EventGridDomainResource{}.basicWithSystemManagedIdentity(data))
}

func (EventGridDomainDataSource) basicWithUserAssignedManagedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_eventgrid_domain" "test" {
  name                = azurerm_eventgrid_domain.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, EventGridDomainResource{}.basicWithUserAssignedManagedIdentity(data))
}
