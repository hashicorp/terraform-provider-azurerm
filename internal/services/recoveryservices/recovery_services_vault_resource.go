package recoveryservices

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2021-08-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2021-12-01/backupresourcestorageconfigsnoncrr"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2021-12-01/backupresourcevaultconfigs"
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
			Create: pluginsdk.DefaultTimeout(120 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
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
						// We must use system assigned identity for now since recovery vault only support system assigned for now.
						// We can remove this property, but in that way when we enable user assigned identity in the future
						// , many users might be surprised at update in place. So we use an anonymous function to restrict this value to `true`
						"use_system_assigned_identity": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							ValidateFunc: func(i interface{}, s string) ([]string, []error) {
								use := i.(bool)
								if !use {
									return nil, []error{fmt.Errorf(" at this time `use_system_assigned_identity` only support `true`")}
								}
								return nil, nil
							},
							Default: true,
						},
					},
				},
			},

			// TODO: the API for this also supports UserAssigned & SystemAssigned, UserAssigned
			"identity": commonschema.SystemAssignedIdentityOptional(),

			"tags": commonschema.Tags(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(vaults.SkuNameRSZero),
					string(vaults.SkuNameStandard),
				}, false),
			},

			"storage_mode_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  backupresourcestorageconfigsnoncrr.StorageTypeGeoRedundant,
				ValidateFunc: validation.StringInSlice([]string{
					string(backupresourcestorageconfigsnoncrr.StorageTypeGeoRedundant),
					string(backupresourcestorageconfigsnoncrr.StorageTypeLocallyRedundant),
					string(backupresourcestorageconfigsnoncrr.StorageTypeZoneRedundant),
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
	storageId := backupresourcestorageconfigsnoncrr.VaultId{
		SubscriptionId:    id.SubscriptionId,
		ResourceGroupName: id.ResourceGroupName,
		VaultName:         id.VaultName,
	}
	cfgId := backupresourcevaultconfigs.VaultId{
		SubscriptionId:    id.SubscriptionId,
		ResourceGroupName: id.ResourceGroupName,
		VaultName:         id.VaultName,
	}

	storageMode := d.Get("storage_mode_type").(string)
	crossRegionRestore := d.Get("cross_region_restore_enabled").(bool)

	if crossRegionRestore && storageMode != string(backupresourcestorageconfigsnoncrr.StorageTypeGeoRedundant) {
		return fmt.Errorf("cannot enable cross region restore when storage mode type is not %s. %s", string(backupresourcestorageconfigsnoncrr.StorageTypeGeoRedundant), id.String())
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
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_recovery_services_vault", id.ID())
	}

	expandedIdentity, err := expandVaultIdentity(d.Get("identity").([]interface{}))
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
		return fmt.Errorf("creating %s: %+v", id.String(), err)
	}

	storageType := backupresourcestorageconfigsnoncrr.StorageType(d.Get("storage_mode_type").(string))
	storageCfg := backupresourcestorageconfigsnoncrr.BackupResourceConfigResource{
		Properties: &backupresourcestorageconfigsnoncrr.BackupResourceConfig{
			StorageModelType:       &storageType,
			CrossRegionRestoreFlag: utils.Bool(d.Get("cross_region_restore_enabled").(bool)),
		},
	}

	err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		if resp, err := storageCfgsClient.Update(ctx, storageId, storageCfg); err != nil {
			if response.WasNotFound(resp.HttpResponse) {
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
	err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		if resp, err := storageCfgsClient.Get(ctx, storageId); err == nil {
			if resp.Model == nil {
				return pluginsdk.NonRetryableError(fmt.Errorf("updating %s Storage Config: `model` was nil", id))
			}
			if resp.Model.Properties == nil {
				return pluginsdk.NonRetryableError(fmt.Errorf("updating %s Storage Config: `properties` was nil", id))
			}
			if *resp.Model.Properties.StorageType != *storageCfg.Properties.StorageModelType {
				return pluginsdk.RetryableError(fmt.Errorf("updating Storage Config: %+v", err))
			}
			if *resp.Model.Properties.CrossRegionRestoreFlag != *storageCfg.Properties.CrossRegionRestoreFlag {
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

	// recovery vault's encryption config cannot be set while creation, so a standalone update is required.
	if _, ok := d.GetOk("encryption"); ok {
		err = client.UpdateThenPoll(ctx, id, vaults.PatchVault{
			Properties: &vaults.VaultProperties{
				Encryption: expandEncryption(d),
			},
		})
		if err != nil {
			return fmt.Errorf("updating Recovery Service Encryption %s: %+v, but recovery vault was created, a manually import might be required", id.String(), err)
		}
	}

	// an update on the vault will reset the vault config to default, so we handle it at last.
	enhancedSecurityState := backupresourcevaultconfigs.EnhancedSecurityStateEnabled
	cfg := backupresourcevaultconfigs.BackupResourceVaultConfigResource{
		Properties: &backupresourcevaultconfigs.BackupResourceVaultConfig{
			EnhancedSecurityState: &enhancedSecurityState, // always enabled
		},
	}

	var StateRefreshPendingStrings []string
	var StateRefreshTargetStrings []string
	if sd := d.Get("soft_delete_enabled").(bool); sd {
		state := backupresourcevaultconfigs.SoftDeleteFeatureStateEnabled
		cfg.Properties.SoftDeleteFeatureState = &state
		StateRefreshPendingStrings = []string{string(backupresourcevaultconfigs.SoftDeleteFeatureStateDisabled)}
		StateRefreshTargetStrings = []string{string(backupresourcevaultconfigs.SoftDeleteFeatureStateEnabled)}
	} else {
		state := backupresourcevaultconfigs.SoftDeleteFeatureStateDisabled
		cfg.Properties.SoftDeleteFeatureState = &state
		StateRefreshPendingStrings = []string{string(backupresourcevaultconfigs.SoftDeleteFeatureStateEnabled)}
		StateRefreshTargetStrings = []string{string(backupresourcevaultconfigs.SoftDeleteFeatureStateDisabled)}
	}

	_, err = cfgsClient.Update(ctx, cfgId, cfg)
	if err != nil {
		return err
	}

	// sometimes update sync succeed but READ returns with old value, so we refresh till the value is correct.
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   StateRefreshPendingStrings,
		Target:                    StateRefreshTargetStrings,
		MinTimeout:                30 * time.Second,
		ContinuousTargetOccurence: 3,
		Refresh:                   resourceRecoveryServicesVaultSoftDeleteRefreshFunc(ctx, cfgsClient, cfgId),
	}

	stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for on update for Recovery Service %s: %+v", id.String(), err)
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
	storageId := backupresourcestorageconfigsnoncrr.VaultId{
		SubscriptionId:    id.SubscriptionId,
		ResourceGroupName: id.ResourceGroupName,
		VaultName:         id.VaultName,
	}
	cfgId := backupresourcevaultconfigs.VaultId{
		SubscriptionId:    id.SubscriptionId,
		ResourceGroupName: id.ResourceGroupName,
		VaultName:         id.VaultName,
	}

	encryption := expandEncryption(d)
	existing, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("checking for presence of existing Recovery Service %s: %+v", id.String(), err)
	}
	if existing.Model != nil && existing.Model.Properties != nil {
		prop := existing.Model.Properties
		if prop.Encryption != nil {
			if encryption == nil {
				return fmt.Errorf("once encryption with your own key has been enabled it's not possible to disable it")
			}
			if *encryption.InfrastructureEncryption != *prop.Encryption.InfrastructureEncryption {
				return fmt.Errorf("once `infrastructure_encryption_enabled` has been set it's not possible to change it")
			}
			if d.HasChange("sku") {
				// Once encryption has been enabled, calling `CreateOrUpdate` without it is not allowed.
				// But `sku` can only be updated by `CreateOrUpdate` and the support for `encryption` in `CreateOrUpdate` is still under preview (https://docs.microsoft.com/azure/backup/encryption-at-rest-with-cmk?tabs=portal#enable-encryption-using-customer-managed-keys-at-vault-creation-in-preview).
				// TODO remove this restriction and add `encryption` to below `sku` update block when `encryption` in `CreateOrUpdate` is GA
				return fmt.Errorf("`sku` cannot be changed when encryption with your own key has been enabled")
			}
		}
	}

	storageMode := d.Get("storage_mode_type").(string)
	crossRegionRestore := d.Get("cross_region_restore_enabled").(bool)

	if crossRegionRestore && storageMode != string(backupresourcestorageconfigsnoncrr.StorageTypeGeoRedundant) {
		return fmt.Errorf("cannot enable cross region restore when storage mode type is not %s. %s", string(backupresourcestorageconfigsnoncrr.StorageTypeGeoRedundant), id.String())
	}

	expandedIdentity, err := expandVaultIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	enhanchedSecurityState := backupresourcevaultconfigs.EnhancedSecurityStateEnabled
	cfg := backupresourcevaultconfigs.BackupResourceVaultConfigResource{
		Properties: &backupresourcevaultconfigs.BackupResourceVaultConfig{
			EnhancedSecurityState: &enhanchedSecurityState, // always enabled
		},
	}

	if d.HasChanges("storage_mode_type", "cross_region_restore_enabled") {
		storageType := backupresourcestorageconfigsnoncrr.StorageType(storageMode)
		storageCfg := backupresourcestorageconfigsnoncrr.BackupResourceConfigResource{
			Properties: &backupresourcestorageconfigsnoncrr.BackupResourceConfig{
				StorageModelType:       &storageType,
				CrossRegionRestoreFlag: utils.Bool(crossRegionRestore),
			},
		}

		err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutUpdate), func() *pluginsdk.RetryError {
			if resp, err := storageCfgsClient.Update(ctx, storageId, storageCfg); err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return pluginsdk.RetryableError(fmt.Errorf("updating Recovery Service Storage Cfg %s: %+v", id.String(), err))
				}
				if response.WasBadRequest(resp.HttpResponse) {
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
			if resp, err := storageCfgsClient.Get(ctx, storageId); err == nil {
				if resp.Model == nil {
					return pluginsdk.NonRetryableError(fmt.Errorf("updating %s Storage Config: `model` was nil", id))
				}
				if resp.Model.Properties == nil {
					return pluginsdk.NonRetryableError(fmt.Errorf("updating %s Storage Config: `properties` was nil", id))
				}
				if *resp.Model.Properties.StorageType != *storageCfg.Properties.StorageModelType {
					return pluginsdk.RetryableError(fmt.Errorf("updating Storage Config: %+v", err))
				}
				if *resp.Model.Properties.CrossRegionRestoreFlag != *storageCfg.Properties.CrossRegionRestoreFlag {
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

		if vaults.SkuName(sku) == vaults.SkuNameRSZero {
			vault.Sku.Tier = utils.String("Standard")
		}

		err = client.CreateOrUpdateThenPoll(ctx, id, vault)
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

	if d.HasChange("tags") {
		vault.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	err = client.UpdateThenPoll(ctx, id, vault)
	if err != nil {
		return fmt.Errorf("updating  %s: %+v", id, err)
	}

	// an update on vault will cause the vault config reset to default, so whether the config has change or not, it needs to be updated.
	var StateRefreshPendingStrings []string
	var StateRefreshTargetStrings []string
	if sd := d.Get("soft_delete_enabled").(bool); sd {
		state := backupresourcevaultconfigs.SoftDeleteFeatureStateEnabled
		cfg.Properties.SoftDeleteFeatureState = &state
		StateRefreshPendingStrings = []string{string(backupresourcevaultconfigs.SoftDeleteFeatureStateDisabled)}
		StateRefreshTargetStrings = []string{string(backupresourcevaultconfigs.SoftDeleteFeatureStateEnabled)}
	} else {
		state := backupresourcevaultconfigs.SoftDeleteFeatureStateDisabled
		cfg.Properties.SoftDeleteFeatureState = &state
		StateRefreshPendingStrings = []string{string(backupresourcevaultconfigs.SoftDeleteFeatureStateEnabled)}
		StateRefreshTargetStrings = []string{string(backupresourcevaultconfigs.SoftDeleteFeatureStateDisabled)}
	}

	_, err = cfgsClient.Update(ctx, cfgId, cfg)
	if err != nil {
		return err
	}

	// sometimes update sync succeed but READ returns with old value, so we refresh till the value is correct.
	// tracked by https://github.com/Azure/azure-rest-api-specs/issues/21548
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   StateRefreshPendingStrings,
		Target:                    StateRefreshTargetStrings,
		MinTimeout:                30 * time.Second,
		ContinuousTargetOccurence: 3,
		Refresh:                   resourceRecoveryServicesVaultSoftDeleteRefreshFunc(ctx, cfgsClient, cfgId),
	}

	stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for on update for Recovery Service %s: %+v", id.String(), err)
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
	storageId := backupresourcestorageconfigsnoncrr.VaultId{
		SubscriptionId:    id.SubscriptionId,
		ResourceGroupName: id.ResourceGroupName,
		VaultName:         id.VaultName,
	}
	cfgId := backupresourcevaultconfigs.VaultId{
		SubscriptionId:    id.SubscriptionId,
		ResourceGroupName: id.ResourceGroupName,
		VaultName:         id.VaultName,
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

	if resp.Model == nil {
		return fmt.Errorf("recovery Service Vault response %q : model is nil", id.ID())
	}
	model := resp.Model

	d.Set("name", id.VaultName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("location", location.Normalize(model.Location))

	if sku := model.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	cfg, err := cfgsClient.Get(ctx, cfgId)
	if err != nil {
		return fmt.Errorf("reading Recovery Service Vault Cfg %s: %+v", id.String(), err)
	}

	if cfg.Model != nil && cfg.Model.Properties != nil && cfg.Model.Properties.SoftDeleteFeatureState != nil {
		d.Set("soft_delete_enabled", *cfg.Model.Properties.SoftDeleteFeatureState == backupresourcevaultconfigs.SoftDeleteFeatureStateEnabled)
	}

	storageCfg, err := storageCfgsClient.Get(ctx, storageId)
	if err != nil {
		return fmt.Errorf("reading Recovery Service storage Cfg %s: %+v", id.String(), err)
	}

	if storageCfg.Model != nil && storageCfg.Model.Properties != nil {
		props := storageCfg.Model.Properties
		d.Set("storage_mode_type", props.StorageModelType)
		d.Set("cross_region_restore_enabled", props.CrossRegionRestoreFlag)
	}

	flattenIdentity := flattenVaultIdentity(model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", flattenIdentity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	encryption := flattenVaultEncryption(*model)
	if encryption != nil {
		d.Set("encryption", []interface{}{encryption})
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

	_, err = client.Delete(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id.String(), err)
	}

	return nil
}

func expandVaultIdentity(input []interface{}) (*identity.SystemAndUserAssignedMap, error) {
	expanded, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	identityIds := make(map[string]identity.UserAssignedIdentityDetails, 0)
	return &identity.SystemAndUserAssignedMap{
		Type:        expanded.Type,
		PrincipalId: expanded.PrincipalId,
		TenantId:    expanded.TenantId,
		IdentityIds: identityIds,
	}, nil
}

func flattenVaultIdentity(input *identity.SystemAndUserAssignedMap) []interface{} {
	systemAssignedIdentity := &identity.SystemAssigned{
		Type:        input.Type,
		PrincipalId: input.PrincipalId,
		TenantId:    input.TenantId,
	}
	return identity.FlattenSystemAssigned(systemAssignedIdentity)
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

func resourceRecoveryServicesVaultSoftDeleteRefreshFunc(ctx context.Context, cfgsClient *backupresourcevaultconfigs.BackupResourceVaultConfigsClient, id backupresourcevaultconfigs.VaultId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := cfgsClient.Get(ctx, id)
		if err != nil {
			if strings.Contains(err.Error(), "ResourceNotYetSynced") {
				return resp, "syncing", nil
			}
			return resp, "error", fmt.Errorf("refreshing Recovery Service Vault Cfg %s: %+v", id.String(), err)
		}

		if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.SoftDeleteFeatureState != nil {
			return resp.Model, string(*resp.Model.Properties.SoftDeleteFeatureState), nil
		}

		return resp, "error", fmt.Errorf("refreshing Recovery Service Vault Cfg %s: Properties is nil", id.String())
	}
}
