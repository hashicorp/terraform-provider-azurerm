package parse

import (
    "fmt"

    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NetworkExpressRouteCircuitConnectionId struct {
    SubscriptionId string
    ResourceGroup string
    CircuitName string
    PeeringName string
    Name string
}

func NewExpressRouteCircuitConnectionID(subscriptionId string, resourcegroup string, circuitname string, peeringname string, name string) NetworkExpressRouteCircuitConnectionId {
    return NetworkExpressRouteCircuitConnectionId{
        SubscriptionId: subscriptionId,
        ResourceGroup: resourcegroup,
        CircuitName: circuitname,
        PeeringName: peeringname,
        Name: name,
    }
}

func (id NetworkExpressRouteCircuitConnectionId) ID() string {
    fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteCircuits/%s/peerings/%s/connections/%s"
    return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.CircuitName, id.PeeringName, id.Name)
}

func ExpressRouteCircuitConnectionID(input string) (*NetworkExpressRouteCircuitConnectionId, error) {
    id, err := azure.ParseAzureResourceID(input)
    if err != nil {
        return nil, fmt.Errorf("parsing networkExpressRouteCircuitConnection ID %q: %+v", input, err)
    }

    networkExpressRouteCircuitConnection := NetworkExpressRouteCircuitConnectionId{
        SubscriptionId: id.SubscriptionID,
        ResourceGroup: id.ResourceGroup,
    }
    if networkExpressRouteCircuitConnection.CircuitName, err = id.PopSegment("expressRouteCircuits"); err != nil {
        return nil, err
    }
    if networkExpressRouteCircuitConnection.PeeringName, err = id.PopSegment("peerings"); err != nil {
        return nil, err
    }
    if networkExpressRouteCircuitConnection.Name, err = id.PopSegment("connections"); err != nil {
        return nil, err
    }
    if err := id.ValidateNoEmptySegments(input); err != nil {
        return nil, err
    }

    return &networkExpressRouteCircuitConnection, nil
}
