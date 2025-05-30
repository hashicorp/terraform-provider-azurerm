// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signalr_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type WebPubSubSocketIOTestResource struct{}

func TestAccWebPubSubSocketIO_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_socketio", "test")
	r := WebPubSubSocketIOTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "Standard_S1", 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWebPubSubSocketIO_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_socketio", "test")
	r := WebPubSubSocketIOTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "Standard_S1", 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccWebPubSubSocketIO_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_socketio", "test")
	r := WebPubSubSocketIOTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "Standard_S1", 1),
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
			Config: r.basic(data, "Standard_S1", 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWebPubSubSocketIO_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_socketio", "test")
	r := WebPubSubSocketIOTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userAssignedIdentity(data),
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
			Config: r.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWebPubSubSocketIO_skus(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_socketio", "test")
	r := WebPubSubSocketIOTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "Free_F1", 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "Standard_S1", 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "Premium_P1", 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "Premium_P2", 100),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (WebPubSubSocketIOTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := webpubsub.ParseWebPubSubID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.SignalR.WebPubSubClient.WebPubSub.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r WebPubSubSocketIOTestResource) basic(data acceptance.TestData, sku string, capacity int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_web_pubsub_socketio" "test" {
  name                = "acctestWebPubsubSocketIO-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "%s"
    capacity = %d
  }
}
`, r.template(data), data.RandomInteger, sku, capacity)
}

func (r WebPubSubSocketIOTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_web_pubsub_socketio" "import" {
  name                = azurerm_web_pubsub_socketio.test.name
  location            = azurerm_web_pubsub_socketio.test.location
  resource_group_name = azurerm_web_pubsub_socketio.test.resource_group_name

  sku {
    name     = azurerm_web_pubsub_socketio.test.sku.0.name
    capacity = azurerm_web_pubsub_socketio.test.sku.0.capacity
  }
}`, r.basic(data, "Standard_S1", 1))
}

func (r WebPubSubSocketIOTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_web_pubsub_socketio" "test" {
  name                = "acctestWebPubsubSocketIO-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Standard_S1"
    capacity = 2
  }

  aad_auth_enabled = true

  identity {
    type = "SystemAssigned"
  }

  live_trace_enabled                   = true
  live_trace_connectivity_logs_enabled = false
  live_trace_http_request_logs_enabled = true
  live_trace_messaging_logs_enabled    = false

  local_auth_enabled      = false
  public_network_access   = "Disabled"
  service_mode            = "Serverless"
  tls_client_cert_enabled = true

  tags = {
    Hello = "World"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r WebPubSubSocketIOTestResource) userAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_web_pubsub_socketio" "test" {
  name                = "acctestWebPubsubSocketIO-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name = "Standard_S1"
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.template(data), data.RandomInteger)
}

func (WebPubSubSocketIOTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
