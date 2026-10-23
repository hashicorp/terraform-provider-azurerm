// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkInterfaceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseNetworkInterfaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if err := resourceNetworkInterfaceFlatten(d, id, resp.Model); err != nil {
		return fmt.Errorf("encoding %s: %+v", *id, err)
	}

	return nil
}
