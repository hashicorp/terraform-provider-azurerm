package azurerm

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKeyVaultKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKeyVaultKeyCreate,
		Read:   resourceArmKeyVaultKeyRead,
		Update: resourceArmKeyVaultKeyUpdate,
		Delete: resourceArmKeyVaultKeyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceArmKeyVaultChildResourceImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateKeyVaultChildName,
			},

			"key_vault_id": {
				Type:          schema.TypeString,
				Optional:      true, //todo required in 2.0
				Computed:      true, //todo removed in 2.0
				ForceNew:      true,
				ValidateFunc:  azure.ValidateResourceID,
				ConflictsWith: []string{"vault_uri"},
			},

			//todo remove in 2.0
			"vault_uri": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				Deprecated:    "This property has been deprecated in favour of the key_vault_id property. This will prevent a class of bugs as described in https://github.com/terraform-providers/terraform-provider-azurerm/issues/2396 and will be removed in version 2.0 of the provider",
				ValidateFunc:  validate.URLIsHTTPS,
				ConflictsWith: []string{"key_vault_id"},
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
					string(keyvault.ECHSM),
					string(keyvault.RSA),
					string(keyvault.RSAHSM),
				}, false),
			},

			"key_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"curve"},
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

			"curve": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(keyvault.P256),
					string(keyvault.P384),
					string(keyvault.P521),
					string(keyvault.SECP256K1),
				}, false),
				// TODO: the curve name should probably be mandatory for EC in the future,
				// but handle the diff so that we don't break existing configurations and
				// imported EC keys
				ConflictsWith: []string{"key_size"},
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

			"x": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"y": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmKeyVaultKeyCreate(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*ArmClient).keyvault.VaultsClient
	client := meta.(*ArmClient).keyvault.ManagementClient
	ctx := meta.(*ArmClient).StopContext

	log.Print("[INFO] preparing arguments for AzureRM KeyVault Key creation.")
	name := d.Get("name").(string)
	keyVaultBaseUri := d.Get("vault_uri").(string)
	keyVaultId := d.Get("key_vault_id").(string)

	if keyVaultBaseUri == "" {
		if keyVaultId == "" {
			return fmt.Errorf("one of `key_vault_id` or `vault_uri` must be set")
		}

		pKeyVaultBaseUrl, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("Error looking up Key %q vault url form id %q: %+v", name, keyVaultId, err)
		}

		keyVaultBaseUri = pKeyVaultBaseUrl
	} else {
		id, err := azure.GetKeyVaultIDFromBaseUrl(ctx, vaultClient, keyVaultBaseUri)
		if err != nil {
			return fmt.Errorf("Error unable to find key vault ID from URL %q for certificate %q: %+v", keyVaultBaseUri, name, err)
		}
		d.Set("key_vault_id", id)
	}

	if features.ShouldResourcesBeImported() {
		existing, err := client.GetKey(ctx, keyVaultBaseUri, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Key %q (Key Vault %q): %s", name, keyVaultBaseUri, err)
			}
		}

		if existing.Key != nil && existing.Key.Kid != nil && *existing.Key.Kid != "" {
			return tf.ImportAsExistsError("azurerm_key_vault_key", *existing.Key.Kid)
		}
	}

	keyType := d.Get("key_type").(string)
	keyOptions := expandKeyVaultKeyOptions(d)
	t := d.Get("tags").(map[string]interface{})

	// TODO: support Importing Keys once this is fixed:
	// https://github.com/Azure/azure-rest-api-specs/issues/1747
	parameters := keyvault.KeyCreateParameters{
		Kty:    keyvault.JSONWebKeyType(keyType),
		KeyOps: keyOptions,
		KeyAttributes: &keyvault.KeyAttributes{
			Enabled: utils.Bool(true),
		},

		Tags: tags.Expand(t),
	}

	if parameters.Kty == keyvault.EC || parameters.Kty == keyvault.ECHSM {
		curveName := d.Get("curve").(string)
		parameters.Curve = keyvault.JSONWebKeyCurveName(curveName)
	} else if parameters.Kty == keyvault.RSA || parameters.Kty == keyvault.RSAHSM {
		keySize, ok := d.GetOk("key_size")
		if !ok {
			return fmt.Errorf("Key size is required when creating an RSA key")
		}
		parameters.KeySize = utils.Int32(int32(keySize.(int)))
	}
	// TODO: support `oct` once this is fixed
	// https://github.com/Azure/azure-rest-api-specs/issues/1739#issuecomment-332236257

	if _, err := client.CreateKey(ctx, keyVaultBaseUri, name, parameters); err != nil {
		return fmt.Errorf("Error Creating Key: %+v", err)
	}

	// "" indicates the latest version
	read, err := client.GetKey(ctx, keyVaultBaseUri, name, "")
	if err != nil {
		return err
	}

	d.SetId(*read.Key.Kid)

	return resourceArmKeyVaultKeyRead(d, meta)
}

func resourceArmKeyVaultKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*ArmClient).keyvault.VaultsClient
	client := meta.(*ArmClient).keyvault.ManagementClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseKeyVaultChildID(d.Id())
	if err != nil {
		return err
	}

	keyVaultId, err := azure.GetKeyVaultIDFromBaseUrl(ctx, vaultClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("Error retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultId == nil {
		return fmt.Errorf("Unable to determine the Resource ID for the Key Vault at URL %q", id.KeyVaultBaseUrl)
	}

	ok, err := azure.KeyVaultExists(ctx, vaultClient, *keyVaultId)
	if err != nil {
		return fmt.Errorf("Error checking if key vault %q for Key %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Key %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	keyOptions := expandKeyVaultKeyOptions(d)
	t := d.Get("tags").(map[string]interface{})

	parameters := keyvault.KeyUpdateParameters{
		KeyOps: keyOptions,
		KeyAttributes: &keyvault.KeyAttributes{
			Enabled: utils.Bool(true),
		},
		Tags: tags.Expand(t),
	}

	if _, err = client.UpdateKey(ctx, id.KeyVaultBaseUrl, id.Name, id.Version, parameters); err != nil {
		return err
	}

	return resourceArmKeyVaultKeyRead(d, meta)
}

func resourceArmKeyVaultKeyRead(d *schema.ResourceData, meta interface{}) error {
	keyVaultClient := meta.(*ArmClient).keyvault.VaultsClient
	client := meta.(*ArmClient).keyvault.ManagementClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseKeyVaultChildID(d.Id())
	if err != nil {
		return err
	}

	keyVaultId, err := azure.GetKeyVaultIDFromBaseUrl(ctx, keyVaultClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("Error retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultId == nil {
		log.Printf("[DEBUG] Unable to determine the Resource ID for the Key Vault at URL %q - removing from state!", id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	ok, err := azure.KeyVaultExists(ctx, keyVaultClient, *keyVaultId)
	if err != nil {
		return fmt.Errorf("Error checking if key vault %q for Key %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Key %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
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
		d.Set("x", key.X)
		d.Set("y", key.Y)
		if key.N != nil {
			nBytes, err := base64.RawURLEncoding.DecodeString(*key.N)
			if err != nil {
				return fmt.Errorf("Could not decode N: %+v", err)
			}
			d.Set("key_size", len(nBytes)*8)
		}

		d.Set("curve", key.Crv)
	}

	// Computed
	d.Set("version", id.Version)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmKeyVaultKeyDelete(d *schema.ResourceData, meta interface{}) error {
	keyVaultClient := meta.(*ArmClient).keyvault.VaultsClient
	client := meta.(*ArmClient).keyvault.ManagementClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseKeyVaultChildID(d.Id())
	if err != nil {
		return err
	}

	keyVaultId, err := azure.GetKeyVaultIDFromBaseUrl(ctx, keyVaultClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("Error retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultId == nil {
		return fmt.Errorf("Unable to determine the Resource ID for the Key Vault at URL %q", id.KeyVaultBaseUrl)
	}

	ok, err := azure.KeyVaultExists(ctx, keyVaultClient, *keyVaultId)
	if err != nil {
		return fmt.Errorf("Error checking if key vault %q for Key %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Key %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
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
