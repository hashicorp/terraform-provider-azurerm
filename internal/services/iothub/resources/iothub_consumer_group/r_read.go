// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_consumer_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotHubConsumerGroupRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ConsumerGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetEventHubConsumerGroup(ctx, id.ResourceGroup, id.IotHubName, id.EventHubEndpointName, id.Name)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making read request for %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("iothub_name", id.IotHubName)
	d.Set("eventhub_endpoint_name", id.EventHubEndpointName)
	d.Set("resource_group_name", id.ResourceGroup)

	return nil
}
