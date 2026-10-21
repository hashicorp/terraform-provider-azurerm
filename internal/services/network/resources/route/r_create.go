// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package route

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routetables"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRouteCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.Routes
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	addressPrefix := d.Get("address_prefix").(string)
	nextHopType := d.Get("next_hop_type").(string)

	id := routes.NewRouteID(subscriptionId, d.Get("resource_group_name").(string), d.Get("route_table_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_route", id.ID())
		}
	}

	rtId := routetables.NewRouteTableID(id.SubscriptionId, id.ResourceGroupName, id.RouteTableName)
	locks.ByID(rtId.ID())
	defer locks.UnlockByID(rtId.ID())

	route := routes.Route{
		Name: pointer.To(id.RouteName),
		Properties: &routes.RoutePropertiesFormat{
			AddressPrefix: &addressPrefix,
			NextHopType:   routes.RouteNextHopType(nextHopType),
		},
	}

	if v, ok := d.GetOk("next_hop_in_ip_address"); ok {
		route.Properties.NextHopIPAddress = pointer.To(v.(string))
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, route, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceRouteRead(d, meta)
}
