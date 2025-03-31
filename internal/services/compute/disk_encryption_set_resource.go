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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskencryptionsets"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	managedHsmHelpers "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/helpers"
	managedHsmParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	managedHsmValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func resourceDiskEncryptionSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDiskEncryptionSetCreate,
		Read:   resourceDiskEncryptionSetRead,
		Update: resourceDiskEncryptionSetUpdate,
		Delete: resourceDiskEncryptionSetDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DiskEncryptionSetID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DiskEncryptionSetName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			// Issue #22864
			"key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
				ExactlyOneOf: []string{"managed_hsm_key_id", "key_vault_key_id"},
			},

			"managed_hsm_key_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.Any(managedHsmValidate.ManagedHSMDataPlaneVersionedKeyID, managedHsmValidate.ManagedHSMDataPlaneVersionlessKeyID),
				ExactlyOneOf: []string{"managed_hsm_key_id", "key_vault_key_id"},
			},

			"auto_key_rotation_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"encryption_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(diskencryptionsets.DiskEncryptionSetTypeEncryptionAtRestWithCustomerKey),
				ValidateFunc: validation.StringInSlice([]string{
					string(diskencryptionsets.DiskEncryptionSetTypeEncryptionAtRestWithCustomerKey),
					string(diskencryptionsets.DiskEncryptionSetTypeEncryptionAtRestWithPlatformAndCustomerKeys),
					string(diskencryptionsets.DiskEncryptionSetTypeConfidentialVMEncryptedWithCustomerKey),
				}, false),
			},

			"federated_client_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

			"tags": commonschema.Tags(),

			"key_vault_key_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("identity.0.type", func(ctx context.Context, old, new, meta interface{}) bool {
				// cannot change identity type from userAssigned to systemAssigned
				return (old.(string) == string(identity.TypeUserAssigned) || old.(string) == string(identity.TypeSystemAssignedUserAssigned)) && (new.(string) == string(identity.TypeSystemAssigned))
			}),
		),
	}
}

func resourceDiskEncryptionSetCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	keyVaultKeyClient := meta.(*clients.Client).KeyVault.ManagementClient
	managedkeyBundleClient := meta.(*clients.Client).ManagedHSMs.DataPlaneKeysClient
	env := meta.(*clients.Client).Account.Environment
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewDiskEncryptionSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for present of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_disk_encryption_set", id.ID())
	}

	rotationToLatestKeyVersionEnabled := d.Get("auto_key_rotation_enabled").(bool)
	activeKey := &diskencryptionsets.KeyForDiskEncryptionSet{}

	if keyVaultKeyId, ok := d.GetOk("key_vault_key_id"); ok {
		keyVaultKey, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(keyVaultKeyId.(string))
		if err != nil {
			return err
		}

		err = validateKeyAndRotationEnabled(rotationToLatestKeyVersionEnabled, keyVaultKey.Version != "", keyVaultKey.ID())
		if err != nil {
			return err
		}

		keyVaultDetails, err := diskEncryptionSetRetrieveKeyVault(ctx, keyVaultsClient, subscriptionId, *keyVaultKey)
		if err != nil {
			return fmt.Errorf("validating Key Vault Key %q for Disk Encryption Set: %+v", keyVaultKey.ID(), err)
		}

		if keyVaultDetails != nil {
			err = validateKeyVaultDetails(keyVaultDetails)
			if err != nil {
				return err
			}

			activeKey.SourceVault = &diskencryptionsets.SourceVault{
				Id: utils.String(keyVaultDetails.keyVaultId),
			}
		}

		// NOTE: The API requires a versioned key to be sent however if rotationToLatestKeyVersion is enabled this will cause
		// terraform to revert the rotated key to the previous version that is defined in the configuration file...
		// Issue #22864
		if rotationToLatestKeyVersionEnabled {
			// Get the latest version of the key...
			keyBundle, err := keyVaultKeyClient.GetKey(ctx, keyVaultKey.KeyVaultBaseUrl, keyVaultKey.Name, "")
			if err != nil {
				return err
			}

			if keyBundle.Key != nil {
				activeKey.KeyURL = pointer.From(keyBundle.Key.Kid)
			}
		} else {
			// Use the passed version of the key...
			activeKey.KeyURL = keyVaultKey.ID()
		}
	} else if managedHsmKeyId, ok := d.GetOk("managed_hsm_key_id"); ok {
		keyUrl, err := getManagedHsmKeyURL(ctx, managedkeyBundleClient, managedHsmKeyId.(string), rotationToLatestKeyVersionEnabled, env)
		if err != nil {
			return err
		}
		activeKey.KeyURL = keyUrl
	}

	encryptionType := diskencryptionsets.DiskEncryptionSetType(d.Get("encryption_type").(string))
	t := d.Get("tags").(map[string]interface{})

	expandedIdentity, err := expandDiskEncryptionSetIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	params := diskencryptionsets.DiskEncryptionSet{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &diskencryptionsets.EncryptionSetProperties{
			ActiveKey:                         activeKey,
			RotationToLatestKeyVersionEnabled: utils.Bool(rotationToLatestKeyVersionEnabled),
			EncryptionType:                    &encryptionType,
		},
		Identity: expandedIdentity,
		Tags:     tags.Expand(t),
	}

	if v, ok := d.GetOk("federated_client_id"); ok {
		params.Properties.FederatedClientId = utils.String(v.(string))
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, params)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDiskEncryptionSetRead(d, meta)
}

func validateKeyAndRotationEnabled(rotationToLatestKeyVersionEnabled bool, hasVersion bool, keyVaultKeyId string) error {
	if rotationToLatestKeyVersionEnabled {
		if hasVersion {
			return fmt.Errorf("'auto_key_rotation_enabled' field is set to 'true' expected a key vault key with a versionless ID but version information was found: %s", keyVaultKeyId)
		}
	} else {
		if !hasVersion {
			return fmt.Errorf("'auto_key_rotation_enabled' field is set to 'false' expected a key vault key with a versioned ID but no version information was found: %s", keyVaultKeyId)
		}
	}
	return nil
}

func resourceDiskEncryptionSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	env := meta.(*clients.Client).Account.Environment
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseDiskEncryptionSetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Disk Encryption Set %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Disk Encryption Set %q (Resource Group %q): %+v", id.DiskEncryptionSetName, id.ResourceGroupName, err)
	}

	d.Set("name", id.DiskEncryptionSetName)
	d.Set("resource_group_name", id.ResourceGroupName)

	model := resp.Model
	if model == nil {
		return fmt.Errorf("reading Disk Encryption Set : %+v", err)
	}

	if l := model.Location; l != "" {
		d.Set("location", location.Normalize(l))
	}

	if props := model.Properties; props != nil {
		var keyVaultKey *keyVaultParse.NestedItemId

		RotationToLatestKeyVersionEnabled := pointer.From(props.RotationToLatestKeyVersionEnabled)
		d.Set("auto_key_rotation_enabled", RotationToLatestKeyVersionEnabled)

		if props.ActiveKey != nil && props.ActiveKey.KeyURL != "" {
			keyVaultURI := props.ActiveKey.KeyURL

			isHSMURI, err, instanceName, domainSuffix := managedHsmHelpers.IsManagedHSMURI(env, keyVaultURI)
			if err != nil {
				return err
			}

			switch {
			case !isHSMURI:
				{
					keyVaultKey, err = keyVaultParse.ParseOptionallyVersionedNestedItemID(props.ActiveKey.KeyURL)
					if err != nil {
						return err
					}

					// This "should" never happen, but if keyVaultKey does not get assigned above it
					// would cause a panic when referenced below, so check to make sure it was
					// assigned or not...
					if keyVaultKey == nil {
						return fmt.Errorf("`KeyForDiskEncryptionSet.ActiveKey` was nil")
					}
					d.Set("key_vault_key_url", keyVaultKey.ID())

					// NOTE: Since the auto rotation changes the version information when the key is rotated
					// we need to persist the versionless key ID to the state file else terraform will always
					// try to revert to the original version of the key once it has been rotated...
					// Issue #22864
					if RotationToLatestKeyVersionEnabled {
						d.Set("key_vault_key_id", keyVaultKey.VersionlessID())
					} else {
						d.Set("key_vault_key_id", keyVaultKey.ID())
					}
				}

			case isHSMURI:
				{
					keyId, err := managedHsmParse.ManagedHSMDataPlaneVersionedKeyID(keyVaultURI, &domainSuffix)
					if err != nil {
						return err
					}

					// See comment above and issue #22864
					if RotationToLatestKeyVersionEnabled {
						versionlessKeyId := managedHsmParse.NewManagedHSMDataPlaneVersionlessKeyID(instanceName, domainSuffix, keyId.KeyName)
						d.Set("managed_hsm_key_id", versionlessKeyId.ID())
					} else {
						d.Set("managed_hsm_key_id", keyId.ID())
					}
				}
			}
		}

		encryptionType := string(diskencryptionsets.DiskEncryptionSetTypeEncryptionAtRestWithCustomerKey)
		if props.EncryptionType != nil {
			encryptionType = string(*props.EncryptionType)
		}
		d.Set("encryption_type", encryptionType)

		federatedClientId := ""
		if props.FederatedClientId != nil {
			federatedClientId = *props.FederatedClientId
		}
		d.Set("federated_client_id", federatedClientId)
	}

	flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}

	if err := d.Set("identity", flattenedIdentity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, model.Tags)
}

func resourceDiskEncryptionSetUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	keyVaultKeyClient := meta.(*clients.Client).KeyVault.ManagementClient
	managedkeyBundleClient := meta.(*clients.Client).ManagedHSMs.DataPlaneKeysClient
	env := meta.(*clients.Client).Account.Environment
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseDiskEncryptionSetID(d.Id())
	if err != nil {
		return err
	}

	update := diskencryptionsets.DiskEncryptionSetUpdate{}

	if d.HasChange("identity") {
		expandedIdentity, err := expandDiskEncryptionSetIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		update.Identity = expandedIdentity
	}

	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	rotationToLatestKeyVersionEnabled := d.Get("auto_key_rotation_enabled").(bool)

	if keyVaultKeyId, ok := d.GetOk("key_vault_key_id"); ok && d.HasChange("key_vault_key_id") {
		keyVaultKey, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(keyVaultKeyId.(string))
		if err != nil {
			return err
		}

		err = validateKeyAndRotationEnabled(rotationToLatestKeyVersionEnabled, keyVaultKey.Version != "", keyVaultKey.ID())
		if err != nil {
			return err
		}

		if update.Properties == nil {
			update.Properties = &diskencryptionsets.DiskEncryptionSetUpdateProperties{
				ActiveKey: &diskencryptionsets.KeyForDiskEncryptionSet{},
			}
		}

		// NOTE: The API requires a versioned key to be sent however if rotationToLatestKeyVersion is enabled this will cause
		// terraform to revert the rotated key to the previous version that is defined in the configuration file...
		// Issue #22864
		if rotationToLatestKeyVersionEnabled {
			// Get the latest version of the key...
			keyBundle, err := keyVaultKeyClient.GetKey(ctx, keyVaultKey.KeyVaultBaseUrl, keyVaultKey.Name, "")
			if err != nil {
				return err
			}

			if keyBundle.Key != nil {
				update.Properties.ActiveKey.KeyURL = pointer.From(keyBundle.Key.Kid)
			}
		} else {
			// Use the passed version of the key...
			update.Properties.ActiveKey.KeyURL = keyVaultKey.ID()
		}

		keyVaultDetails, err := diskEncryptionSetRetrieveKeyVault(ctx, keyVaultsClient, id.SubscriptionId, *keyVaultKey)
		if err != nil {
			return fmt.Errorf("validating Key Vault Key %q for Disk Encryption Set: %+v", keyVaultKey.ID(), err)
		}

		if keyVaultDetails != nil {
			err = validateKeyVaultDetails(keyVaultDetails)
			if err != nil {
				return err
			}

			update.Properties.ActiveKey.SourceVault = &diskencryptionsets.SourceVault{
				Id: utils.String(keyVaultDetails.keyVaultId),
			}
		}
	} else if managedHsmKeyId, ok := d.GetOk("managed_hsm_key_id"); ok && d.HasChange("managed_hsm_key_id") {
		keyUrl, err := getManagedHsmKeyURL(ctx, managedkeyBundleClient, managedHsmKeyId.(string), rotationToLatestKeyVersionEnabled, env)
		if err != nil {
			return err
		}
		update.Properties.ActiveKey.KeyURL = keyUrl
	}

	if d.HasChange("auto_key_rotation_enabled") {
		if update.Properties == nil {
			update.Properties = &diskencryptionsets.DiskEncryptionSetUpdateProperties{}
		}

		update.Properties.RotationToLatestKeyVersionEnabled = utils.Bool(rotationToLatestKeyVersionEnabled)
	}

	if d.HasChange("federated_client_id") {
		if update.Properties == nil {
			update.Properties = &diskencryptionsets.DiskEncryptionSetUpdateProperties{}
		}

		v, ok := d.GetOk("federated_client_id")
		if ok {
			update.Properties.FederatedClientId = utils.String(v.(string))
		} else {
			update.Properties.FederatedClientId = utils.String("None") // this is the only way to remove the federated client id
		}
	}

	err = client.UpdateThenPoll(ctx, *id, update)
	if err != nil {
		return fmt.Errorf("updating Disk Encryption Set %q (Resource Group %q): %+v", id.DiskEncryptionSetName, id.ResourceGroupName, err)
	}

	return resourceDiskEncryptionSetRead(d, meta)
}

func getManagedHsmKeyURL(ctx context.Context, managedkeyBundleClient *keyvault.BaseClient, managedHsmKeyId string, rotationToLatestKeyVersionEnabled bool, env environments.Environment) (string, error) {
	domainSuffix, ok := env.ManagedHSM.DomainSuffix()
	if !ok {
		return "", fmt.Errorf("managed HSM is not supported in this Environment")
	}
	managedHsmVersionedKey, err := managedHsmParse.ManagedHSMDataPlaneVersionedKeyID(managedHsmKeyId, domainSuffix)
	if err == nil {
		err = validateKeyAndRotationEnabled(rotationToLatestKeyVersionEnabled, true, managedHsmVersionedKey.ID())
		if err != nil {
			return "", err
		}
		return managedHsmVersionedKey.ID(), nil
	} else {
		managedHsmVersionlessKey, err := managedHsmParse.ManagedHSMDataPlaneVersionlessKeyID(managedHsmKeyId, domainSuffix)
		if err != nil {
			return "", err
		}
		err = validateKeyAndRotationEnabled(rotationToLatestKeyVersionEnabled, false, managedHsmVersionlessKey.ID())
		if err != nil {
			return "", err
		}
		keyBundle, err := managedkeyBundleClient.GetKey(ctx, managedHsmVersionlessKey.BaseUri(), managedHsmVersionlessKey.KeyName, "")
		if err != nil {
			return "", err
		}

		if keyBundle.Key.Kid != nil {
			return pointer.From(keyBundle.Key.Kid), nil
		}
		return "", fmt.Errorf("retrieving key version for key %s: Key Vault did not return a version", managedHsmKeyId)
	}
}

func validateKeyVaultDetails(keyVaultDetails *diskEncryptionSetKeyVault) error {
	if !keyVaultDetails.softDeleteEnabled {
		return fmt.Errorf("validating Key Vault %q (Resource Group %q) for Disk Encryption Set: Soft Delete must be enabled", keyVaultDetails.keyVaultName, keyVaultDetails.resourceGroupName)
	}

	return nil
}

func resourceDiskEncryptionSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseDiskEncryptionSetID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting Disk Encryption Set %q (Resource Group %q): %+v", id.DiskEncryptionSetName, id.ResourceGroupName, err)
	}

	return nil
}

func expandDiskEncryptionSetIdentity(input []interface{}) (*identity.SystemAndUserAssignedMap, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	return expanded, nil
}

type diskEncryptionSetKeyVault struct {
	keyVaultId             string
	resourceGroupName      string
	keyVaultName           string
	purgeProtectionEnabled bool
	softDeleteEnabled      bool
}

func diskEncryptionSetRetrieveKeyVault(ctx context.Context, keyVaultsClient *client.Client, subscriptionId string, keyVaultKeyId keyVaultParse.NestedItemId) (*diskEncryptionSetKeyVault, error) {
	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultID, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, keyVaultKeyId.KeyVaultBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", keyVaultKeyId.KeyVaultBaseUrl, err)
	}
	if keyVaultID == nil {
		return nil, nil
	}

	parsedKeyVaultID, err := commonids.ParseKeyVaultID(*keyVaultID)
	if err != nil {
		return nil, err
	}

	resp, err := keyVaultsClient.VaultsClient.Get(ctx, *parsedKeyVaultID)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *parsedKeyVaultID, err)
	}

	purgeProtectionEnabled := false
	softDeleteEnabled := false
	if model := resp.Model; model != nil {
		if model.Properties.EnableSoftDelete != nil {
			softDeleteEnabled = *model.Properties.EnableSoftDelete
		}

		if model.Properties.EnablePurgeProtection != nil {
			purgeProtectionEnabled = *model.Properties.EnablePurgeProtection
		}
	}

	return &diskEncryptionSetKeyVault{
		keyVaultId:             *keyVaultID,
		resourceGroupName:      parsedKeyVaultID.ResourceGroupName,
		keyVaultName:           parsedKeyVaultID.VaultName,
		purgeProtectionEnabled: purgeProtectionEnabled,
		softDeleteEnabled:      softDeleteEnabled,
	}, nil
}
