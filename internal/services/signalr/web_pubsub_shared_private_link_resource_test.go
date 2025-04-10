// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signalr_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WebPubsubSharedPrivateLinkResource struct{}

func TestAccWebPubsubSharedPrivateLinkResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_shared_private_link_resource", "test")
	r := WebPubsubSharedPrivateLinkResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccWebPubsubSharedPrivateLinkResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_shared_private_link_resource", "test")
	r := WebPubsubSharedPrivateLinkResource{}

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

func (r WebPubsubSharedPrivateLinkResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := webpubsub.ParseSharedPrivateLinkResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SignalR.WebPubSubClient.WebPubSub.SharedPrivateLinkResourcesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (r WebPubsubSharedPrivateLinkResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault" "test" {
  name                       = "vault%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "ManageContacts",
    ]

    key_permissions = [
      "Create",
    ]

    secret_permissions = [
      "Set",
    ]
  }
}

resource "azurerm_web_pubsub_shared_private_link_resource" "test" {
  name               = "acctest-%d"
  web_pubsub_id      = azurerm_web_pubsub.test.id
  subresource_name   = "vault"
  target_resource_id = azurerm_key_vault.test.id
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r WebPubsubSharedPrivateLinkResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_web_pubsub_shared_private_link_resource" "import" {
  name               = azurerm_web_pubsub_shared_private_link_resource.test.name
  web_pubsub_id      = azurerm_web_pubsub_shared_private_link_resource.test.web_pubsub_id
  subresource_name   = azurerm_web_pubsub_shared_private_link_resource.test.subresource_name
  target_resource_id = azurerm_web_pubsub_shared_private_link_resource.test.target_resource_id
}
`, config)
}

func (r WebPubsubSharedPrivateLinkResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-wps-%d"
  location = "%s"
}

resource "azurerm_web_pubsub" "test" {
  name                = "acctestWebPubsub-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_S1"
  capacity            = 1
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
