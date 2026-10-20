// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit_authorization_test

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutecircuitauthorizations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ExpressRouteCircuitAuthorizationResource struct{}

func (t ExpressRouteCircuitAuthorizationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := expressroutecircuitauthorizations.ParseAuthorizationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ExpressRoute.ExpressRouteCircuitAuthorizations.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}
