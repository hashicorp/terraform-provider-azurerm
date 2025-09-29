// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid_test

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

type EventGridPartnerNamespaceChannelTestResource struct{}

func TestAccEventGridPartnerNamespaceChannel_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_namespace_channel", "test")
	r := EventGridPartnerNamespaceChannelTestResource{}

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

func TestAccEventGridPartnerNamespaceChannel_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_namespace_channel", "test")
	r := EventGridPartnerNamespaceChannelTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_eventgrid_partner_namespace_channel"),
		},
	})
}

func TestAccEventGridPartnerNamespaceChannel_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_namespace_channel", "test")
	r := EventGridPartnerNamespaceChannelTestResource{}

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

func TestAccEventGridPartnerNamespaceChannel_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_namespace_channel", "test")
	r := EventGridPartnerNamespaceChannelTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (EventGridPartnerNamespaceChannelTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := channels.ParseChannelID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.EventGrid.Channels.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r EventGridPartnerNamespaceChannelTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_eventgrid_partner_namespace_channel" "test" {
  name                 = "acctest-egpnc-%[2]d"
  partner_namespace_id = azurerm_eventgrid_partner_namespace.test.id
  partner_topic {
    subscription_id     = "%[3]s"
    resource_group_name = azurerm_resource_group.test.name
    name                = "acctest-egpt-%[2]d"
    source              = "https://example.com/partner-topic"
  }
}
`, r.template(data), data.RandomInteger, data.Subscriptions.Primary)
}

func (r EventGridPartnerNamespaceChannelTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_eventgrid_partner_namespace_channel" "import" {
  name                                    = azurerm_eventgrid_partner_namespace_channel.test.name
  partner_namespace_id                    = azurerm_eventgrid_partner_namespace_channel.test.partner_namespace_id
  expiration_time_if_not_activated_in_utc = azurerm_eventgrid_partner_namespace_channel.test.expiration_time_if_not_activated_in_utc
  partner_topic {
    subscription_id     = azurerm_eventgrid_partner_namespace_channel.test.partner_topic[0].subscription_id
    resource_group_name = azurerm_eventgrid_partner_namespace_channel.test.partner_topic[0].resource_group_name
    name                = azurerm_eventgrid_partner_namespace_channel.test.partner_topic[0].name
    source              = azurerm_eventgrid_partner_namespace_channel.test.partner_topic[0].source
  }
}
`, r.basic(data))
}

func (r EventGridPartnerNamespaceChannelTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_eventgrid_partner_namespace_channel" "test" {
  name                                    = "acctest-egpnc-%[2]d"
  partner_namespace_id                    = azurerm_eventgrid_partner_namespace.test.id
  channel_type                            = "PartnerTopic"
  expiration_time_if_not_activated_in_utc = azurerm_eventgrid_partner_configuration.test.partner_authorization.0.authorization_expiration_time_in_utc
  partner_topic {
    subscription_id     = "%[3]s"
    resource_group_name = azurerm_resource_group.test.name
    name                = "acctest-egpt-%[2]d"
    source              = "https://example.com/partner-topic"
    event_type_definitions {
      inline_event_type {
        name              = "SampleEventType"
        display_name      = "Sample Event Type"
        description       = "This is a sample event type"
        data_schema_url   = "https://example.com/sample-event-type-schema"
        documentation_url = "https://example.com/sample-event-type-docs"
      }
      kind = "Inline"
    }
  }
}`, r.template(data), data.RandomInteger, data.Subscriptions.Primary)
}

func (r EventGridPartnerNamespaceChannelTestResource) update(data acceptance.TestData) string {
	expiryTime := time.Now().In(time.UTC).Add(7 * 24 * time.Hour).Format(time.RFC3339)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_eventgrid_partner_namespace_channel" "test" {
  name                                    = "acctest-egpnc-%[2]d"
  partner_namespace_id                    = azurerm_eventgrid_partner_namespace.test.id
  channel_type                            = "PartnerTopic"
  expiration_time_if_not_activated_in_utc = "%[4]s"
  partner_topic {
    subscription_id     = "%[3]s"
    resource_group_name = azurerm_resource_group.test.name
    name                = "acctest-egpt-%[2]d"
    source              = "https://example.com/partner-topic"
    event_type_definitions {
      inline_event_type {
        name              = "SampleEventType"
        display_name      = "Sample Event Type"
        description       = "This is a sample event type update"
        data_schema_url   = "https://example.com/sample-event-type-schema"
        documentation_url = "https://example.com/sample-event-type-docs"
      }
      inline_event_type {
        name              = "SampleEventType2"
        display_name      = "Sample Event Type 2"
        description       = "This is a sample event type 2"
        data_schema_url   = "https://example.com/sample-event-type-schema"
        documentation_url = "https://example.com/sample-event-type-docs"
      }
      kind = "Inline"
    }
  }
}`, r.template(data), data.RandomInteger, data.Subscriptions.Primary, expiryTime)
}

func (r EventGridPartnerNamespaceChannelTestResource) template(data acceptance.TestData) interface{} {
	expiryTime := time.Now().In(time.UTC).Add(3 * 24 * time.Hour).Format(time.RFC3339)
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventgrid_partner_registration" "test" {
  name                = "acctest-egpr-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_eventgrid_partner_namespace" "test" {
  name                    = "acctest-egpn-%[1]d"
  location                = "%[2]s"
  resource_group_name     = azurerm_resource_group.test.name
  partner_registration_id = azurerm_eventgrid_partner_registration.test.id
}

resource "azurerm_eventgrid_partner_configuration" "test" {
  resource_group_name                     = azurerm_resource_group.test.name
  default_maximum_expiration_time_in_days = 180

  partner_authorization {
    partner_registration_id              = azurerm_eventgrid_partner_registration.test.partner_registration_id
    partner_name                         = ""
    authorization_expiration_time_in_utc = "%[3]s"
  }
}
`, data.RandomInteger, data.Locations.Primary, expiryTime)
}
