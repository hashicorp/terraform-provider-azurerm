// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_port_authorization

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteportauthorizations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteports"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	expressrouteClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/expressroute/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRoutePortAuthorizationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := expressrouteClient.From(meta).ExpressRoutePortAuthorizations
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := expressrouteportauthorizations.ParseExpressRoutePortAuthorizationID(d.Id())
	if err != nil {
		return err
	}

	portID := expressrouteports.NewExpressRoutePortID(id.SubscriptionId, id.ResourceGroupName, id.ExpressRoutePortName)
	locks.ByID(portID.ID())
	defer locks.UnlockByID(portID.ID())

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
