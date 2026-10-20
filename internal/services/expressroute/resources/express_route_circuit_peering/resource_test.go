// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit_peering_test

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ExpressRouteCircuitPeeringResource struct{}

func (r ExpressRouteCircuitPeeringResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseExpressRouteCircuitPeeringID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ExpressRoute.ExpressRouteCircuitPeerings.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}
