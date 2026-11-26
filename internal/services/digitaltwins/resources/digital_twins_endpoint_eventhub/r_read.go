// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package digital_twins_endpoint_eventhub

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/digitaltwinsinstance"
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDigitalTwinsEndpointEventHubRead(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DigitalTwins.EndpointClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := endpoints.ParseEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.DigitalTwinsEndpointGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	d.Set("name", id.EndpointName)
	d.Set("digital_twins_id", digitaltwinsinstance.NewDigitalTwinsInstanceID(subscriptionId, id.ResourceGroupName, id.DigitalTwinsInstanceName).ID())

	if model := resp.Model; model != nil {
		if _, ok := model.Properties.(endpoints.EventHub); !ok {
			return fmt.Errorf("retrieving %s: expected an EventHub type but got: %+v", *id, model.Properties)
		}
	}

	return nil
}
