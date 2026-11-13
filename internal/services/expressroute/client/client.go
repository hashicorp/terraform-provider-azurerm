// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/expressroutecircuitauthorizations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/expressroutecircuitconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/expressroutecircuitpeerings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/expressroutecircuits"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/expressrouteportauthorizations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/expressrouteports"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/routefilters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ExpressRouteCircuitAuthorizations *expressroutecircuitauthorizations.ExpressRouteCircuitAuthorizationsClient
	ExpressRouteCircuitConnections    *expressroutecircuitconnections.ExpressRouteCircuitConnectionsClient
	ExpressRouteCircuitPeerings       *expressroutecircuitpeerings.ExpressRouteCircuitPeeringsClient
	ExpressRouteCircuits              *expressroutecircuits.ExpressRouteCircuitsClient
	ExpressRoutePortAuthorizations    *expressrouteportauthorizations.ExpressRoutePortAuthorizationsClient
	ExpressRoutePorts                 *expressrouteports.ExpressRoutePortsClient
	RouteFilters                      *routefilters.RouteFiltersClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	expressRouteCircuitAuthorizationsClient, err := expressroutecircuitauthorizations.NewExpressRouteCircuitAuthorizationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ExpressRouteCircuitAuthorizations client: %+v", err)
	}
	o.Configure(expressRouteCircuitAuthorizationsClient.Client, o.Authorizers.ResourceManager)

	expressRouteCircuitConnectionsClient, err := expressroutecircuitconnections.NewExpressRouteCircuitConnectionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ExpressRouteCircuitConnections client: %+v", err)
	}
	o.Configure(expressRouteCircuitConnectionsClient.Client, o.Authorizers.ResourceManager)

	expressRouteCircuitPeeringsClient, err := expressroutecircuitpeerings.NewExpressRouteCircuitPeeringsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ExpressRouteCircuitPeerings client: %+v", err)
	}
	o.Configure(expressRouteCircuitPeeringsClient.Client, o.Authorizers.ResourceManager)

	expressRouteCircuitsClient, err := expressroutecircuits.NewExpressRouteCircuitsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ExpressRouteCircuits client: %+v", err)
	}
	o.Configure(expressRouteCircuitsClient.Client, o.Authorizers.ResourceManager)

	expressRoutePortAuthorizationsClient, err := expressrouteportauthorizations.NewExpressRoutePortAuthorizationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ExpressRoutePortAuthorizations client: %+v", err)
	}
	o.Configure(expressRoutePortAuthorizationsClient.Client, o.Authorizers.ResourceManager)

	expressRoutePortsClient, err := expressrouteports.NewExpressRoutePortsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ExpressRoutePorts client: %+v", err)
	}
	o.Configure(expressRoutePortsClient.Client, o.Authorizers.ResourceManager)

	routeFiltersClient, err := routefilters.NewRouteFiltersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building RouteFilters client: %+v", err)
	}
	o.Configure(routeFiltersClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ExpressRouteCircuitAuthorizations: expressRouteCircuitAuthorizationsClient,
		ExpressRouteCircuitConnections:    expressRouteCircuitConnectionsClient,
		ExpressRouteCircuitPeerings:       expressRouteCircuitPeeringsClient,
		ExpressRouteCircuits:              expressRouteCircuitsClient,
		ExpressRoutePortAuthorizations:    expressRoutePortAuthorizationsClient,
		ExpressRoutePorts:                 expressRoutePortsClient,
		RouteFilters:                      routeFiltersClient,
	}, nil
}
