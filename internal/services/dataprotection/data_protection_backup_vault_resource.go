// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backupvaults"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataProtectionBackupVault() *pluginsdk.Resource {
	return &pluginsdk.Resource{
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

			"datastore_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(backupvaults.StorageSettingStoreTypesArchiveStore),
					string(backupvaults.StorageSettingStoreTypesOperationalStore),
					string(backupvaults.StorageSettingStoreTypesVaultStore),
				}, false),
			},

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

			"immutability": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      backupvaults.ImmutabilityStateDisabled,
				ValidateFunc: validation.StringInSlice(backupvaults.PossibleValuesForImmutabilityState(), false),
			},

			"encryption_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"identity_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: commonids.ValidateUserAssignedIdentityID,
						},

						"infrastructure_encryption_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							ForceNew: true,
							RequiredWith: []string{
								"encryption_settings.0.identity_id",
								"encryption_settings.0.key_vault_key_id",
							},
						},

						"key_vault_key_id": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							RequiredWith: []string{
								"encryption_settings.0.identity_id",
							},
							ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
						},
					},
				},
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(

			// Once `cross_region_restore_enabled` is enabled it cannot be disabled.
			pluginsdk.ForceNewIfChange("cross_region_restore_enabled", func(ctx context.Context, old, new, meta interface{}) bool {
				return old.(bool) && new.(bool) != old.(bool)
			}),

			// Once `immutability` is enabled it cannot be disabled.
			pluginsdk.ForceNewIfChange("immutability", func(ctx context.Context, old, new, meta interface{}) bool {
				return old.(string) == string(backupvaults.ImmutabilityStateLocked) && new.(string) != string(backupvaults.ImmutabilityStateLocked)
			}),

			pluginsdk.ForceNewIfChange("soft_delete", func(ctx context.Context, old, new, meta interface{}) bool {
				return old.(string) == string(backupvaults.SoftDeleteStateAlwaysOn) && new.(string) != string(backupvaults.SoftDeleteStateAlwaysOn)
			}),

			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
				redundancy := d.Get("redundancy").(string)
				crossRegionRestore := d.GetRawConfig().AsValueMap()["cross_region_restore_enabled"]
				if !crossRegionRestore.IsNull() && redundancy != string(backupvaults.StorageSettingTypesGeoRedundant) {
					// Cross region restore is only allowed on `GeoRedundant` vault.
					return fmt.Errorf("`cross_region_restore_enabled` can only be specified when `redundancy` is specified for `GeoRedundant`")
				}
				return nil
			}),
		),
	}
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

	parameters := backupvaults.BackupVaultResource{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: backupvaults.BackupVault{
			StorageSettings: []backupvaults.StorageSetting{
				{
					DatastoreType: pointer.To(backupvaults.StorageSettingStoreTypes(d.Get("datastore_type").(string))),
					Type:          pointer.To(backupvaults.StorageSettingTypes(d.Get("redundancy").(string))),
				},
			},
			SecuritySettings: &backupvaults.SecuritySettings{
				SoftDeleteSettings: &backupvaults.SoftDeleteSettings{
					State: pointer.To(backupvaults.SoftDeleteState(d.Get("soft_delete").(string))),
				},
				ImmutabilitySettings: &backupvaults.ImmutabilitySettings{
					State: pointer.To(backupvaults.ImmutabilityState(d.Get("immutability").(string))),
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

	if v, ok := d.GetOk("encryption_settings"); ok {
		encryptionSettings, err := expandBackupVaultEncryptionSettings(v.([]interface{}))

		if err != nil {
			return err
		}

		parameters.Properties.SecuritySettings.EncryptionSettings = encryptionSettings
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

		if len(props.StorageSettings) > 0 {
			d.Set("datastore_type", string(pointer.From((props.StorageSettings)[0].DatastoreType)))
			d.Set("redundancy", string(pointer.From((props.StorageSettings)[0].Type)))
		}

		immutability := backupvaults.ImmutabilityStateDisabled
		if securitySetting := model.Properties.SecuritySettings; securitySetting != nil {
			if immutabilitySettings := securitySetting.ImmutabilitySettings; immutabilitySettings != nil {
				if immutabilitySettings.State != nil {
					immutability = *immutabilitySettings.State
				}
			}
			if softDelete := securitySetting.SoftDeleteSettings; softDelete != nil {
				d.Set("soft_delete", string(pointer.From(softDelete.State)))
				d.Set("retention_duration_in_days", pointer.From(softDelete.RetentionDurationInDays))
			}
			if securitySetting.EncryptionSettings != nil {
				d.Set("encryption_settings", *flattenBackupVaultEncryptionSettings(securitySetting.EncryptionSettings))
			}
		}
		d.Set("immutability", string(immutability))

		crossRegionStoreEnabled := false
		if featureSetting := model.Properties.FeatureSettings; featureSetting != nil {
			if crossRegionRestore := featureSetting.CrossRegionRestoreSettings; crossRegionRestore != nil {
				if pointer.From(crossRegionRestore.State) == backupvaults.CrossRegionRestoreStateEnabled {
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

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting DataProtection BackupVault (%q): %+v", id, err)
	}
	return nil
}

func expandBackupVaultDppIdentityDetails(input []interface{}) (*backupvaults.DppIdentityDetails, error) {
	config, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	identity := backupvaults.DppIdentityDetails{
		Type: pointer.To(string(config.Type)),
	}

	if len(config.IdentityIds) > 0 {
		identityIds := make(map[string]backupvaults.UserAssignedIdentity, len(config.IdentityIds))
		for id := range config.IdentityIds {
			identityIds[id] = backupvaults.UserAssignedIdentity{}
		}
		identity.UserAssignedIdentities = pointer.To(identityIds)
	}

	return &identity, nil
}

func flattenBackupVaultDppIdentityDetails(input *backupvaults.DppIdentityDetails) (*[]interface{}, error) {
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

func expandBackupVaultEncryptionSettings(input []interface{}) (*backupvaults.EncryptionSettings, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	v := input[0].(map[string]interface{})
	output := &backupvaults.EncryptionSettings{
		KekIdentity: &backupvaults.CmkKekIdentity{},
	}

	if v["identity_id"].(string) != "" {
		output.KekIdentity.IdentityId = pointer.To(v["identity_id"].(string))
		output.KekIdentity.IdentityType = pointer.To(backupvaults.IdentityTypeUserAssigned)
		output.State = pointer.To(backupvaults.EncryptionStateEnabled)

		if v["infrastructure_encryption_enabled"].(bool) {
			output.InfrastructureEncryption = pointer.To(backupvaults.InfrastructureEncryptionStateEnabled)
		} else {
			output.InfrastructureEncryption = pointer.To(backupvaults.InfrastructureEncryptionStateDisabled)
		}

		if v["key_vault_key_id"].(string) != "" {
			keyId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(v["key_vault_key_id"].(string))

			if err != nil {
				return nil, err
			}

			output.KeyVaultProperties = &backupvaults.CmkKeyVaultProperties{
				KeyUri: pointer.To(keyId.ID()),
			}
		}
	}

	return output, nil
}

func flattenBackupVaultEncryptionSettings(input *backupvaults.EncryptionSettings) *[]interface{} {
	output := make(map[string]interface{})

	if input.InfrastructureEncryption != nil {
		output["infrastructure_encryption_enabled"] = pointer.From(input.InfrastructureEncryption) == backupvaults.InfrastructureEncryptionStateEnabled
	}

	if input.KekIdentity != nil && input.KekIdentity.IdentityId != nil {
		output["identity_id"] = pointer.From(input.KekIdentity.IdentityId)
	}

	if input.KeyVaultProperties != nil && input.KeyVaultProperties.KeyUri != nil {
		output["key_vault_key_id"] = pointer.From(input.KeyVaultProperties.KeyUri)
	}

	return &[]interface{}{output}
}