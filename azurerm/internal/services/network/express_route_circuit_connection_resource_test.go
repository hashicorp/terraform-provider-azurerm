package network_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ExpressRouteCircuitConnectionResource struct{}

func TestAccExpressRouteCircuitConnection_basic(t *testing.T) {
	if os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" || os.Getenv("ARM_TEST_CIRCUIT_NAME_FIRST") == "" || os.Getenv("ARM_TEST_CIRCUIT_NAME_SECOND") == "" {
		t.Skip("Skipping as ARM_TEST_DATA_RESOURCE_GROUP and/or ARM_TEST_CIRCUIT_NAME_FIRST and/or ARM_TEST_CIRCUIT_NAME_SECOND are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
	r := ExpressRouteCircuitConnectionResource{}
	data.ResourceSequentialTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccExpressRouteCircuitConnection_requiresImport(t *testing.T) {
	if os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" || os.Getenv("ARM_TEST_CIRCUIT_NAME_FIRST") == "" ||
		os.Getenv("ARM_TEST_CIRCUIT_NAME_SECOND") == "" {
		t.Skip("Skipping as ARM_TEST_DATA_RESOURCE_GROUP and/or ARM_TEST_CIRCUIT_NAME_FIRST and/or ARM_TEST_CIRCUIT_NAME_SECOND are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
	r := ExpressRouteCircuitConnectionResource{}
	data.ResourceSequentialTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccExpressRouteCircuitConnection_complete(t *testing.T) {
	if os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" || os.Getenv("ARM_TEST_CIRCUIT_NAME_FIRST") == "" ||
		os.Getenv("ARM_TEST_CIRCUIT_NAME_SECOND") == "" {
		t.Skip("Skipping as ARM_TEST_DATA_RESOURCE_GROUP and/or ARM_TEST_CIRCUIT_NAME_FIRST and/or ARM_TEST_CIRCUIT_NAME_SECOND are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
	r := ExpressRouteCircuitConnectionResource{}
	data.ResourceSequentialTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authorization_key"),
	})
}

func TestAccExpressRouteCircuitConnection_update(t *testing.T) {
	if os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" || os.Getenv("ARM_TEST_CIRCUIT_NAME_FIRST") == "" || os.Getenv("ARM_TEST_CIRCUIT_NAME_SECOND") == "" {
		t.Skip("Skipping as ARM_TEST_DATA_RESOURCE_GROUP and/or ARM_TEST_CIRCUIT_NAME_FIRST and/or ARM_TEST_CIRCUIT_NAME_SECOND are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
	r := ExpressRouteCircuitConnectionResource{}
	data.ResourceSequentialTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authorization_key"),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authorization_key"),
	})
}

func TestAccExpressRouteCircuitConnection_updateAddressPrefixIPv6(t *testing.T) {
	if os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" || os.Getenv("ARM_TEST_CIRCUIT_NAME_FIRST") == "" || os.Getenv("ARM_TEST_CIRCUIT_NAME_SECOND") == "" {
		t.Skip("Skipping as ARM_TEST_DATA_RESOURCE_GROUP and/or ARM_TEST_CIRCUIT_NAME_FIRST and/or ARM_TEST_CIRCUIT_NAME_SECOND are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
	r := ExpressRouteCircuitConnectionResource{}
	data.ResourceSequentialTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateAddressPrefixIPv6(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ExpressRouteCircuitConnectionResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ExpressRouteCircuitConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.ExpressRouteCircuitConnectionClient.Get(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, id.PeeringName, id.ConnectionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
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
`, r.template(), data.RandomInteger)
}

func (r ExpressRouteCircuitConnectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit_connection" "import" {
  name                = azurerm_express_route_circuit_connection.test.name
  peering_id          = azurerm_express_route_circuit_peering.test.id
  peer_peering_id     = azurerm_express_route_circuit_peering.peer_test.id
  address_prefix_ipv4 = azurerm_express_route_circuit_connection.test.address_prefix
}
`, r.basic(data))
}

func (r ExpressRouteCircuitConnectionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit_connection" "test" {
  name                = "acctest-ExpressRouteCircuitConn-%d"
  peering_id          = azurerm_express_route_circuit_peering.test.id
  peer_peering_id     = azurerm_express_route_circuit_peering.peer_test.id
  address_prefix_ipv4 = "192.169.8.0/29"
  authorization_key   = "946a1918-b7a2-4917-b43c-8c4cdaee006a"
  address_prefix_ipv6 = "aa:bb::/125"
}
`, r.template(), data.RandomInteger)
}

func (r ExpressRouteCircuitConnectionResource) updateAddressPrefixIPv6(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit_connection" "test" {
  name                = "acctest-ExpressRouteCircuitConn-%d"
  peering_id          = azurerm_express_route_circuit_peering.test.id
  peer_peering_id     = azurerm_express_route_circuit_peering.peer_test.id
  address_prefix_ipv4 = "192.169.8.0/29"
  authorization_key   = "946a1918-b7a2-4917-b43c-8c4cdaee006a"
  address_prefix_ipv6 = "aa:bb::/125"
}
`, r.template(), data.RandomInteger)
}

func (r ExpressRouteCircuitConnectionResource) template() string {
	rg := os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP")
	circuitName1 := os.Getenv("ARM_TEST_CIRCUIT_NAME_FIRST")
	circuitName2 := os.Getenv("ARM_TEST_CIRCUIT_NAME_SECOND")

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = "%[1]s"
  resource_group_name           = "%[2]s"
  shared_key                    = "SSSSsssssshhhhhItsASecret"
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.1.0/30"
  vlan_id                       = 100
}

resource "azurerm_express_route_circuit_peering" "peer_test" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = "%[3]s"
  resource_group_name           = "%[2]s"
  shared_key                    = "SSSSsssssshhhhhItsASecret"
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.1.0/30"
  vlan_id                       = 100
}
`, circuitName1, rg, circuitName2)
}
