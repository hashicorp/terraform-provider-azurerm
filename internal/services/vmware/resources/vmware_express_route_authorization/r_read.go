// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vmware_express_route_authorization

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/authorizations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/privateclouds"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	vmwareClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/vmware/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVmwareExpressRouteAuthorizationRead(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := account.From(meta).SubscriptionId
	client := vmwareClient.From(meta).AuthorizationClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := authorizations.ParseAuthorizationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", *id, err)
	}

	d.Set("name", id.AuthorizationName)
	d.Set("private_cloud_id", privateclouds.NewPrivateCloudID(subscriptionId, id.ResourceGroupName, id.PrivateCloudName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("express_route_authorization_id", props.ExpressRouteAuthorizationId)
			d.Set("express_route_authorization_key", props.ExpressRouteAuthorizationKey)
		}
	}

	return nil
}
