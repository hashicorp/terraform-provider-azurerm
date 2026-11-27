// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	keyvaultClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	webClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceCertificateCreate(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := keyvaultClient.From(meta)
	client := webClient.From(meta).CertificatesClient

	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	id := certificates.NewCertificateID(account.From(meta).SubscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if !response.WasNotFound(existing.HttpResponse) {
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %w", id, err)
			}
			return tf.ImportAsExistsError("azurerm_app_service_certificate", id.ID())
		}
	}

	certificate := certificates.Certificate{
		Properties: &certificates.CertificateProperties{
			Password: pointer.To(d.Get("password").(string)),
		},
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v := d.Get("app_service_plan_id").(string); v != "" {
		certificate.Properties.ServerFarmId = &v
	}

	if v := d.Get("pfx_blob").(string); v != "" {
		certificate.Properties.PfxBlob = pointer.To(v)
	}

	if v := d.Get("key_vault_secret_id").(string); v != "" {
		parsedSecretId, err := keyvault.ParseNestedItemID(v, keyvault.VersionTypeAny, keyvault.NestedItemTypeAny)
		if err != nil {
			return err
		}

		var keyVaultId *string
		if v := d.Get("key_vault_id").(string); v != "" {
			keyVaultId = pointer.To(v)
		} else {
			keyVaultBaseUrl := parsedSecretId.KeyVaultBaseURL

			subscriptionResourceId := commonids.NewSubscriptionID(id.SubscriptionId)
			keyVaultId, err = keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, keyVaultBaseUrl)
			if err != nil {
				return fmt.Errorf("retrieving the Resource ID for the Key Vault at URL %q: %s", keyVaultBaseUrl, err)
			}
			if keyVaultId == nil {
				return fmt.Errorf("unable to determine the Resource ID for the Key Vault at URL %q", keyVaultBaseUrl)
			}
		}

		certificate.Properties.KeyVaultId = keyVaultId
		certificate.Properties.KeyVaultSecretName = pointer.To(parsedSecretId.Name)
	}

	if _, err := client.CreateOrUpdate(ctx, id, certificate); err != nil {
		return fmt.Errorf("creating %s: %s", id, err)
	}

	d.SetId(id.ID())

	return resourceAppServiceCertificateRead(d, meta)
}
