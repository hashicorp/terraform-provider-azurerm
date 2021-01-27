package keyvault

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceKeyVaultSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeyVaultSecretCreate,
		Read:   resourceKeyVaultSecretRead,
		Update: resourceKeyVaultSecretUpdate,
		Delete: resourceKeyVaultSecretDelete,
		Importer: &schema.ResourceImporter{
			State: nestedItemResourceImporter,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: keyVaultValidate.NestedItemName,
			},

			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: keyVaultValidate.VaultID,
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

			"not_before_date": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"expiration_date": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceKeyVaultSecretCreate(d *schema.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Print("[INFO] preparing arguments for AzureRM KeyVault Secret creation.")

	name := d.Get("name").(string)
	keyVaultId, err := parse.VaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUrl, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("Error looking up Secret %q vault url from id %q: %+v", name, *keyVaultId, err)
	}

	existing, err := client.GetSecret(ctx, *keyVaultBaseUrl, name, "")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing Secret %q (Key Vault %q): %s", name, *keyVaultBaseUrl, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_key_vault_secret", *existing.ID)
	}

	value := d.Get("value").(string)
	contentType := d.Get("content_type").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := keyvault.SecretSetParameters{
		Value:            utils.String(value),
		ContentType:      utils.String(contentType),
		Tags:             tags.Expand(t),
		SecretAttributes: &keyvault.SecretAttributes{},
	}

	if v, ok := d.GetOk("not_before_date"); ok {
		notBeforeDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		notBeforeUnixTime := date.UnixTime(notBeforeDate)
		parameters.SecretAttributes.NotBefore = &notBeforeUnixTime
	}

	if v, ok := d.GetOk("expiration_date"); ok {
		expirationDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		expirationUnixTime := date.UnixTime(expirationDate)
		parameters.SecretAttributes.Expires = &expirationUnixTime
	}

	if resp, err := client.SetSecret(ctx, *keyVaultBaseUrl, name, parameters); err != nil {
		// In the case that the Secret already exists in a Soft Deleted / Recoverable state we check if `recover_soft_deleted_key_vaults` is set
		// and attempt recovery where appropriate
		if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedKeyVaults && utils.ResponseWasConflict(resp.Response) {
			recoveredSecret, err := client.RecoverDeletedSecret(ctx, *keyVaultBaseUrl, name)
			if err != nil {
				return err
			}
			log.Printf("[DEBUG] Recovering Secret %q with ID: %q", name, *recoveredSecret.ID)
			// We need to wait for consistency, recovered Key Vault Child items are not as readily available as newly created
			if secret := recoveredSecret.ID; secret != nil {
				stateConf := &resource.StateChangeConf{
					Pending:                   []string{"pending"},
					Target:                    []string{"available"},
					Refresh:                   keyVaultChildItemRefreshFunc(*secret),
					Delay:                     30 * time.Second,
					PollInterval:              10 * time.Second,
					ContinuousTargetOccurence: 10,
					Timeout:                   d.Timeout(schema.TimeoutCreate),
				}

				if _, err := stateConf.WaitForState(); err != nil {
					return fmt.Errorf("Error waiting for Key Vault Secret %q to become available: %s", name, err)
				}
				log.Printf("[DEBUG] Secret %q recovered with ID: %q", name, *recoveredSecret.ID)
			}
		} else {
			// If the error response was anything else, or `recover_soft_deleted_key_vaults` is `false` just return the error
			return err
		}
	}

	// "" indicates the latest version
	read, err := client.GetSecret(ctx, *keyVaultBaseUrl, name, "")
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read KeyVault Secret '%s' (in key vault '%s')", name, *keyVaultBaseUrl)
	}

	d.SetId(*read.ID)

	return resourceKeyVaultSecretRead(d, meta)
}

func resourceKeyVaultSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	resourcesClient := meta.(*clients.Client).Resource
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Print("[INFO] preparing arguments for AzureRM KeyVault Secret update.")

	id, err := parse.ParseNestedItemID(d.Id())
	if err != nil {
		return err
	}

	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("Error retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultIdRaw == nil {
		return fmt.Errorf("Unable to determine the Resource ID for the Key Vault at URL %q", id.KeyVaultBaseUrl)
	}
	keyVaultId, err := parse.VaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
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

	secretAttributes := &keyvault.SecretAttributes{}

	if v, ok := d.GetOk("not_before_date"); ok {
		notBeforeDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		notBeforeUnixTime := date.UnixTime(notBeforeDate)
		secretAttributes.NotBefore = &notBeforeUnixTime
	}

	if v, ok := d.GetOk("expiration_date"); ok {
		expirationDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		expirationUnixTime := date.UnixTime(expirationDate)
		secretAttributes.Expires = &expirationUnixTime
	}

	if d.HasChange("value") {
		// for changing the value of the secret we need to create a new version
		parameters := keyvault.SecretSetParameters{
			Value:            utils.String(value),
			ContentType:      utils.String(contentType),
			Tags:             tags.Expand(t),
			SecretAttributes: secretAttributes,
		}

		if _, err = client.SetSecret(ctx, id.KeyVaultBaseUrl, id.Name, parameters); err != nil {
			return err
		}
	} else {
		parameters := keyvault.SecretUpdateParameters{
			ContentType:      utils.String(contentType),
			Tags:             tags.Expand(t),
			SecretAttributes: secretAttributes,
		}

		if _, err = client.UpdateSecret(ctx, id.KeyVaultBaseUrl, id.Name, "", parameters); err != nil {
			return err
		}
	}

	// "" indicates the latest version
	read, err := client.GetSecret(ctx, id.KeyVaultBaseUrl, id.Name, "")
	if err != nil {
		return fmt.Errorf("getting Key Vault Secret %q : %+v", id.Name, err)
	}

	if _, err = parse.ParseNestedItemID(*read.ID); err != nil {
		return err
	}

	// the ID is suffixed with the secret version
	d.SetId(*read.ID)

	return resourceKeyVaultSecretRead(d, meta)
}

func resourceKeyVaultSecretRead(d *schema.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	resourcesClient := meta.(*clients.Client).Resource
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id())
	if err != nil {
		return err
	}

	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("Error retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultIdRaw == nil {
		log.Printf("[DEBUG] Unable to determine the Resource ID for the Key Vault at URL %q - removing from state!", id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}
	keyVaultId, err := parse.VaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
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
	respID, err := parse.ParseNestedItemID(*resp.ID)
	if err != nil {
		return err
	}

	d.Set("name", respID.Name)
	d.Set("value", resp.Value)
	d.Set("version", respID.Version)
	d.Set("content_type", resp.ContentType)

	if attributes := resp.Attributes; attributes != nil {
		if v := attributes.NotBefore; v != nil {
			d.Set("not_before_date", time.Time(*v).Format(time.RFC3339))
		}

		if v := attributes.Expires; v != nil {
			d.Set("expiration_date", time.Time(*v).Format(time.RFC3339))
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceKeyVaultSecretDelete(d *schema.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	resourcesClient := meta.(*clients.Client).Resource
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id())
	if err != nil {
		return err
	}

	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("Error retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultIdRaw == nil {
		return fmt.Errorf("Unable to determine the Resource ID for the Key Vault at URL %q", id.KeyVaultBaseUrl)
	}
	keyVaultId, err := parse.VaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("Error checking if key vault %q for Secret %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Secret %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	shouldPurge := meta.(*clients.Client).Features.KeyVault.PurgeSoftDeleteOnDestroy
	description := fmt.Sprintf("Secret %q (Key Vault %q)", id.Name, id.KeyVaultBaseUrl)
	deleter := deleteAndPurgeSecret{
		client:      client,
		keyVaultUri: id.KeyVaultBaseUrl,
		name:        id.Name,
	}
	if err := deleteAndOptionallyPurge(ctx, description, shouldPurge, deleter); err != nil {
		return err
	}

	return nil
}

var _ deleteAndPurgeNestedItem = deleteAndPurgeSecret{}

type deleteAndPurgeSecret struct {
	client      *keyvault.BaseClient
	keyVaultUri string
	name        string
}

func (d deleteAndPurgeSecret) DeleteNestedItem(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.DeleteSecret(ctx, d.keyVaultUri, d.name)
	return resp.Response, err
}

func (d deleteAndPurgeSecret) NestedItemHasBeenDeleted(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.GetSecret(ctx, d.keyVaultUri, d.name, "")
	return resp.Response, err
}

func (d deleteAndPurgeSecret) PurgeNestedItem(ctx context.Context) (autorest.Response, error) {
	return d.client.PurgeDeletedSecret(ctx, d.keyVaultUri, d.name)
}

func (d deleteAndPurgeSecret) NestedItemHasBeenPurged(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.GetDeletedSecret(ctx, d.keyVaultUri, d.name)
	return resp.Response, err
}
