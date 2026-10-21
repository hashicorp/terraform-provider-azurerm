// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bastion_host

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/bastionhosts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBastionHostUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.BastionHostsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := bastionhosts.ParseBastionHostID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	sku := bastionhosts.BastionHostSkuName(d.Get("sku").(string))

	if d.HasChange("sku") {
		payload.Sku = &bastionhosts.Sku{
			Name: pointer.To(sku),
		}
	}

	if d.HasChange("copy_paste_enabled") {
		payload.Properties.DisableCopyPaste = pointer.To(!d.Get("copy_paste_enabled").(bool))
	}

	if d.HasChange("file_copy_enabled") {
		fileCopyEnabled := d.Get("file_copy_enabled").(bool)
		if fileCopyEnabled && (sku != bastionhosts.BastionHostSkuNameStandard && sku != bastionhosts.BastionHostSkuNamePremium) {
			return fmt.Errorf("`file_copy_enabled` is only supported when `sku` is `Standard` or `Premium`")
		}
		payload.Properties.EnableFileCopy = pointer.To(fileCopyEnabled)
	}

	if d.HasChange("ip_connect_enabled") {
		ipConnectEnabled := d.Get("ip_connect_enabled").(bool)
		if ipConnectEnabled && (sku != bastionhosts.BastionHostSkuNameStandard && sku != bastionhosts.BastionHostSkuNamePremium) {
			return fmt.Errorf("`ip_connect_enabled` is only supported when `sku` is `Standard` or `Premium`")
		}
		payload.Properties.EnableIPConnect = pointer.To(ipConnectEnabled)
	}

	if d.HasChange("scale_units") {
		scaleUnits := d.Get("scale_units").(int)
		if scaleUnits > 2 && (sku != bastionhosts.BastionHostSkuNameStandard && sku != bastionhosts.BastionHostSkuNamePremium) {
			return fmt.Errorf("`scale_units` only can be changed when `sku` is `Standard` or `Premium`. `scale_units` is always `2` when `sku` is `Basic`")
		}
		payload.Properties.ScaleUnits = pointer.To(int64(scaleUnits))
	}

	if d.HasChange("shareable_link_enabled") {
		shareableLinkEnabled := d.Get("shareable_link_enabled").(bool)
		if shareableLinkEnabled && (sku != bastionhosts.BastionHostSkuNameStandard && sku != bastionhosts.BastionHostSkuNamePremium) {
			return fmt.Errorf("`shareable_link_enabled` is only supported when `sku` is `Standard` or `Premium`")
		}
		payload.Properties.EnableShareableLink = pointer.To(shareableLinkEnabled)
	}

	if d.HasChange("tunneling_enabled") {
		tunnelingEnabled := d.Get("tunneling_enabled").(bool)
		if tunnelingEnabled && (sku != bastionhosts.BastionHostSkuNameStandard && sku != bastionhosts.BastionHostSkuNamePremium) {
			return fmt.Errorf("`tunneling_enabled` is only supported when `sku` is `Standard` or `Premium`")
		}
		payload.Properties.EnableTunneling = pointer.To(tunnelingEnabled)
	}

	if d.HasChange("session_recording_enabled") {
		sessionRecordingEnabled := d.Get("session_recording_enabled").(bool)
		if sessionRecordingEnabled && sku != bastionhosts.BastionHostSkuNamePremium {
			return fmt.Errorf("`session_recording_enabled` is only supported when `sku` is `Premium`")
		}
		payload.Properties.EnableSessionRecording = pointer.To(sessionRecordingEnabled)
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if d.HasChange("zones") {
		payload.Zones = pointer.To(zones.ExpandUntyped(d.Get("zones").(*schema.Set).List()))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceBastionHostRead(d, meta)
}
