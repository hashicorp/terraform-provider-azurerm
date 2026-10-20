// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package express_route_port

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteports"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmExpressRoutePortCreate,
		Read:   resourceArmExpressRoutePortRead,
		Update: resourceArmExpressRoutePortUpdate,
		Delete: resourceArmExpressRoutePortDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := expressrouteports.ParseExpressRoutePortID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(120 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(120 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(120 * time.Minute),
		},

		Schema: expressRoutePortResourceSchema(),
	}
}
