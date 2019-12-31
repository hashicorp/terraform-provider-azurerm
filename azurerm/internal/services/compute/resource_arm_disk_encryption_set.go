package compute

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
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

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDiskEncryptionSetName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"key_vault_key_uri": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.SystemAssigned),
							}, false),
							Default: string(compute.SystemAssigned),
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
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
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

	location := azure.NormalizeLocation(d.Get("location").(string))
	//activeKey := d.Get("active_key").([]interface{})
	t := d.Get("tags").(map[string]interface{})

	keyURL := d.Get("key_vault_key_uri").(string)
	vaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	keyID, err := azure.ParseKeyVaultChildID(keyURL)
	if err != nil {
		return fmt.Errorf("Error creating Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	keyVaultID, err := azure.GetKeyVaultIDFromBaseUrl(ctx, vaultClient, keyID.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("Error retrieving the Resource ID the Key Vault at URL %q: %s", keyID.KeyVaultBaseUrl, err)
	}
	if keyVaultID == nil {
		return fmt.Errorf("Error creating Disk Encryption Set %q (Resource Group %q): Unable to determine the Resource ID for the Key Vault at URL %q", name, resourceGroup, keyID.KeyVaultBaseUrl)
	}

	diskEncryptionSet := compute.DiskEncryptionSet{
		Location: utils.String(location),
		EncryptionSetProperties: &compute.EncryptionSetProperties{
			ActiveKey: &compute.KeyVaultAndKeyReference{
				KeyURL: &keyURL,
				SourceVault: &compute.SourceVault{
					ID: keyVaultID,
				},
			},
		},
		Tags: tags.Expand(t),
	}

	diskEncryptionSet.Identity = expandArmDiskEncryptionSetIdentity(d)

	// validate whether the keyvault has soft-delete and purge-protection enabled
	if err := validateKeyVault(ctx, meta.(*clients.Client), resourceGroup, *keyVaultID); err != nil {
		return fmt.Errorf("Error creating Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, diskEncryptionSet)
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

func validateKeyVault(ctx context.Context, armClient *clients.Client, resourceGroup string, keyVaultID string) error {
	client := armClient.KeyVault.VaultsClient
	parsedId, err := azure.ParseAzureResourceID(keyVaultID)
	if err != nil {
		return fmt.Errorf("Error parsing ID for keyvault in Disk Encryption Set: %+v", err)
	}
	keyVaultName := parsedId.Path["vaults"]
	log.Printf("[INFO] Keyvault name input in Disk Encryption Set: %s", keyVaultName)
	resp, err := client.Get(ctx, resourceGroup, keyVaultName)
	if err != nil {
		return fmt.Errorf("Error reading keyvault %q (Resource Group %q): %+v", keyVaultName, resourceGroup, err)
	}
	if props := resp.Properties; props != nil {
		if softDelete := props.EnableSoftDelete; softDelete != nil {
			if !*softDelete {
				return fmt.Errorf("the keyvault in Disk Encryption Set must enable soft delete")
			}
		}
		if purgeProtection := props.EnablePurgeProtection; purgeProtection != nil {
			if !*purgeProtection {
				return fmt.Errorf("the keyvault in Disk Encryption Set must enable purge protection")
			}
		}
	}
	return nil
}

func resourceArmDiskEncryptionSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["diskEncryptionSets"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Disk Encryption Set %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if encryptionSetProperties := resp.EncryptionSetProperties; encryptionSetProperties != nil {
		if err := d.Set("key_vault_key_uri", encryptionSetProperties.ActiveKey.KeyURL); err != nil {
			return fmt.Errorf("Error setting `active_key`: %+v", err)
		}
	}
	if identity := resp.Identity; identity != nil {
		if err := d.Set("identity", flattenArmDiskEncryptionSetIdentity(identity)); err != nil {
			return fmt.Errorf("Error setting `identity`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDiskEncryptionSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DiskEncryptionSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["diskEncryptionSets"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deleting Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func expandArmDiskEncryptionSetIdentity(d *schema.ResourceData) *compute.EncryptionSetIdentity {
	if v, ok := d.GetOk("identity"); ok {
		input := v.([]interface{})[0].(map[string]interface{})
		t := input["type"].(string)
		return &compute.EncryptionSetIdentity{
			Type: compute.DiskEncryptionSetIdentityType(t),
		}
	}
	return &compute.EncryptionSetIdentity{
		Type: compute.SystemAssigned,
	}
}

func flattenArmDiskEncryptionSetIdentity(input *compute.EncryptionSetIdentity) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["type"] = string(input.Type)
	if principalId := input.PrincipalID; principalId != nil {
		result["principal_id"] = *principalId
	}
	if tenantId := input.TenantID; tenantId != nil {
		result["tenant_id"] = *tenantId
	}
	return []interface{}{result}
}
