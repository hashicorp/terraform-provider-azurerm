package recoveryservices

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-08-01/recoveryservices"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-12-01/backup"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2022-03-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyvaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceRecoveryServicesVault() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRecoveryServicesVaultCreate,
		Read:   resourceRecoveryServicesVaultRead,
		Update: resourceRecoveryServicesVaultUpdate,
		Delete: resourceRecoveryServicesVaultDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := vaults.ParseVaultID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RecoveryServicesVaultName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"encryption": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				RequiredWith: []string{"identity"},
				MaxItems:     1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: keyvaultValidate.NestedItemIdWithOptionalVersion,
						},
						"infrastructure_encryption_enabled": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
						"use_system_assigned_identity": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},

			"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

			"tags": commonschema.Tags(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(recoveryservices.SkuNameRS0),
					string(recoveryservices.SkuNameStandard),
				}, false),
			},

			"storage_mode_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  backup.StorageTypeGeoRedundant,
				ValidateFunc: validation.StringInSlice([]string{
					string(backup.StorageTypeGeoRedundant),
					string(backup.StorageTypeLocallyRedundant),
					string(backup.StorageTypeZoneRedundant),
				}, false),
			},

			"cross_region_restore_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"soft_delete_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"monitor_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"azure_monitor_alert_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},
						"classic_alert_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
		},
	}
}

func resourceRecoveryServicesVaultCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.VaultsClient
	cfgsClient := meta.(*clients.Client).RecoveryServices.VaultsConfigsClient
	storageCfgsClient := meta.(*clients.Client).RecoveryServices.StorageConfigsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := vaults.NewVaultID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	storageMode := d.Get("storage_mode_type").(string)
	crossRegionRestore := d.Get("cross_region_restore_enabled").(bool)

	if crossRegionRestore && storageMode != string(backup.StorageTypeGeoRedundant) {
		return fmt.Errorf("cannot enable cross region restore when storage mode type is not %s. %s", string(backup.StorageTypeGeoRedundant), id.String())
	}

	location := d.Get("location").(string)
	t := d.Get("tags").(map[string]interface{})

	log.Printf("[DEBUG] Creating Recovery Service %s", id.String())

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing Recovery Service %s: %+v", id.String(), err)
		}
	}
	if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
		return tf.ImportAsExistsError("azurerm_recovery_services_vault", *existing.Model.Id)
	}

	expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	sku := d.Get("sku").(string)
	vault := vaults.Vault{
		Location: location,
		Tags:     tags.Expand(t),
		Identity: expandedIdentity,
		Sku: &vaults.Sku{
			Name: vaults.SkuName(sku),
		},
		Properties: &vaults.VaultProperties{},
	}

	if vaults.SkuName(sku) == vaults.SkuNameRSZero {
		vault.Sku.Tier = utils.String("Standard")
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, vault)
	if err != nil {
		return fmt.Errorf("creating Recovery Service %s: %+v", id.String(), err)
	}

	cfg := backup.ResourceVaultConfigResource{
		Properties: &backup.ResourceVaultConfig{
			EnhancedSecurityState: backup.EnhancedSecurityStateEnabled, // always enabled
		},
	}

	if sd := d.Get("soft_delete_enabled").(bool); sd {
		cfg.Properties.SoftDeleteFeatureState = backup.SoftDeleteFeatureStateEnabled
	} else {
		cfg.Properties.SoftDeleteFeatureState = backup.SoftDeleteFeatureStateDisabled
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"syncing"},
		Target:     []string{"success"},
		MinTimeout: 30 * time.Second,
		Refresh: func() (interface{}, string, error) {
			resp, err := cfgsClient.Update(ctx, id.VaultName, id.ResourceGroupName, cfg)
			if err != nil {
				if strings.Contains(err.Error(), "ResourceNotYetSynced") {
					return resp, "syncing", nil
				}
				return resp, "error", fmt.Errorf("updating Recovery Service Vault Cfg %s: %+v", id.String(), err)
			}

			return resp, "success", nil
		},
	}

	stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for on update for Recovery Service  %s: %+v", id.String(), err)
	}

	storageCfg := backup.ResourceConfigResource{
		Properties: &backup.ResourceConfig{
			StorageModelType:       backup.StorageType(d.Get("storage_mode_type").(string)),
			CrossRegionRestoreFlag: utils.Bool(d.Get("cross_region_restore_enabled").(bool)),
		},
	}

	err = pluginsdk.Retry(stateConf.Timeout, func() *pluginsdk.RetryError {
		if resp, err := storageCfgsClient.Update(ctx, id.VaultName, id.ResourceGroupName, storageCfg); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return pluginsdk.RetryableError(fmt.Errorf("updating Recovery Service Storage Cfg %s: %+v", id.String(), err))
			}
			if utils.ResponseWasBadRequest(resp.Response) {
				return pluginsdk.RetryableError(fmt.Errorf("updating Recovery Service Storage Cfg %s: %+v", id.String(), err))
			}

			return pluginsdk.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// storage type is not updated instantaneously, so we wait until storage type is correct
	err = pluginsdk.Retry(stateConf.Timeout, func() *pluginsdk.RetryError {
		if resp, err := storageCfgsClient.Get(ctx, id.VaultName, id.ResourceGroupName); err == nil {
			if resp.Properties == nil {
				return pluginsdk.NonRetryableError(fmt.Errorf("updating %s Storage Config: `properties` was nil", id))
			}
			if resp.Properties.StorageType != storageCfg.Properties.StorageModelType {
				return pluginsdk.RetryableError(fmt.Errorf("updating Storage Config: %+v", err))
			}
			if *resp.Properties.CrossRegionRestoreFlag != *storageCfg.Properties.CrossRegionRestoreFlag {
				return pluginsdk.RetryableError(fmt.Errorf("updating Storage Config: %+v", err))
			}
		} else {
			return pluginsdk.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// recovery vault's encryption config and monitor settings cannot be set while creation, so a standalone update is required.
	needUpdate := false
	prop := vaults.VaultProperties{}

	if _, ok := d.GetOk("encryption"); ok {
		prop.Encryption = expandEncryption(d)
		needUpdate = true
	}

	if _, ok := d.GetOk("monitor_settings"); ok {
		prop.MonitoringSettings = expandMonitorSettings(d)
		needUpdate = true
	}

	if needUpdate {
		err := client.UpdateThenPoll(ctx, id, vaults.PatchVault{
			Properties: &prop,
		})
		if err != nil {
			return fmt.Errorf("updating Recovery Service Encryption %s: %+v, but recovery vault was created, a manually import might be required", id.String(), err)
		}
	}

	d.SetId(id.ID())
	return resourceRecoveryServicesVaultRead(d, meta)
}

func resourceRecoveryServicesVaultUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.VaultsClient
	cfgsClient := meta.(*clients.Client).RecoveryServices.VaultsConfigsClient
	storageCfgsClient := meta.(*clients.Client).RecoveryServices.StorageConfigsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := vaults.NewVaultID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	encryption := expandEncryption(d)
	existing, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("checking for presence of existing Recovery Service %s: %+v", id.String(), err)
	}
	model := existing.Model
	if model == nil {
		return fmt.Errorf("checking for presence of existing Recovery Service %s: `model` was nil", id.String())
	}
	if model.Properties != nil && model.Properties.Encryption != nil {
		if encryption == nil {
			return fmt.Errorf("once encryption with your own key has been enabled it's not possible to disable it")
		}
		if encryption.InfrastructureEncryption != model.Properties.Encryption.InfrastructureEncryption {
			return fmt.Errorf("once `infrastructure_encryption_enabled` has been set it's not possible to change it")
		}
		if d.HasChange("sku") {
			// Once encryption has been enabled, calling `CreateOrUpdate` without it is not allowed.
			// But `sku` can only be updated by `CreateOrUpdate` and the support for `encryption` in `CreateOrUpdate` is still under preview (https://docs.microsoft.com/azure/backup/encryption-at-rest-with-cmk?tabs=portal#enable-encryption-using-customer-managed-keys-at-vault-creation-in-preview).
			// TODO remove this restriction and add `encryption` to below `sku` update block when `encryption` in `CreateOrUpdate` is GA
			return fmt.Errorf("`sku` cannot be changed when encryption with your own key has been enabled")
		}
	}

	storageMode := d.Get("storage_mode_type").(string)
	crossRegionRestore := d.Get("cross_region_restore_enabled").(bool)

	if crossRegionRestore && storageMode != string(backup.StorageTypeGeoRedundant) {
		return fmt.Errorf("cannot enable cross region restore when storage mode type is not %s. %s", string(backup.StorageTypeGeoRedundant), id.String())
	}

	expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	cfg := backup.ResourceVaultConfigResource{
		Properties: &backup.ResourceVaultConfig{
			EnhancedSecurityState: backup.EnhancedSecurityStateEnabled, // always enabled
		},
	}

	if d.HasChange("soft_delete_enabled") {
		if sd := d.Get("soft_delete_enabled").(bool); sd {
			cfg.Properties.SoftDeleteFeatureState = backup.SoftDeleteFeatureStateEnabled
		} else {
			cfg.Properties.SoftDeleteFeatureState = backup.SoftDeleteFeatureStateDisabled
		}

		stateConf := &pluginsdk.StateChangeConf{
			Pending:    []string{"syncing"},
			Target:     []string{"success"},
			MinTimeout: 30 * time.Second,
			Refresh: func() (interface{}, string, error) {
				resp, err := cfgsClient.Update(ctx, id.VaultName, id.ResourceGroupName, cfg)
				if err != nil {
					if strings.Contains(err.Error(), "ResourceNotYetSynced") {
						return resp, "syncing", nil
					}
					return resp, "error", fmt.Errorf("updating Recovery Service Vault Cfg %s: %+v", id.String(), err)
				}

				return resp, "success", nil
			},
		}

		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for on update for Recovery Service  %s: %+v", id.String(), err)
		}
	}

	if d.HasChanges("storage_mode_type", "cross_region_restore_enabled") {
		storageCfg := backup.ResourceConfigResource{
			Properties: &backup.ResourceConfig{
				StorageModelType:       backup.StorageType(storageMode),
				CrossRegionRestoreFlag: utils.Bool(crossRegionRestore),
			},
		}

		err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutUpdate), func() *pluginsdk.RetryError {
			if resp, err := storageCfgsClient.Update(ctx, id.VaultName, id.ResourceGroupName, storageCfg); err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return pluginsdk.RetryableError(fmt.Errorf("updating Recovery Service Storage Cfg %s: %+v", id.String(), err))
				}
				if utils.ResponseWasBadRequest(resp.Response) {
					return pluginsdk.RetryableError(fmt.Errorf("updating Recovery Service Storage Cfg %s: %+v", id.String(), err))
				}

				return pluginsdk.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}

		// storage type is not updated instantaneously, so we wait until storage type is correct
		err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutUpdate), func() *pluginsdk.RetryError {
			if resp, err := storageCfgsClient.Get(ctx, id.VaultName, id.ResourceGroupName); err == nil {
				if resp.Properties == nil {
					return pluginsdk.NonRetryableError(fmt.Errorf("updating %s Storage Config: `properties` was nil", id))
				}
				if resp.Properties.StorageType != storageCfg.Properties.StorageModelType {
					return pluginsdk.RetryableError(fmt.Errorf("updating Storage Config: %+v", err))
				}
				if *resp.Properties.CrossRegionRestoreFlag != *storageCfg.Properties.CrossRegionRestoreFlag {
					return pluginsdk.RetryableError(fmt.Errorf("updating Storage Config: %+v", err))
				}
			} else {
				return pluginsdk.NonRetryableError(err)
			}

			return nil
		})
		if err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	// `sku` can only be updated by `CreateOrUpdate` but not `Update`, so use `CreateOrUpdate` with required and unchangeable properties
	if d.HasChange("sku") {
		sku := d.Get("sku").(string)
		vault := vaults.Vault{
			Location: d.Get("location").(string),
			Identity: expandedIdentity,
			Sku: &vaults.Sku{
				Name: vaults.SkuName(sku),
			},
			Properties: &vaults.VaultProperties{},
		}

		if recoveryservices.SkuName(sku) == recoveryservices.SkuNameRS0 {
			vault.Sku.Tier = utils.String("Standard")
		}

		err := client.CreateOrUpdateThenPoll(ctx, id, vault)
		if err != nil {
			return fmt.Errorf("updating Recovery Service %s: %+v", id.String(), err)
		}
	}

	vault := vaults.PatchVault{}

	if d.HasChange("identity") {
		vault.Identity = expandedIdentity
	}

	if d.HasChange("encryption") {
		if vault.Properties == nil {
			vault.Properties = &vaults.VaultProperties{}
		}

		vault.Properties.Encryption = encryption
	}

	if d.HasChange("monitor_settings") {
		vault.Properties.MonitoringSettings = expandMonitorSettings(d)
	}

	if d.HasChange("tags") {
		vault.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	err = client.UpdateThenPoll(ctx, id, vault)
	if err != nil {
		return fmt.Errorf("updating Recovery Service Encryption %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceRecoveryServicesVaultRead(d, meta)
}

func resourceRecoveryServicesVaultRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.VaultsClient
	cfgsClient := meta.(*clients.Client).RecoveryServices.VaultsConfigsClient
	storageCfgsClient := meta.(*clients.Client).RecoveryServices.StorageConfigsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := vaults.ParseVaultID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading Recovery Service %s", id.String())

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Recovery Service %s: %+v", id.String(), err)
	}

	d.Set("name", id.VaultName)
	d.Set("resource_group_name", id.ResourceGroupName)

	model := resp.Model
	if model == nil {
		return fmt.Errorf("retrieving Recovery Service %s: `model` was nil", id)
	}

	d.Set("location", location.Normalize(model.Location))
	if sku := model.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	cfg, err := cfgsClient.Get(ctx, id.VaultName, id.ResourceGroupName)
	if err != nil {
		return fmt.Errorf("reading Recovery Service Vault Cfg %s: %+v", id.String(), err)
	}

	if props := cfg.Properties; props != nil {
		d.Set("soft_delete_enabled", props.SoftDeleteFeatureState == backup.SoftDeleteFeatureStateEnabled)
	}

	storageCfg, err := storageCfgsClient.Get(ctx, id.VaultName, id.ResourceGroupName)
	if err != nil {
		return fmt.Errorf("reading Recovery Service storage Cfg %s: %+v", id.String(), err)
	}

	if props := storageCfg.Properties; props != nil {
		d.Set("storage_mode_type", string(props.StorageModelType))
		d.Set("cross_region_restore_enabled", props.CrossRegionRestoreFlag)
	}

	flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", flattenedIdentity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	encryption := flattenVaultEncryption(*model)
	if encryption != nil {
		d.Set("encryption", []interface{}{encryption})
	}

	monitorSettings := flattenVaultMonitorSettings(*model)
	if monitorSettings != nil {
		d.Set("monitor_settings", []interface{}{monitorSettings})
	}

	return tags.FlattenAndSet(d, model.Tags)
}

func resourceRecoveryServicesVaultDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.VaultsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := vaults.ParseVaultID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting Recovery Service  %s", id.String())

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("issuing delete request for Recovery Service %s: %+v", id.String(), err)
		}
	}

	return nil
}

func expandEncryption(d *pluginsdk.ResourceData) *vaults.VaultPropertiesEncryption {
	encryptionRaw := d.Get("encryption")
	if encryptionRaw == nil {
		return nil
	}
	settings := encryptionRaw.([]interface{})
	if len(settings) == 0 {
		return nil
	}
	encryptionMap := settings[0].(map[string]interface{})
	keyUri := encryptionMap["key_id"].(string)
	enabledInfraEncryption := encryptionMap["infrastructure_encryption_enabled"].(bool)
	infraEncryptionState := vaults.InfrastructureEncryptionStateEnabled
	if !enabledInfraEncryption {
		infraEncryptionState = vaults.InfrastructureEncryptionStateDisabled
	}
	encryption := &vaults.VaultPropertiesEncryption{
		KeyVaultProperties: &vaults.CmkKeyVaultProperties{
			KeyUri: utils.String(keyUri),
		},
		KekIdentity: &vaults.CmkKekIdentity{
			UseSystemAssignedIdentity: utils.Bool(encryptionMap["use_system_assigned_identity"].(bool)),
		},
		InfrastructureEncryption: &infraEncryptionState,
	}
	return encryption
}

func flattenVaultEncryption(model vaults.Vault) interface{} {
	if model.Properties == nil || model.Properties.Encryption == nil {
		return nil
	}
	encryption := model.Properties.Encryption
	if encryption.KeyVaultProperties == nil || encryption.KeyVaultProperties.KeyUri == nil {
		return nil
	}
	if encryption.KekIdentity == nil || encryption.KekIdentity.UseSystemAssignedIdentity == nil {
		return nil
	}
	encryptionMap := make(map[string]interface{})

	encryptionMap["key_id"] = encryption.KeyVaultProperties.KeyUri
	encryptionMap["use_system_assigned_identity"] = *encryption.KekIdentity.UseSystemAssignedIdentity
	encryptionMap["infrastructure_encryption_enabled"] = *encryption.InfrastructureEncryption == vaults.InfrastructureEncryptionStateEnabled
	return encryptionMap
}

func expandMonitorSettings(d *pluginsdk.ResourceData) *vaults.MonitoringSettings {
	monitorSettingsRaw := d.Get("monitor_settings")
	if monitorSettingsRaw == nil {
		return nil
	}
	settings := monitorSettingsRaw.([]interface{})
	if len(settings) == 0 {
		return nil
	}
	monitorSettingsMap := settings[0].(map[string]interface{})
	monitorSettings := &vaults.MonitoringSettings{
		AzureMonitorAlertSettings: &vaults.AzureMonitorAlertSettings{
			AlertsForAllJobFailures: expandAlertState(monitorSettingsMap["azure_monitor_alert_enabled"].(bool)),
		},
		ClassicAlertSettings: &vaults.ClassicAlertSettings{
			AlertsForCriticalOperations: expandAlertState(monitorSettingsMap["classic_alert_enabled"].(bool)),
		},
	}
	return monitorSettings
}

func flattenVaultMonitorSettings(model vaults.Vault) interface{} {
	if model.Properties == nil || model.Properties.MonitoringSettings == nil {
		return nil
	}
	monitorSettings := model.Properties.MonitoringSettings
	if monitorSettings.AzureMonitorAlertSettings == nil || monitorSettings.AzureMonitorAlertSettings.AlertsForAllJobFailures == nil {
		return nil
	}
	if monitorSettings.ClassicAlertSettings == nil || monitorSettings.ClassicAlertSettings.AlertsForCriticalOperations == nil {
		return nil
	}
	monitorSettingsMap := make(map[string]interface{})
	monitorSettingsMap["azure_monitor_alert_enabled"] = *monitorSettings.AzureMonitorAlertSettings.AlertsForAllJobFailures == vaults.AlertsStateEnabled
	monitorSettingsMap["classic_alert_enabled"] = *monitorSettings.ClassicAlertSettings.AlertsForCriticalOperations == vaults.AlertsStateEnabled
	return monitorSettingsMap
}

func expandAlertState(enabled bool) *vaults.AlertsState {
	out := vaults.AlertsStateDisabled
	if enabled {
		out = vaults.AlertsStateEnabled
	}
	return &out
}
