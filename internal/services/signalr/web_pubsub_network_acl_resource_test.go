// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signalr_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2023-02-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WebPubsubNetworkACLResource struct{}

func TestAccWebPubsubNetworkACL_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_network_acl", "test")
	r := WebPubsubNetworkACLResource{}

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

func TestAccWebPubsubNetworkACL_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_network_acl", "test")
	r := WebPubsubNetworkACLResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccWebPubsubNetworkACL_complete_withoutPrivateEndpoint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_network_acl", "test")
	r := WebPubsubNetworkACLResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete_withoutPrivateEndpoint(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccWebPubsubNetworkACL_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_network_acl", "test")
	r := WebPubsubNetworkACLResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccWebPubsubNetworkACL_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_network_acl", "test")
	r := WebPubsubNetworkACLResource{}

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

func TestAccWebPubsubNetworkACL_updateMultiplePrivateEndpoints(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_network_acl", "test")
	r := WebPubsubNetworkACLResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
		{
			Config: r.multiplePrivateEndpoints(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func (r WebPubsubNetworkACLResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := webpubsub.ParseWebPubSubID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SignalR.WebPubSubClient.WebPubSub.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	isDefaultConfiguration := false
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if acls := props.NetworkACLs; acls != nil {
				hasDefaultAction := false
				if acls.DefaultAction != nil && *acls.DefaultAction != "" {
					hasDefaultAction = *acls.DefaultAction == webpubsub.ACLActionDeny
				}

				hasDefaultMatches := false
				if acls.PublicNetwork != nil && acls.PublicNetwork.Allow != nil {
					hasDefaultMatches = len(*acls.PublicNetwork.Allow) == 4
				}

				isDefaultConfiguration = hasDefaultAction && hasDefaultMatches
			}
		}
	}

	return utils.Bool(!isDefaultConfiguration), nil
}

func (r WebPubsubNetworkACLResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_web_pubsub_network_acl" "test" {
  web_pubsub_id  = azurerm_web_pubsub.test.id
  default_action = "Allow"
  public_network {
    denied_request_types = ["RESTAPI"]
  }
  depends_on = [azurerm_web_pubsub.test]
}
`, r.template(data))
}

func (r WebPubsubNetworkACLResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  address_space       = ["10.5.0.0/16"]
}
resource "azurerm_subnet" "test" {
  name                                          = "acctest-subnet-%d"
  resource_group_name                           = azurerm_resource_group.test.name
  virtual_network_name                          = azurerm_virtual_network.test.name
  address_prefixes                              = ["10.5.2.0/24"]
  private_link_service_network_policies_enabled = true
}
resource "azurerm_private_endpoint" "test" {
  name                = "acctest-pe-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.test.id

  private_service_connection {
    name                           = "psc-sig-test"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_web_pubsub.test.id
    subresource_names              = ["webpubsub"]
  }
}

resource "azurerm_web_pubsub_network_acl" "test" {
  web_pubsub_id  = azurerm_web_pubsub.test.id
  default_action = "Deny"
  public_network {
    allowed_request_types = ["ClientConnection"]
  }

  private_endpoint {
    id                    = azurerm_private_endpoint.test.id
    allowed_request_types = ["ClientConnection"]
  }

  depends_on = [azurerm_web_pubsub.test]
}
`, r.template(data), data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r WebPubsubNetworkACLResource) complete_withoutPrivateEndpoint(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_web_pubsub_network_acl" "test" {
  web_pubsub_id  = azurerm_web_pubsub.test.id
  default_action = "Deny"
  public_network {
    allowed_request_types = ["ClientConnection", "RESTAPI"]
  }

  depends_on = [azurerm_web_pubsub.test]
}
`, r.template(data))
}

func (r WebPubsubNetworkACLResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_web_pubsub_network_acl" "import" {
  web_pubsub_id  = azurerm_web_pubsub_network_acl.test.web_pubsub_id
  default_action = azurerm_web_pubsub_network_acl.test.default_action
  public_network {
    denied_request_types = azurerm_web_pubsub_network_acl.test.public_network[0].denied_request_types
  }
  depends_on = [azurerm_web_pubsub.test]
}
`, r.basic(data))
}

func (r WebPubsubNetworkACLResource) multiplePrivateEndpoints(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}
resource "azurerm_subnet" "test" {
  name                                          = "acctest-subnet-%d"
  resource_group_name                           = azurerm_resource_group.test.name
  virtual_network_name                          = azurerm_virtual_network.test.name
  address_prefixes                              = ["10.5.2.0/24"]
  private_link_service_network_policies_enabled = true
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-pe-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.test.id

  private_service_connection {
    name                           = "psc-sig-test"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_web_pubsub.test.id
    subresource_names              = ["webpubsub"]
  }
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctest-vnet2-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "test2" {
  name                 = "acctest-subnet2-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test2.name
  address_prefixes     = ["10.5.2.0/24"]

  private_link_service_network_policies_enabled = true
}

resource "azurerm_private_endpoint" "test2" {
  name                = "acctest-pe2-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.test2.id

  private_service_connection {
    name                           = "psc-sig-test2"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_web_pubsub.test.id
    subresource_names              = ["webpubsub"]
  }
}

resource "azurerm_web_pubsub_network_acl" "test" {
  web_pubsub_id  = azurerm_web_pubsub.test.id
  default_action = "Allow"

  public_network {
    denied_request_types = ["ClientConnection"]
  }

  private_endpoint {
    id                   = azurerm_private_endpoint.test.id
    denied_request_types = ["ClientConnection"]
  }

  private_endpoint {
    id                   = azurerm_private_endpoint.test2.id
    denied_request_types = ["ServerConnection"]
  }

  depends_on = [azurerm_web_pubsub.test]
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r WebPubsubNetworkACLResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-webpubsub-%d"
  location = "%s"
}

resource "azurerm_web_pubsub" "test" {
  name                = "acctestRG-webpubsub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_S1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
