// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_route

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotHubRouteRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.RouteID(d.Id())
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

	d.Set("name", id.Name)
	d.Set("iothub_name", id.IotHubName)
	d.Set("resource_group_name", id.ResourceGroup)

	if iothub.Properties == nil || iothub.Properties.Routing == nil {
		d.SetId("")
		return nil
	}

	exist := false

	if routes := iothub.Properties.Routing.Routes; routes != nil {
		for _, route := range *routes {
			if route.Name != nil {
				if strings.EqualFold(*route.Name, id.Name) {
					exist = true
					d.Set("source", route.Source)
					d.Set("condition", route.Condition)
					d.Set("enabled", route.IsEnabled)
					d.Set("endpoint_names", route.EndpointNames)
				}
			}
		}
	}

	if !exist {
		d.SetId("")
	}

	return nil
}
