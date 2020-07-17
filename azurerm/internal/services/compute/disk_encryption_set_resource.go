package compute

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDiskEncryptionSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDiskEncryptionSetCreateUpdate,
		Read:   resourceArmDiskEncryptionSetRead,
		Update: resourceArmDiskEncryptionSetCreateUpdate,
		Delete: resourceArmDiskEncryptionSetDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DiskEncryptionSetID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDiskEncryptionSetName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"key_vault_key_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateKeyVaultChildId,
			},

			"identity": {
				Type: schema.TypeList,
				// whilst the API Documentation shows optional - attempting to send nothing returns:
				// `Required parameter 'ResourceIdentity' is missing (null)`
				// hence this is required
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.SystemAssigned),
							}, false),
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDiskEncryptionSetCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	vaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_disk_encryption_set", *existing.ID)
		}
	}

	keyVaultKeyId := d.Get("key_vault_key_id").(string)
	keyVaultDetails, err := diskEncryptionSetRetrieveKeyVault(ctx, vaultClient, keyVaultKeyId)
	if err != nil {
		return fmt.Errorf("Error validating Key Vault Key %q for Disk Encryption Set: %+v", keyVaultKeyId, err)
	}
	if !keyVaultDetails.softDeleteEnabled {
		return fmt.Errorf("Error validating Key Vault %q (Resource Group %q) for Disk Encryption Set: Soft Delete must be enabled but it isn't!", keyVaultDetails.keyVaultName, keyVaultDetails.resourceGroupName)
	}
	if !keyVaultDetails.purgeProtectionEnabled {
		return fmt.Errorf("Error validating Key Vault %q (Resource Group %q) for Disk Encryption Set: Purge Protection must be enabled but it isn't!", keyVaultDetails.keyVaultName, keyVaultDetails.resourceGroupName)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	identityRaw := d.Get("identity").([]interface{})
	t := d.Get("tags").(map[string]interface{})

	params := compute.DiskEncryptionSet{
		Location: utils.String(location),
		EncryptionSetProperties: &compute.EncryptionSetProperties{
			ActiveKey: &compute.KeyVaultAndKeyReference{
				KeyURL: utils.String(keyVaultKeyId),
				SourceVault: &compute.SourceVault{
					ID: utils.String(keyVaultDetails.keyVaultId),
				},
			},
		},
		Identity: expandArmDiskEncryptionSetIdentity(identityRaw),
		Tags:     tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, params)
	if err != nil {
		return fmt.Errorf("Error creating Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Disk Encryption Set %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmDiskEncryptionSetRead(d, meta)
}

func resourceArmDiskEncryptionSetRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error reading Disk Encryption Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
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
	}

	if err := d.Set("identity", flattenArmDiskEncryptionSetIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("Error setting `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDiskEncryptionSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DiskEncryptionSetID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Disk Encryption Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deleting Disk Encryption Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandArmDiskEncryptionSetIdentity(input []interface{}) *compute.EncryptionSetIdentity {
	val := input[0].(map[string]interface{})
	return &compute.EncryptionSetIdentity{
		Type: compute.DiskEncryptionSetIdentityType(val["type"].(string)),
	}
}

func flattenArmDiskEncryptionSetIdentity(input *compute.EncryptionSetIdentity) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	identityType := string(input.Type)
	principalId := ""
	if input.PrincipalID != nil {
		principalId = *input.PrincipalID
	}
	tenantId := ""
	if input.TenantID != nil {
		tenantId = *input.TenantID
	}

	return []interface{}{
		map[string]interface{}{
			"type":         identityType,
			"principal_id": principalId,
			"tenant_id":    tenantId,
		},
	}
}

type diskEncryptionSetKeyVault struct {
	keyVaultId             string
	resourceGroupName      string
	keyVaultName           string
	purgeProtectionEnabled bool
	softDeleteEnabled      bool
}

func diskEncryptionSetRetrieveKeyVault(ctx context.Context, meta interface{}, id string) (*diskEncryptionSetKeyVault, error) {
	vaultClient := meta.(*clients.Client).KeyVault

	keyVaultKeyId, err := azure.ParseKeyVaultChildID(id)
	if err != nil {
		return nil, err
	}

	vault, err := vaultClient.FindKeyVault(ctx, keyVaultKeyId.KeyVaultBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("retrieving the Resource ID for the Key Vault at URL %q: %s", keyVaultKeyId.KeyVaultBaseUrl, err)
	}
	if vault == nil {
		return nil, fmt.Errorf("retrieving key vault %q", keyVaultKeyId.KeyVaultBaseUrl)
	}

	purgeProtectionEnabled := false
	softDeleteEnabled := false

	if props := vault.Properties; props != nil {
		if props.EnableSoftDelete != nil {
			softDeleteEnabled = *props.EnableSoftDelete
		}

		if props.EnablePurgeProtection != nil {
			purgeProtectionEnabled = *props.EnablePurgeProtection
		}
	}

	return &diskEncryptionSetKeyVault{
		keyVaultId:             vault.ID,
		resourceGroupName:      vault.ResourceGroup,
		keyVaultName:           vault.Name,
		purgeProtectionEnabled: purgeProtectionEnabled,
		softDeleteEnabled:      softDeleteEnabled,
	}, nil
}
