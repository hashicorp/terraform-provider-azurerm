// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dedicated_host

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/dedicatedhosts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDedicatedHostUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.DedicatedHostsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseDedicatedHostID(d.Id())
	if err != nil {
		return err
	}

	payload := dedicatedhosts.DedicatedHostUpdate{}

	if d.HasChanges("auto_replace_on_failure", "license_type") {
		payload.Properties = &dedicatedhosts.DedicatedHostProperties{}
		if d.HasChange("auto_replace_on_failure") {
			payload.Properties.AutoReplaceOnFailure = pointer.To(d.Get("auto_replace_on_failure").(bool))
		}
		if d.HasChange("license_type") {
			licenseType := dedicatedhosts.DedicatedHostLicenseTypesNone
			if v := d.Get("license_type").(string); v != "" {
				licenseType = dedicatedhosts.DedicatedHostLicenseTypes(v)
			}
			payload.Properties.LicenseType = pointer.To(licenseType)
		}
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceDedicatedHostRead(d, meta)
}
