package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKeyVaultKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKeyVaultKeyCreate,
		Read:   resourceArmKeyVaultKeyRead,
		Update: resourceArmKeyVaultKeyUpdate,
		Delete: resourceArmKeyVaultKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateKeyVaultChildName,
			},

			"vault_uri": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.URLIsHTTPS,
			},

			"key_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// turns out Azure's *really* sensitive about the casing of these
				// issue: https://github.com/Azure/azure-rest-api-specs/issues/1739
				ValidateFunc: validation.StringInSlice([]string{
					// TODO: add `oct` back in once this is fixed
					// https://github.com/Azure/azure-rest-api-specs/issues/1739#issuecomment-332236257
					string(keyvault.EC),
					string(keyvault.RSA),
					string(keyvault.RSAHSM),
				}, false),
			},

			"key_size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"key_opts": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					// turns out Azure's *really* sensitive about the casing of these
					// issue: https://github.com/Azure/azure-rest-api-specs/issues/1739
					ValidateFunc: validation.StringInSlice([]string{
						string(keyvault.Decrypt),
						string(keyvault.Encrypt),
						string(keyvault.Sign),
						string(keyvault.UnwrapKey),
						string(keyvault.Verify),
						string(keyvault.WrapKey),
					}, false),
				},
			},

			// Computed
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"n": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"e": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmKeyVaultKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultManagementClient
	ctx := meta.(*ArmClient).StopContext

	log.Print("[INFO] preparing arguments for AzureRM KeyVault Key creation.")
	name := d.Get("name").(string)
	keyVaultBaseUrl := d.Get("vault_uri").(string)

	keyType := d.Get("key_type").(string)
	keyOptions := expandKeyVaultKeyOptions(d)
	tags := d.Get("tags").(map[string]interface{})

	// TODO: support Importing Keys once this is fixed:
	// https://github.com/Azure/azure-rest-api-specs/issues/1747
	parameters := keyvault.KeyCreateParameters{
		Kty:    keyvault.JSONWebKeyType(keyType),
		KeyOps: keyOptions,
		KeyAttributes: &keyvault.KeyAttributes{
			Enabled: utils.Bool(true),
		},
		KeySize: utils.Int32(int32(d.Get("key_size").(int))),
		Tags:    expandTags(tags),
	}

	if _, err := client.CreateKey(ctx, keyVaultBaseUrl, name, parameters); err != nil {
		return fmt.Errorf("Error Creating Key: %+v", err)
	}

	// "" indicates the latest version
	read, err := client.GetKey(ctx, keyVaultBaseUrl, name, "")
	if err != nil {
		return err
	}

	d.SetId(*read.Key.Kid)

	return resourceArmKeyVaultKeyRead(d, meta)
}

func resourceArmKeyVaultKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultManagementClient
	ctx := meta.(*ArmClient).StopContext

	log.Print("[INFO] preparing arguments for AzureRM KeyVault Key update.")
	id, err := azure.ParseKeyVaultChildID(d.Id())
	if err != nil {
		return err
	}

	keyOptions := expandKeyVaultKeyOptions(d)
	tags := d.Get("tags").(map[string]interface{})

	parameters := keyvault.KeyUpdateParameters{
		KeyOps: keyOptions,
		KeyAttributes: &keyvault.KeyAttributes{
			Enabled: utils.Bool(true),
		},
		Tags: expandTags(tags),
	}

	_, err = client.UpdateKey(ctx, id.KeyVaultBaseUrl, id.Name, id.Version, parameters)
	if err != nil {
		return err
	}

	return resourceArmKeyVaultKeyRead(d, meta)
}

func resourceArmKeyVaultKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultManagementClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseKeyVaultChildID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetKey(ctx, id.KeyVaultBaseUrl, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Key %q was not found in Key Vault at URI %q - removing from state", id.Name, id.KeyVaultBaseUrl)
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", id.Name)
	d.Set("vault_uri", id.KeyVaultBaseUrl)
	if key := resp.Key; key != nil {
		d.Set("key_type", string(key.Kty))

		options := flattenKeyVaultKeyOptions(key.KeyOps)
		if err := d.Set("key_opts", options); err != nil {
			return err
		}

		d.Set("n", key.N)
		d.Set("e", key.E)
	}

	// Computed
	d.Set("version", id.Version)

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmKeyVaultKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultManagementClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseKeyVaultChildID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.DeleteKey(ctx, id.KeyVaultBaseUrl, id.Name)
	return err
}

func expandKeyVaultKeyOptions(d *schema.ResourceData) *[]keyvault.JSONWebKeyOperation {
	options := d.Get("key_opts").([]interface{})
	results := make([]keyvault.JSONWebKeyOperation, 0, len(options))

	for _, option := range options {
		results = append(results, keyvault.JSONWebKeyOperation(option.(string)))
	}

	return &results
}

func flattenKeyVaultKeyOptions(input *[]string) []interface{} {
	results := make([]interface{}, 0, len(*input))

	for _, option := range *input {
		results = append(results, option)
	}

	return results
}
