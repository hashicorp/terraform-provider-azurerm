package cognitive

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/mgmt/2021-04-30/cognitiveservices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			"key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
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
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if d.IsNewResource() {
		if resp.Properties != nil && resp.Properties.Encryption != nil && resp.Properties.Encryption.KeySource != cognitiveservices.KeySourceMicrosoftCognitiveServices {
			return tf.ImportAsExistsError("azurerm_cognitive_account_customer_managed_key", id.ID())
		}
	}

	props := cognitiveservices.Account{
		Properties: &cognitiveservices.AccountProperties{
			Encryption: &cognitiveservices.Encryption{
				KeySource: cognitiveservices.KeySourceMicrosoftKeyVault,
			},
		},
	}

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

	if _, err = client.Update(ctx, id.ResourceGroup, id.Name, props); err != nil {
		return fmt.Errorf("adding Customer Managed Key for %s: %+v", id, err)
	}

	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Accepted"},
		Target:     []string{"Succeeded"},
		Refresh:    cognitiveAccountStateRefreshFunc(ctx, client, *id),
		MinTimeout: 15 * time.Second,
		Timeout:    time.Until(timeout),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
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
	if resp.Properties == nil || resp.Properties.Encryption == nil {
		d.SetId("")
		return nil
	}
	if resp.Properties.Encryption.KeySource == cognitiveservices.KeySourceMicrosoftCognitiveServices {
		d.SetId("")
		return nil
	}

	d.Set("cognitive_account_id", id.ID())
	if props := resp.Properties.Encryption.KeyVaultProperties; props != nil {
		keyVaultKeyId, err := keyVaultParse.NewNestedItemID(*props.KeyVaultURI, "keys", *props.KeyName, *props.KeyVersion)
		if err != nil {
			return fmt.Errorf("parsing `key_vault_key_id`: %+v", err)
		}
		d.Set("key_vault_key_id", keyVaultKeyId.ID())
		d.Set("identity_client_id", props.IdentityClientID)
	}

	return nil
}

func resourceCognitiveAccountCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, "azurerm_cognitive_account")
	defer locks.UnlockByName(id.Name, "azurerm_cognitive_account")

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.Properties == nil || resp.Properties.Encryption == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	// set key source to Microsoft.CognitiveServices to disable customer managed key
	props := cognitiveservices.Account{
		Properties: &cognitiveservices.AccountProperties{
			Encryption: &cognitiveservices.Encryption{
				KeySource: cognitiveservices.KeySourceMicrosoftCognitiveServices,
			},
		},
	}

	if _, err = client.Update(ctx, id.ResourceGroup, id.Name, props); err != nil {
		return fmt.Errorf("removing Customer Managed Key for %s: %+v", *id, err)
	}

	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Accepted"},
		Target:     []string{"Succeeded"},
		Refresh:    cognitiveAccountStateRefreshFunc(ctx, client, *id),
		MinTimeout: 15 * time.Second,
		Timeout:    time.Until(timeout),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for removal of Customer Managed Key for %s: %+v", *id, err)
	}

	return nil
}
