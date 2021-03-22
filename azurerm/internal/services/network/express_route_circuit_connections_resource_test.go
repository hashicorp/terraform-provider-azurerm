package network_test

import (
    "context"
    "fmt"
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
    data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
    r := ExpressRouteCircuitConnectionResource{}
    data.ResourceTest(t, r, []resource.TestStep{
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
    data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
    r := ExpressRouteCircuitConnectionResource{}
    data.ResourceTest(t, r, []resource.TestStep{
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
    data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
    r := ExpressRouteCircuitConnectionResource{}
    data.ResourceTest(t, r, []resource.TestStep{
        {
            Config: r.complete(data),
            Check: resource.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(),
    })
}

func TestAccExpressRouteCircuitConnection_update(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
    r := ExpressRouteCircuitConnectionResource{}
    data.ResourceTest(t, r, []resource.TestStep{
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
        data.ImportStep(),
        {
            Config: r.basic(data),
            Check: resource.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(),
    })
}

func TestAccExpressRouteCircuitConnection_updateIpv6CircuitConnectionConfig(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
    r := ExpressRouteCircuitConnectionResource{}
    data.ResourceTest(t, r, []resource.TestStep{
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
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-network-%d"
  location = "%s"
}

resource "azurerm_express_route_circuit" "test" {
  name = "acctest-nerc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
}

resource "azurerm_express_route_circuit_peering" "test" {
  name = "acctest-nercp-%d"
  resource_group_name = azurerm_resource_group.test.name
  circuit_name = azurerm_express_route_circuit.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ExpressRouteCircuitConnectionResource) basic(data acceptance.TestData) string {
    template := r.template(data)
    return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit_connection" "test" {
  name = "acctest-nercc-%d"
  resource_group_name = azurerm_resource_group.test.name
  circuit_name = azurerm_express_route_circuit.test.name
  peering_name = azurerm_express_route_circuit_peering.test.name
}
`, template, data.RandomInteger)
}

func (r ExpressRouteCircuitConnectionResource) requiresImport(data acceptance.TestData) string {
    config := r.basic(data)
    return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit_connection" "import" {
  name = azurerm_express_route_circuit_connection.test.name
  resource_group_name = azurerm_express_route_circuit_connection.test.resource_group_name
  circuit_name = azurerm_express_route_circuit_connection.test.circuit_name
  peering_name = azurerm_express_route_circuit_connection.test.peering_name
}
`, config)
}

func (r ExpressRouteCircuitConnectionResource) complete(data acceptance.TestData) string {
    template := r.template(data)
    return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit_connection" "test" {
  name = "acctest-nercc-%d"
  resource_group_name = azurerm_resource_group.test.name
  circuit_name = azurerm_express_route_circuit.test.name
  peering_name = azurerm_express_route_circuit_peering.test.name
  address_prefix = "10.0.0.0/29"
  authorization_key = "946a1918-b7a2-4917-b43c-8c4cdaee006a"
  ipv6circuit_connection_config {
    address_prefix = "aa:bb::/125"
  }
}
`, template, data.RandomInteger)
}

func (r ExpressRouteCircuitConnectionResource) updateIpv6CircuitConnectionConfig(data acceptance.TestData) string {
    template := r.template(data)
    return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit_connection" "test" {
  name = "acctest-nercc-%d"
  resource_group_name = azurerm_resource_group.test.name
  circuit_name = azurerm_express_route_circuit.test.name
  peering_name = azurerm_express_route_circuit_peering.test.name
  address_prefix = "10.0.0.0/29"
  authorization_key = "946a1918-b7a2-4917-b43c-8c4cdaee006a"
  ipv6circuit_connection_config {
    address_prefix = "aa:bb::/125"
  }
}
`, template, data.RandomInteger)
}
