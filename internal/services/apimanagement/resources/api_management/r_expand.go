// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management

import (
	"errors"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/delegationsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/signinsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/signupsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/tenantaccess"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandAzureRmApiManagementHostnameConfigurations(d *pluginsdk.ResourceData) *[]apimanagementservice.HostnameConfiguration {
	results := make([]apimanagementservice.HostnameConfiguration, 0)
	vs := d.Get("hostname_configuration")
	if vs == nil {
		return &results
	}
	hostnameVs := vs.([]any)

	for _, hostnameRawVal := range hostnameVs {
		// hostnameRawVal is guaranteed to be non-nil as there is AtLeastOneOf constraint on its containing properties.
		hostnameV := hostnameRawVal.(map[string]any)

		managementVs := hostnameV["management"].([]any)
		for _, managementV := range managementVs {
			v := managementV.(map[string]any)
			output := helpers.ExpandApiManagementCommonHostnameConfiguration(v, apimanagementservice.HostnameTypeManagement)
			results = append(results, output)
		}

		portalVs := hostnameV["portal"].([]any)
		for _, portalV := range portalVs {
			v := portalV.(map[string]any)
			output := helpers.ExpandApiManagementCommonHostnameConfiguration(v, apimanagementservice.HostnameTypePortal)
			results = append(results, output)
		}

		developerPortalVs := hostnameV["developer_portal"].([]any)
		for _, developerPortalV := range developerPortalVs {
			v := developerPortalV.(map[string]any)
			output := helpers.ExpandApiManagementCommonHostnameConfiguration(v, apimanagementservice.HostnameTypeDeveloperPortal)
			results = append(results, output)
		}

		proxyVs := hostnameV["proxy"].([]any)
		for _, proxyV := range proxyVs {
			v := proxyV.(map[string]any)
			output := helpers.ExpandApiManagementCommonHostnameConfiguration(v, apimanagementservice.HostnameTypeProxy)

			if value, ok := v["default_ssl_binding"]; ok {
				output.DefaultSslBinding = pointer.To(value.(bool))
			}

			results = append(results, output)
		}

		scmVs := hostnameV["scm"].([]any)
		for _, scmV := range scmVs {
			v := scmV.(map[string]any)
			output := helpers.ExpandApiManagementCommonHostnameConfiguration(v, apimanagementservice.HostnameTypeScm)
			results = append(results, output)
		}
	}

	return &results
}

func expandAzureRmApiManagementCertificates(d *pluginsdk.ResourceData) *[]apimanagementservice.CertificateConfiguration {
	vs := d.Get("certificate").([]any)

	results := make([]apimanagementservice.CertificateConfiguration, 0)

	for _, v := range vs {
		config := v.(map[string]any)

		certBase64 := config["encoded_certificate"].(string)
		storeName := apimanagementservice.StoreName(config["store_name"].(string))

		cert := apimanagementservice.CertificateConfiguration{
			EncodedCertificate: pointer.To(certBase64),
			StoreName:          storeName,
		}

		cert.CertificatePassword = pointer.To(config["certificate_password"].(string))

		results = append(results, cert)
	}

	return &results
}

func expandAzureRmApiManagementAdditionalLocations(d *pluginsdk.ResourceData, sku apimanagementservice.ApiManagementServiceSkuProperties) (*[]apimanagementservice.AdditionalLocation, error) {
	inputLocations := d.Get("additional_location").([]any)
	parentVnetConfig := d.Get("virtual_network_configuration").([]any)

	additionalLocations := make([]apimanagementservice.AdditionalLocation, 0)

	for _, v := range inputLocations {
		config := v.(map[string]any)
		location := location.Normalize(config["location"].(string))

		if config["capacity"].(int) > 0 {
			sku.Capacity = int64(config["capacity"].(int))
		}

		additionalLocation := apimanagementservice.AdditionalLocation{
			Location:       location,
			Sku:            sku,
			DisableGateway: pointer.To(config["gateway_disabled"].(bool)),
		}

		childVnetConfig := config["virtual_network_configuration"].([]any)
		switch {
		case len(childVnetConfig) == 0 && len(parentVnetConfig) > 0:
			return nil, errors.New("`virtual_network_configuration` must be specified in any `additional_location` block when top-level `virtual_network_configuration` is supplied")
		case len(childVnetConfig) > 0 && len(parentVnetConfig) == 0:
			return nil, errors.New("`virtual_network_configuration` must be empty in all `additional_location` blocks when top-level `virtual_network_configuration` is not supplied")
		case len(childVnetConfig) > 0 && len(parentVnetConfig) > 0:
			v := childVnetConfig[0].(map[string]any)
			subnetResourceId := v["subnet_id"].(string)
			additionalLocation.VirtualNetworkConfiguration = &apimanagementservice.VirtualNetworkConfiguration{
				SubnetResourceId: pointer.To(subnetResourceId),
			}
		}

		publicIPAddressID := config["public_ip_address_id"].(string)
		if publicIPAddressID != "" {
			if sku.Name != apimanagementservice.SkuTypePremium {
				if len(childVnetConfig) == 0 {
					return nil, errors.New("`public_ip_address_id` for an additional location is only supported when sku type is `Premium`, and the APIM instance is deployed in a virtual network")
				}
			}
			additionalLocation.PublicIPAddressId = &publicIPAddressID
		}

		zones := zones.ExpandUntyped(config["zones"].(*schema.Set).List())
		if len(zones) > 0 {
			additionalLocation.Zones = &zones
		}

		additionalLocations = append(additionalLocations, additionalLocation)
	}

	return &additionalLocations, nil
}

func expandAzureRmApiManagementSkuName(input string) apimanagementservice.ApiManagementServiceSkuProperties {
	// "sku_name" is validated to be in this format above, and is required
	skuParts := strings.Split(input, "_")
	name := skuParts[0]
	capacity, _ := strconv.Atoi(skuParts[1])
	return apimanagementservice.ApiManagementServiceSkuProperties{
		Name:     apimanagementservice.SkuType(name),
		Capacity: int64(capacity),
	}
}

func expandApiManagementCustomProperties(d *pluginsdk.ResourceData, skuIsConsumption bool) (map[string]string, error) {
	backendProtocolSsl3 := false
	backendProtocolTls10 := false
	backendProtocolTls11 := false
	frontendProtocolSsl3 := false
	frontendProtocolTls10 := false
	frontendProtocolTls11 := false
	tripleDesCiphers := false
	tlsEcdheEcdsaWithAes256CbcShaCiphers := false
	tlsEcdheEcdsaWithAes128CbcShaCiphers := false
	tlsEcdheRsaWithAes256CbcShaCiphers := false
	tlsEcdheRsaWithAes128CbcShaCiphers := false
	tlsRsaWithAes128GcmSha256Ciphers := false
	tlsRsaWithAes256GcmSha384Ciphers := false
	tlsRsaWithAes256CbcSha256Ciphers := false
	tlsRsaWithAes128CbcSha256Ciphers := false
	tlsRsaWithAes256CbcShaCiphers := false
	tlsRsaWithAes128CbcShaCiphers := false

	if vs := d.Get("security").([]any); len(vs) > 0 {
		v := vs[0].(map[string]any)
		backendProtocolSsl3 = v["backend_ssl30_enabled"].(bool)
		backendProtocolTls10 = v["backend_tls10_enabled"].(bool)
		backendProtocolTls11 = v["backend_tls11_enabled"].(bool)
		frontendProtocolSsl3 = v["frontend_ssl30_enabled"].(bool)
		frontendProtocolTls10 = v["frontend_tls10_enabled"].(bool)
		frontendProtocolTls11 = v["frontend_tls11_enabled"].(bool)

		if v, exists := v["triple_des_ciphers_enabled"]; exists {
			tripleDesCiphers = v.(bool)
		}

		tlsEcdheEcdsaWithAes256CbcShaCiphers = v["tls_ecdhe_ecdsa_with_aes256_cbc_sha_ciphers_enabled"].(bool)
		tlsEcdheEcdsaWithAes128CbcShaCiphers = v["tls_ecdhe_ecdsa_with_aes128_cbc_sha_ciphers_enabled"].(bool)
		tlsEcdheRsaWithAes256CbcShaCiphers = v["tls_ecdhe_rsa_with_aes256_cbc_sha_ciphers_enabled"].(bool)
		tlsEcdheRsaWithAes128CbcShaCiphers = v["tls_ecdhe_rsa_with_aes128_cbc_sha_ciphers_enabled"].(bool)
		tlsRsaWithAes128GcmSha256Ciphers = v["tls_rsa_with_aes128_gcm_sha256_ciphers_enabled"].(bool)
		tlsRsaWithAes256GcmSha384Ciphers = v["tls_rsa_with_aes256_gcm_sha384_ciphers_enabled"].(bool)
		tlsRsaWithAes256CbcSha256Ciphers = v["tls_rsa_with_aes256_cbc_sha256_ciphers_enabled"].(bool)
		tlsRsaWithAes128CbcSha256Ciphers = v["tls_rsa_with_aes128_cbc_sha256_ciphers_enabled"].(bool)
		tlsRsaWithAes256CbcShaCiphers = v["tls_rsa_with_aes256_cbc_sha_ciphers_enabled"].(bool)
		tlsRsaWithAes128CbcShaCiphers = v["tls_rsa_with_aes128_cbc_sha_ciphers_enabled"].(bool)

		if skuIsConsumption && frontendProtocolSsl3 {
			return nil, errors.New("`frontend_ssl30_enabled` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tripleDesCiphers {
			return nil, errors.New("`enable_triple_des_ciphers` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tlsEcdheEcdsaWithAes256CbcShaCiphers {
			return nil, errors.New("`tls_ecdhe_ecdsa_with_aes256_cbc_sha_ciphers_enabled` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tlsEcdheEcdsaWithAes128CbcShaCiphers {
			return nil, errors.New("`tls_ecdhe_ecdsa_with_aes128_cbc_sha_ciphers_enabled` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tlsEcdheRsaWithAes256CbcShaCiphers {
			return nil, errors.New("`tls_ecdhe_rsa_with_aes256_cbc_sha_ciphers_enabled` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tlsEcdheRsaWithAes128CbcShaCiphers {
			return nil, errors.New("`tls_ecdhe_rsa_with_aes128_cbc_sha_ciphers_enabled` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tlsRsaWithAes128GcmSha256Ciphers {
			return nil, errors.New("`tls_rsa_with_aes128_gcm_sha256_ciphers_enabled` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tlsRsaWithAes256CbcSha256Ciphers {
			return nil, errors.New("`tls_rsa_with_aes256_cbc_sha256_ciphers_enabled` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tlsRsaWithAes128CbcSha256Ciphers {
			return nil, errors.New("`tls_rsa_with_aes128_cbc_sha256_ciphers_enabled` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tlsRsaWithAes256CbcShaCiphers {
			return nil, errors.New("`tls_rsa_with_aes256_cbc_sha_ciphers_enabled` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tlsRsaWithAes128CbcShaCiphers {
			return nil, errors.New("`tls_rsa_with_aes128_cbc_sha_ciphers_enabled` is not supported for Sku Tier `Consumption`")
		}
	}

	customProperties := map[string]string{
		apimBackendProtocolSsl3:   strconv.FormatBool(backendProtocolSsl3),
		apimBackendProtocolTls10:  strconv.FormatBool(backendProtocolTls10),
		apimBackendProtocolTls11:  strconv.FormatBool(backendProtocolTls11),
		apimFrontendProtocolTls10: strconv.FormatBool(frontendProtocolTls10),
		apimFrontendProtocolTls11: strconv.FormatBool(frontendProtocolTls11),
	}

	if !skuIsConsumption {
		customProperties[apimFrontendProtocolSsl3] = strconv.FormatBool(frontendProtocolSsl3)
		customProperties[apimTripleDesCiphers] = strconv.FormatBool(tripleDesCiphers)
		customProperties[apimTlsEcdheEcdsaWithAes256CbcShaCiphers] = strconv.FormatBool(tlsEcdheEcdsaWithAes256CbcShaCiphers)
		customProperties[apimTlsEcdheEcdsaWithAes128CbcShaCiphers] = strconv.FormatBool(tlsEcdheEcdsaWithAes128CbcShaCiphers)
		customProperties[apimTlsEcdheRsaWithAes256CbcShaCiphers] = strconv.FormatBool(tlsEcdheRsaWithAes256CbcShaCiphers)
		customProperties[apimTlsEcdheRsaWithAes128CbcShaCiphers] = strconv.FormatBool(tlsEcdheRsaWithAes128CbcShaCiphers)
		customProperties[apimTlsRsaWithAes128GcmSha256Ciphers] = strconv.FormatBool(tlsRsaWithAes128GcmSha256Ciphers)
		customProperties[apimTlsRsaWithAes256GcmSha384Ciphers] = strconv.FormatBool(tlsRsaWithAes256GcmSha384Ciphers)
		customProperties[apimTlsRsaWithAes256CbcSha256Ciphers] = strconv.FormatBool(tlsRsaWithAes256CbcSha256Ciphers)
		customProperties[apimTlsRsaWithAes128CbcSha256Ciphers] = strconv.FormatBool(tlsRsaWithAes128CbcSha256Ciphers)
		customProperties[apimTlsRsaWithAes256CbcShaCiphers] = strconv.FormatBool(tlsRsaWithAes256CbcShaCiphers)
		customProperties[apimTlsRsaWithAes128CbcShaCiphers] = strconv.FormatBool(tlsRsaWithAes128CbcShaCiphers)
	}

	if vp := d.Get("protocols").([]any); len(vp) > 0 {
		vpr := vp[0].(map[string]any)
		enableHttp2 := vpr["http2_enabled"].(bool)
		customProperties[apimHttp2Protocol] = strconv.FormatBool(enableHttp2)
	}

	return customProperties, nil
}

func expandAzureRmApiManagementVirtualNetworkConfigurations(d *pluginsdk.ResourceData) *apimanagementservice.VirtualNetworkConfiguration {
	vs := d.Get("virtual_network_configuration").([]any)
	if len(vs) == 0 {
		return nil
	}

	v := vs[0].(map[string]any)

	return &apimanagementservice.VirtualNetworkConfiguration{
		SubnetResourceId: pointer.To(v["subnet_id"].(string)),
	}
}

func expandApiManagementSignInSettings(input []any) signinsettings.PortalSigninSettings {
	enabled := false

	if len(input) > 0 {
		vs := input[0].(map[string]any)
		enabled = vs["enabled"].(bool)
	}

	return signinsettings.PortalSigninSettings{
		Properties: &signinsettings.PortalSigninSettingProperties{
			Enabled: pointer.To(enabled),
		},
	}
}

func expandApiManagementDelegationSettings(input []any) delegationsettings.PortalDelegationSettings {
	if len(input) == 0 {
		return delegationsettings.PortalDelegationSettings{}
	}

	vs := input[0].(map[string]any)

	props := delegationsettings.PortalDelegationSettingsProperties{
		UserRegistration: &delegationsettings.RegistrationDelegationSettingsProperties{
			Enabled: pointer.To(vs["user_registration_enabled"].(bool)),
		},
		Subscriptions: &delegationsettings.SubscriptionsDelegationSettingsProperties{
			Enabled: pointer.To(vs["subscriptions_enabled"].(bool)),
		},
	}

	validationKey := vs["validation_key"].(string)
	if !vs["user_registration_enabled"].(bool) && !vs["subscriptions_enabled"].(bool) && validationKey == "" {
		// for some reason we cannot leave this empty
		props.ValidationKey = pointer.To("cGxhY2Vob2xkZXIxCg==")
	}
	if validationKey != "" {
		props.ValidationKey = pointer.To(validationKey)
	}

	url := vs["url"].(string)
	if !vs["user_registration_enabled"].(bool) && !vs["subscriptions_enabled"].(bool) && url == "" {
		// for some reason we cannot leave this empty
		props.Url = pointer.To("https://www.placeholder.com")
	}
	if url != "" {
		props.Url = pointer.To(url)
	}

	return delegationsettings.PortalDelegationSettings{
		Properties: &props,
	}
}

func expandApiManagementSignUpSettings(input []any) signupsettings.PortalSignupSettings {
	if len(input) == 0 {
		return signupsettings.PortalSignupSettings{
			Properties: &signupsettings.PortalSignupSettingsProperties{
				Enabled: pointer.To(false),
				TermsOfService: &signupsettings.TermsOfServiceProperties{
					ConsentRequired: pointer.To(false),
					Enabled:         pointer.To(false),
					Text:            pointer.To(""),
				},
			},
		}
	}

	vs := input[0].(map[string]any)

	props := signupsettings.PortalSignupSettingsProperties{
		Enabled: pointer.To(vs["enabled"].(bool)),
	}

	termsOfServiceRaw := vs["terms_of_service"].([]any)
	if len(termsOfServiceRaw) > 0 {
		termsOfServiceVs := termsOfServiceRaw[0].(map[string]any)
		props.TermsOfService = &signupsettings.TermsOfServiceProperties{
			Enabled:         pointer.To(termsOfServiceVs["enabled"].(bool)),
			ConsentRequired: pointer.To(termsOfServiceVs["consent_required"].(bool)),
			Text:            pointer.To(termsOfServiceVs["text"].(string)),
		}
	}

	return signupsettings.PortalSignupSettings{
		Properties: &props,
	}
}

func expandApiManagementTenantAccessSettings(input []any) tenantaccess.AccessInformationUpdateParameters {
	enabled := false

	if len(input) > 0 {
		vs := input[0].(map[string]any)
		enabled = vs["enabled"].(bool)
	}

	return tenantaccess.AccessInformationUpdateParameters{
		Properties: &tenantaccess.AccessInformationUpdateParameterProperties{
			Enabled: pointer.To(enabled),
		},
	}
}
