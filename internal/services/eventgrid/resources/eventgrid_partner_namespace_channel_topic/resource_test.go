// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_namespace_channel_topic_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/channels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type EventgridPartnerNamespaceChannelTopicResource struct{}

func TestAccEventGridPartnerNamespaceChannelTopic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_namespace_channel_topic", "test")
	r := EventgridPartnerNamespaceChannelTopicResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridPartnerNamespaceChannelTopic_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_namespace_channel_topic", "test")
	r := EventgridPartnerNamespaceChannelTopicResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccEventGridPartnerNamespaceChannelTopic_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_namespace_channel_topic", "test")
	r := EventgridPartnerNamespaceChannelTopicResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridPartnerNamespaceChannelTopic_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_namespace_channel_topic", "test")
	r := EventgridPartnerNamespaceChannelTopicResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (EventgridPartnerNamespaceChannelTopicResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := channels.ParseChannelID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.EventGrid.Channels.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil && resp.Model.Properties != nil && pointer.From(resp.Model.Properties.ChannelType) == channels.ChannelTypePartnerTopic), nil
}

func (EventgridPartnerNamespaceChannelTopicResource) template(data acceptance.TestData) string {
	expiryTime := time.Now().UTC().Add(7 * 24 * time.Hour).Format(time.RFC3339)

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "partner" {
  name     = "acctestRG-partner-%[1]d"
  location = "%[2]s"
}

resource "azurerm_resource_group" "subscriber" {
  name     = "acctestRG-subscriber-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventgrid_partner_registration" "test" {
  name                = "acctest-egpr-%[1]d"
  resource_group_name = azurerm_resource_group.partner.name
}

resource "azurerm_eventgrid_partner_namespace" "test" {
  name                       = "acctest-egpn-%[1]d"
  resource_group_name        = azurerm_resource_group.partner.name
  location                   = azurerm_resource_group.partner.location
  partner_registration_id    = azurerm_eventgrid_partner_registration.test.id
  partner_topic_routing_mode = "ChannelNameHeader"
}

resource "azurerm_eventgrid_partner_configuration" "test" {
  resource_group_name = azurerm_resource_group.subscriber.name

  partner_authorization {
    partner_registration_id              = azurerm_eventgrid_partner_registration.test.partner_registration_id
    authorization_expiration_time_in_utc = "%[3]s"
  }
}
`, data.RandomInteger, data.Locations.Primary, expiryTime)
}

func (r EventgridPartnerNamespaceChannelTopicResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_eventgrid_partner_namespace_channel_topic" "test" {
  name                 = "acctest-channel-%d"
  partner_namespace_id = azurerm_eventgrid_partner_namespace.test.id

  partner_topic {
    name                = "acctest-topic-%d"
    resource_group_name = azurerm_resource_group.subscriber.name
    subscription_id     = "%s"
    source              = "Contoso.Accounts.User1"
  }

  depends_on = [azurerm_eventgrid_partner_configuration.test]
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.Subscriptions.Primary)
}

func (r EventgridPartnerNamespaceChannelTopicResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_eventgrid_partner_namespace_channel_topic" "test" {
  name                 = "acctest-channel-%d"
  partner_namespace_id = azurerm_eventgrid_partner_namespace.test.id

  partner_topic {
    name                = "acctest-topic-%d"
    resource_group_name = azurerm_resource_group.subscriber.name
    subscription_id     = "%s"
    source              = "Contoso.Accounts.User1"
  }

  event_type {
    name              = "Contoso.OrderCreated"
    display_name      = "Order created"
    description       = "An order was created."
    documentation_url = "https://example.com/events/order-created"
    data_schema_url   = "https://example.com/schemas/order-created.json"
  }

  depends_on = [azurerm_eventgrid_partner_configuration.test]
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.Subscriptions.Primary)
}

func (r EventgridPartnerNamespaceChannelTopicResource) completeUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_eventgrid_partner_namespace_channel_topic" "test" {
  name                 = "acctest-channel-%d"
  partner_namespace_id = azurerm_eventgrid_partner_namespace.test.id

  partner_topic {
    name                = "acctest-topic-%d"
    resource_group_name = azurerm_resource_group.subscriber.name
    subscription_id     = "%s"
    source              = "Contoso.Accounts.User1"
  }

  event_type {
    name              = "Contoso.OrderCreated"
    display_name      = "Order created"
    description       = "An order was created."
    documentation_url = "https://example.com/events/order-created"
    data_schema_url   = "https://example.com/schemas/order-created.json"
  }

  event_type {
    name              = "Contoso.OrderUpdated"
    display_name      = "Order updated"
    description       = "An order was updated."
    documentation_url = "https://example.com/events/order-updated"
    data_schema_url   = "https://example.com/schemas/order-updated.json"
  }

  depends_on = [azurerm_eventgrid_partner_configuration.test]
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.Subscriptions.Primary)
}

func (r EventgridPartnerNamespaceChannelTopicResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_eventgrid_partner_namespace_channel_topic" "import" {
  name                 = azurerm_eventgrid_partner_namespace_channel_topic.test.name
  partner_namespace_id = azurerm_eventgrid_partner_namespace_channel_topic.test.partner_namespace_id

  partner_topic {
    name                = azurerm_eventgrid_partner_namespace_channel_topic.test.partner_topic[0].name
    resource_group_name = azurerm_eventgrid_partner_namespace_channel_topic.test.partner_topic[0].resource_group_name
    subscription_id     = azurerm_eventgrid_partner_namespace_channel_topic.test.partner_topic[0].subscription_id
    source              = azurerm_eventgrid_partner_namespace_channel_topic.test.partner_topic[0].source
  }
}
`, r.basic(data))
}
