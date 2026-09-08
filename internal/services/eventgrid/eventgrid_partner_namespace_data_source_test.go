// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type EventGridPartnerNamespaceDataSource struct{}

func TestAccEventGridPartnerNamespaceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventgrid_partner_namespace", "test")
	r := EventGridPartnerNamespaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("endpoint").Exists(),
				check.That(data.ResourceName).Key("local_authentication_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("partner_registration_id").Exists(),
				check.That(data.ResourceName).Key("partner_topic_routing_mode").HasValue("ChannelNameHeader"),
				check.That(data.ResourceName).Key("public_network_access").HasValue("Enabled"),
			),
		},
	})
}

func TestAccEventGridPartnerNamespaceDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventgrid_partner_namespace", "test")
	r := EventGridPartnerNamespaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("endpoint").Exists(),
				check.That(data.ResourceName).Key("inbound_ip_rule.0.ip_mask").HasValue("10.0.0.0/16"),
				check.That(data.ResourceName).Key("inbound_ip_rule.0.action").HasValue("Allow"),
				check.That(data.ResourceName).Key("inbound_ip_rule.1.ip_mask").HasValue("10.1.0.0/16"),
				check.That(data.ResourceName).Key("inbound_ip_rule.1.action").HasValue("Allow"),
				check.That(data.ResourceName).Key("local_authentication_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("partner_registration_id").Exists(),
				check.That(data.ResourceName).Key("partner_topic_routing_mode").HasValue("ChannelNameHeader"),
				check.That(data.ResourceName).Key("public_network_access").HasValue("Enabled"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.foo").HasValue("bar"),
			),
		},
	})
}

func (EventGridPartnerNamespaceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_eventgrid_partner_namespace" "test" {
  name                = azurerm_eventgrid_partner_namespace.test.name
  resource_group_name = azurerm_eventgrid_partner_namespace.test.resource_group_name
}
`, EventGridPartnerNamespaceTestResource{}.basic(data))
}

func (EventGridPartnerNamespaceDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_eventgrid_partner_namespace" "test" {
  name                = azurerm_eventgrid_partner_namespace.test.name
  resource_group_name = azurerm_eventgrid_partner_namespace.test.resource_group_name
}
`, EventGridPartnerNamespaceTestResource{}.complete(data))
}
