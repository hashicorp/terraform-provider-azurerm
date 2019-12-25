package compute

import (
	"context"
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
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
				ValidateFunc: ValidateDiskEncryptionSetName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"active_key": {
				Type:     schema.TypeList,
				Required: true,
				// the forceNew is enabled because currently key rotation is not supported, you cannot change the key vault and key associated with disk encryption set.
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"source_vault_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
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

			"previous_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"source_vault_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
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
	activeKey := d.Get("active_key").([]interface{})
	t := d.Get("tags").(map[string]interface{})

	diskEncryptionSet := compute.DiskEncryptionSet{
		Location: utils.String(location),
		EncryptionSetProperties: &compute.EncryptionSetProperties{
			ActiveKey: expandArmDiskEncryptionSetKeyVaultAndKeyReference(activeKey),
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("identity"); ok {
		diskEncryptionSet.Identity = expandArmDiskEncryptionSetIdentity(v.([]interface{}))
	} else {
		diskEncryptionSet.Identity = &compute.EncryptionSetIdentity{
			Type: compute.SystemAssigned,
		}
	}

	// validate whether the keyvault has soft-delete and purge-protection enabled
	err := validateKeyVaultAndKey(ctx, meta, resourceGroup, diskEncryptionSet.ActiveKey)
	if err != nil {
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

func validateKeyVaultAndKey(ctx context.Context, meta interface{}, resourceGroup string, keyVaultAndKey *compute.KeyVaultAndKeyReference) error {
	armClient := meta.(*clients.Client)
	if keyVaultAndKey == nil {
		return nil
	}
	if keyVault := keyVaultAndKey.SourceVault; keyVault != nil {
		if keyVaultId := keyVault.ID; keyVaultId != nil {
			client := armClient.KeyVault.VaultsClient
			parsedId, err := azure.ParseAzureResourceID(*keyVaultId)
			if err != nil {
				return fmt.Errorf("Error parsing ID for keyvault in Disk Encryption Set: %+v", err)
			}
			keyVaultName := parsedId.Path["name"]
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
		if err := d.Set("active_key", flattenArmDiskEncryptionSetKeyVaultAndKeyReference(encryptionSetProperties.ActiveKey)); err != nil {
			return fmt.Errorf("Error setting `active_key`: %+v", err)
		}
		if err := d.Set("previous_keys", flattenArmDiskEncryptionSetKeyVaultAndKeyReferenceArray(encryptionSetProperties.PreviousKeys)); err != nil {
			return fmt.Errorf("Error setting `previous_keys`: %+v", err)
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
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Disk Encryption Set %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandArmDiskEncryptionSetKeyVaultAndKeyReference(input []interface{}) *compute.KeyVaultAndKeyReference {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	sourceVaultId := v["source_vault_id"].(string)
	keyUrl := v["key_url"].(string)

	result := compute.KeyVaultAndKeyReference{
		KeyURL: utils.String(keyUrl),
		SourceVault: &compute.SourceVault{
			ID: utils.String(sourceVaultId),
		},
	}
	return &result
}

func expandArmDiskEncryptionSetIdentity(input []interface{}) *compute.EncryptionSetIdentity {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	t := v["type"].(string)
	result := compute.EncryptionSetIdentity{
		Type: compute.DiskEncryptionSetIdentityType(t),
	}

	return &result
}

func flattenArmDiskEncryptionSetKeyVaultAndKeyReference(input *compute.KeyVaultAndKeyReference) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if keyUrl := input.KeyURL; keyUrl != nil {
		result["key_url"] = *keyUrl
	}
	if sourceVault := input.SourceVault; sourceVault != nil {
		if sourceVaultId := sourceVault.ID; sourceVaultId != nil {
			result["source_vault_id"] = *sourceVaultId
		}
	}

	return []interface{}{result}
}

func flattenArmDiskEncryptionSetKeyVaultAndKeyReferenceArray(input *[]compute.KeyVaultAndKeyReference) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		if sourceVault := item.SourceVault; sourceVault != nil {
			if sourceVaultId := sourceVault.ID; sourceVaultId != nil {
				v["source_vault_id"] = *sourceVaultId
			}
		}

		results = append(results, v)
	}

	return results
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
