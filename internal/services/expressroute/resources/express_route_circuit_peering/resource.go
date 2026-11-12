// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit_peering

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var expressRouteCircuitResourceName = "azurerm_express_route_circuit"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceExpressRouteCircuitPeeringCreate,
		Read:   resourceExpressRouteCircuitPeeringRead,
		Update: resourceExpressRouteCircuitPeeringUpdate,
		Delete: resourceExpressRouteCircuitPeeringDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseExpressRouteCircuitPeeringID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: expressRouteCircuitPeeringSchema(),
	}
}
