// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dedicated_host_group

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/dedicatedhostgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDedicatedHostGroupRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.DedicatedHostGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseDedicatedHostGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, dedicatedhostgroups.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %q was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return resourceDedicatedHostGroupFlatten(d, id, resp.Model)
}
