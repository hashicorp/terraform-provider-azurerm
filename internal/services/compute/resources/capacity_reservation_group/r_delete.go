// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package capacity_reservation_group

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCapacityReservationGroupDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.CapacityReservationGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacityreservationgroups.ParseCapacityReservationGroupID(d.Id())
	if err != nil {
		return err
	}

	// It takes several seconds to sync the cache of reservations list in Capacity Reservation Group. Delete operation requires the list to be empty, and fails before the cache sync is completed.
	// Retry the delete operation after a minute as a workaround. Issue is tracked by: https://github.com/Azure/azure-rest-api-specs/issues/18767
	if _, err := client.Delete(ctx, *id); err != nil {
		stateConf := &pluginsdk.StateChangeConf{
			Pending: []string{"Deleting"},
			Target:  []string{"Deleted"},
			Refresh: func() (any, string, error) {
				res, err := client.Delete(ctx, *id)
				if err != nil {
					return res, "Deleting", nil // lint:ignore nilerr Returning nil error is intentional as we will retry the delete operation
				}
				return res, "Deleted", nil
			},
			MinTimeout: 15 * time.Second,
			Timeout:    1 * time.Minute,
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
		}
	}

	return nil
}
