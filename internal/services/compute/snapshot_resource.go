// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskaccesses"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSnapshot() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSnapshotCreateUpdate,
		Read:   resourceSnapshotRead,
		Update: resourceSnapshotCreateUpdate,
		Delete: resourceSnapshotDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := snapshots.ParseSnapshotID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SnapshotV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SnapshotName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"create_option": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(snapshots.DiskCreateOptionCopy),
					string(snapshots.DiskCreateOptionImport),
				}, false),
			},

			"incremental_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"source_uri": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"network_access_policy": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(snapshots.PossibleValuesForNetworkAccessPolicy(), false),
				Default:      string(snapshots.NetworkAccessPolicyAllowAll),
			},

			"disk_access_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: diskaccesses.ValidateDiskAccessID,
				// TODO:
				// the snapshot API is broken and returns the Resource Group name in UPPERCASE
				// tracked by https://github.com/Azure/azure-rest-api-specs/issues/29187
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"source_resource_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"storage_account_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"disk_size_gb": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Computed: true,
			},

			"encryption_settings": encryptionSettingsSchema(),

			"trusted_launch_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},

		// Encryption Settings cannot be disabled once enabled
		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("encryption_settings", func(ctx context.Context, old, new, meta interface{}) bool {
				if !features.FourPointOhBeta() {
					return false
				}
				return len(old.([]interface{})) > 0 && len(new.([]interface{})) == 0
			}),
		),
	}
}

func resourceSnapshotCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SnapshotsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := snapshots.NewSnapshotID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := azure.NormalizeLocation(d.Get("location").(string))
	createOption := d.Get("create_option").(string)
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
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

	properties := snapshots.Snapshot{
		Location: location,
		Properties: &snapshots.SnapshotProperties{
			CreationData: snapshots.CreationData{
				CreateOption: snapshots.DiskCreateOption(createOption),
			},
			Incremental: utils.Bool(d.Get("incremental_enabled").(bool)),
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("source_uri"); ok {
		properties.Properties.CreationData.SourceUri = utils.String(v.(string))
	}

	if v, ok := d.GetOk("source_resource_id"); ok {
		properties.Properties.CreationData.SourceResourceId = utils.String(v.(string))
	}

	if v, ok := d.GetOk("storage_account_id"); ok {
		properties.Properties.CreationData.StorageAccountId = utils.String(v.(string))
	}

	if v, ok := d.GetOk("network_access_policy"); ok {
		properties.Properties.NetworkAccessPolicy = pointer.To(snapshots.NetworkAccessPolicy(v.(string)))
	}

	if v, ok := d.GetOk("disk_access_id"); ok {
		properties.Properties.DiskAccessId = utils.String(v.(string))
	}

	properties.Properties.PublicNetworkAccess = pointer.To(snapshots.PublicNetworkAccessEnabled)
	if !d.Get("public_network_access_enabled").(bool) {
		properties.Properties.PublicNetworkAccess = pointer.To(snapshots.PublicNetworkAccessDisabled)
	}

	diskSizeGB := d.Get("disk_size_gb").(int)
	if diskSizeGB > 0 {
		properties.Properties.DiskSizeGB = utils.Int64(int64(diskSizeGB))
	}

	properties.Properties.EncryptionSettingsCollection = expandSnapshotDiskEncryptionSettings(d.Get("encryption_settings").([]interface{}))

	if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSnapshotRead(d, meta)
}

func resourceSnapshotRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		d.Set("location", azure.NormalizeLocation(model.Location))

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

			if err := d.Set("encryption_settings", flattenSnapshotDiskEncryptionSettings(props.EncryptionSettingsCollection)); err != nil {
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

			incrementalEnabled := false
			if props.Incremental != nil {
				incrementalEnabled = *props.Incremental
			}
			d.Set("incremental_enabled", incrementalEnabled)

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

func resourceSnapshotDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SnapshotsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := snapshots.ParseSnapshotID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
