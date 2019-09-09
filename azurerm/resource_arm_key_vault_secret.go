package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKeyVaultSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKeyVaultSecretCreate,
		Read:   resourceArmKeyVaultSecretRead,
		Update: resourceArmKeyVaultSecretUpdate,
		Delete: resourceArmKeyVaultSecretDelete,
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

			"value": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmKeyVaultSecretCreate(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*ArmClient).keyvault.VaultsClient
	client := meta.(*ArmClient).keyvault.ManagementClient
	ctx := meta.(*ArmClient).StopContext

	log.Print("[INFO] preparing arguments for AzureRM KeyVault Secret creation.")

	name := d.Get("name").(string)
	keyVaultBaseUrl := d.Get("vault_uri").(string)
	keyVaultId := d.Get("key_vault_id").(string)

	if keyVaultBaseUrl == "" {
		if keyVaultId == "" {
			return fmt.Errorf("one of `key_vault_id` or `vault_uri` must be set")
		}

		pKeyVaultBaseUrl, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("Error looking up Secret %q vault url form id %q: %+v", name, keyVaultId, err)
		}

		keyVaultBaseUrl = pKeyVaultBaseUrl
	} else {
		id, err := azure.GetKeyVaultIDFromBaseUrl(ctx, vaultClient, keyVaultBaseUrl)
		if err != nil {
			return fmt.Errorf("Error unable to find key vault ID from URL %q for certificate %q: %+v", keyVaultBaseUrl, name, err)
		}
		d.Set("key_vault_id", id)
	}

	if features.ShouldResourcesBeImported() {
		existing, err := client.GetSecret(ctx, keyVaultBaseUrl, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Secret %q (Key Vault %q): %s", name, keyVaultBaseUrl, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_key_vault_secret", *existing.ID)
		}
	}

	value := d.Get("value").(string)
	contentType := d.Get("content_type").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := keyvault.SecretSetParameters{
		Value:       utils.String(value),
		ContentType: utils.String(contentType),
		Tags:        tags.Expand(t),
	}

	if _, err := client.SetSecret(ctx, keyVaultBaseUrl, name, parameters); err != nil {
		return err
	}

	// "" indicates the latest version
	read, err := client.GetSecret(ctx, keyVaultBaseUrl, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read KeyVault Secret '%s' (in key vault '%s')", name, keyVaultBaseUrl)
	}

	d.SetId(*read.ID)

	return resourceArmKeyVaultSecretRead(d, meta)
}

func resourceArmKeyVaultSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	keyVaultClient := meta.(*ArmClient).keyvault.VaultsClient
	client := meta.(*ArmClient).keyvault.ManagementClient
	ctx := meta.(*ArmClient).StopContext
	log.Print("[INFO] preparing arguments for AzureRM KeyVault Secret update.")

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
		return fmt.Errorf("Error checking if key vault %q for Secret %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Secret %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	value := d.Get("value").(string)
	contentType := d.Get("content_type").(string)
	t := d.Get("tags").(map[string]interface{})

	if d.HasChange("value") {
		// for changing the value of the secret we need to create a new version
		parameters := keyvault.SecretSetParameters{
			Value:       utils.String(value),
			ContentType: utils.String(contentType),
			Tags:        tags.Expand(t),
		}

		if _, err = client.SetSecret(ctx, id.KeyVaultBaseUrl, id.Name, parameters); err != nil {
			return err
		}

		// "" indicates the latest version
		read, err2 := client.GetSecret(ctx, id.KeyVaultBaseUrl, id.Name, "")
		if err2 != nil {
			return fmt.Errorf("Error getting Key Vault Secret %q : %+v", id.Name, err2)
		}

		if _, err = azure.ParseKeyVaultChildID(*read.ID); err != nil {
			return err
		}

		// the ID is suffixed with the secret version
		d.SetId(*read.ID)
	} else {
		parameters := keyvault.SecretUpdateParameters{
			ContentType: utils.String(contentType),
			Tags:        tags.Expand(t),
		}

		if _, err = client.UpdateSecret(ctx, id.KeyVaultBaseUrl, id.Name, id.Version, parameters); err != nil {
			return err
		}
	}

	return resourceArmKeyVaultSecretRead(d, meta)
}

func resourceArmKeyVaultSecretRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error checking if key vault %q for Secret %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Secret %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	// we always want to get the latest version
	resp, err := client.GetSecret(ctx, id.KeyVaultBaseUrl, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Secret %q was not found in Key Vault at URI %q - removing from state", id.Name, id.KeyVaultBaseUrl)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure KeyVault Secret %s: %+v", id.Name, err)
	}

	// the version may have changed, so parse the updated id
	respID, err := azure.ParseKeyVaultChildID(*resp.ID)
	if err != nil {
		return err
	}

	d.Set("name", respID.Name)
	d.Set("vault_uri", respID.KeyVaultBaseUrl)
	d.Set("value", resp.Value)
	d.Set("version", respID.Version)
	d.Set("content_type", resp.ContentType)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmKeyVaultSecretDelete(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error checking if key vault %q for Secret %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Secret %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	_, err = client.DeleteSecret(ctx, id.KeyVaultBaseUrl, id.Name)
	return err
}
