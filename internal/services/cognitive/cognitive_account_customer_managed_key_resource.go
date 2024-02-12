// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2023-05-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
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

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := cognitiveservicesaccounts.ParseAccountID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"cognitive_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: cognitiveservicesaccounts.ValidateAccountID,
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

	id, err := cognitiveservicesaccounts.ParseAccountID(d.Get("cognitive_account_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(id.AccountName, "azurerm_cognitive_account")
	defer locks.UnlockByName(id.AccountName, "azurerm_cognitive_account")

	resp, err := client.AccountsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if d.IsNewResource() {
		if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.Encryption != nil && resp.Model.Properties.Encryption.KeySource != nil && *resp.Model.Properties.Encryption.KeySource != cognitiveservicesaccounts.KeySourceMicrosoftPointCognitiveServices {
			return tf.ImportAsExistsError("azurerm_cognitive_account_customer_managed_key", id.ID())
		}
	}

	keySource := cognitiveservicesaccounts.KeySourceMicrosoftPointKeyVault

	props := cognitiveservicesaccounts.Account{
		Properties: &cognitiveservicesaccounts.AccountProperties{
			Encryption: &cognitiveservicesaccounts.Encryption{
				KeySource: &keySource,
			},
		},
	}

	keyId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(d.Get("key_vault_key_id").(string))
	if err != nil {
		return err
	}
	props.Properties.Encryption.KeyVaultProperties = &cognitiveservicesaccounts.KeyVaultProperties{
		KeyName:     utils.String(keyId.Name),
		KeyVersion:  utils.String(keyId.Version),
		KeyVaultUri: utils.String(keyId.KeyVaultBaseUrl),
	}

	if identityClientId := d.Get("identity_client_id").(string); identityClientId != "" {
		props.Properties.Encryption.KeyVaultProperties.IdentityClientId = pointer.To(identityClientId)
	}

	// todo check if poll works in all the resources
	if _, err = client.AccountsUpdate(ctx, *id, props); err != nil {
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

	id, err := cognitiveservicesaccounts.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.AccountsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Encryption == nil || resp.Model.Properties.Encryption.KeySource == nil {
		d.SetId("")
		return nil
	}
	if *resp.Model.Properties.Encryption.KeySource == cognitiveservicesaccounts.KeySourceMicrosoftPointCognitiveServices {
		d.SetId("")
		return nil
	}

	d.Set("cognitive_account_id", id.ID())
	if props := resp.Model.Properties.Encryption.KeyVaultProperties; props != nil {
		keyVaultKeyId, err := keyVaultParse.NewNestedItemID(*props.KeyVaultUri, keyVaultParse.NestedItemTypeKey, *props.KeyName, *props.KeyVersion)
		if err != nil {
			return fmt.Errorf("parsing `key_vault_key_id`: %+v", err)
		}
		d.Set("key_vault_key_id", keyVaultKeyId.ID())
		d.Set("identity_client_id", props.IdentityClientId)
	}

	return nil
}

func resourceCognitiveAccountCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cognitiveservicesaccounts.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.AccountName, "azurerm_cognitive_account")
	defer locks.UnlockByName(id.AccountName, "azurerm_cognitive_account")

	resp, err := client.AccountsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Encryption == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	keySource := cognitiveservicesaccounts.KeySourceMicrosoftPointCognitiveServices
	// set key source to Microsoft.CognitiveServices to disable customer managed key
	props := cognitiveservicesaccounts.Account{
		Properties: &cognitiveservicesaccounts.AccountProperties{
			Encryption: &cognitiveservicesaccounts.Encryption{
				KeySource: &keySource,
			},
		},
	}

	if _, err = client.AccountsUpdate(ctx, *id, props); err != nil {
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
