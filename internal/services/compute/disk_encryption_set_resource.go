package compute

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskencryptionsets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	resourcesClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			"key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.NestedItemId,
			},

			"auto_key_rotation_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
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
		},
	}
}

func resourceDiskEncryptionSetCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	resourcesClient := meta.(*clients.Client).Resource
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := diskencryptionsets.NewDiskEncryptionSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for present of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_disk_encryption_set", id.ID())
	}

	keyVaultKeyId := d.Get("key_vault_key_id").(string)
	keyVaultDetails, err := diskEncryptionSetRetrieveKeyVault(ctx, keyVaultsClient, resourcesClient, keyVaultKeyId)
	if err != nil {
		return fmt.Errorf("validating Key Vault Key %q for Disk Encryption Set: %+v", keyVaultKeyId, err)
	}

	if keyVaultDetails != nil {
		if !keyVaultDetails.softDeleteEnabled {
			return fmt.Errorf("validating Key Vault %q (Resource Group %q) for Disk Encryption Set: Soft Delete must be enabled but it isn't!", keyVaultDetails.keyVaultName, keyVaultDetails.resourceGroupName)
		}
	}

	rotationToLatestKeyVersionEnabled := d.Get("auto_key_rotation_enabled").(bool)
	encryptionType := diskencryptionsets.DiskEncryptionSetType(d.Get("encryption_type").(string))
	t := d.Get("tags").(map[string]interface{})

	expandedIdentity, err := expandDiskEncryptionSetIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	params := diskencryptionsets.DiskEncryptionSet{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &diskencryptionsets.EncryptionSetProperties{
			ActiveKey: &diskencryptionsets.KeyForDiskEncryptionSet{
				KeyUrl: keyVaultKeyId,
			},
			RotationToLatestKeyVersionEnabled: utils.Bool(rotationToLatestKeyVersionEnabled),
			EncryptionType:                    &encryptionType,
		},
		Identity: expandedIdentity,
		Tags:     tags.Expand(t),
	}

	if v, ok := d.GetOk("federated_client_id"); ok {
		params.Properties.FederatedClientId = utils.String(v.(string))
	}

	if keyVaultDetails != nil {
		params.Properties.ActiveKey.SourceVault = &diskencryptionsets.SourceVault{
			Id: utils.String(keyVaultDetails.keyVaultId),
		}
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, params)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDiskEncryptionSetRead(d, meta)
}

func resourceDiskEncryptionSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := diskencryptionsets.ParseDiskEncryptionSetID(d.Id())
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
		keyVaultKeyId := ""
		if props.ActiveKey != nil && props.ActiveKey.KeyUrl != "" {
			keyVaultKeyId = props.ActiveKey.KeyUrl
		}
		d.Set("key_vault_key_id", keyVaultKeyId)
		d.Set("auto_key_rotation_enabled", props.RotationToLatestKeyVersionEnabled)

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
	resourcesClient := meta.(*clients.Client).Resource
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := diskencryptionsets.ParseDiskEncryptionSetID(d.Id())
	if err != nil {
		return err
	}

	update := diskencryptionsets.DiskEncryptionSetUpdate{}
	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("key_vault_key_id") {
		keyVaultKeyId := d.Get("key_vault_key_id").(string)
		keyVaultDetails, err := diskEncryptionSetRetrieveKeyVault(ctx, keyVaultsClient, resourcesClient, keyVaultKeyId)
		if err != nil {
			return fmt.Errorf("validating Key Vault Key %q for Disk Encryption Set: %+v", keyVaultKeyId, err)
		}

		if keyVaultDetails != nil {
			if !keyVaultDetails.softDeleteEnabled {
				return fmt.Errorf("validating Key Vault %q (Resource Group %q) for Disk Encryption Set: Soft Delete must be enabled but it isn't!", keyVaultDetails.keyVaultName, keyVaultDetails.resourceGroupName)
			}
			if !keyVaultDetails.purgeProtectionEnabled {
				return fmt.Errorf("validating Key Vault %q (Resource Group %q) for Disk Encryption Set: Purge Protection must be enabled but it isn't!", keyVaultDetails.keyVaultName, keyVaultDetails.resourceGroupName)
			}
		}

		update.Properties = &diskencryptionsets.DiskEncryptionSetUpdateProperties{
			ActiveKey: &diskencryptionsets.KeyForDiskEncryptionSet{
				KeyUrl: keyVaultKeyId,
			},
		}

		if keyVaultDetails != nil {
			update.Properties.ActiveKey.SourceVault = &diskencryptionsets.SourceVault{
				Id: utils.String(keyVaultDetails.keyVaultId),
			}
		}
	}

	if d.HasChange("auto_key_rotation_enabled") {
		if update.Properties == nil {
			update.Properties = &diskencryptionsets.DiskEncryptionSetUpdateProperties{}
		}

		update.Properties.RotationToLatestKeyVersionEnabled = utils.Bool(d.Get("auto_key_rotation_enabled").(bool))
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

func resourceDiskEncryptionSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := diskencryptionsets.ParseDiskEncryptionSetID(d.Id())
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

func diskEncryptionSetRetrieveKeyVault(ctx context.Context, keyVaultsClient *client.Client, resourcesClient *resourcesClient.Client, id string) (*diskEncryptionSetKeyVault, error) {
	keyVaultKeyId, err := keyVaultParse.ParseNestedItemID(id)
	if err != nil {
		return nil, err
	}
	keyVaultID, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, keyVaultKeyId.KeyVaultBaseUrl)
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

	resp, err := keyVaultsClient.VaultsClient.Get(ctx, parsedKeyVaultID.ResourceGroupName, parsedKeyVaultID.VaultName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *parsedKeyVaultID, err)
	}

	purgeProtectionEnabled := false
	softDeleteEnabled := false

	if props := resp.Properties; props != nil {
		if props.EnableSoftDelete != nil {
			softDeleteEnabled = *props.EnableSoftDelete
		}

		if props.EnablePurgeProtection != nil {
			purgeProtectionEnabled = *props.EnablePurgeProtection
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
