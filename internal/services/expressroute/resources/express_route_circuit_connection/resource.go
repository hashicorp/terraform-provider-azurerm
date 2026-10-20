// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit_connection

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutecircuitconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceExpressRouteCircuitConnectionCreate,
		Read:   resourceExpressRouteCircuitConnectionRead,
		Update: resourceExpressRouteCircuitConnectionUpdate,
		Delete: resourceExpressRouteCircuitConnectionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := expressroutecircuitconnections.ParsePeeringConnectionID(id)
			return err
		}),

		Schema: expressRouteCircuitConnectionSchema(),
	}
}
