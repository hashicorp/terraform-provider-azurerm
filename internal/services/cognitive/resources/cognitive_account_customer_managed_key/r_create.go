// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_customer_managed_key

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCognitiveAccountCustomerManagedKeyCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cognitiveservicesaccounts.ParseAccountID(d.Get("cognitive_account_id").(string))
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			resp, err := client.AccountsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.Encryption != nil && resp.Model.Properties.Encryption.KeySource != nil && *resp.Model.Properties.Encryption.KeySource != cognitiveservicesaccounts.KeySourceMicrosoftPointCognitiveServices {
				return tf.ImportAsExistsError("azurerm_cognitive_account_customer_managed_key", id.ID())
			}
		}
	}

	props := cognitiveservicesaccounts.Account{
		Properties: &cognitiveservicesaccounts.AccountProperties{
			Encryption: &cognitiveservicesaccounts.Encryption{
				KeySource: pointer.To(cognitiveservicesaccounts.KeySourceMicrosoftPointKeyVault),
			},
		},
	}

	keyId, err := keyvault.ParseNestedItemID(d.Get("key_vault_key_id").(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
	if err != nil {
		return err
	}
	props.Properties.Encryption.KeyVaultProperties = &cognitiveservicesaccounts.KeyVaultProperties{
		KeyName:     pointer.To(keyId.Name),
		KeyVersion:  pointer.To(keyId.Version),
		KeyVaultUri: pointer.To(keyId.KeyVaultBaseURL),
	}

	if identityClientId := d.Get("identity_client_id").(string); identityClientId != "" {
		props.Properties.Encryption.KeyVaultProperties.IdentityClientId = pointer.To(identityClientId)
	}

	if err := client.AccountsUpdateCallbackThenPoll(ctx, *id, props, sdk.SetIDCallback(meta, id, d)); err != nil {
		return fmt.Errorf("adding Customer Managed Key for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCognitiveAccountCustomerManagedKeyRead(d, meta)
}
