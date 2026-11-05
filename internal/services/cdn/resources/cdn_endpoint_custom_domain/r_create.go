// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_endpoint_custom_domain

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/customdomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmCdnEndpointCustomDomainCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.CustomDomainsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// todo 6.0 - move to the case-sensitive parser when validation.AsGeneratedID is removed: this parses a config
	// value which the paired AsGeneratedID validator accepts with legacy casing, and configs cannot be migrated.
	cdnEndpointId, err := endpoints.ParseEndpointIDInsensitively(d.Get("cdn_endpoint_id").(string))
	if err != nil {
		return err
	}

	id := customdomains.NewEndpointCustomDomainID(cdnEndpointId.SubscriptionId, cdnEndpointId.ResourceGroupName, cdnEndpointId.ProfileName, cdnEndpointId.EndpointName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id.ResourceGroupName, id.ProfileName, id.EndpointName, id.CustomDomainName)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for existing %q: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_cdn_endpoint_custom_domain", id.ID())
		}
	}

	props := cdn.CustomDomainParameters{
		CustomDomainPropertiesParameters: &cdn.CustomDomainPropertiesParameters{
			HostName: pointer.To(d.Get("host_name").(string)),
		},
	}

	future, err := client.Create(ctx, id.ResourceGroupName, id.ProfileName, id.EndpointName, id.CustomDomainName, props)
	if err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %q: %+v", id, err)
	}

	d.SetId(id.ID())

	// Enable https if specified
	var params cdn.BasicCustomDomainHTTPSParameters
	if v, ok := d.GetOk("user_managed_https"); ok {
		// User managed certificate is only available for Azure CDN from Microsoft and Azure CDN from Verizon profiles.
		// https://docs.microsoft.com/azure/cdn/cdn-custom-ssl?tabs=option-2-enable-https-with-your-own-certificate#tlsssl-certificates
		pfClient := meta.(*clients.Client).Cdn.ProfilesClient
		cdnEndpointResp, err := pfClient.Get(ctx, id.ResourceGroupName, id.ProfileName)
		if err != nil {
			return fmt.Errorf("retrieving Cdn Profile %q (Resource Group %q): %+v",
				id.ResourceGroupName, id.ProfileName, err)
		}
		supportedSku := map[cdn.SkuName]bool{
			cdn.SkuNamePremiumVerizon:    true,
			cdn.SkuNameStandardVerizon:   true,
			cdn.SkuNameStandardMicrosoft: true,
		}
		if cdnEndpointResp.Sku != nil && !supportedSku[cdnEndpointResp.Sku.Name] {
			return fmt.Errorf("user managed HTTPS certificate is only available for Azure CDN from Microsoft or Azure CDN from Verizon profiles")
		}
		params, err = expandArmCdnEndpointCustomDomainUserManagedHttpsSettings(ctx, v.([]any), meta.(*clients.Client))
		if err != nil {
			return err
		}
	} else if v, ok := d.GetOk("cdn_managed_https"); ok {
		params = expandArmCdnEndpointCustomDomainCdnManagedHttpsSettings(v.([]any))
	}

	if params != nil {
		if err := enableArmCdnEndpointCustomDomainHttps(ctx, client, id, params); err != nil {
			return err
		}
	}

	return resourceArmCdnEndpointCustomDomainRead(d, meta)
}

func expandArmCdnEndpointCustomDomainCdnManagedHttpsSettings(input []any) cdn.BasicCustomDomainHTTPSParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)
	return &cdn.ManagedHTTPSParameters{
		CertificateSourceParameters: &cdn.CertificateSourceParameters{
			OdataType:       pointer.To("#Microsoft.Azure.Cdn.Models.CdnCertificateSourceParameters"),
			CertificateType: cdn.CertificateType(raw["certificate_type"].(string)),
		},
		CertificateSource: cdn.CertificateSourceCdn,
		ProtocolType:      cdn.ProtocolType(raw["protocol_type"].(string)),
		MinimumTLSVersion: cdn.MinimumTLSVersion(raw["tls_version"].(string)),
	}
}

func expandArmCdnEndpointCustomDomainUserManagedHttpsSettings(ctx context.Context, input []any, clients *clients.Client) (cdn.BasicCustomDomainHTTPSParameters, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	raw := input[0].(map[string]any)

	idLiteral := raw["key_vault_secret_id"].(string)

	keyVaultSecretId, err := keyvault.ParseNestedItemID(idLiteral, keyvault.VersionTypeAny, keyvault.NestedItemTypeAny)
	if err != nil {
		return nil, err
	}

	subscriptionId := commonids.NewSubscriptionID(clients.Account.SubscriptionId)
	keyVaultIdRaw, err := clients.KeyVault.KeyVaultIDFromBaseUrl(ctx, subscriptionId, keyVaultSecretId.KeyVaultBaseURL)
	if err != nil {
		return nil, fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", keyVaultSecretId.KeyVaultBaseURL, err)
	}
	if keyVaultIdRaw == nil {
		return nil, fmt.Errorf("unable to find the Resource Manager ID for the Key Vault URI %q in %s", keyVaultSecretId.KeyVaultBaseURL, subscriptionId)
	}
	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return nil, err
	}

	// Fix for issue #20772
	var SecretVersion *string
	if keyVaultSecretId.Version != "" {
		SecretVersion = pointer.To(keyVaultSecretId.Version)
	}

	output := &cdn.UserManagedHTTPSParameters{
		CertificateSourceParameters: &cdn.KeyVaultCertificateSourceParameters{
			OdataType:         pointer.To("#Microsoft.Azure.Cdn.Models.KeyVaultCertificateSourceParameters"),
			SubscriptionID:    pointer.To(keyVaultId.SubscriptionId),
			ResourceGroupName: pointer.To(keyVaultId.ResourceGroupName),
			VaultName:         pointer.To(keyVaultId.VaultName),
			SecretName:        pointer.To(keyVaultSecretId.Name),
			SecretVersion:     SecretVersion,
			UpdateRule:        pointer.To("NoAction"),
			DeleteRule:        pointer.To("NoAction"),
		},
		CertificateSource: cdn.CertificateSourceAzureKeyVault,
		ProtocolType:      cdn.ProtocolTypeServerNameIndication,
		MinimumTLSVersion: cdn.MinimumTLSVersion(raw["tls_version"].(string)),
	}

	return output, nil
}
