package azurerm

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/dataplane/keyvault"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
				ValidateFunc: validateKeyVaultChildName,
			},

			"vault_uri": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

	_, err := client.CreateKey(keyVaultBaseUrl, name, parameters)
	if err != nil {
		return fmt.Errorf("Error Creating Key: %+v", err)
	}

	// "" indicates the latest version
	read, err := client.GetKey(keyVaultBaseUrl, name, "")
	if err != nil {
		return err
	}

	d.SetId(*read.Key.Kid)

	return resourceArmKeyVaultKeyRead(d, meta)
}

func resourceArmKeyVaultKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultManagementClient

	log.Print("[INFO] preparing arguments for AzureRM KeyVault Key update.")
	id, err := parseKeyVaultKeyID(d.Id())
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

	_, err = client.UpdateKey(id.KeyVaultBaseUrl, id.Name, id.Version, parameters)
	if err != nil {
		return err
	}

	return resourceArmKeyVaultKeyRead(d, meta)
}

func resourceArmKeyVaultKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultManagementClient

	id, err := parseKeyVaultKeyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetKey(id.KeyVaultBaseUrl, id.Name, "")
	if err != nil {
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

	id, err := parseKeyVaultKeyID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.DeleteKey(id.KeyVaultBaseUrl, id.Name)

	return err
}

func parseKeyVaultKeyID(id string) (*KeyVaultKey, error) {
	// example: https://tharvey-keyvault.vault.azure.net/keys/bird/fdf067c93bbb4b22bff4d8b7a9a56217
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse Azure KeyVault Key Id: %q", err)
	}

	path := idURL.Path

	path = strings.TrimSpace(path)
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	components := strings.Split(path, "/")

	if len(components) != 3 {
		return nil, fmt.Errorf("Azure KeyVault Key Id should have 3 segments, got %d: %q", len(components), path)
	}

	key := KeyVaultKey{
		KeyVaultBaseUrl: fmt.Sprintf("%s://%s/", idURL.Scheme, idURL.Host),
		Name:            components[1],
		Version:         components[2],
	}

	return &key, nil
}

type KeyVaultKey struct {
	KeyVaultBaseUrl string
	Name            string
	Version         string
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
