// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signalr_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2023-02-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WebPubsubHubResource struct{}

func TestAccWebPubsubHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_hub", "test")
	r := WebPubsubHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccWebPubsubHub_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_hub", "test")
	r := WebPubsubHubResource{}

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

func TestAccWebPubsubHub_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_hub", "test")
	r := WebPubsubHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccWebPubsubHub_usingAuthGuid(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_hub", "test")
	r := WebPubsubHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.usingAuthGuid(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccWebPubsubHub_withAuthUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_hub", "test")
	r := WebPubsubHubResource{}

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
	})
}

func TestAccWebPubsubHub_withPropertyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_hub", "test")
	r := WebPubsubHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withPropertyUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWebPubsubHub_withMultipleEventHandlerSettingsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_hub", "test")
	r := WebPubsubHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withMultipleEventhandlerSettingsAndNoAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
		{
			Config: r.withMultipleEventHandlerSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccWebPubsubHub_withMultipleEventListenerSettingsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_hub", "test")
	r := WebPubsubHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withMultipleEventListenerSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
		{
			Config: r.withMultipleEventListenerSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func (r WebPubsubHubResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := webpubsub.ParseHubID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SignalR.WebPubSubClient.WebPubSub.HubsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r WebPubsubHubResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_web_pubsub_hub" "test" {
  name          = "acctestwpsh%d"
  web_pubsub_id = azurerm_web_pubsub.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r WebPubsubHubResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_web_pubsub_hub" "test" {
  name          = "acctestwpsh%d"
  web_pubsub_id = azurerm_web_pubsub.test.id
  event_handler {
    url_template       = "https://test.com/api/{hub}/{event}"
    user_event_pattern = "*"
    system_events      = ["connect", "connected"]

    auth {
      managed_identity_id = azurerm_user_assigned_identity.test.id
    }
  }
  anonymous_connections_enabled = true

  depends_on = [
    azurerm_web_pubsub.test
  ]
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r WebPubsubHubResource) usingAuthGuid(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_web_pubsub_hub" "test" {
  name          = "acctestwpsh%d"
  web_pubsub_id = azurerm_web_pubsub.test.id
  event_handler {
    url_template       = "https://test.com/api/{hub}/{event}"
    user_event_pattern = "*"
    system_events      = ["connect", "connected"]

    auth {
      managed_identity_id = "12345678-9012-3456-7890-123456789012"
    }
  }
  anonymous_connections_enabled = true

  depends_on = [
    azurerm_web_pubsub.test
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r WebPubsubHubResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_web_pubsub_hub" "import" {
  name          = azurerm_web_pubsub_hub.test.name
  web_pubsub_id = azurerm_web_pubsub.test.id

  event_handler {
    url_template       = "https://test.com/api/{hub}/{event}"
    user_event_pattern = "*"
    system_events      = ["connect", "connected"]
  }
}
`, r.basic(data))
}

func (r WebPubsubHubResource) withMultipleEventHandlerSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test1" {
  name                = "acctest-uai1-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest-uai2-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_web_pubsub_hub" "test" {
  name          = "acctestwpsh%d"
  web_pubsub_id = azurerm_web_pubsub.test.id
  event_handler {
    url_template       = "https://test.com/api/{hub1}/{event2}"
    user_event_pattern = "*"
    system_events      = ["connect", "connected"]
    auth {
      managed_identity_id = azurerm_user_assigned_identity.test1.id
    }
  }
  event_handler {
    url_template       = "https://test.com/api/{hub2}/{event1}"
    user_event_pattern = "event1, event2"
    system_events      = ["connected"]
    auth {
      managed_identity_id = azurerm_user_assigned_identity.test2.id
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r WebPubsubHubResource) withMultipleEventListenerSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 1
  message_retention   = 1
}

resource "azurerm_eventhub" "test1" {
  name                = "acctesteventhub-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 1
  message_retention   = 1
}

resource "azurerm_web_pubsub_hub" "test" {
  name          = "acctestwpsh%d"
  web_pubsub_id = azurerm_web_pubsub.test.id

  event_listener {
    system_event_name_filter = ["disconnected", "connected"]
    user_event_name_filter   = ["event1"]
    eventhub_namespace_name  = azurerm_eventhub_namespace.test.name
    eventhub_name            = azurerm_eventhub.test.name
  }

  event_listener {
    system_event_name_filter = ["connected"]
    user_event_name_filter   = ["event1", "event2"]
    eventhub_namespace_name  = azurerm_eventhub_namespace.test.name
    eventhub_name            = azurerm_eventhub.test1.name
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r WebPubsubHubResource) withMultipleEventhandlerSettingsAndNoAuth(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test1" {
  name                = "acctest-uai1-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_web_pubsub_hub" "test" {
  name          = "acctestwpsh%d"
  web_pubsub_id = azurerm_web_pubsub.test.id
  event_handler {
    url_template       = "https://test.com/api/{hub1}/{event2}"
    user_event_pattern = "*"
    system_events      = ["connect", "connected"]
  }
  event_handler {
    url_template       = "https://test.com/api/{hub2}/{event1}"
    user_event_pattern = "event1, event2"
    system_events      = ["connected"]
    auth {
      managed_identity_id = azurerm_user_assigned_identity.test1.id
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r WebPubsubHubResource) withPropertyUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test1" {
  name                = "acctest-uai1-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_web_pubsub_hub" "test" {
  name          = "acctestwpsh%d"
  web_pubsub_id = azurerm_web_pubsub.test.id
  event_handler {
    url_template       = "https://test.com/api/{testhub}/{testevent1}"
    user_event_pattern = "event1, event2"
    system_events      = ["disconnected", "connect", "connected"]
    auth {
      managed_identity_id = azurerm_user_assigned_identity.test1.id
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r WebPubsubHubResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-wps-%d"
  location = "%s"
}

resource "azurerm_web_pubsub" "test" {
  name                = "acctest-webpubsub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_S1"

  identity {
    type = "SystemAssigned"
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
