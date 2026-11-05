// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_endpoint

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorEndpointUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.AFDEndpointsClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdendpoints.ParseAfdEndpointID(d.Id())
	if err != nil {
		return err
	}

	// TODO: Check if this can be migrated to a PUT upgrade path for consistency with the rest of the provider
	props := afdendpoints.AFDEndpointUpdateParameters{}

	if d.HasChange("enabled") {
		enabledState := afdendpoints.EnabledStateEnabled
		if !d.Get("enabled").(bool) {
			enabledState = afdendpoints.EnabledStateDisabled
		}

		props.Properties = &afdendpoints.AFDEndpointPropertiesUpdateParameters{
			EnabledState: &enabledState,
		}
	}

	if d.HasChange("tags") {
		props.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.UpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceCdnFrontDoorEndpointRead(d, meta)
}
