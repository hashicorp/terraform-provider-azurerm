// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nat_gateway

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/natgateways"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNatGatewayCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NatGateways
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := natgateways.NewNatGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.NatGatewayName, natGatewayResourceName)
	defer locks.UnlockByName(id.NatGatewayName, natGatewayResourceName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		resp, err := client.Get(ctx, id, natgateways.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for present of existing %s: %+v", id, err)
			}
		}
		if resp.Model != nil && resp.Model.Id != nil && *resp.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_nat_gateway", id.ID())
		}
	}

	parameters := natgateways.NatGateway{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &natgateways.NatGatewayPropertiesFormat{
			IdleTimeoutInMinutes: pointer.To(int64(d.Get("idle_timeout_in_minutes").(int))),
		},
		Sku: &natgateways.NatGatewaySku{
			Name: pointer.ToEnum[natgateways.NatGatewaySkuName](d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		parameters.Zones = &zones
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceNatGatewayRead(d, meta)
}
