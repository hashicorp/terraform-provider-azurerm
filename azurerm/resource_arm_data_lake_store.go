package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/datalake/store/mgmt/2016-11-01/account"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"strings"
)

func resourceArmDataLakeStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDateLakeStoreCreate,
		Read:   resourceArmDateLakeStoreRead,
		Update: resourceArmDateLakeStoreUpdate,
		Delete: resourceArmDateLakeStoreDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`\A([a-z0-9]{3,24})\z`),
					"Name can only consist of lowercase letters and numbers, and must be between 3 and 24 characters long",
				),
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"tier": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          string(account.Consumption),
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(account.Consumption),
					string(account.Commitment1TB),
					string(account.Commitment10TB),
					string(account.Commitment100TB),
					string(account.Commitment500TB),
					string(account.Commitment1PB),
					string(account.Commitment5PB),
				}, true),
			},

			"encryption": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true, //true by default, so allow read to set this
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
							ForceNew: true,
						},

						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(account.ServiceManaged),
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(account.ServiceManaged),
								string(account.UserManaged),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},

						//the follow are required if UserManaged is selected
						"key_vault_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"key_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.NoZeroValues,
						},

						"key_version": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.NoZeroValues,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {

			encryptionType, hasEncryptionType := d.GetOk("encryption.0.type")

			if hasEncryptionType && strings.EqualFold(encryptionType.(string), string(account.UserManaged)) {
				if _, hasKeyValueId := d.GetOk("encryption.0.key_vault_id"); !hasKeyValueId {
					//return fmt.Errorf("encryption key_vault_id must be specified if encryption type is UserManaged")
				}
				if _, hasKeyName := d.GetOk("encryption.0.key_name"); !hasKeyName {
					return fmt.Errorf("encryption key_name must be specified if encryption type is UserManaged")
				}
			}

			return nil
		},
	}
}

func resourceArmDateLakeStoreCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataLakeStoreAccountClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM Date Lake Store creation.")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)
	tier := d.Get("tier").(string)
	tags := d.Get("tags").(map[string]interface{})

	dateLakeStore := account.CreateDataLakeStoreAccountParameters{
		Location: &location,
		Tags:     expandTags(tags),
		CreateDataLakeStoreAccountProperties: &account.CreateDataLakeStoreAccountProperties{
			NewTier:          account.TierType(tier),
			EncryptionConfig: expandAzureRmDataLakeStoreEncryptionConfig(d),
		},
	}

	if encryptionEnabled, ok := d.GetOk("encryption.0.enabled"); ok {
		if encryptionEnabled.(bool) {
			dateLakeStore.CreateDataLakeStoreAccountProperties.EncryptionState = account.Enabled
		} else {
			dateLakeStore.CreateDataLakeStoreAccountProperties.EncryptionState = account.Disabled
		}
	}

	future, err := client.Create(ctx, resourceGroup, name, dateLakeStore)
	if err != nil {
		return fmt.Errorf("Error issuing create request for Data Lake Store %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error creating Data Lake Store %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Data Lake Store %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Data Lake Store %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmDateLakeStoreRead(d, meta)
}

func resourceArmDateLakeStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataLakeStoreAccountClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	newTags := d.Get("tags").(map[string]interface{})
	newTier := d.Get("tier").(string)

	props := account.UpdateDataLakeStoreAccountParameters{
		Tags: expandTags(newTags),
		UpdateDataLakeStoreAccountProperties: &account.UpdateDataLakeStoreAccountProperties{
			NewTier:          account.TierType(newTier),
			EncryptionConfig: expandAzureRmDataLakeStoreUpdateEncryptionConfig(d),
		},
	}

	future, err := client.Update(ctx, resourceGroup, name, props)
	if err != nil {
		return fmt.Errorf("Error issuing update request for Data Lake Store %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the update of Data Lake Store %q (Resource Group %q) to commplete: %+v", name, resourceGroup, err)
	}

	return resourceArmDateLakeStoreRead(d, meta)
}

func resourceArmDateLakeStoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataLakeStoreAccountClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["accounts"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] DataLakeStoreAccount '%s' was not found (resource group '%s')", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Data Lake Store %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if properties := resp.DataLakeStoreAccountProperties; properties != nil {
		d.Set("tier", string(properties.CurrentTier))
		d.Set("encryption", flattenAzureRmDataLakeStoreEncryption(properties))

	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmDateLakeStoreDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataLakeStoreAccountClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["accounts"]
	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing delete request for Data Lake Store %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Data Lake Store %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func expandAzureRmDataLakeStoreEncryptionConfig(d *schema.ResourceData) *account.EncryptionConfig {
	blocks, ok := d.GetOk("encryption")
	if !ok {
		return nil
	}

	block := blocks.([]interface{})[0].(map[string]interface{})
	config := account.EncryptionConfig{
		Type: account.EncryptionConfigType(block["type"].(string)),
	}

	if config.Type == account.UserManaged {
		config.KeyVaultMetaInfo = &account.KeyVaultMetaInfo{
			KeyVaultResourceID: utils.String(block["key_vault_id"].(string)),
			EncryptionKeyName:  utils.String(block["key_name"].(string)),
		}

		if v, ok := block["key_version"]; ok {
			config.KeyVaultMetaInfo.EncryptionKeyVersion = utils.String(v.(string))
		}
	}

	return &config
}

func expandAzureRmDataLakeStoreUpdateEncryptionConfig(d *schema.ResourceData) *account.UpdateEncryptionConfig {
	blocks, ok := d.GetOk("encryption")
	if !ok {
		return nil
	}

	block := blocks.([]interface{})[0].(map[string]interface{})
	config := account.UpdateEncryptionConfig{}
	if itemType, ok := block["type"]; ok && strings.EqualFold(itemType.(string), string(account.UserManaged)) {
		if v, ok := block["key_version"]; ok {
			config.KeyVaultMetaInfo.EncryptionKeyVersion = utils.String(v.(string))
		}
	}

	return &config
}

func flattenAzureRmDataLakeStoreEncryption(properties *account.DataLakeStoreAccountProperties) interface{} {
	block := map[string]interface{}{
		"enabled": bool(properties.EncryptionState == account.Enabled),
	}

	if config := properties.EncryptionConfig; config != nil {
		block["type"] = string(config.Type)
		if keyVault := config.KeyVaultMetaInfo; keyVault != nil {
			if v := keyVault.KeyVaultResourceID; v != nil {
				block["key_vault_id"] = *v
			}
			if v := keyVault.EncryptionKeyName; v != nil {
				block["key_name"] = *v
			}
			if v := keyVault.EncryptionKeyName; v != nil {
				block["key_version"] = *v
			}
		}
	}

	return []interface{}{block}
}
