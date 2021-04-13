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

func TestAccExpressRouteCircuitConnection_updateIpv6CircuitConnectionConfig(t *testing.T) {
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
		data.ImportStep(),
		{
			Config: r.updateIpv6CircuitConnectionConfig(data),
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
	resp, err := client.Network.ExpressRouteCircuitConnectionClient.Get(ctx, id.ResourceGroup, id.CircuitName, id.PeeringName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving ExpressRouteCircuitConnection %q (Resource Group %q / circuitName %q / peeringName %q): %+v", id.Name, id.ResourceGroup, id.CircuitName, id.PeeringName, err)
	}
	return utils.Bool(true), nil
}

func (r ExpressRouteCircuitConnectionResource) template(data acceptance.TestData) string {
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

func (r ExpressRouteCircuitConnectionResource) basic(data acceptance.TestData) string {
	rg := os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP")
	circuitName := os.Getenv("ARM_TEST_CIRCUIT_NAME_FIRST")
	template := r.template(data)

	return fmt.Sprintf(`%s

resource "azurerm_express_route_circuit_connection" "test" {
  name = "acctest-nercc-%d"
  resource_group_name = "%s"
  circuit_name = "%s"
  peering_id = azurerm_express_route_circuit_peering.test.id
  peer_peering_id = azurerm_express_route_circuit_peering.peer_test.id
  address_prefix = "192.169.8.0/29"
}
`, template, data.RandomInteger, rg, circuitName)
}

func (r ExpressRouteCircuitConnectionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`%s

resource "azurerm_express_route_circuit_connection" "import" {
  name = azurerm_express_route_circuit_connection.test.name
  resource_group_name = azurerm_express_route_circuit_connection.test.resource_group_name
  circuit_name = azurerm_express_route_circuit_connection.test.circuit_name
  peering_id = azurerm_express_route_circuit_peering.test.id
  peer_peering_id = azurerm_express_route_circuit_peering.peer_test.id
  address_prefix = "192.169.8.0/29"
}
`, config)
}

func (r ExpressRouteCircuitConnectionResource) complete(data acceptance.TestData) string {
	rg := os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP")
	circuitName := os.Getenv("ARM_TEST_CIRCUIT_NAME_FIRST")

	template := r.template(data)
	return fmt.Sprintf(`%s

resource "azurerm_express_route_circuit_connection" "test" {
  name = "acctest-nercc-%d"
  resource_group_name = "%s"
  circuit_name = "%s"
  peering_id = azurerm_express_route_circuit_peering.test.id
  peer_peering_id = azurerm_express_route_circuit_peering.peer_test.id
  address_prefix = "192.169.8.0/29"
  authorization_key = "946a1918-b7a2-4917-b43c-8c4cdaee006a"
  ipv6circuit_connection_config {
    address_prefix = "aa:bb::/125"
  }
}
`, template, data.RandomInteger, rg, circuitName)
}

func (r ExpressRouteCircuitConnectionResource) updateIpv6CircuitConnectionConfig(data acceptance.TestData) string {
	rg := os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP")
	circuitName := os.Getenv("ARM_TEST_CIRCUIT_NAME_FIRST")

	template := r.template(data)
	return fmt.Sprintf(`%s

resource "azurerm_express_route_circuit_connection" "test" {
  name = "acctest-nercc-%d"
  resource_group_name = "%s"
  circuit_name = "%s"
  peering_id = azurerm_express_route_circuit_peering.test.id
  peer_peering_id = azurerm_express_route_circuit_peering.peer_test.id
  address_prefix = "192.169.8.0/29"
  authorization_key = "946a1918-b7a2-4917-b43c-8c4cdaee006a"
  ipv6circuit_connection_config {
    address_prefix = "aa:bb::/125"
  }
}
`, template, data.RandomInteger, rg, circuitName)
}
