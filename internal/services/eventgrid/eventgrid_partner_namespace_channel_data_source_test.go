// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type EventGridPartnerNamespaceChannelDataSource struct{}

func TestAccEventGridPartnerNamespaceChannelDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventgrid_partner_namespace_channel", "test")
	r := EventGridPartnerNamespaceChannelDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("partner_namespace_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("channel_type").HasValue("PartnerTopic"),
				check.That(data.ResourceName).Key("expiration_time_if_not_activated_in_utc").Exists(),
				check.That(data.ResourceName).Key("readiness_state").Exists(),
			),
		},
	})
}

func TestAccEventGridPartnerNamespaceChannelDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventgrid_partner_namespace_channel", "test")
	r := EventGridPartnerNamespaceChannelDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("partner_namespace_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("channel_type").HasValue("PartnerTopic"),
				check.That(data.ResourceName).Key("expiration_time_if_not_activated_in_utc").Exists(),
				check.That(data.ResourceName).Key("partner_topic.#").HasValue("1"),
				check.That(data.ResourceName).Key("partner_topic.0.name").Exists(),
				check.That(data.ResourceName).Key("partner_topic.0.subscription_id").Exists(),
				check.That(data.ResourceName).Key("partner_topic.0.resource_group_name").Exists(),
				check.That(data.ResourceName).Key("partner_topic.0.source").HasValue("https://example.com/partner-topic"),
				check.That(data.ResourceName).Key("partner_topic.0.event_type_definitions.#").HasValue("1"),
				check.That(data.ResourceName).Key("partner_topic.0.event_type_definitions.0.kind").HasValue("Inline"),
				check.That(data.ResourceName).Key("partner_topic.0.event_type_definitions.0.inline_event_type.#").HasValue("1"),
				check.That(data.ResourceName).Key("partner_topic.0.event_type_definitions.0.inline_event_type.0.name").HasValue("SampleEventType"),
				check.That(data.ResourceName).Key("partner_topic.0.event_type_definitions.0.inline_event_type.0.display_name").HasValue("Sample Event Type"),
				check.That(data.ResourceName).Key("partner_topic.0.event_type_definitions.0.inline_event_type.0.data_schema_url").HasValue("https://example.com/sample-event-type-schema"),
				check.That(data.ResourceName).Key("partner_topic.0.event_type_definitions.0.inline_event_type.0.description").HasValue("This is a sample event type"),
				check.That(data.ResourceName).Key("partner_topic.0.event_type_definitions.0.inline_event_type.0.documentation_url").HasValue("https://example.com/sample-event-type-docs"),
				check.That(data.ResourceName).Key("readiness_state").Exists(),
			),
		},
	})
}

func (EventGridPartnerNamespaceChannelDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_eventgrid_partner_namespace_channel" "test" {
  name                   = azurerm_eventgrid_partner_namespace_channel.test.name
  partner_namespace_name = azurerm_eventgrid_partner_namespace_channel.test.partner_namespace_name
  resource_group_name    = azurerm_eventgrid_partner_namespace_channel.test.resource_group_name
}
`, EventGridPartnerNamespaceChannelTestResource{}.basic(data))
}

func (EventGridPartnerNamespaceChannelDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_eventgrid_partner_namespace_channel" "test" {
  name                   = azurerm_eventgrid_partner_namespace_channel.test.name
  partner_namespace_name = azurerm_eventgrid_partner_namespace_channel.test.partner_namespace_name
  resource_group_name    = azurerm_eventgrid_partner_namespace_channel.test.resource_group_name
}
`, EventGridPartnerNamespaceChannelTestResource{}.complete(data))
}
