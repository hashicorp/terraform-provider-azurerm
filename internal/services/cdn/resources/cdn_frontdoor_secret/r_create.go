// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_secret

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/secrets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorSecretCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorSecretsClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profile, err := profiles.ParseProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	id := secrets.NewSecretID(profile.SubscriptionId, profile.ResourceGroupName, profile.ProfileName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_secret", id.ID())
		}
	}

	secretParams, err := expandCdnFrontDoorSecretParameters(ctx, d.Get("secret").([]any), meta.(*clients.Client))
	if err != nil {
		return fmt.Errorf("expanding `secret`: %+v", err)
	}

	props := secrets.Secret{
		Properties: &secrets.SecretProperties{
			Parameters: secretParams,
		},
	}

	if err := client.CreateCallbackThenPoll(ctx, id, props, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnFrontDoorSecretRead(d, meta)
}

func expandCdnFrontDoorSecretParameters(ctx context.Context, input []any, clients *clients.Client) (secrets.SecretParameters, error) {
	v := input[0].(map[string]any)

	cc := v["customer_certificate"].([]any)
	cc0 := cc[0].(map[string]any)

	certificateId, err := keyvault.ParseNestedItemID(cc0["key_vault_certificate_id"].(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeCertificate)
	if err != nil {
		return nil, err
	}

	keyVaultBaseId, err := clients.KeyVault.KeyVaultIDFromBaseUrl(ctx, commonids.NewSubscriptionID(clients.Account.SubscriptionId), certificateId.KeyVaultBaseURL)
	if err != nil {
		return nil, fmt.Errorf("retrieving the Key Vault Resource ID from the Key Vault Base URL (`%s`): %w", certificateId.KeyVaultBaseURL, err)
	}

	if keyVaultBaseId == nil {
		return nil, fmt.Errorf("retrieving the Key Vault Resource ID from the Key Vault Base URL (`%s`): id was nil", certificateId.KeyVaultBaseURL)
	}

	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultBaseId)
	if err != nil {
		return nil, err
	}

	useLatest := certificateId.Version == ""
	customerCertificate := &secrets.CustomerCertificateParameters{
		Type: secrets.SecretTypeCustomerCertificate,
		SecretSource: secrets.ResourceReference{
			Id: pointer.To(parse.NewSecretVersionlessID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, certificateId.Name).ID()),
		},
		UseLatestVersion: pointer.To(useLatest),
	}

	if !useLatest {
		customerCertificate.SecretVersion = pointer.To(certificateId.Version)
	}

	return customerCertificate, nil
}
