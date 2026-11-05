// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package capacity_reservation_group

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCapacityReservationGroupRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.CapacityReservationGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacityreservationgroups.ParseCapacityReservationGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, capacityreservationgroups.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return resourceCapacityReservationGroupFlatten(d, id, resp.Model)
}
