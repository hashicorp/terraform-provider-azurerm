// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/delegationsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/signinsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/signupsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/tenantaccess"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenApiManagementHostnameConfigurations(input *[]apimanagementservice.HostnameConfiguration, d *pluginsdk.ResourceData) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	managementResults := make([]any, 0)
	portalResults := make([]any, 0)
	developerPortalResults := make([]any, 0)
	proxyResults := make([]any, 0)
	scmResults := make([]any, 0)

	for _, config := range *input {
		output := make(map[string]any)

		output["host_name"] = config.HostName
		output["negotiate_client_certificate"] = pointer.From(config.NegotiateClientCertificate)
		output["key_vault_certificate_id"] = pointer.From(config.KeyVaultId)
		output["ssl_keyvault_identity_client_id"] = pointer.From(config.IdentityClientId)

		if config.Certificate != nil {
			if config.Certificate.Expiry != "" {
				output["expiry"] = config.Certificate.Expiry
			}
			output["thumbprint"] = config.Certificate.Thumbprint
			output["subject"] = config.Certificate.Subject
		}

		output["certificate_source"] = pointer.From(config.CertificateSource)
		output["certificate_status"] = pointer.From(config.CertificateStatus)

		var configType string
		switch strings.ToLower(string(config.Type)) {
		case strings.ToLower(string(apimanagementservice.HostnameTypeProxy)):
			// only set SSL binding for proxy types
			output["default_ssl_binding"] = pointer.From(config.DefaultSslBinding)
			proxyResults = append(proxyResults, output)
			configType = "proxy"

		case strings.ToLower(string(apimanagementservice.HostnameTypeManagement)):
			managementResults = append(managementResults, output)
			configType = "management"

		case strings.ToLower(string(apimanagementservice.HostnameTypePortal)):
			portalResults = append(portalResults, output)
			configType = "portal"

		case strings.ToLower(string(apimanagementservice.HostnameTypeDeveloperPortal)):
			developerPortalResults = append(developerPortalResults, output)
			configType = "developer_portal"

		case strings.ToLower(string(apimanagementservice.HostnameTypeScm)):
			scmResults = append(scmResults, output)
			configType = "scm"
		}

		existingHostnames := d.Get("hostname_configuration").([]any)
		if len(existingHostnames) > 0 && configType != "" {
			v := existingHostnames[0].(map[string]any)

			if valsRaw, ok := v[configType]; ok {
				vals := valsRaw.([]any)
				schemaz.CopyCertificateAndPassword(vals, config.HostName, output)
			}
		}
	}

	if len(managementResults) == 0 && len(portalResults) == 0 && len(developerPortalResults) == 0 && len(proxyResults) == 0 && len(scmResults) == 0 {
		return []any{}
	}

	return []any{
		map[string]any{
			"management":       managementResults,
			"portal":           portalResults,
			"developer_portal": developerPortalResults,
			"proxy":            proxyResults,
			"scm":              scmResults,
		},
	}
}

func flattenApiManagementAdditionalLocations(input *[]apimanagementservice.AdditionalLocation) ([]any, error) {
	results := make([]any, 0)
	if input == nil {
		return results, nil
	}

	for _, prop := range *input {
		virtualNetworkConfiguration, err := flattenApiManagementVirtualNetworkConfiguration(prop.VirtualNetworkConfiguration)
		if err != nil {
			return results, err
		}

		results = append(results, map[string]any{
			"capacity":                      int32(prop.Sku.Capacity),
			"gateway_regional_url":          pointer.From(prop.GatewayRegionalURL),
			"location":                      location.NormalizeNilable(pointer.To(prop.Location)),
			"private_ip_addresses":          pointer.From(prop.PrivateIPAddresses),
			"public_ip_address_id":          pointer.From(prop.PublicIPAddressId),
			"public_ip_addresses":           pointer.From(prop.PublicIPAddresses),
			"virtual_network_configuration": virtualNetworkConfiguration,
			"zones":                         zones.FlattenUntyped(prop.Zones),
			"gateway_disabled":              pointer.From(prop.DisableGateway),
		})
	}

	return results, nil
}

func flattenApiManagementServiceSkuName(input *apimanagementservice.ApiManagementServiceSkuProperties) string {
	if input == nil {
		return ""
	}

	return fmt.Sprintf("%s_%d", string(input.Name), input.Capacity)
}

func flattenApiManagementSecurityCustomProperties(input map[string]string, skuIsConsumption bool) []any {
	output := make(map[string]any)

	output["backend_ssl30_enabled"] = parseApiManagementNilableDictionary(input, apimBackendProtocolSsl3)
	output["backend_tls10_enabled"] = parseApiManagementNilableDictionary(input, apimBackendProtocolTls10)
	output["backend_tls11_enabled"] = parseApiManagementNilableDictionary(input, apimBackendProtocolTls11)
	output["frontend_tls10_enabled"] = parseApiManagementNilableDictionary(input, apimFrontendProtocolTls10)
	output["frontend_tls11_enabled"] = parseApiManagementNilableDictionary(input, apimFrontendProtocolTls11)

	if !skuIsConsumption {
		output["frontend_ssl30_enabled"] = parseApiManagementNilableDictionary(input, apimFrontendProtocolSsl3)

		output["triple_des_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTripleDesCiphers)
		output["tls_ecdhe_ecdsa_with_aes256_cbc_sha_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsEcdheEcdsaWithAes256CbcShaCiphers)
		output["tls_ecdhe_ecdsa_with_aes128_cbc_sha_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsEcdheEcdsaWithAes128CbcShaCiphers)
		output["tls_ecdhe_rsa_with_aes256_cbc_sha_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsEcdheRsaWithAes256CbcShaCiphers)
		output["tls_ecdhe_rsa_with_aes128_cbc_sha_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsEcdheRsaWithAes128CbcShaCiphers)
		output["tls_rsa_with_aes256_gcm_sha384_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsRsaWithAes256GcmSha384Ciphers)
		output["tls_rsa_with_aes128_gcm_sha256_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsRsaWithAes128GcmSha256Ciphers)
		output["tls_rsa_with_aes256_cbc_sha256_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsRsaWithAes256CbcSha256Ciphers)
		output["tls_rsa_with_aes128_cbc_sha256_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsRsaWithAes128CbcSha256Ciphers)
		output["tls_rsa_with_aes256_cbc_sha_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsRsaWithAes256CbcShaCiphers)
		output["tls_rsa_with_aes128_cbc_sha_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsRsaWithAes128CbcShaCiphers)
	}

	return []any{output}
}

func flattenApiManagementProtocolsCustomProperties(input map[string]string) []any {
	output := make(map[string]any)

	output["http2_enabled"] = parseApiManagementNilableDictionary(input, apimHttp2Protocol)

	return []any{output}
}

func flattenApiManagementVirtualNetworkConfiguration(input *apimanagementservice.VirtualNetworkConfiguration) ([]any, error) {
	if input == nil {
		return []any{}, nil
	}

	virtualNetworkConfiguration := make(map[string]any)

	if input.SubnetResourceId != nil {
		subnetId, err := commonids.ParseSubnetIDInsensitively(*input.SubnetResourceId)
		if err != nil {
			return []any{}, err
		}
		virtualNetworkConfiguration["subnet_id"] = subnetId.ID()
	}

	return []any{virtualNetworkConfiguration}, nil
}

func flattenApiManagementSignInSettings(input signinsettings.PortalSigninSettings) []any {
	enabled := false

	if props := input.Properties; props != nil {
		if props.Enabled != nil {
			enabled = pointer.From(props.Enabled)
		}
	}

	return []any{
		map[string]any{
			"enabled": enabled,
		},
	}
}

func flattenApiManagementDelegationSettings(input delegationsettings.PortalDelegationSettings, keyContract delegationsettings.PortalSettingValidationKeyContract) []any {
	url := ""
	subscriptionsEnabled := false
	userRegistrationEnabled := false

	if props := input.Properties; props != nil {
		url = pointer.From(props.Url)
		if props.Subscriptions != nil {
			subscriptionsEnabled = pointer.From(props.Subscriptions.Enabled)
		}
		if props.UserRegistration != nil {
			userRegistrationEnabled = pointer.From(props.UserRegistration.Enabled)
		}
	}

	return []any{
		map[string]any{
			"url":                       url,
			"subscriptions_enabled":     subscriptionsEnabled,
			"user_registration_enabled": userRegistrationEnabled,
			"validation_key":            pointer.From(keyContract.ValidationKey),
		},
	}
}

func flattenApiManagementSignUpSettings(input signupsettings.PortalSignupSettings) []any {
	enabled := false
	termsOfService := make([]any, 0)

	if props := input.Properties; props != nil {
		enabled = pointer.From(props.Enabled)

		if tos := props.TermsOfService; tos != nil {
			output := make(map[string]any)

			output["enabled"] = pointer.From(tos.Enabled)
			output["consent_required"] = pointer.From(tos.ConsentRequired)
			output["text"] = pointer.From(tos.Text)

			termsOfService = append(termsOfService, output)
		}
	}

	return []any{
		map[string]any{
			"enabled":          enabled,
			"terms_of_service": termsOfService,
		},
	}
}

func flattenApiManagementTenantAccessSettings(input tenantaccess.AccessInformationSecretsContract) []any {
	result := make(map[string]any)

	result["enabled"] = pointer.From(input.Enabled)
	result["tenant_id"] = pointer.From(input.Id)
	result["primary_key"] = pointer.From(input.PrimaryKey)
	result["secondary_key"] = pointer.From(input.SecondaryKey)

	return []any{result}
}

func flattenAPIManagementCertificates(d *pluginsdk.ResourceData, inputs *[]apimanagementservice.CertificateConfiguration) []any {
	if inputs == nil || len(*inputs) == 0 {
		return []any{}
	}

	outputs := []any{}
	for i, input := range *inputs {
		var pwd, encodedCertificate string
		if v, ok := d.GetOk(fmt.Sprintf("certificate.%d.certificate_password", i)); ok {
			pwd = v.(string)
		}

		if v, ok := d.GetOk(fmt.Sprintf("certificate.%d.encoded_certificate", i)); ok {
			encodedCertificate = v.(string)
		}

		output := map[string]any{
			"certificate_password": pwd,
			"encoded_certificate":  encodedCertificate,
			"store_name":           string(input.StoreName),
			"expiry":               input.Certificate.Expiry,
			"subject":              input.Certificate.Subject,
			"thumbprint":           input.Certificate.Thumbprint,
		}
		outputs = append(outputs, output)
	}
	return outputs
}
