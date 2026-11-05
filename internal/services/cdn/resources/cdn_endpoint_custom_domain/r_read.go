// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_endpoint_custom_domain

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/customdomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyvaultClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmCdnEndpointCustomDomainRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.CustomDomainsClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := customdomains.ParseEndpointCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.ProfileName, id.EndpointName, id.CustomDomainName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[DEBUG] %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	cdnEndpointId := endpoints.NewEndpointID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.EndpointName)

	d.Set("name", resp.Name)
	d.Set("cdn_endpoint_id", cdnEndpointId.ID())
	if props := resp.CustomDomainProperties; props != nil {
		d.Set("host_name", props.HostName)

		switch params := props.CustomHTTPSParameters.(type) {
		case cdn.ManagedHTTPSParameters:
			if err := d.Set("cdn_managed_https", flattenArmCdnEndpointCustomDomainCdnManagedHttpsSettings(params)); err != nil {
				return fmt.Errorf("setting `cdn_managed_https`: %+v", err)
			}
		case cdn.UserManagedHTTPSParameters:
			var isVersioned bool
			if b := d.Get("user_managed_https").([]any); len(b) == 1 {
				b := b[0].(map[string]any)

				secretIdRaw := b["key_vault_secret_id"].(string)
				if secretIdRaw != "" {
					id, err := keyvault.ParseNestedItemID(secretIdRaw, keyvault.VersionTypeAny, keyvault.NestedItemTypeAny)
					if err != nil {
						return err
					}
					isVersioned = id.Version != ""
				}
			}
			settings, err := flattenArmCdnEndpointCustomDomainUserManagedHttpsSettings(ctx, params, keyVaultsClient, isVersioned)
			if err != nil {
				return err
			}
			if err := d.Set("user_managed_https", settings); err != nil {
				return fmt.Errorf("setting `user_managed_https`: %+v", err)
			}
		default:
			d.Set("cdn_managed_https", nil)
			d.Set("user_managed_https", nil)
		}
	}

	return nil
}

func flattenArmCdnEndpointCustomDomainCdnManagedHttpsSettings(input cdn.ManagedHTTPSParameters) []any {
	certificateType := ""
	if params := input.CertificateSourceParameters; params != nil {
		certificateType = string(params.CertificateType)
	}

	return []any{
		map[string]any{
			"certificate_type": certificateType,
			"protocol_type":    string(input.ProtocolType),
			"tls_version":      string(input.MinimumTLSVersion),
		},
	}
}

func flattenArmCdnEndpointCustomDomainUserManagedHttpsSettings(ctx context.Context, input cdn.UserManagedHTTPSParameters, keyVaultsClient *keyvaultClient.Client, isVersioned bool) ([]any, error) {
	params := input.CertificateSourceParameters
	if params == nil {
		return nil, fmt.Errorf("unexpected nil Certificate Source Parameters from API")
	}

	if params.SubscriptionID == nil {
		return nil, fmt.Errorf("unexpected nil `subscriptionId` in the Certificate Source Parameters from API")
	}
	subscriptionId := *params.SubscriptionID

	if params.ResourceGroupName == nil {
		return nil, fmt.Errorf("unexpected nil `resourceGroupName` in the Certificate Source Parameters from API")
	}
	resourceGroupName := *params.ResourceGroupName

	if params.VaultName == nil {
		return nil, fmt.Errorf("unexpected nil `vaultName` in the Certificate Source Parameters from API")
	}
	vaultName := *params.VaultName

	if params.SecretName == nil {
		return nil, fmt.Errorf("unexpected nil `secretName` in the Certificate Source Parameters from API")
	}
	secretName := *params.SecretName

	secretVersion := pointer.From(params.SecretVersion)

	keyVaultId := commonids.NewKeyVaultID(subscriptionId, resourceGroupName, vaultName)
	keyVaultBaseUrl, err := keyVaultsClient.BaseUriForKeyVault(ctx, keyVaultId)
	if err != nil {
		return nil, fmt.Errorf("looking up Key Vault Secret %q vault url from id %q: %+v", vaultName, keyVaultId, err)
	}

	secret, err := keyVaultsClient.ManagementClient.GetSecret(ctx, *keyVaultBaseUrl, secretName, secretVersion)
	if err != nil {
		return nil, err
	}

	if secret.ID == nil {
		return nil, fmt.Errorf("unexpected null Key Vault Secret retrieved for Key Vault %s / Secret Name %s / Secret Version %s", keyVaultId, secretName, secretVersion)
	}

	secretId, err := keyvault.ParseNestedItemID(*secret.ID, keyvault.VersionTypeAny, keyvault.NestedItemTypeAny)
	if err != nil {
		return nil, err
	}

	secretIdLiteral := secretId.ID()
	if !isVersioned {
		secretIdLiteral = secretId.VersionlessID()
	}

	return []any{map[string]any{
		"key_vault_secret_id": secretIdLiteral,
		"tls_version":         string(input.MinimumTLSVersion),
	}}, nil
}

func enableArmCdnEndpointCustomDomainHttps(ctx context.Context, client *cdn.CustomDomainsClient, id customdomains.EndpointCustomDomainId, params cdn.BasicCustomDomainHTTPSParameters) error {
	future, err := client.EnableCustomHTTPS(ctx, id.ResourceGroupName, id.ProfileName, id.EndpointName, id.CustomDomainName, &params)
	if err != nil {
		return fmt.Errorf("sending enable request: %+v", err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for enabling HTTPS: %+v", err)
	}
	return nil
}

func disableArmCdnEndpointCustomDomainHttps(ctx context.Context, client *cdn.CustomDomainsClient, id customdomains.EndpointCustomDomainId) error {
	future, err := client.DisableCustomHTTPS(ctx, id.ResourceGroupName, id.ProfileName, id.EndpointName, id.CustomDomainName)
	if err != nil {
		return fmt.Errorf("sending disable request: %+v", err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for disabling HTTPS: %+v", err)
	}
	return nil
}
