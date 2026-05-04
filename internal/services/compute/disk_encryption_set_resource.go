// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskencryptionsets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDiskEncryptionSet() *pluginsdk.Resource {
	r := &pluginsdk.Resource{
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

			"key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
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

	if !features.FivePointOh() {
		r.Schema["key_vault_key_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeAny),
			ExactlyOneOf: []string{"managed_hsm_key_id", "key_vault_key_id"},
		}

		r.Schema["managed_hsm_key_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeAny),
			ExactlyOneOf: []string{"managed_hsm_key_id", "key_vault_key_id"},
			Deprecated:   "`managed_hsm_key_id` has been deprecated in favour of `key_vault_key_id` and will be removed in v5.0 of the AzureRM Provider",
		}
	}

	return r
}

func resourceDiskEncryptionSetCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
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

	if !features.FivePointOh() && !d.GetRawConfig().AsValueMap()["managed_hsm_key_id"].IsNull() {
		key, err := keyvault.ParseNestedItemID(d.Get("managed_hsm_key_id").(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeAny)
		if err != nil {
			return err
		}

		keyURL, err := getKeyURL(ctx, key, rotationToLatestKeyVersionEnabled, meta)
		if err != nil {
			return err
		}
		activeKey.KeyURL = keyURL
	} else {
		nestedItemType := keyvault.NestedItemTypeKey
		if !features.FivePointOh() {
			nestedItemType = keyvault.NestedItemTypeAny
		}

		key, err := keyvault.ParseNestedItemID(d.Get("key_vault_key_id").(string), keyvault.VersionTypeAny, nestedItemType)
		if err != nil {
			return err
		}

		keyURL, err := getKeyURL(ctx, key, rotationToLatestKeyVersionEnabled, meta)
		if err != nil {
			return err
		}
		activeKey.KeyURL = keyURL
	}

	expandedIdentity, err := expandDiskEncryptionSetIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	params := diskencryptionsets.DiskEncryptionSet{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &diskencryptionsets.EncryptionSetProperties{
			ActiveKey:                         activeKey,
			RotationToLatestKeyVersionEnabled: pointer.To(rotationToLatestKeyVersionEnabled),
			EncryptionType:                    pointer.ToEnum[diskencryptionsets.DiskEncryptionSetType](d.Get("encryption_type").(string)),
		},
		Identity: expandedIdentity,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("federated_client_id"); ok {
		params.Properties.FederatedClientId = pointer.To(v.(string))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, params); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDiskEncryptionSetRead(d, meta)
}

func resourceDiskEncryptionSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
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
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.DiskEncryptionSetName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			rotationToLatestKeyVersionEnabled := pointer.From(props.RotationToLatestKeyVersionEnabled)
			d.Set("auto_key_rotation_enabled", rotationToLatestKeyVersionEnabled)

			encryptionType := string(diskencryptionsets.DiskEncryptionSetTypeEncryptionAtRestWithCustomerKey)
			if props.EncryptionType != nil {
				encryptionType = string(*props.EncryptionType)
			}
			d.Set("encryption_type", encryptionType)
			d.Set("federated_client_id", pointer.From(props.FederatedClientId))

			if props.ActiveKey != nil && props.ActiveKey.KeyURL != "" {
				nestedItemType := keyvault.NestedItemTypeKey
				if !features.FivePointOh() {
					nestedItemType = keyvault.NestedItemTypeAny
				}

				key, err := keyvault.ParseNestedItemID(props.ActiveKey.KeyURL, keyvault.VersionTypeAny, nestedItemType)
				if err != nil {
					return err
				}
				d.Set("key_vault_key_url", key.ID())

				if rotationToLatestKeyVersionEnabled {
					key.Version = ""
				}

				d.Set("key_vault_key_id", key.ID())
				if !features.FivePointOh() {
					if key.IsManagedHSM() {
						d.Set("managed_hsm_key_id", key.ID())
					} else {
						d.Set("managed_hsm_key_id", "")
					}
				}
			}
		}

		flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceDiskEncryptionSetUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
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

	if !features.FivePointOh() && d.HasChange("managed_hsm_key_id") {
		if update.Properties == nil {
			update.Properties = &diskencryptionsets.DiskEncryptionSetUpdateProperties{
				ActiveKey: &diskencryptionsets.KeyForDiskEncryptionSet{},
			}
		}

		key, err := keyvault.ParseNestedItemID(d.Get("managed_hsm_key_id").(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeAny)
		if err != nil {
			return err
		}

		keyURL, err := getKeyURL(ctx, key, rotationToLatestKeyVersionEnabled, meta)
		if err != nil {
			return err
		}
		update.Properties.ActiveKey.KeyURL = keyURL
	}

	if d.HasChange("key_vault_key_id") {
		if update.Properties == nil {
			update.Properties = &diskencryptionsets.DiskEncryptionSetUpdateProperties{
				ActiveKey: &diskencryptionsets.KeyForDiskEncryptionSet{},
			}
		}

		nestedItemType := keyvault.NestedItemTypeKey
		if !features.FivePointOh() {
			nestedItemType = keyvault.NestedItemTypeAny
		}

		key, err := keyvault.ParseNestedItemID(d.Get("key_vault_key_id").(string), keyvault.VersionTypeAny, nestedItemType)
		if err != nil {
			return err
		}

		keyURL, err := getKeyURL(ctx, key, rotationToLatestKeyVersionEnabled, meta)
		if err != nil {
			return err
		}
		update.Properties.ActiveKey.KeyURL = keyURL
	}

	if d.HasChange("auto_key_rotation_enabled") {
		if update.Properties == nil {
			update.Properties = &diskencryptionsets.DiskEncryptionSetUpdateProperties{}
		}

		update.Properties.RotationToLatestKeyVersionEnabled = pointer.To(rotationToLatestKeyVersionEnabled)
	}

	if d.HasChange("federated_client_id") {
		if update.Properties == nil {
			update.Properties = &diskencryptionsets.DiskEncryptionSetUpdateProperties{}
		}

		update.Properties.FederatedClientId = pointer.To("None")
		if v, ok := d.GetOk("federated_client_id"); ok {
			update.Properties.FederatedClientId = pointer.To(v.(string))
		}
	}

	if err := client.UpdateThenPoll(ctx, *id, update); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceDiskEncryptionSetRead(d, meta)
}

func resourceDiskEncryptionSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseDiskEncryptionSetID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func getKeyURL(ctx context.Context, id *keyvault.NestedItemID, rotationToLatestKeyVersionEnabled bool, meta any) (string, error) {
	keyVaultClient := meta.(*clients.Client).KeyVault.ManagementClient
	managedHSMClient := meta.(*clients.Client).ManagedHSMs.DataPlaneKeysClient

	keyURL := ""

	if rotationToLatestKeyVersionEnabled {
		if id.Version != "" {
			return keyURL, fmt.Errorf("`auto_key_rotation_enabled` field is set to `true` expected a key vault key with a versionless ID but version information was found: %s", id)
		}

		if id.IsManagedHSM() {
			keyBundle, err := managedHSMClient.GetKey(ctx, id.KeyVaultBaseURL, id.Name, "")
			if err != nil {
				return keyURL, err
			}

			if keyBundle.Key != nil {
				keyURL = pointer.From(keyBundle.Key.Kid)
			}
		} else {
			keyBundle, err := keyVaultClient.GetKey(ctx, id.KeyVaultBaseURL, id.Name, "")
			if err != nil {
				return keyURL, err
			}

			if keyBundle.Key != nil {
				keyURL = pointer.From(keyBundle.Key.Kid)
			}
		}
	} else {
		if id.Version == "" {
			return keyURL, fmt.Errorf("`auto_key_rotation_enabled` field is set to `false` expected a key vault key with a versioned ID but no version information was found: %s", id)
		}
		keyURL = id.ID()
	}

	if keyURL == "" {
		return keyURL, errors.New("internal-error: received an unexpected empty key URL")
	}

	return keyURL, nil
}

func expandDiskEncryptionSetIdentity(input []interface{}) (*identity.SystemAndUserAssignedMap, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	return expanded, nil
}
