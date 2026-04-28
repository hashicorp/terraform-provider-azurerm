// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/deletedsecrets"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/secrets"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultSecret() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKeyVaultSecretCreate,
		Read:   resourceKeyVaultSecretRead,
		Update: resourceKeyVaultSecretUpdate,
		Delete: resourceKeyVaultSecretDelete,
		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.ParseNestedItemID(id)
			return err
		}, nestedItemResourceImporter),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			// TODO: Change this back to 5min, once https://github.com/hashicorp/terraform-provider-azurerm/issues/11059 is addressed.
			Read:   pluginsdk.DefaultTimeout(30 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: keyVaultValidate.NestedItemName,
			},

			"key_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KeyVaultId{}),

			"value": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ExactlyOneOf: []string{"value", "value_wo"},
			},

			"value_wo": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				WriteOnly:    true,
				RequiredWith: []string{"value_wo_version"},
				ExactlyOneOf: []string{"value", "value_wo"},
			},

			"value_wo_version": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				RequiredWith: []string{"value_wo"},
				ValidateFunc: validation.IntAtLeast(1),
			},

			"content_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"not_before_date": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"expiration_date": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"versionless_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"resource_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"resource_versionless_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsWithMaximumElements(15),
		},
	}
}

func resourceKeyVaultSecretCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Print("[INFO] preparing arguments for AzureRM KeyVault Secret creation.")

	name := d.Get("name").(string)
	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUrl, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("looking up Secret %q vault url from id %q: %+v", name, *keyVaultId, err)
	}

	client := meta.(*clients.Client).KeyVault.DataPlaneKeyVaultClient.Secrets.Clone(*keyVaultBaseUrl)
	secretId := secrets.NewSecretID(*keyVaultBaseUrl, name)
	secretVersionId := secrets.NewSecretversionID(secretId.BaseURI, secretId.SecretName, "")

	existing, err := client.GetSecret(ctx, secretVersionId)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing Secret %q (Key Vault %q): %s", name, *keyVaultBaseUrl, err)
		}
	}

	if model := existing.Model; model != nil && pointer.From(model.Id) != "" {
		return tf.ImportAsExistsError("azurerm_key_vault_secret", *model.Id)
	}

	value := d.Get("value").(string)

	valueWo, err := pluginsdk.GetWriteOnly(d, "value_wo", cty.String)
	if err != nil {
		return err
	}
	if !valueWo.IsNull() {
		value = valueWo.AsString()
	}

	contentType := d.Get("content_type").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := secrets.SecretSetParameters{
		Value:       value,
		ContentType: pointer.To(contentType),
		Tags:        pointer.To(tags.ToTypedObject(tags.Expand(t))),
		Attributes:  &secrets.SecretAttributes{},
	}

	if v, ok := d.GetOk("not_before_date"); ok {
		notBeforeDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		parameters.Attributes.Nbf = pointer.To(notBeforeDate.Unix())
	}

	if v, ok := d.GetOk("expiration_date"); ok {
		expirationDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		parameters.Attributes.Exp = pointer.To(expirationDate.Unix())
	}

	if resp, err := client.SetSecret(ctx, secretId, parameters); err != nil {
		// In the case that the Secret already exists in a Soft Deleted / Recoverable state we check if `recover_soft_deleted_key_vaults` is set
		// and attempt recovery where appropriate
		if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedSecrets && response.WasConflict(resp.HttpResponse) {
			deletedSecretsClient := meta.(*clients.Client).KeyVault.DataPlaneKeyVaultClient.DeletedSecrets.Clone(*keyVaultBaseUrl)
			deletedSecretId := deletedsecrets.NewDeletedsecretID(*keyVaultBaseUrl, name)
			recoveredSecret, err := deletedSecretsClient.RecoverDeletedSecret(ctx, deletedSecretId)
			if err != nil {
				return err
			}
			if recoveredSecret.Model != nil {
				secretIdStr := recoveredSecret.Model.Id
				log.Printf("[DEBUG] Recovering Secret %q with ID: %q", name, pointer.From(secretIdStr))
				// We need to wait for consistency, recovered Key Vault Child items are not as readily available as newly created
				if secretIdStr != nil {
					stateConf := &pluginsdk.StateChangeConf{
						Pending:                   []string{"pending"},
						Target:                    []string{"available"},
						Refresh:                   keyVaultChildItemRefreshFunc(*secretIdStr),
						Delay:                     30 * time.Second,
						PollInterval:              10 * time.Second,
						ContinuousTargetOccurence: 10,
						Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
					}

					if _, err := stateConf.WaitForStateContext(ctx); err != nil {
						return fmt.Errorf("waiting for Key Vault Secret %q to become available: %s", name, err)
					}
					log.Printf("[DEBUG] Secret %q recovered with ID: %q", name, *secretIdStr)

					if _, err := client.SetSecret(ctx, secretId, parameters); err != nil {
						return err
					}
				}
			}
		} else {
			// If the error response was anything else, or `recover_soft_deleted_key_vaults` is `false` just return the error
			return err
		}
	}

	// "" indicates the latest version
	read, err := client.GetSecret(ctx, secretVersionId)
	if err != nil {
		return err
	}

	if read.Model == nil || read.Model.Id == nil {
		return fmt.Errorf("cannot read KeyVault Secret '%s' (in key vault '%s')", name, *keyVaultBaseUrl)
	}

	parsedSecretId, err := parse.ParseNestedItemID(*read.Model.Id)
	if err != nil {
		return err
	}

	d.SetId(parsedSecretId.ID())

	return resourceKeyVaultSecretRead(d, meta)
}

func resourceKeyVaultSecretUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Print("[INFO] preparing arguments for AzureRM KeyVault Secret update.")

	id, err := parse.ParseNestedItemID(d.Id())
	if err != nil {
		return err
	}

	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	meta.(*clients.Client).KeyVault.AddToCache(*keyVaultId, id.KeyVaultBaseUrl)

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking if key vault %q for Secret %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Secret %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	client := meta.(*clients.Client).KeyVault.DataPlaneKeyVaultClient.Secrets.Clone(id.KeyVaultBaseUrl)
	secretId := secrets.NewSecretID(id.KeyVaultBaseUrl, id.Name)
	secretVersionId := secrets.NewSecretversionID(id.KeyVaultBaseUrl, id.Name, "")

	value := d.Get("value").(string)
	contentType := d.Get("content_type").(string)
	t := d.Get("tags").(map[string]interface{})

	secretAttributes := &secrets.SecretAttributes{}

	if v, ok := d.GetOk("not_before_date"); ok {
		notBeforeDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		secretAttributes.Nbf = pointer.To(notBeforeDate.Unix())
	}

	if v, ok := d.GetOk("expiration_date"); ok {
		expirationDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		secretAttributes.Exp = pointer.To(expirationDate.Unix())
	}

	if d.HasChanges("value", "value_wo_version") {
		valueWo, err := pluginsdk.GetWriteOnly(d, "value_wo", cty.String)
		if err != nil {
			return err
		}
		if !valueWo.IsNull() {
			value = valueWo.AsString()
		}
		// for changing the value of the secret we need to create a new version
		parameters := secrets.SecretSetParameters{
			Value:       value,
			ContentType: pointer.To(contentType),
			Tags:        pointer.To(tags.ToTypedObject(tags.Expand(t))),
			Attributes:  secretAttributes,
		}

		if _, err = client.SetSecret(ctx, secretId, parameters); err != nil {
			return err
		}
	} else {
		parameters := secrets.SecretUpdateParameters{
			ContentType: pointer.To(contentType),
			Tags:        pointer.To(tags.ToTypedObject(tags.Expand(t))),
			Attributes:  secretAttributes,
		}

		if _, err = client.UpdateSecret(ctx, secretVersionId, parameters); err != nil {
			return err
		}
	}

	// "" indicates the latest version
	read, err := client.GetSecret(ctx, secretVersionId)
	if err != nil {
		return fmt.Errorf("getting Key Vault Secret %q : %+v", id.Name, err)
	}

	if read.Model == nil || read.Model.Id == nil {
		return fmt.Errorf("cannot read KeyVault Secret '%s' (in key vault '%s')", id.Name, id.KeyVaultBaseUrl)
	}

	parsedSecretId, err := parse.ParseNestedItemID(*read.Model.Id)
	if err != nil {
		return err
	}

	// the ID is suffixed with the secret version
	d.SetId(parsedSecretId.ID())

	return resourceKeyVaultSecretRead(d, meta)
}

func resourceKeyVaultSecretRead(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id())
	if err != nil {
		return err
	}

	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultIdRaw == nil {
		log.Printf("[DEBUG] Unable to determine the Resource ID for the Key Vault at URL %q - removing from state!", id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}
	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking if key vault %q for Secret %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Secret %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	client := meta.(*clients.Client).KeyVault.DataPlaneKeyVaultClient.Secrets.Clone(id.KeyVaultBaseUrl)
	secretVersionId := secrets.NewSecretversionID(id.KeyVaultBaseUrl, id.Name, "")

	// we always want to get the latest version
	resp, err := client.GetSecret(ctx, secretVersionId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Secret %q was not found in Key Vault at URI %q - removing from state", id.Name, id.KeyVaultBaseUrl)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Azure KeyVault Secret %s: %+v", id.Name, err)
	}

	if resp.Model == nil || resp.Model.Id == nil {
		return fmt.Errorf("reading KeyVault Secret %q: response model was nil", id.Name)
	}

	// the version may have changed, so parse the updated id
	respID, err := parse.ParseNestedItemID(*resp.Model.Id)
	if err != nil {
		return err
	}

	d.Set("name", respID.Name)
	d.Set("value", resp.Model.Value)
	// Unset value if is a write-only value
	if _, ok := d.GetOk("value_wo_version"); ok {
		d.Set("value", nil)
	}
	d.Set("version", respID.Version)
	d.Set("content_type", resp.Model.ContentType)
	d.Set("versionless_id", id.VersionlessID())
	d.Set("value_wo_version", d.Get("value_wo_version").(int))

	if attributes := resp.Model.Attributes; attributes != nil {
		notBeforeDate := ""
		if v := attributes.Nbf; v != nil {
			notBeforeDate = time.Unix(*v, 0).UTC().Format(time.RFC3339)
		}
		d.Set("not_before_date", notBeforeDate)

		expirationDate := ""
		if v := attributes.Exp; v != nil {
			expirationDate = time.Unix(*v, 0).UTC().Format(time.RFC3339)
		}
		d.Set("expiration_date", expirationDate)
	}

	d.Set("resource_id", parse.NewSecretID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, id.Name, id.Version).ID())
	d.Set("resource_versionless_id", parse.NewSecretVersionlessID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, id.Name).ID())

	return tags.FlattenAndSet(d, tags.FromTypedObject(pointer.From(resp.Model.Tags)))
}

func resourceKeyVaultSecretDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id())
	if err != nil {
		return err
	}

	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultIdRaw == nil {
		return fmt.Errorf("unable to determine the Resource ID for the Key Vault at URL %q", id.KeyVaultBaseUrl)
	}
	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	kv, err := keyVaultsClient.VaultsClient.Get(ctx, *keyVaultId)
	if err != nil {
		if response.WasNotFound(kv.HttpResponse) {
			log.Printf("[DEBUG] Secret %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("checking if key vault %q for Secret %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}

	shouldPurge := meta.(*clients.Client).Features.KeyVault.PurgeSoftDeletedSecretsOnDestroy
	if shouldPurge && kv.Model != nil && pointer.From(kv.Model.Properties.EnablePurgeProtection) {
		log.Printf("[DEBUG] cannot purge secret %q because %s has purge protection enabled", id.Name, *keyVaultId)
		shouldPurge = false
	}

	description := fmt.Sprintf("Secret %q (Key Vault %q)", id.Name, id.KeyVaultBaseUrl)
	deleter := deleteAndPurgeSecret{
		secretsClient:        keyVaultsClient.DataPlaneKeyVaultClient.Secrets.Clone(id.KeyVaultBaseUrl),
		deletedSecretsClient: keyVaultsClient.DataPlaneKeyVaultClient.DeletedSecrets.Clone(id.KeyVaultBaseUrl),
		keyVaultUri:          id.KeyVaultBaseUrl,
		name:                 id.Name,
	}
	if err := deleteAndOptionallyPurge(ctx, description, shouldPurge, deleter); err != nil {
		return err
	}

	return nil
}

var _ deleteAndPurgeNestedItem = deleteAndPurgeSecret{}

type deleteAndPurgeSecret struct {
	secretsClient        *secrets.SecretsClient
	deletedSecretsClient *deletedsecrets.DeletedSecretsClient
	keyVaultUri          string
	name                 string
}

func (d deleteAndPurgeSecret) DeleteNestedItem(ctx context.Context) (autorest.Response, error) {
	secretId := secrets.NewSecretID(d.keyVaultUri, d.name)
	resp, err := d.secretsClient.DeleteSecret(ctx, secretId)
	return autorest.Response{Response: resp.HttpResponse}, err
}

func (d deleteAndPurgeSecret) NestedItemHasBeenDeleted(ctx context.Context) (autorest.Response, error) {
	secretVersionId := secrets.NewSecretversionID(d.keyVaultUri, d.name, "")
	resp, err := d.secretsClient.GetSecret(ctx, secretVersionId)
	return autorest.Response{Response: resp.HttpResponse}, err
}

func (d deleteAndPurgeSecret) PurgeNestedItem(ctx context.Context) (autorest.Response, error) {
	deletedSecretId := deletedsecrets.NewDeletedsecretID(d.keyVaultUri, d.name)
	resp, err := d.deletedSecretsClient.PurgeDeletedSecret(ctx, deletedSecretId)
	return autorest.Response{Response: resp.HttpResponse}, err
}

func (d deleteAndPurgeSecret) NestedItemHasBeenPurged(ctx context.Context) (autorest.Response, error) {
	deletedSecretId := deletedsecrets.NewDeletedsecretID(d.keyVaultUri, d.name)
	resp, err := d.deletedSecretsClient.GetDeletedSecret(ctx, deletedSecretId)
	return autorest.Response{Response: resp.HttpResponse}, err
}
