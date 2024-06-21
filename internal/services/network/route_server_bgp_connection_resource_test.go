// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type RouteServerBGPConnectionResource struct{}

func TestAccRouteServerBgpConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_server_bgp_connection", "test")
	r := RouteServerBGPConnectionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccRouteServerBgpConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_server_bgp_connection", "test")
	r := RouteServerBGPConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r RouteServerBGPConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseVirtualHubBGPConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.VirtualWANs.VirtualHubBgpConnectionGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r RouteServerBGPConnectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_route_server_bgp_connection" "test" {
  name            = "acctest-rs-bgp-%d"
  route_server_id = azurerm_route_server.test.id
  peer_asn        = 65501
  peer_ip         = "169.254.21.5"

}
`, RouteServerResource{}.basic(data), data.RandomInteger)
}

func (r RouteServerBGPConnectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_route_server_bgp_connection" "import" {
  name            = azurerm_route_server_bgp_connection.test.name
  route_server_id = azurerm_route_server_bgp_connection.test.route_server_id
  peer_asn        = azurerm_route_server_bgp_connection.test.peer_asn
  peer_ip         = azurerm_route_server_bgp_connection.test.peer_ip
}
`, r.basic(data))
}
