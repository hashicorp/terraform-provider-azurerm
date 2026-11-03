// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_fallback_route

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotHubFallbackRouteRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FallbackRouteID(d.Id())
	if err != nil {
		return err
	}

	iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		if response.WasNotFound(iothub.Response.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	d.Set("iothub_name", id.IotHubName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := iothub.Properties; props != nil {
		if routing := props.Routing; routing != nil {
			if fallbackRoute := routing.FallbackRoute; fallbackRoute != nil {
				d.Set("condition", fallbackRoute.Condition)
				d.Set("enabled", fallbackRoute.IsEnabled)
				d.Set("endpoint_names", fallbackRoute.EndpointNames)
				d.Set("source", fallbackRoute.Source)
			}
		}
	}

	return nil
}
