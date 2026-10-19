// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package snapshot

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSnapshotRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.SnapshotsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := snapshots.ParseSnapshotID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Error reading Snapshot %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.SnapshotName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			data := props.CreationData
			d.Set("create_option", string(data.CreateOption))
			d.Set("storage_account_id", data.StorageAccountId)
			d.Set("disk_access_id", pointer.From(props.DiskAccessId))

			diskSizeGb := 0
			if props.DiskSizeGB != nil {
				diskSizeGb = int(*props.DiskSizeGB)
			}
			d.Set("disk_size_gb", diskSizeGb)

			if err := d.Set("encryption_settings", helpers.FlattenSnapshotDiskEncryptionSettings(props.EncryptionSettingsCollection)); err != nil {
				return fmt.Errorf("setting `encryption_settings`: %+v", err)
			}

			networkAccessPolicy := snapshots.NetworkAccessPolicyAllowAll
			if props.NetworkAccessPolicy != nil {
				networkAccessPolicy = *props.NetworkAccessPolicy
			}
			d.Set("network_access_policy", string(networkAccessPolicy))

			publicNetworkAccessEnabled := true
			if v := props.PublicNetworkAccess; v != nil && *v != snapshots.PublicNetworkAccessEnabled {
				publicNetworkAccessEnabled = false
			}
			d.Set("public_network_access_enabled", publicNetworkAccessEnabled)

			d.Set("incremental_enabled", pointer.From(props.Incremental))

			trustedLaunchEnabled := false
			if securityProfile := props.SecurityProfile; securityProfile != nil && securityProfile.SecurityType != nil {
				trustedLaunchEnabled = *securityProfile.SecurityType == snapshots.DiskSecurityTypesTrustedLaunch
			}
			d.Set("trusted_launch_enabled", trustedLaunchEnabled)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
