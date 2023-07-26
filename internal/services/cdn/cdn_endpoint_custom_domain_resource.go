// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	keyvaultClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	keyvaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyvaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmCdnEndpointCustomDomain() *pluginsdk.Resource {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CdnEndpointCustomDomainName(),
		},

		"cdn_endpoint_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.EndpointID,
		},

		"host_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"cdn_managed_https": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"certificate_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(cdn.CertificateTypeShared),
							string(cdn.CertificateTypeDedicated),
						}, false),
					},
					"protocol_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(cdn.ProtocolTypeServerNameIndication),
							string(cdn.ProtocolTypeIPBased),
						}, false),
					},
					"tls_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(cdn.MinimumTLSVersionTLS10),
							string(cdn.MinimumTLSVersionTLS12),
							string(cdn.MinimumTLSVersionNone),
						}, false),
						Default: string(cdn.MinimumTLSVersionTLS12),
					},
				},
			},
			ConflictsWith: []string{"user_managed_https"},
		},

		"user_managed_https": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"tls_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(cdn.MinimumTLSVersionTLS10),
							string(cdn.MinimumTLSVersionTLS12),
							string(cdn.MinimumTLSVersionNone),
						}, false),
						Default: string(cdn.MinimumTLSVersionTLS12),
					},
				},
			},
			ConflictsWith: []string{"cdn_managed_https"},
		},
	}
	if !features.FourPointOhBeta() {
		schema["user_managed_https"].Elem.(*pluginsdk.Resource).Schema["key_vault_certificate_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: keyvaultValidate.NestedItemIdWithOptionalVersion,
			ExactlyOneOf: []string{"user_managed_https.0.key_vault_certificate_id", "user_managed_https.0.key_vault_secret_id"},
			Deprecated:   "This is deprecated in favor of `key_vault_secret_id` as the service is actually looking for a secret, not a certificate",
		}

		schema["user_managed_https"].Elem.(*pluginsdk.Resource).Schema["key_vault_secret_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: keyvaultValidate.NestedItemIdWithOptionalVersion,
			ExactlyOneOf: []string{"user_managed_https.0.key_vault_certificate_id", "user_managed_https.0.key_vault_secret_id"},
		}
	} else {
		schema["user_managed_https"].Elem.(*pluginsdk.Resource).Schema["key_vault_secret_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: keyvaultValidate.NestedItemIdWithOptionalVersion,
		}
	}

	return &pluginsdk.Resource{
		Create: resourceArmCdnEndpointCustomDomainCreate,
		Read:   resourceArmCdnEndpointCustomDomainRead,
		Update: resourceArmCdnEndpointCustomDomainUpdate,
		Delete: resourceArmCdnEndpointCustomDomainDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CustomDomainID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(12 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(24 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(12 * time.Hour),
		},

		Schema: schema,
	}
}

func resourceArmCdnEndpointCustomDomainCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.CustomDomainsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	epid := d.Get("cdn_endpoint_id").(string)

	cdnEndpointId, err := parse.EndpointID(epid)
	if err != nil {
		return err
	}

	id := parse.NewCustomDomainID(cdnEndpointId.SubscriptionId, cdnEndpointId.ResourceGroup, cdnEndpointId.ProfileName, cdnEndpointId.Name, name)

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %q: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cdn_endpoint_custom_domain", id.ID())
	}

	props := cdn.CustomDomainParameters{
		CustomDomainPropertiesParameters: &cdn.CustomDomainPropertiesParameters{
			HostName: utils.String(d.Get("host_name").(string)),
		},
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name, props)
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
		// https://docs.microsoft.com/en-us/azure/cdn/cdn-custom-ssl?tabs=option-2-enable-https-with-your-own-certificate#tlsssl-certificates
		pfClient := meta.(*clients.Client).Cdn.ProfilesClient
		cdnEndpointResp, err := pfClient.Get(ctx, id.ResourceGroup, id.ProfileName)
		if err != nil {
			return fmt.Errorf("retrieving Cdn Profile %q (Resource Group %q): %+v",
				id.ResourceGroup, id.ProfileName, err)
		}
		supportedSku := map[cdn.SkuName]bool{
			cdn.SkuNamePremiumVerizon:    true,
			cdn.SkuNameStandardVerizon:   true,
			cdn.SkuNameStandardMicrosoft: true,
		}
		if cdnEndpointResp.Sku != nil && !supportedSku[cdnEndpointResp.Sku.Name] {
			return fmt.Errorf("user managed HTTPS certificate is only available for Azure CDN from Microsoft or Azure CDN from Verizon profiles")
		}
		params, err = expandArmCdnEndpointCustomDomainUserManagedHttpsSettings(ctx, v.([]interface{}), meta.(*clients.Client))
		if err != nil {
			return err
		}
	} else if v, ok := d.GetOk("cdn_managed_https"); ok {
		params = expandArmCdnEndpointCustomDomainCdnManagedHttpsSettings(v.([]interface{}))
	}

	if params != nil {
		if err := enableArmCdnEndpointCustomDomainHttps(ctx, client, id, params); err != nil {
			return err
		}
	}

	return resourceArmCdnEndpointCustomDomainRead(d, meta)
}

func resourceArmCdnEndpointCustomDomainUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.CustomDomainsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CustomDomainID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	const (
		turnOn = iota
		turnOff
		update
		noChange
	)

	var (
		cdnManagedHTTPSStatus = noChange
		cdnManagedHTTPSParams cdn.BasicCustomDomainHTTPSParameters

		userManagedHTTPSStatus = noChange
		userManagedHTTPSParams cdn.BasicCustomDomainHTTPSParameters
	)

	if d.HasChange("cdn_managed_https") {
		props := resp.CustomDomainProperties
		if props == nil {
			return fmt.Errorf("unexpected nil of `CustomDomainProperties` in response")
		}

		cdnManagedHTTPSParams = expandArmCdnEndpointCustomDomainCdnManagedHttpsSettings(d.Get("cdn_managed_https").([]interface{}))

		if props.CustomHTTPSParameters == nil {
			cdnManagedHTTPSStatus = turnOn
		} else {
			if cdnManagedHTTPSParams == nil {
				cdnManagedHTTPSStatus = turnOff
			} else {
				cdnManagedHTTPSStatus = update
			}
		}
	}

	if d.HasChange("user_managed_https") {
		props := resp.CustomDomainProperties
		if props == nil {
			return fmt.Errorf("unexpected nil of `CustomDomainProperties` in response")
		}

		var err error
		userManagedHTTPSParams, err = expandArmCdnEndpointCustomDomainUserManagedHttpsSettings(ctx, d.Get("user_managed_https").([]interface{}), meta.(*clients.Client))
		if err != nil {
			return err
		}

		if props.CustomHTTPSParameters == nil {
			userManagedHTTPSStatus = turnOn
		} else {
			if userManagedHTTPSParams == nil {
				userManagedHTTPSStatus = turnOff
			} else {
				userManagedHTTPSStatus = update
			}
		}
	}

	// There are theoretically 16 (4x4) combinations of the cdn/user managed https status combinations.
	// While actually there are only following 8 combinations due to the exclusive nature of both settings.
	// +-----------------------------------+
	// |     	| n/a | on | off | update  |
	// |--------|--------------------------|
	// | n/a 	|     |  x |  x  |    x    |
	// | on  	|  x  |    |  x  |         |
	// | off    |  x  |  x |     |         |
	// | update |  x  |    |     |         |
	// +-----------------------------------+

	switch {
	case cdnManagedHTTPSStatus == turnOff || cdnManagedHTTPSStatus == update:
		if err := disableArmCdnEndpointCustomDomainHttps(ctx, client, *id); err != nil {
			return fmt.Errorf("disable CDN Managed HTTPS on %q: %+v", *id, err)
		}
	case userManagedHTTPSStatus == turnOff || userManagedHTTPSStatus == update:
		if err := disableArmCdnEndpointCustomDomainHttps(ctx, client, *id); err != nil {
			return fmt.Errorf("disable User Managed HTTPS on %q: %+v", *id, err)
		}
	}

	switch {
	case cdnManagedHTTPSStatus == turnOn || cdnManagedHTTPSStatus == update:
		if err := enableArmCdnEndpointCustomDomainHttps(ctx, client, *id, cdnManagedHTTPSParams); err != nil {
			return fmt.Errorf("enable CDN Managed HTTPS on %q: %+v", *id, err)
		}
	case userManagedHTTPSStatus == turnOn || userManagedHTTPSStatus == update:
		if err := enableArmCdnEndpointCustomDomainHttps(ctx, client, *id, userManagedHTTPSParams); err != nil {
			return fmt.Errorf("enable User Managed HTTPS on %q: %+v", *id, err)
		}
	}

	return resourceArmCdnEndpointCustomDomainRead(d, meta)
}

func resourceArmCdnEndpointCustomDomainRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.CustomDomainsClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CustomDomainID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	cdnEndpointId := parse.NewEndpointID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.EndpointName)

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
			if b := d.Get("user_managed_https").([]interface{}); len(b) == 1 {
				b := b[0].(map[string]interface{})

				secretIdRaw := b["key_vault_secret_id"].(string)
				if !features.FourPointOhBeta() {
					if secretIdRaw == "" {
						secretIdRaw = b["key_vault_certificate_id"].(string)
					}
				}

				if secretIdRaw != "" {
					id, err := keyvaultParse.ParseOptionallyVersionedNestedItemID(secretIdRaw)
					if err != nil {
						return fmt.Errorf("parsing Key Vault Secret Id %q: %v", secretIdRaw, err)
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

func resourceArmCdnEndpointCustomDomainDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.CustomDomainsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CustomDomainID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %q: %+v", id, err)
	}

	return nil
}

func expandArmCdnEndpointCustomDomainCdnManagedHttpsSettings(input []interface{}) cdn.BasicCustomDomainHTTPSParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	output := &cdn.ManagedHTTPSParameters{
		CertificateSourceParameters: &cdn.CertificateSourceParameters{
			OdataType:       utils.String("#Microsoft.Azure.Cdn.Models.CdnCertificateSourceParameters"),
			CertificateType: cdn.CertificateType(raw["certificate_type"].(string)),
		},
		CertificateSource: cdn.CertificateSourceCdn,
		ProtocolType:      cdn.ProtocolType(raw["protocol_type"].(string)),
		MinimumTLSVersion: cdn.MinimumTLSVersion(raw["tls_version"].(string)),
	}

	return output
}

func expandArmCdnEndpointCustomDomainUserManagedHttpsSettings(ctx context.Context, input []interface{}, clients *clients.Client) (cdn.BasicCustomDomainHTTPSParameters, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	raw := input[0].(map[string]interface{})

	idLiteral := raw["key_vault_secret_id"].(string)
	if !features.FourPointOhBeta() {
		if idLiteral == "" {
			idLiteral = raw["key_vault_certificate_id"].(string)
		}
	}

	keyVaultSecretId, err := keyvaultParse.ParseOptionallyVersionedNestedItemID(idLiteral)
	if err != nil {
		return nil, err
	}

	keyVaultIdRaw, err := clients.KeyVault.KeyVaultIDFromBaseUrl(ctx, clients.Resource, keyVaultSecretId.KeyVaultBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", keyVaultSecretId.KeyVaultBaseUrl, err)
	}
	if keyVaultIdRaw == nil {
		return nil, fmt.Errorf("unexpected nil Key Vault ID retrieved at URL %q", keyVaultSecretId.KeyVaultBaseUrl)
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
			OdataType:         utils.String("#Microsoft.Azure.Cdn.Models.KeyVaultCertificateSourceParameters"),
			SubscriptionID:    pointer.To(keyVaultId.SubscriptionId),
			ResourceGroupName: pointer.To(keyVaultId.ResourceGroupName),
			VaultName:         pointer.To(keyVaultId.VaultName),
			SecretName:        pointer.To(keyVaultSecretId.Name),
			SecretVersion:     SecretVersion,
			UpdateRule:        utils.String("NoAction"),
			DeleteRule:        utils.String("NoAction"),
		},
		CertificateSource: cdn.CertificateSourceAzureKeyVault,
		ProtocolType:      cdn.ProtocolTypeServerNameIndication,
		MinimumTLSVersion: cdn.MinimumTLSVersion(raw["tls_version"].(string)),
	}

	return output, nil
}

func flattenArmCdnEndpointCustomDomainCdnManagedHttpsSettings(input cdn.ManagedHTTPSParameters) []interface{} {
	certificateType := ""
	if params := input.CertificateSourceParameters; params != nil {
		certificateType = string(params.CertificateType)
	}

	return []interface{}{
		map[string]interface{}{
			"certificate_type": certificateType,
			"protocol_type":    string(input.ProtocolType),
			"tls_version":      string(input.MinimumTLSVersion),
		},
	}
}

func flattenArmCdnEndpointCustomDomainUserManagedHttpsSettings(ctx context.Context, input cdn.UserManagedHTTPSParameters, keyVaultsClient *keyvaultClient.Client, isVersioned bool) ([]interface{}, error) {
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

	var secretVersion string
	if params.SecretVersion != nil {
		secretVersion = *params.SecretVersion
	}

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

	secretId, err := keyvaultParse.ParseOptionallyVersionedNestedItemID(*secret.ID)
	if err != nil {
		return nil, err
	}

	secretIdLiteral := secretId.ID()
	if !isVersioned {
		secretIdLiteral = secretId.VersionlessID()
	}

	m := map[string]interface{}{
		"key_vault_secret_id": secretIdLiteral,
		"tls_version":         string(input.MinimumTLSVersion),
	}

	if features.FourPointOhBeta() {
		return []interface{}{m}, nil
	}

	// We try to retrieve the certificate with the given secret name and version. If it returns error, then we tolerate the error and simply setting empty string for the certificate id.
	// As in this case, the users might be using a secret rather than a certificate.
	var certIdLiteral string
	cert, err := keyVaultsClient.ManagementClient.GetCertificate(ctx, *keyVaultBaseUrl, secretName, secretVersion)
	if err == nil && cert.ID != nil {
		certId, _ := keyvaultParse.ParseOptionallyVersionedNestedItemID(*cert.ID)
		certIdLiteral = certId.ID()
		if !isVersioned {
			certIdLiteral = certId.VersionlessID()
		}
	}

	m["key_vault_certificate_id"] = certIdLiteral

	return []interface{}{m}, nil
}

func enableArmCdnEndpointCustomDomainHttps(ctx context.Context, client *cdn.CustomDomainsClient, id parse.CustomDomainId, params cdn.BasicCustomDomainHTTPSParameters) error {
	future, err := client.EnableCustomHTTPS(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name, &params)
	if err != nil {
		return fmt.Errorf("sending enable request: %+v", err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for enabling HTTPS: %+v", err)
	}
	return nil
}

func disableArmCdnEndpointCustomDomainHttps(ctx context.Context, client *cdn.CustomDomainsClient, id parse.CustomDomainId) error {
	future, err := client.DisableCustomHTTPS(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name)
	if err != nil {
		return fmt.Errorf("sending disable request: %+v", err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for disabling HTTPS: %+v", err)
	}
	return nil
}
