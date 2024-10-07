// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backupvaults"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataProtectionBackupVault() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceDataProtectionBackupVaultCreateUpdate,
		Read:   resourceDataProtectionBackupVaultRead,
		Update: resourceDataProtectionBackupVaultCreateUpdate,
		Delete: resourceDataProtectionBackupVaultDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := backupvaults.ParseBackupVaultIDInsensitively(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{2,50}$"),
					"DataProtection BackupVault name must be 2 - 50 characters long, contain only letters, numbers and hyphens.).",
				),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"redundancy": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(backupvaults.StorageSettingTypesGeoRedundant),
					string(backupvaults.StorageSettingTypesLocallyRedundant),
					string(backupvaults.StorageSettingTypesZoneRedundant),
				}, false),
			},

			"cross_region_restore_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"identity": commonschema.SystemAssignedIdentityOptional(),

			"retention_duration_in_days": {
				Type:         pluginsdk.TypeFloat,
				Optional:     true,
				Default:      14,
				ValidateFunc: validation.FloatBetween(14, 180),
			},

			"soft_delete": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      backupvaults.SoftDeleteStateOn,
				ValidateFunc: validation.StringInSlice(backupvaults.PossibleValuesForSoftDeleteState(), false),
			},

			"tags": tags.Schema(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("soft_delete", func(ctx context.Context, old, new, meta interface{}) bool {
				return old.(string) == string(backupvaults.SoftDeleteStateAlwaysOn) && new.(string) != string(backupvaults.SoftDeleteStateAlwaysOn)
			}),

			// Once `cross_region_restore_enabled` is enabled it cannot be disabled.
			pluginsdk.ForceNewIfChange("cross_region_restore_enabled", func(ctx context.Context, old, new, meta interface{}) bool {
				return old.(bool) && new.(bool) != old.(bool)
			}),

			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
				redundancy := d.Get("redundancy").(string)
				crossRegionRestore := d.GetRawConfig().AsValueMap()["cross_region_restore_enabled"]
				if !crossRegionRestore.IsNull() && redundancy != string(backupvaults.StorageSettingTypesGeoRedundant) {
					// Cross region restore is only allowed on `GeoRedundant` vault.
					return fmt.Errorf("`cross_region_restore_enabled` can only be specified when `redundancy` is specified for `GeoRedundant`.")
				}
				return nil
			}),
		),
	}

	// Confirmed with the service team that `SnapshotStore` has been replaced with `OperationalStore`.
	if !features.FourPointOhBeta() {
		resource.Schema["datastore_type"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(backupvaults.StorageSettingStoreTypesArchiveStore),
				"SnapshotStore",
				string(backupvaults.StorageSettingStoreTypesOperationalStore),
				string(backupvaults.StorageSettingStoreTypesVaultStore),
			}, false),
		}
	} else {
		resource.Schema["datastore_type"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(backupvaults.StorageSettingStoreTypesArchiveStore),
				string(backupvaults.StorageSettingStoreTypesOperationalStore),
				string(backupvaults.StorageSettingStoreTypesVaultStore),
			}, false),
		}
	}
	return resource
}

func resourceDataProtectionBackupVaultCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := backupvaults.NewBackupVaultID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing DataProtection BackupVault (%q): %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_data_protection_backup_vault", id.ID())
		}
	}

	expandedIdentity, err := expandBackupVaultDppIdentityDetails(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	datastoreType := backupvaults.StorageSettingStoreTypes(d.Get("datastore_type").(string))
	storageSettingType := backupvaults.StorageSettingTypes(d.Get("redundancy").(string))

	parameters := backupvaults.BackupVaultResource{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: backupvaults.BackupVault{
			StorageSettings: []backupvaults.StorageSetting{
				{
					DatastoreType: &datastoreType,
					Type:          &storageSettingType,
				},
			},
			SecuritySettings: &backupvaults.SecuritySettings{
				SoftDeleteSettings: &backupvaults.SoftDeleteSettings{
					State: pointer.To(backupvaults.SoftDeleteState(d.Get("soft_delete").(string))),
				},
			},
		},
		Identity: expandedIdentity,
		Tags:     expandTags(d.Get("tags").(map[string]interface{})),
	}

	if !pluginsdk.IsExplicitlyNullInConfig(d, "cross_region_restore_enabled") {
		parameters.Properties.FeatureSettings = &backupvaults.FeatureSettings{
			CrossRegionRestoreSettings: &backupvaults.CrossRegionRestoreSettings{},
		}
		if d.Get("cross_region_restore_enabled").(bool) {
			parameters.Properties.FeatureSettings.CrossRegionRestoreSettings.State = pointer.To(backupvaults.CrossRegionRestoreStateEnabled)
		} else {
			parameters.Properties.FeatureSettings.CrossRegionRestoreSettings.State = pointer.To(backupvaults.CrossRegionRestoreStateDisabled)
		}
	}

	if v, ok := d.GetOk("retention_duration_in_days"); ok {
		parameters.Properties.SecuritySettings.SoftDeleteSettings.RetentionDurationInDays = pointer.To(v.(float64))
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, parameters, backupvaults.DefaultCreateOrUpdateOperationOptions())
	if err != nil {
		return fmt.Errorf("creating DataProtection BackupVault (%q): %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDataProtectionBackupVaultRead(d, meta)
}

func resourceDataProtectionBackupVaultRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backupvaults.ParseBackupVaultID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] DataProtection BackupVault %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataProtection BackupVault (%q): %+v", id, err)
	}
	d.Set("name", id.BackupVaultName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		props := model.Properties
		if props.StorageSettings != nil && len(props.StorageSettings) > 0 {
			d.Set("datastore_type", string(pointer.From((props.StorageSettings)[0].DatastoreType)))
			d.Set("redundancy", string(pointer.From((props.StorageSettings)[0].Type)))
		}
		if securitySetting := model.Properties.SecuritySettings; securitySetting != nil {
			if softDelete := securitySetting.SoftDeleteSettings; softDelete != nil {
				d.Set("soft_delete", string(pointer.From(softDelete.State)))
				d.Set("retention_duration_in_days", pointer.From(softDelete.RetentionDurationInDays))
			}
		}
		crossRegionStoreEnabled := false
		if featureSetting := model.Properties.FeatureSettings; featureSetting != nil {
			if featureSetting := model.Properties.FeatureSettings; featureSetting != nil {
				if pointer.From(featureSetting.CrossRegionRestoreSettings.State) == backupvaults.CrossRegionRestoreStateEnabled {
					crossRegionStoreEnabled = true
				}
			}
		}

		d.Set("cross_region_restore_enabled", crossRegionStoreEnabled)

		if err = d.Set("identity", flattenBackupVaultDppIdentityDetails(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}
		if err = tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}

	return nil
}

func resourceDataProtectionBackupVaultDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backupvaults.ParseBackupVaultID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id); err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}
		return fmt.Errorf("deleting DataProtection BackupVault (%q): %+v", id, err)
	}
	return nil
}

func expandBackupVaultDppIdentityDetails(input []interface{}) (*backupvaults.DppIdentityDetails, error) {
	config, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	return &backupvaults.DppIdentityDetails{
		Type: utils.String(string(config.Type)),
	}, nil
}

func flattenBackupVaultDppIdentityDetails(input *backupvaults.DppIdentityDetails) []interface{} {
	var config *identity.SystemAssigned
	if input != nil {
		principalId := ""
		if input.PrincipalId != nil {
			principalId = *input.PrincipalId
		}

		tenantId := ""
		if input.TenantId != nil {
			tenantId = *input.TenantId
		}
		config = &identity.SystemAssigned{
			Type:        identity.Type(*input.Type),
			PrincipalId: principalId,
			TenantId:    tenantId,
		}
	}
	return identity.FlattenSystemAssigned(config)
}
