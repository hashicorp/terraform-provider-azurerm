// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_customer_managed_key

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCognitiveAccountCustomerManagedKeyRead(d *pluginsdk.ResourceData, meta any) error {
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
		keyVaultKeyId, err := keyvault.NewNestedItemID(pointer.From(props.KeyVaultUri), keyvault.NestedItemTypeKey, pointer.From(props.KeyName), pointer.From(props.KeyVersion))
		if err != nil {
			return fmt.Errorf("parsing `key_vault_key_id`: %+v", err)
		}
		d.Set("key_vault_key_id", keyVaultKeyId.ID())
		d.Set("identity_client_id", props.IdentityClientId)
	}

	return nil
}
