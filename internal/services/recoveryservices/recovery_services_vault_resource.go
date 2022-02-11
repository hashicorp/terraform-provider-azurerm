package recoveryservices

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-07-01/backup"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-08-01/recoveryservices"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	keyvaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceRecoveryServicesVault() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRecoveryServicesVaultCreateUpdate,
		Read:   resourceRecoveryServicesVaultRead,
		Update: resourceRecoveryServicesVaultCreateUpdate,
		Delete: resourceRecoveryServicesVaultDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VaultID(id)
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

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

			"tags": tags.Schema(),

			"sku": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifferenceV2Only,
				ValidateFunc: validation.StringInSlice([]string{
					string(recoveryservices.SkuNameRS0),
					string(recoveryservices.SkuNameStandard),
				}, !features.ThreePointOh()),
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

			"soft_delete_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceRecoveryServicesVaultCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.VaultsClient
	cfgsClient := meta.(*clients.Client).RecoveryServices.VaultsConfigsClient
	storageCfgsClient := meta.(*clients.Client).RecoveryServices.StorageConfigsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewVaultID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	location := d.Get("location").(string)
	t := d.Get("tags").(map[string]interface{})

	log.Printf("[DEBUG] Creating/updating Recovery Service %s", id.String())

	encryption := expandEncryption(d)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Recovery Service %s: %+v", id.String(), err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_recovery_services_vault", *existing.ID)
		}
	} else {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("checking for presence of existing Recovery Service %s: %+v", id.String(), err)
		}
		if existing.Properties != nil && existing.Properties.Encryption != nil {
			if encryption == nil {
				return fmt.Errorf("once encryption with your own key has been enabled it's not possible to disable it")
			}
			if encryption.InfrastructureEncryption != existing.Properties.Encryption.InfrastructureEncryption {
				return fmt.Errorf("once `infrastructure_encryption_enabled` has been set it's not possible to change it")
			}
		}
	}

	expandedIdentity, err := expandVaultIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	vault := recoveryservices.Vault{
		Location: utils.String(location),
		Tags:     tags.Expand(t),
		Identity: expandedIdentity,
		Sku: &recoveryservices.Sku{
			Name: recoveryservices.SkuName(d.Get("sku").(string)),
		},
		Properties: &recoveryservices.VaultProperties{},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, vault); err != nil {
		return fmt.Errorf("creating/updating Recovery Service %s: %+v", id.String(), err)
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
			resp, err := cfgsClient.Update(ctx, id.Name, id.ResourceGroup, cfg)
			if err != nil {
				if strings.Contains(err.Error(), "ResourceNotYetSynced") {
					return resp, "syncing", nil
				}
				return resp, "error", fmt.Errorf("updating Recovery Service Vault Cfg %s: %+v", id.String(), err)
			}

			return resp, "success", nil
		},
	}

	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for on update for Recovery Service  %s: %+v", id.String(), err)
	}

	storageCfg := backup.ResourceConfigResource{
		Properties: &backup.ResourceConfig{
			StorageModelType: backup.StorageType(d.Get("storage_mode_type").(string)),
		},
	}

	err = pluginsdk.Retry(stateConf.Timeout, func() *pluginsdk.RetryError {
		if resp, err := storageCfgsClient.Update(ctx, id.Name, id.ResourceGroup, storageCfg); err != nil {
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
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	// storage type is not updated instantaneously, so we wait until storage type is correct
	err = pluginsdk.Retry(stateConf.Timeout, func() *pluginsdk.RetryError {
		if resp, err := storageCfgsClient.Get(ctx, id.Name, id.ResourceGroup); err == nil {
			if resp.Properties == nil {
				return pluginsdk.NonRetryableError(fmt.Errorf("updating %s Storage Config: `properties` was nil", id))
			}
			if resp.Properties.StorageType != storageCfg.Properties.StorageModelType {
				return pluginsdk.RetryableError(fmt.Errorf("updating Storage Config: %+v", err))
			}
		} else {
			return pluginsdk.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	// recovery vault's encryption config cannot be set while creation, so a standalone update is required.
	if _, ok := d.GetOk("encryption"); ok {
		updateFuture, err := client.Update(ctx, id.ResourceGroup, id.Name, recoveryservices.PatchVault{
			Properties: &recoveryservices.VaultProperties{
				Encryption: encryption,
			},
		})
		if err != nil {
			return fmt.Errorf("updating Recovery Service Encryption %s: %+v, but recovery vault was created, a manually import might be required", id.String(), err)
		}
		if err = updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for update encryption of %s: %+v, but recovery vault was created, a manually import might be required", id.String(), err)
		}
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

	id, err := parse.VaultID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading Recovery Service %s", id.String())

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Recovery Service %s: %+v", id.String(), err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	cfg, err := cfgsClient.Get(ctx, id.Name, id.ResourceGroup)
	if err != nil {
		return fmt.Errorf("reading Recovery Service Vault Cfg %s: %+v", id.String(), err)
	}

	if props := cfg.Properties; props != nil {
		d.Set("soft_delete_enabled", props.SoftDeleteFeatureState == backup.SoftDeleteFeatureStateEnabled)
	}

	storageCfg, err := storageCfgsClient.Get(ctx, id.Name, id.ResourceGroup)
	if err != nil {
		return fmt.Errorf("reading Recovery Service storage Cfg %s: %+v", id.String(), err)
	}

	if props := storageCfg.Properties; props != nil {
		d.Set("storage_mode_type", string(props.StorageModelType))
	}

	if err := d.Set("identity", flattenVaultIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	encryption := flattenVaultEncryption(resp)
	if encryption != nil {
		d.Set("encryption", []interface{}{encryption})
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceRecoveryServicesVaultDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.VaultsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VaultID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting Recovery Service  %s", id.String())

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("issuing delete request for Recovery Service %s: %+v", id.String(), err)
		}
	}

	return nil
}

func expandVaultIdentity(input []interface{}) (*recoveryservices.IdentityData, error) {
	expanded, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	return &recoveryservices.IdentityData{
		Type: recoveryservices.ResourceIdentityType(string(expanded.Type)),
	}, nil
}

func flattenVaultIdentity(input *recoveryservices.IdentityData) []interface{} {
	var transition *identity.SystemAssigned

	if input != nil {
		transition = &identity.SystemAssigned{
			Type: identity.Type(string(input.Type)),
		}
		if input.PrincipalID != nil {
			transition.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transition.TenantId = *input.TenantID
		}
	}

	return identity.FlattenSystemAssigned(transition)
}

func expandEncryption(d *pluginsdk.ResourceData) *recoveryservices.VaultPropertiesEncryption {
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
	infraEncryptionState := recoveryservices.InfrastructureEncryptionStateEnabled
	if !enabledInfraEncryption {
		infraEncryptionState = recoveryservices.InfrastructureEncryptionStateDisabled
	}
	encryption := &recoveryservices.VaultPropertiesEncryption{
		KeyVaultProperties: &recoveryservices.CmkKeyVaultProperties{
			KeyURI: utils.String(keyUri),
		},
		KekIdentity: &recoveryservices.CmkKekIdentity{
			UseSystemAssignedIdentity: utils.Bool(encryptionMap["use_system_assigned_identity"].(bool)),
		},
		InfrastructureEncryption: infraEncryptionState,
	}
	return encryption
}

func flattenVaultEncryption(resp recoveryservices.Vault) interface{} {
	if resp.Properties == nil || resp.Properties.Encryption == nil {
		return nil
	}
	encryption := resp.Properties.Encryption
	if encryption.KeyVaultProperties == nil || encryption.KeyVaultProperties.KeyURI == nil {
		return nil
	}
	if encryption.KekIdentity == nil || encryption.KekIdentity.UseSystemAssignedIdentity == nil {
		return nil
	}
	encryptionMap := make(map[string]interface{})

	encryptionMap["key_id"] = encryption.KeyVaultProperties.KeyURI
	encryptionMap["use_system_assigned_identity"] = *encryption.KekIdentity.UseSystemAssignedIdentity
	encryptionMap["infrastructure_encryption_enabled"] = encryption.InfrastructureEncryption == recoveryservices.InfrastructureEncryptionStateEnabled
	return encryptionMap
}
