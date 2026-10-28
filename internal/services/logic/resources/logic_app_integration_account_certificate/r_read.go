// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_integration_account_certificate

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountcertificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogicAppIntegrationAccountCertificateRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountCertificateClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationaccountcertificates.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.CertificateName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("integration_account_name", id.IntegrationAccountName)

	if model := resp.Model; model != nil {
		props := model.Properties
		if err := d.Set("key_vault_key", flattenIntegrationAccountCertificateKeyVaultKey(props.Key)); err != nil {
			return fmt.Errorf("setting `key_vault_key`: %+v", err)
		}

		if props.Metadata != nil {
			d.Set("metadata", props.Metadata)
		}

		d.Set("public_certificate", props.PublicCertificate)
	}

	return nil
}

func flattenIntegrationAccountCertificateKeyVaultKey(input *integrationaccountcertificates.KeyVaultKeyReference) []any {
	if input == nil {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"key_name":     input.KeyName,
			"key_vault_id": pointer.From(input.KeyVault.Id),
			"key_version":  pointer.From(input.KeyVersion),
		},
	}
}
