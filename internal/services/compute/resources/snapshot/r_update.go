// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package snapshot

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/custompoller"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSnapshotCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.SnapshotsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := snapshots.NewSnapshotID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := location.Normalize(d.Get("location").(string))
	createOption := d.Get("create_option").(string)
	t := d.Get("tags").(map[string]any)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_snapshot", id.ID())
			}
		}
	}

	properties := snapshots.Snapshot{
		Location: location,
		Properties: &snapshots.SnapshotProperties{
			CreationData: snapshots.CreationData{
				CreateOption: snapshots.DiskCreateOption(createOption),
			},
			Incremental: pointer.To(d.Get("incremental_enabled").(bool)),
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("source_uri"); ok {
		properties.Properties.CreationData.SourceUri = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("source_resource_id"); ok {
		properties.Properties.CreationData.SourceResourceId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("storage_account_id"); ok {
		properties.Properties.CreationData.StorageAccountId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("network_access_policy"); ok {
		properties.Properties.NetworkAccessPolicy = pointer.ToEnum[snapshots.NetworkAccessPolicy](v.(string))
	}

	if v, ok := d.GetOk("disk_access_id"); ok {
		properties.Properties.DiskAccessId = pointer.To(v.(string))
	}

	properties.Properties.PublicNetworkAccess = pointer.To(snapshots.PublicNetworkAccessEnabled)
	if !d.Get("public_network_access_enabled").(bool) {
		properties.Properties.PublicNetworkAccess = pointer.To(snapshots.PublicNetworkAccessDisabled)
	}

	diskSizeGB := d.Get("disk_size_gb").(int)
	if diskSizeGB > 0 {
		properties.Properties.DiskSizeGB = pointer.To(int64(diskSizeGB))
	}

	properties.Properties.EncryptionSettingsCollection = helpers.ExpandSnapshotDiskEncryptionSettings(d.Get("encryption_settings").([]any))

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, properties, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())

		// When `create_option` is `CopyStart` the API returns as soon as the copy
		// operation has been initiated, however the resulting snapshot is not usable
		// until the background copy has completed. Wait for `CompletionPercent` to
		// reach 100 so that downstream resources (e.g. `azurerm_managed_disk`) don't
		// consume an incomplete snapshot.
		if createOption == string(snapshots.DiskCreateOptionCopyStart) {
			log.Printf("[DEBUG] Waiting for the copy of %s to complete", id)

			pollerType := custompoller.NewSnapshotCopyStartPoller(client, id)
			poller := pollers.NewPoller(pollerType, 30*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for the copy of %s to complete: %+v", id, err)
			}
		}
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceSnapshotRead(d, meta)
}
