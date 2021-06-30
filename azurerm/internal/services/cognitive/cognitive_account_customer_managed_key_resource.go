package cognitive

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/mgmt/2021-04-30/cognitiveservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cognitive/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cognitive/validate"
	keyVaultParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceCognitiveAccountCustomerManagedKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCognitiveAccountCustomerManagedKeyCreateUpdate,
		Read:   resourceCognitiveAccountCustomerManagedKeyRead,
		Update: resourceCognitiveAccountCustomerManagedKeyCreateUpdate,
		Delete: resourceCognitiveAccountCustomerManagedKeyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.DefaultImporter(),

		Schema: map[string]*pluginsdk.Schema{
			"cognitive_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountID,
			},

			"key_source": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cognitiveservices.KeySourceMicrosoftKeyVault),
					string(cognitiveservices.KeySourceMicrosoftCognitiveServices),
				}, false),
			},

			"key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
			},

			"identity_client_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceCognitiveAccountCustomerManagedKeyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Get("cognitive_account_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(id.Name, "azurerm_cognitive_account")
	defer locks.UnlockByName(id.Name, "azurerm_cognitive_account")

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if d.IsNewResource() {
		if resp.Properties != nil && resp.Properties.Encryption != nil {
			return tf.ImportAsExistsError("azurerm_cognitive_account_customer_managed_key", id.ID())
		}
	}

	keySource := cognitiveservices.KeySource(d.Get("key_source").(string))
	props := cognitiveservices.Account{
		Properties: &cognitiveservices.AccountProperties{
			Encryption: &cognitiveservices.Encryption{
				KeySource: keySource,
			},
		},
	}
	if keySource == cognitiveservices.KeySourceMicrosoftCognitiveServices {
		if _, ok := d.GetOk("key_vault_key_id"); ok {
			return fmt.Errorf("can't specify key_vault_key_id when using Microsoft.CognitiveServices key source")
		}
		if _, ok := d.GetOk("identity_client_id"); ok {
			return fmt.Errorf("can't specify key_name when using Microsoft.CognitiveServices key source")
		}
	} else {
		keyId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(d.Get("key_vault_key_id").(string))
		if err != nil {
			return err
		}
		props.Properties.Encryption.KeyVaultProperties = &cognitiveservices.KeyVaultProperties{
			KeyName:          utils.String(keyId.Name),
			KeyVersion:       utils.String(keyId.Version),
			KeyVaultURI:      utils.String(keyId.KeyVaultBaseUrl),
			IdentityClientID: utils.String(d.Get("identity_client_id").(string)),
		}
	}

	if _, err = client.Update(ctx, id.ResourceGroup, id.Name, props); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Accepted"},
		Target:     []string{"Succeeded"},
		Refresh:    cognitiveAccountStateRefreshFunc(ctx, client, *id),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", *id, err)
	}
	d.SetId(id.ID())

	return resourceCognitiveAccountCustomerManagedKeyRead(d, meta)
}

func resourceCognitiveAccountCustomerManagedKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if resp.Properties == nil && resp.Properties.Encryption == nil {
		d.SetId("")
		return nil
	}

	d.Set("cognitive_account_id", id.ID())
	d.Set("key_source", resp.Properties.Encryption.KeySource)
	if props := resp.Properties.Encryption.KeyVaultProperties; props != nil {
		var keyName string
		if props.KeyName != nil {
			keyName = *props.KeyName
		}

		var keyVaultUri string
		if props.KeyVaultURI != nil {
			keyVaultUri = *props.KeyVaultURI
		}

		var keyVersion string
		if props.KeyVersion != nil {
			keyVersion = *props.KeyVersion
		}
		keyVaultKeyId, err := keyVaultParse.NewNestedItemID(keyVaultUri, "keys", keyName, keyVersion)
		if err != nil {
			return fmt.Errorf("parsing `key_vault_key_id`: %+v", err)
		}
		d.Set("key_vault_key_id", keyVaultKeyId.ID())
		if props.IdentityClientID != nil {
			d.Set("identity_client_id", *props.IdentityClientID)
		}
	} else {
		d.Set("key_vault_key_id", nil)
		d.Set("identity_client_id", nil)
	}

	return nil
}

func resourceCognitiveAccountCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	accountsClient := meta.(*clients.Client).Cognitive.AccountsClient
	deletedAccountsClient := meta.(*clients.Client).Cognitive.DeletedAccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, "azurerm_cognitive_account")
	defer locks.UnlockByName(id.Name, "azurerm_cognitive_account")

	account, err := accountsClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	// Since this isn't a real object and it cannot be disabled once Customer Managed Key at rest has been enabled
	// And it must keep at least one key once Customer Managed Key is enabled
	// So for the delete operation, it has to recreate the Cognitive Account with disabled Customer Managed Key
	log.Printf("[DEBUG] Deleting %s..", *id)
	deleteFuture, err := accountsClient.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	if err := deleteFuture.WaitForCompletionRef(ctx, accountsClient.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	log.Printf("[DEBUG] Purging %s..", *id)
	purgeFuture, err := deletedAccountsClient.Purge(ctx, *account.Location, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("purging %s: %+v", *id, err)
	}
	if err := purgeFuture.WaitForCompletionRef(ctx, deletedAccountsClient.Client); err != nil {
		return fmt.Errorf("waiting for purge of %s: %+v", *id, err)
	}

	account.SystemData = nil
	account.Properties.Encryption = nil
	if _, err := accountsClient.Create(ctx, id.ResourceGroup, id.Name, account); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    cognitiveAccountStateRefreshFunc(ctx, accountsClient, *id),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	return nil
}
