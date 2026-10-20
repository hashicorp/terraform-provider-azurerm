// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity

func resourceDataProtectionBackupVault() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceDataProtectionBackupVaultCreateUpdate,
		Read:   resourceDataProtectionBackupVaultRead,
		Update: resourceDataProtectionBackupVaultCreateUpdate,
		Delete: resourceDataProtectionBackupVaultDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DataProtectionBackupVaultV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&backupvaultresources.BackupVaultId{}),
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&backupvaultresources.BackupVaultId{}),
		},

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

			"storage_setting": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"datastore_type": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice(backupvaultresources.PossibleValuesForStorageSettingStoreTypes(), false),
						},
						"redundancy": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice(backupvaultresources.PossibleValuesForStorageSettingTypes(), false),
						},
					},
				},
			},

			"cross_region_restore_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"retention_duration_in_days": {
				Type:         pluginsdk.TypeFloat,
				Optional:     true,
				Default:      14,
				ValidateFunc: validation.FloatBetween(14, 180),
			},

			"soft_delete": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      backupvaultresources.SoftDeleteStateOn,
				ValidateFunc: validation.StringInSlice(backupvaultresources.PossibleValuesForSoftDeleteState(), false),
			},

			"immutability": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      backupvaultresources.ImmutabilityStateDisabled,
				ValidateFunc: validation.StringInSlice(backupvaultresources.PossibleValuesForImmutabilityState(), false),
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(

			// Once `cross_region_restore_enabled` is enabled it cannot be disabled.
			pluginsdk.ForceNewIfChange("cross_region_restore_enabled", func(ctx context.Context, old, new, meta any) bool {
				return old.(bool) && new.(bool) != old.(bool)
			}),

			// Once `immutability` is enabled it cannot be disabled.
			pluginsdk.ForceNewIfChange("immutability", func(ctx context.Context, old, new, meta any) bool {
				return old.(string) == string(backupvaultresources.ImmutabilityStateLocked) && new.(string) != string(backupvaultresources.ImmutabilityStateLocked)
			}),

			pluginsdk.ForceNewIfChange("soft_delete", func(ctx context.Context, old, new, meta any) bool {
				return old.(string) == string(backupvaultresources.SoftDeleteStateAlwaysOn) && new.(string) != string(backupvaultresources.SoftDeleteStateAlwaysOn)
			}),

			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
				crossRegionRestore := d.GetRawConfig().AsValueMap()["cross_region_restore_enabled"]
				if !features.SixPointOh() {
					if !crossRegionRestore.IsNull() && d.Get("redundancy").(string) != string(backupvaultresources.StorageSettingTypesGeoRedundant) {
						return fmt.Errorf("`cross_region_restore_enabled` can only be specified when `redundancy` is specified for `GeoRedundant`")
					}
					return nil
				}

				if !d.GetRawConfig().GetAttr("storage_setting").IsWhollyKnown() {
					return nil
				}
				var redundancy string
				hasArchiveStore, hasVaultStore := false, false
				for _, raw := range d.Get("storage_setting").([]any) {
					setting := raw.(map[string]any)
					datastoreType := setting["datastore_type"].(string)
					if strings.EqualFold(datastoreType, string(backupvaultresources.StorageSettingStoreTypesArchiveStore)) {
						hasArchiveStore = true
					}
					if strings.EqualFold(datastoreType, string(backupvaultresources.StorageSettingStoreTypesVaultStore)) {
						hasVaultStore = true
						redundancy = setting["redundancy"].(string)
					} else if redundancy == "" {
						redundancy = setting["redundancy"].(string)
					}
				}
				if hasArchiveStore && !hasVaultStore {
					return fmt.Errorf("`storage_setting` must include `VaultStore` when `ArchiveStore` is specified")
				}
				if !crossRegionRestore.IsNull() && redundancy != string(backupvaultresources.StorageSettingTypesGeoRedundant) {
					// Cross region restore is only allowed on `GeoRedundant` vault.
					return fmt.Errorf("`cross_region_restore_enabled` can only be specified when `redundancy` is specified for `GeoRedundant`")
				}
				return nil
			}),
		),
	}

	if !features.SixPointOh() {
		resource.SchemaVersion = 0
		resource.StateUpgraders = nil
		delete(resource.Schema, "storage_setting")
		resource.Schema["datastore_type"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(backupvaultresources.PossibleValuesForStorageSettingStoreTypes(), false),
		}
		resource.Schema["redundancy"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(backupvaultresources.PossibleValuesForStorageSettingTypes(), false),
		}
	}

	return resource
}

func resourceDataProtectionBackupVaultCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := backupvaultresources.NewBackupVaultID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.BackupVaultsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing DataProtection BackupVault (%q): %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_data_protection_backup_vault", id.ID())
			}
		}
	}

	expandedIdentity, err := expandBackupVaultDppIdentityDetails(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	parameters := backupvaultresources.BackupVaultResource{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: backupvaultresources.BackupVault{
			SecuritySettings: &backupvaultresources.SecuritySettings{
				SoftDeleteSettings: &backupvaultresources.SoftDeleteSettings{
					State: pointer.ToEnum[backupvaultresources.SoftDeleteState](d.Get("soft_delete").(string)),
				},
				ImmutabilitySettings: &backupvaultresources.ImmutabilitySettings{
					State: pointer.ToEnum[backupvaultresources.ImmutabilityState](d.Get("immutability").(string)),
				},
			},
		},
		Identity: expandedIdentity,
		Tags:     expandTags(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("storage_setting"); ok {
		parameters.Properties.StorageSettings = expandBackupVaultStorageSettings(v.([]any))
	}
	if !features.SixPointOh() {
		parameters.Properties.StorageSettings = []backupvaultresources.StorageSetting{
			{
				DatastoreType: pointer.ToEnum[backupvaultresources.StorageSettingStoreTypes](d.Get("datastore_type").(string)),
				Type:          pointer.ToEnum[backupvaultresources.StorageSettingTypes](d.Get("redundancy").(string)),
			},
		}
		// Azure requires a VaultStore alongside ArchiveStore; use the configured redundancy for both.
		if strings.EqualFold(d.Get("datastore_type").(string), string(backupvaultresources.StorageSettingStoreTypesArchiveStore)) {
			parameters.Properties.StorageSettings = append(parameters.Properties.StorageSettings, backupvaultresources.StorageSetting{
				DatastoreType: pointer.To(backupvaultresources.StorageSettingStoreTypesVaultStore),
				Type:          pointer.ToEnum[backupvaultresources.StorageSettingTypes](d.Get("redundancy").(string)),
			})
		}
	}

	if !pluginsdk.IsExplicitlyNullInConfig(d, "cross_region_restore_enabled") {
		parameters.Properties.FeatureSettings = &backupvaultresources.FeatureSettings{
			CrossRegionRestoreSettings: &backupvaultresources.CrossRegionRestoreSettings{},
		}
		if d.Get("cross_region_restore_enabled").(bool) {
			parameters.Properties.FeatureSettings.CrossRegionRestoreSettings.State = pointer.To(backupvaultresources.CrossRegionRestoreStateEnabled)
		} else {
			parameters.Properties.FeatureSettings.CrossRegionRestoreSettings.State = pointer.To(backupvaultresources.CrossRegionRestoreStateDisabled)
		}
	}

	if v, ok := d.GetOk("retention_duration_in_days"); ok {
		parameters.Properties.SecuritySettings.SoftDeleteSettings.RetentionDurationInDays = pointer.To(v.(float64))
	}

	if d.IsNewResource() {
		if err := client.BackupVaultsCreateOrUpdateCallbackThenPoll(ctx, id, parameters, backupvaultresources.DefaultBackupVaultsCreateOrUpdateOperationOptions(), sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	} else {
		if err := client.BackupVaultsCreateOrUpdateThenPoll(ctx, id, parameters, backupvaultresources.DefaultBackupVaultsCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceDataProtectionBackupVaultRead(d, meta)
}

func resourceDataProtectionBackupVaultRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backupvaultresources.ParseBackupVaultID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.BackupVaultsGet(ctx, *id)
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
		d.Set("location", location.NormalizeNilable(pointer.To(model.Location)))
		props := model.Properties

		immutability := backupvaultresources.ImmutabilityStateDisabled
		if securitySetting := model.Properties.SecuritySettings; securitySetting != nil {
			if immutabilitySettings := securitySetting.ImmutabilitySettings; immutabilitySettings != nil {
				if immutabilitySettings.State != nil {
					immutability = *immutabilitySettings.State
				}
			}
			if softDelete := securitySetting.SoftDeleteSettings; softDelete != nil {
				d.Set("soft_delete", pointer.FromEnum(softDelete.State))
				d.Set("retention_duration_in_days", pointer.From(softDelete.RetentionDurationInDays))
			}
		}
		d.Set("immutability", string(immutability))

		crossRegionStoreEnabled := false
		if featureSetting := model.Properties.FeatureSettings; featureSetting != nil {
			if crossRegionRestore := featureSetting.CrossRegionRestoreSettings; crossRegionRestore != nil {
				if pointer.From(crossRegionRestore.State) == backupvaultresources.CrossRegionRestoreStateEnabled {
					crossRegionStoreEnabled = true
				}
			}
		}
		d.Set("cross_region_restore_enabled", crossRegionStoreEnabled)

		identity, err := flattenBackupVaultDppIdentityDetails(model.Identity)
		if err != nil {
			return err
		}
		d.Set("identity", identity)

		if err = tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}

		if !features.SixPointOh() {
			datastoreType, redundancy := flattenBackupVaultStorageSettingsLegacy(props.StorageSettings, d.Get("datastore_type").(string))
			d.Set("datastore_type", datastoreType)
			d.Set("redundancy", redundancy)
			return pluginsdk.SetResourceIdentityData(d, id)
		}

		if err := d.Set("storage_setting", flattenBackupVaultStorageSettings(props.StorageSettings)); err != nil {
			return fmt.Errorf("setting `storage_setting`: %+v", err)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourceDataProtectionBackupVaultDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backupvaultresources.ParseBackupVaultID(d.Id())
	if err != nil {
		return err
	}

	if err := client.BackupVaultsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting DataProtection BackupVault (%q): %+v", id, err)
	}

	// API has bug, which appears API returns before the resource is fully deleted. Tracked by this issue: https://github.com/Azure/azure-rest-api-specs/issues/38944
	pollerType := custompollers.NewDataProtectionBackupVaultPoller(client, *id)
	poller := pollers.NewPoller(pollerType, 30*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return err
	}

	return nil
}

func expandBackupVaultStorageSettings(input []any) []backupvaultresources.StorageSetting {
	result := make([]backupvaultresources.StorageSetting, 0, len(input))
	for _, raw := range input {
		setting := raw.(map[string]any)
		result = append(result, backupvaultresources.StorageSetting{
			DatastoreType: pointer.ToEnum[backupvaultresources.StorageSettingStoreTypes](setting["datastore_type"].(string)),
			Type:          pointer.ToEnum[backupvaultresources.StorageSettingTypes](setting["redundancy"].(string)),
		})
	}
	return result
}

func flattenBackupVaultStorageSettings(input []backupvaultresources.StorageSetting) []any {
	result := make([]any, 0, len(input))
	for _, setting := range input {
		result = append(result, map[string]any{
			"datastore_type": pointer.FromEnum(setting.DatastoreType),
			"redundancy":     pointer.FromEnum(setting.Type),
		})
	}
	return result
}

func flattenBackupVaultStorageSettingsLegacy(input []backupvaultresources.StorageSetting, currentDatastoreType string) (string, string) {
	if len(input) == 0 {
		return "", ""
	}

	setting := input[0]
	// Preserve an existing store selection when it is still present in Azure.
	// Without prior state, prefer ArchiveStore over its required VaultStore companion.
	for _, candidate := range input {
		if currentDatastoreType != "" && pointer.FromEnum(candidate.DatastoreType) == currentDatastoreType {
			return pointer.FromEnum(candidate.DatastoreType), pointer.FromEnum(candidate.Type)
		}
		if pointer.From(candidate.DatastoreType) == backupvaultresources.StorageSettingStoreTypesArchiveStore {
			setting = candidate
		}
	}

	return pointer.FromEnum(setting.DatastoreType), pointer.FromEnum(setting.Type)
}

func expandBackupVaultDppIdentityDetails(input []any) (*backupvaultresources.DppIdentityDetails, error) {
	config, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	identity := backupvaultresources.DppIdentityDetails{
		Type: pointer.To(string(config.Type)),
	}

	if len(config.IdentityIds) > 0 {
		identityIds := make(map[string]backupvaultresources.UserAssignedIdentity, len(config.IdentityIds))
		for id := range config.IdentityIds {
			identityIds[id] = backupvaultresources.UserAssignedIdentity{}
		}
		identity.UserAssignedIdentities = pointer.To(identityIds)
	}

	return &identity, nil
}

func flattenBackupVaultDppIdentityDetails(input *backupvaultresources.DppIdentityDetails) (*[]any, error) {
	var config *identity.SystemAndUserAssignedMap
	if input != nil {
		config = &identity.SystemAndUserAssignedMap{
			Type: identity.Type(*input.Type),
		}

		config.PrincipalId = pointer.From(input.PrincipalId)
		config.TenantId = pointer.From(input.TenantId)

		if len(pointer.From(input.UserAssignedIdentities)) > 0 {
			config.IdentityIds = make(map[string]identity.UserAssignedIdentityDetails, len(pointer.From(input.UserAssignedIdentities)))
			for k, v := range *input.UserAssignedIdentities {
				config.IdentityIds[k] = identity.UserAssignedIdentityDetails{
					ClientId:    v.ClientId,
					PrincipalId: v.PrincipalId,
				}
			}
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(config)
}
