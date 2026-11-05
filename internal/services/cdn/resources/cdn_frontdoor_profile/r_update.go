// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_profile

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorProfileUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorProfilesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := profiles.ParseProfileID(d.Id())
	if err != nil {
		return err
	}

	props := profiles.ProfileUpdateParameters{
		Tags:       tags.Expand(d.Get("tags").(map[string]any)),
		Properties: &profiles.ProfilePropertiesUpdateParameters{},
	}

	if d.HasChange("response_timeout_seconds") {
		props.Properties.OriginResponseTimeoutSeconds = pointer.To(int64(d.Get("response_timeout_seconds").(int)))
	}

	if d.HasChange("log_scrubbing_rule") {
		props.Properties.LogScrubbing = expandCdnFrontDoorProfileLogScrubbing(d.Get("log_scrubbing_rule").(*pluginsdk.Set).List())
	}

	if d.HasChange("identity") {
		i, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		props.Identity = i
	}

	if err = client.UpdateThenPoll(ctx, pointer.From(id), props); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceCdnFrontDoorProfileRead(d, meta)
}
