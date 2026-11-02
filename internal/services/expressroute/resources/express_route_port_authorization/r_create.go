// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_port_authorization

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteportauthorizations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteports"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRoutePortAuthorizationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ExpressRoute.ExpressRoutePortAuthorizations
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	id := expressrouteportauthorizations.NewExpressRoutePortAuthorizationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("express_route_port_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_express_route_port_authorization", id.ID())
		}
	}

	properties := expressrouteportauthorizations.ExpressRoutePortAuthorization{
		Properties: &expressrouteportauthorizations.ExpressRoutePortAuthorizationPropertiesFormat{},
	}

	// can run only one create/update/delete operation of expressRoutePort at the same time
	portID := expressrouteports.NewExpressRoutePortID(id.SubscriptionId, id.ResourceGroupName, id.ExpressRoutePortName)
	locks.ByID(portID.ID())
	defer locks.UnlockByID(portID.ID())

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, properties, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRoutePortAuthorizationRead(d, meta)
}
