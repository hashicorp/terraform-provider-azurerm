// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nat_gateway

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/natgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNatGatewayUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NatGateways
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := natgateways.ParseNatGatewayID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.NatGatewayName, natGatewayResourceName)
	defer locks.UnlockByName(id.NatGatewayName, natGatewayResourceName)

	existing, err := client.Get(ctx, *id, natgateways.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("%s was not found", *id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}
	props := *existing.Model.Properties

	// intentionally building a new object rather than reusing due to the additional read-only fields
	payload := natgateways.NatGateway{
		Location: existing.Model.Location,
		Properties: &natgateways.NatGatewayPropertiesFormat{
			IdleTimeoutInMinutes: props.IdleTimeoutInMinutes,
			PublicIPAddresses:    props.PublicIPAddresses, // note: these can be managed via the separate resource
			PublicIPPrefixes:     props.PublicIPPrefixes,
		},
		Sku:   existing.Model.Sku,
		Tags:  existing.Model.Tags,
		Zones: existing.Model.Zones,
	}

	if d.HasChange("idle_timeout_in_minutes") {
		timeout := d.Get("idle_timeout_in_minutes").(int)
		payload.Properties.IdleTimeoutInMinutes = pointer.To(int64(timeout))
	}

	if d.HasChange("sku_name") {
		payload.Sku = &natgateways.NatGatewaySku{
			Name: pointer.ToEnum[natgateways.NatGatewaySkuName](d.Get("sku_name").(string)),
		}
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceNatGatewayRead(d, meta)
}
