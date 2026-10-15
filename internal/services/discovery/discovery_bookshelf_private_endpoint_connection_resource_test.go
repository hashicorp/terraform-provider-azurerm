// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package discovery_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelfprivateendpointconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DiscoveryBookshelfPrivateEndpointConnectionResource struct{}

func TestAccDiscoveryBookshelfPrivateEndpointConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_bookshelf_private_endpoint_connection", "test")
	r := DiscoveryBookshelfPrivateEndpointConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("status").HasValue("Approved"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDiscoveryBookshelfPrivateEndpointConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_bookshelf_private_endpoint_connection", "test")
	r := DiscoveryBookshelfPrivateEndpointConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "Approved", "Approved by the acceptance test"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("status").HasValue("Approved"),
				check.That(data.ResourceName).Key("description").HasValue("Approved by the acceptance test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDiscoveryBookshelfPrivateEndpointConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_discovery_bookshelf_private_endpoint_connection", "test")
	r := DiscoveryBookshelfPrivateEndpointConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "Approved", "Approved by the acceptance test"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("status").HasValue("Approved"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, "Rejected", "Rejected by the acceptance test"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("status").HasValue("Rejected"),
				check.That(data.ResourceName).Key("description").HasValue("Rejected by the acceptance test"),
			),
		},
		data.ImportStep(),
	})
}

func (r DiscoveryBookshelfPrivateEndpointConnectionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := bookshelfprivateendpointconnections.ParseBookshelfPrivateEndpointConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Discovery.BookshelfPrivateEndpointConnections.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r DiscoveryBookshelfPrivateEndpointConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_subnet" "endpoint" {
  name                 = "endpointSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.4.0/24"]
}

data "azurerm_discovery_bookshelf_private_link_resources" "test" {
  bookshelf_id = azurerm_discovery_bookshelf.test.id
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctestpe-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = "acctestpe-%[2]d"
    private_connection_resource_id = azurerm_discovery_bookshelf.test.id
    subresource_names              = [data.azurerm_discovery_bookshelf_private_link_resources.test.private_link_resources[0].group_id]
    is_manual_connection           = true
    request_message                = "Requested by the acceptance test"
  }
}
`, DiscoveryBookshelfResource{}.basic(data), data.RandomInteger)
}

func (r DiscoveryBookshelfPrivateEndpointConnectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_discovery_bookshelf_private_endpoint_connection" "test" {
  bookshelf_id        = azurerm_discovery_bookshelf.test.id
  private_endpoint_id = azurerm_private_endpoint.test.id
  status              = "Approved"
}
`, r.template(data))
}

func (r DiscoveryBookshelfPrivateEndpointConnectionResource) complete(data acceptance.TestData, status string, description string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_discovery_bookshelf_private_endpoint_connection" "test" {
  bookshelf_id        = azurerm_discovery_bookshelf.test.id
  private_endpoint_id = azurerm_private_endpoint.test.id
  status              = "%[2]s"
  description         = "%[3]s"
}
`, r.template(data), status, description)
}
