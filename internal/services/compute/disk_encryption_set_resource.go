package compute

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-11-01/compute"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	resourcesClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

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
				Default:  string(compute.DiskEncryptionSetTypeEncryptionAtRestWithCustomerKey),
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.DiskEncryptionSetTypeEncryptionAtRestWithCustomerKey),
					string(compute.DiskEncryptionSetTypeEncryptionAtRestWithPlatformAndCustomerKeys),
					string(compute.DiskEncryptionSetTypeConfidentialVMEncryptedWithCustomerKey),
				}, false),
			},

			// whilst the API Documentation shows optional - attempting to send nothing returns:
			// `Required parameter 'ResourceIdentity' is missing (null)`
			// hence this is required
			"identity": commonschema.SystemAssignedIdentityRequired(),

			"tags": tags.Schema(),
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

	id := parse.NewDiskEncryptionSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_disk_encryption_set", id.ID())
	}

	keyVaultKeyId := d.Get("key_vault_key_id").(string)
	keyVaultDetails, err := diskEncryptionSetRetrieveKeyVault(ctx, keyVaultsClient, resourcesClient, keyVaultKeyId)
	if err != nil {
		return fmt.Errorf("validating Key Vault Key %q for Disk Encryption Set: %+v", keyVaultKeyId, err)
	}
	if !keyVaultDetails.softDeleteEnabled {
		return fmt.Errorf("validating Key Vault %q (Resource Group %q) for Disk Encryption Set: Soft Delete must be enabled but it isn't!", keyVaultDetails.keyVaultName, keyVaultDetails.resourceGroupName)
	}
	if !keyVaultDetails.purgeProtectionEnabled {
		return fmt.Errorf("validating Key Vault %q (Resource Group %q) for Disk Encryption Set: Purge Protection must be enabled but it isn't!", keyVaultDetails.keyVaultName, keyVaultDetails.resourceGroupName)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	rotationToLatestKeyVersionEnabled := d.Get("auto_key_rotation_enabled").(bool)
	encryptionType := d.Get("encryption_type").(string)
	t := d.Get("tags").(map[string]interface{})

	expandedIdentity, err := expandDiskEncryptionSetIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	params := compute.DiskEncryptionSet{
		Location: utils.String(location),
		EncryptionSetProperties: &compute.EncryptionSetProperties{
			ActiveKey: &compute.KeyForDiskEncryptionSet{
				KeyURL: utils.String(keyVaultKeyId),
				SourceVault: &compute.SourceVault{
					ID: utils.String(keyVaultDetails.keyVaultId),
				},
			},
			RotationToLatestKeyVersionEnabled: utils.Bool(rotationToLatestKeyVersionEnabled),
			EncryptionType:                    compute.DiskEncryptionSetType(encryptionType),
		},
		Identity: expandedIdentity,
		Tags:     tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, params)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDiskEncryptionSetRead(d, meta)
}

func resourceDiskEncryptionSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DiskEncryptionSetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Disk Encryption Set %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Disk Encryption Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.EncryptionSetProperties; props != nil {
		keyVaultKeyId := ""
		if props.ActiveKey != nil && props.ActiveKey.KeyURL != nil {
			keyVaultKeyId = *props.ActiveKey.KeyURL
		}
		d.Set("key_vault_key_id", keyVaultKeyId)
		d.Set("auto_key_rotation_enabled", props.RotationToLatestKeyVersionEnabled)

		encryptionType := string(compute.DiskEncryptionSetTypeEncryptionAtRestWithCustomerKey)
		if props.EncryptionType != "" {
			encryptionType = string(props.EncryptionType)
		}
		d.Set("encryption_type", encryptionType)
	}

	if err := d.Set("identity", flattenDiskEncryptionSetIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceDiskEncryptionSetUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	resourcesClient := meta.(*clients.Client).Resource
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DiskEncryptionSetID(d.Id())
	if err != nil {
		return err
	}

	update := compute.DiskEncryptionSetUpdate{}
	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("key_vault_key_id") {
		keyVaultKeyId := d.Get("key_vault_key_id").(string)
		keyVaultDetails, err := diskEncryptionSetRetrieveKeyVault(ctx, keyVaultsClient, resourcesClient, keyVaultKeyId)
		if err != nil {
			return fmt.Errorf("validating Key Vault Key %q for Disk Encryption Set: %+v", keyVaultKeyId, err)
		}
		if !keyVaultDetails.softDeleteEnabled {
			return fmt.Errorf("validating Key Vault %q (Resource Group %q) for Disk Encryption Set: Soft Delete must be enabled but it isn't!", keyVaultDetails.keyVaultName, keyVaultDetails.resourceGroupName)
		}
		if !keyVaultDetails.purgeProtectionEnabled {
			return fmt.Errorf("validating Key Vault %q (Resource Group %q) for Disk Encryption Set: Purge Protection must be enabled but it isn't!", keyVaultDetails.keyVaultName, keyVaultDetails.resourceGroupName)
		}
		update.DiskEncryptionSetUpdateProperties = &compute.DiskEncryptionSetUpdateProperties{
			ActiveKey: &compute.KeyForDiskEncryptionSet{
				KeyURL: utils.String(keyVaultKeyId),
				SourceVault: &compute.SourceVault{
					ID: utils.String(keyVaultDetails.keyVaultId),
				},
			},
		}
	}

	if d.HasChange("auto_key_rotation_enabled") {
		if update.DiskEncryptionSetUpdateProperties == nil {
			update.DiskEncryptionSetUpdateProperties = &compute.DiskEncryptionSetUpdateProperties{}
		}

		update.DiskEncryptionSetUpdateProperties.RotationToLatestKeyVersionEnabled = utils.Bool(d.Get("auto_key_rotation_enabled").(bool))
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.Name, update)
	if err != nil {
		return fmt.Errorf("updating Disk Encryption Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of Disk Encryption Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return resourceDiskEncryptionSetRead(d, meta)
}

func resourceDiskEncryptionSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DiskEncryptionSetID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Disk Encryption Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deleting Disk Encryption Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandDiskEncryptionSetIdentity(input []interface{}) (*compute.EncryptionSetIdentity, error) {
	expanded, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	return &compute.EncryptionSetIdentity{
		Type: compute.DiskEncryptionSetIdentityType(string(expanded.Type)),
	}, nil
}

func flattenDiskEncryptionSetIdentity(input *compute.EncryptionSetIdentity) []interface{} {
	var transform *identity.SystemAssigned

	if input != nil {
		transform = &identity.SystemAssigned{
			Type: identity.Type(string(input.Type)),
		}
		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
	}

	return identity.FlattenSystemAssigned(transform)
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
		return nil, fmt.Errorf("Unable to determine the Resource ID for the Key Vault at URL %q", keyVaultKeyId.KeyVaultBaseUrl)
	}

	parsedKeyVaultID, err := keyVaultParse.VaultID(*keyVaultID)
	if err != nil {
		return nil, err
	}

	resp, err := keyVaultsClient.VaultsClient.Get(ctx, parsedKeyVaultID.ResourceGroup, parsedKeyVaultID.Name)
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
		resourceGroupName:      parsedKeyVaultID.ResourceGroup,
		keyVaultName:           parsedKeyVaultID.Name,
		purgeProtectionEnabled: purgeProtectionEnabled,
		softDeleteEnabled:      softDeleteEnabled,
	}, nil
}
