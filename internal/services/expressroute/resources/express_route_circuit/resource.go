// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutecircuits"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var expressRouteCircuitResourceName = "azurerm_express_route_circuit"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceExpressRouteCircuitCreate,
		Read:   resourceExpressRouteCircuitRead,
		Update: resourceExpressRouteCircuitUpdate,
		Delete: resourceExpressRouteCircuitDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := expressroutecircuits.ParseExpressRouteCircuitID(id)
			return err
		}),

		CustomizeDiff: pluginsdk.CustomDiffInSequence(
			// If bandwidth is reduced force a new resource
			pluginsdk.ForceNewIfChange("bandwidth_in_mbps", func(ctx context.Context, old, new, meta any) bool {
				return new.(int) < old.(int)
			}),
		),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: expressRouteCircuitSchema(),
	}
}
