// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package capacity_reservation

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCapacityReservationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.CapacityReservationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacityreservations.ParseCapacityReservationID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
