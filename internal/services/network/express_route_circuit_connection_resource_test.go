// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/expressroutecircuitconnections"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ExpressRouteCircuitConnectionResource struct{}

func TestAccExpressRouteCircuitConnection(t *testing.T) {
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"Resource": {
			"basic":          testAccExpressRouteCircuitConnection_basic,
			"requiresImport": testAccExpressRouteCircuitConnection_requiresImport,
			"complete":       testAccExpressRouteCircuitConnection_complete,
			"update":         testAccExpressRouteCircuitConnection_update,
		},
	})
}

func testAccExpressRouteCircuitConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
	r := ExpressRouteCircuitConnectionResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccExpressRouteCircuitConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
	r := ExpressRouteCircuitConnectionResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccExpressRouteCircuitConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
	r := ExpressRouteCircuitConnectionResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "846a1918-b7a2-4917-b43c-8c4cdaee006a"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccExpressRouteCircuitConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
	r := ExpressRouteCircuitConnectionResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, "846a1918-b7a2-4917-b43c-8c4cdaee006a"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, "946a1918-b7a2-4917-b43c-8c4cdaee006a"),
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

func TestAccExpressRouteCircuitConnection_writeOnlyAuthorizationKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
	r := ExpressRouteCircuitConnectionResource{}

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.11.0"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.writeOnlyAuthorizationKey(data, "846a1918-b7a2-4917-b43c-8c4cdaee006a", 1),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("authorization_key_wo_version"),
			{
				Config: r.writeOnlyAuthorizationKey(data, "946a1918-b7a2-4917-b43c-8c4cdaee006a", 2),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("authorization_key_wo_version"),
		},
	})
}

func TestAccExpressRouteCircuitConnection_updateToWriteOnlyAuthorizationKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
	r := ExpressRouteCircuitConnectionResource{}

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.11.0"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.complete(data, "846a1918-b7a2-4917-b43c-8c4cdaee006a"),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("authorization_key"),
			{
				Config: r.writeOnlyAuthorizationKey(data, "846a1918-b7a2-4917-b43c-8c4cdaee006a", 1),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("authorization_key_wo_version"),
			{
				Config: r.complete(data, "846a1918-b7a2-4917-b43c-8c4cdaee006a"),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("authorization_key"),
		},
	})
}

func (r ExpressRouteCircuitConnectionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := expressroutecircuitconnections.ParsePeeringConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.ExpressRouteCircuitConnections.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ExpressRouteCircuitConnectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit_connection" "test" {
  name                = "acctest-ExpressRouteCircuitConn-%d"
  peering_id          = azurerm_express_route_circuit_peering.test.id
  peer_peering_id     = azurerm_express_route_circuit_peering.peer_test.id
  address_prefix_ipv4 = "192.169.8.0/29"
}
`, r.template(data), data.RandomInteger)
}

func (r ExpressRouteCircuitConnectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit_connection" "import" {
  name                = azurerm_express_route_circuit_connection.test.name
  peering_id          = azurerm_express_route_circuit_connection.test.peering_id
  peer_peering_id     = azurerm_express_route_circuit_connection.test.peer_peering_id
  address_prefix_ipv4 = azurerm_express_route_circuit_connection.test.address_prefix_ipv4
}
`, r.basic(data))
}

func (r ExpressRouteCircuitConnectionResource) complete(data acceptance.TestData, authorizationKey string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit_connection" "test" {
  name                = "acctest-ExpressRouteCircuitConn-%d"
  peering_id          = azurerm_express_route_circuit_peering.test.id
  peer_peering_id     = azurerm_express_route_circuit_peering.peer_test.id
  address_prefix_ipv4 = "192.169.8.0/29"
  authorization_key   = "%s"

  lifecycle {
    ignore_changes = [
      "authorization_key"
    ]
  }
}
`, r.template(data), data.RandomInteger, authorizationKey)
}

func (r ExpressRouteCircuitConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ercircuitconn-%d"
  location = "%s"
}

resource "azurerm_express_route_port" "test" {
  name                = "acctest-erp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  peering_location    = "Equinix-London-LD5"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  express_route_port_id = azurerm_express_route_port.test.id
  bandwidth_in_gbps     = 5

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }
}

resource "azurerm_express_route_port" "peer_test" {
  name                = "acctest-erp2-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  peering_location    = "Equinix-Sydney-SY2"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
}

resource "azurerm_express_route_circuit" "peer_test" {
  name                  = "acctest-erc2-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  express_route_port_id = azurerm_express_route_port.peer_test.id
  bandwidth_in_gbps     = 5

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = azurerm_express_route_circuit.test.name
  resource_group_name           = azurerm_resource_group.test.name
  shared_key                    = "ItsASecret"
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.1.0/30"
  vlan_id                       = 100
}

resource "azurerm_express_route_circuit_peering" "peer_test" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = azurerm_express_route_circuit.peer_test.name
  resource_group_name           = azurerm_resource_group.test.name
  shared_key                    = "ItsASecret"
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.1.0/30"
  vlan_id                       = 100
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ExpressRouteCircuitConnectionResource) writeOnlyAuthorizationKey(data acceptance.TestData, secret string, version int) string {
	return fmt.Sprintf(`
%s

%s

resource "azurerm_express_route_circuit_connection" "test" {
  name                         = "acctest-ExpressRouteCircuitConn-%d"
  peering_id                   = azurerm_express_route_circuit_peering.test.id
  peer_peering_id              = azurerm_express_route_circuit_peering.peer_test.id
  address_prefix_ipv4          = "192.169.8.0/29"
  authorization_key_wo         = ephemeral.azurerm_key_vault_secret.test.value
  authorization_key_wo_version = %d
}
`, r.template(data), acceptance.WriteOnlyKeyVaultSecretTemplate(data, secret), data.RandomInteger, version)
}
