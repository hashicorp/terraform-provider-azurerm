// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_ddos_protection_plan

import (
	"fmt"
	"log"
	"slices"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ddosprotectionplans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkDDoSProtectionPlanRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.DdosProtectionPlans
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ddosprotectionplans.ParseDdosProtectionPlanID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}
	return resourceNetworkDDoSProtectionPlanFlatten(d, id, resp.Model)
}

func extractVnetIDs(input *[]ddosprotectionplans.SubResource) (*[]string, error) {
	vnetIDs := make([]string, 0)

	if input != nil {
		for _, subresource := range *input {
			if subresource.Id == nil {
				continue
			}

			id, err := commonids.ParseVirtualNetworkIDInsensitively(*subresource.Id)
			if err != nil {
				return nil, err
			}

			if !slices.Contains(vnetIDs, id.ID()) {
				vnetIDs = append(vnetIDs, id.ID())
			}
		}
	}

	return &vnetIDs, nil
}
