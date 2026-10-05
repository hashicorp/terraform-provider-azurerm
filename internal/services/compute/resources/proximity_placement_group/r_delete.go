// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package proximity_placement_group

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceProximityPlacementGroupDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.ProximityPlacementGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := proximityplacementgroups.ParseProximityPlacementGroupID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id); err != nil {
		if response.WasNotFound(resp.HttpResponse) || response.WasStatusCode(resp.HttpResponse, http.StatusNoContent) { // API quirk that delete on non-existent resource returns 204 ¯\_(ツ)_/¯
			return nil
		}
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
