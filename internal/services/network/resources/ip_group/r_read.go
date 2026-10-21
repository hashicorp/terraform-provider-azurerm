// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ip_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ipgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIpGroupRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.IPGroups
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ipgroups.ParseIPGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, ipgroups.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return resourceIpGroupFlatten(d, id, resp.Model)
}

func getIds(subResource *[]ipgroups.SubResource) []string {
	if subResource == nil {
		return nil
	}

	ids := make([]string, 0)
	for _, v := range *subResource {
		if v.Id == nil {
			continue
		}

		ids = append(ids, *v.Id)
	}

	return ids
}
