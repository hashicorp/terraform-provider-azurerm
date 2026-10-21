// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package route_table

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routetables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRouteTableRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.RouteTables
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := routetables.ParseRouteTableID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, routetables.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if err := resourceRouteTableFlatten(d, id, resp.Model); err != nil {
		return fmt.Errorf("flattening %s: %+v", id, err)
	}

	return nil
}
