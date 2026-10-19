// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_origin

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigins"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorOriginDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdorigins.ParseOriginGroupOriginID(d.Id())
	if err != nil {
		return err
	}

	// @tombuildsstuff: JC/WS to dig into if we need to conditionally remove the Private Link
	// via an Update before deletion - presumably we'd also need a Lock on the private link resource

	/*
		original:
			// TODO: Check to see if there is a Load Balancer Private Link connected,
			// if so disconnect the Private Link association with the Frontdoor Origin
			// else the destroy will fail because the Private Link Service has an active
			// Private Link Endpoint connection...

			// It looks like Frontdoor does remove the Private link, I just need to poll here until it is removed...
			// Investigate this further...
			// WS: There is a bug in the service code, for only the load balancer scenario, the private link connection is not removed until the
			// origin is totally destroyed. The workaround for this issue is to put a depends_on the private link service to the origin so the origin
			// will be deleted first before the private link service is destroyed.
	*/

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
