// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_integration_account_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountcertificates"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogicAppIntegrationAccountCertificateCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Logic.IntegrationAccountCertificateClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := integrationaccountcertificates.NewCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("integration_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_logic_app_integration_account_certificate", id.ID())
			}
		}
	}

	parameters := integrationaccountcertificates.IntegrationAccountCertificate{
		Properties: integrationaccountcertificates.IntegrationAccountCertificateProperties{},
	}

	if v, ok := d.GetOk("key_vault_key"); ok {
		parameters.Properties.Key = expandIntegrationAccountCertificateKeyVaultKey(v.([]any))
	}

	if v, ok := d.GetOk("metadata"); ok {
		parameters.Properties.Metadata = &v
	}

	if v, ok := d.GetOk("public_certificate"); ok {
		parameters.Properties.PublicCertificate = pointer.To(v.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogicAppIntegrationAccountCertificateRead(d, meta)
}

func expandIntegrationAccountCertificateKeyVaultKey(input []any) *integrationaccountcertificates.KeyVaultKeyReference {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	result := integrationaccountcertificates.KeyVaultKeyReference{
		KeyVault: integrationaccountcertificates.KeyVaultKeyReferenceKeyVault{
			Id: pointer.To(v["key_vault_id"].(string)),
		},
		KeyName: v["key_name"].(string),
	}

	if keyVersion := v["key_version"].(string); keyVersion != "" {
		result.KeyVersion = pointer.To(keyVersion)
	}

	return &result
}
