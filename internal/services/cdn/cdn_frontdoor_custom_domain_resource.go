// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-04-15/afdcustomdomains"
	dnsValidate "github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorCustomDomain() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceCdnFrontDoorCustomDomainCreate,
		Read:   resourceCdnFrontDoorCustomDomainRead,
		Update: resourceCdnFrontDoorCustomDomainUpdate,
		Delete: resourceCdnFrontDoorCustomDomainDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			// NOTE: These timeouts are extremely long due to the manual
			// step of approving the private link if defined.
			Create: pluginsdk.DefaultTimeout(12 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(24 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(12 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontDoorCustomDomainID(id)
			return err
		}),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(validateCipherSuiteConfiguration),
		),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorCustomDomainName,
			},

			// NOTE: I need this fake field because I need the profile name during the create operation
			"cdn_frontdoor_profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorProfileID,
			},

			"dns_zone_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: dnsValidate.ValidateDnsZoneID,
			},

			"host_name": {
				Type:     pluginsdk.TypeString,
				ForceNew: true,
				Required: true,
			},

			"tls": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"certificate_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(cdn.AfdCertificateTypeManagedCertificate),
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.AfdCertificateTypeCustomerCertificate),
								string(cdn.AfdCertificateTypeManagedCertificate),
							}, false),
						},

						"minimum_tls_version": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(cdn.AfdMinimumTLSVersionTLS12),
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.AfdMinimumTLSVersionTLS12),
							}, false),
						},

						// NOTE: If the secret is managed by FrontDoor this will cause a perpetual diff,
						// so this has to be an optional computed field.
						"cdn_frontdoor_secret_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.FrontDoorSecretID,
						},

						"cipher_suite_set_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(afdcustomdomains.AfdCipherSuiteSetTypeCustomized),
								string(afdcustomdomains.AfdCipherSuiteSetTypeTLSOneTwoTwoZeroTwoThree),
								string(afdcustomdomains.AfdCipherSuiteSetTypeTLSOneTwoTwoZeroTwoTwo),
								// Explicitly exclude TLS10_2019 - TLS 1.0/1.1 support retired March 1, 2025
							}, false),
						},

						"customized_cipher_suite": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"tls12_cipher_suites": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice(
												afdcustomdomains.PossibleValuesForAfdCustomizedCipherSuiteForTls12(),
												false,
											),
										},
									},

									"tls13_cipher_suites": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice(
												afdcustomdomains.PossibleValuesForAfdCustomizedCipherSuiteForTls13(),
												false,
											),
										},
									},
								},
							},
						},
					},
				},
			},

			"expiration_date": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"validation_token": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}

	if !features.FivePointOh() {
		resource.Schema["tls"].Elem.(*pluginsdk.Resource).Schema["minimum_tls_version"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeString,
			Optional:   true,
			Default:    string(cdn.AfdMinimumTLSVersionTLS12),
			Deprecated: "As of March 1, 2025, support for 'TLS10' will be retired from Azure Front Door, therefore the 'TLS10' property value will be removed in v5.0 of the provider.",
			ValidateFunc: validation.StringInSlice([]string{
				string(cdn.AfdMinimumTLSVersionTLS12),
				string(cdn.AfdMinimumTLSVersionTLS10),
			}, false),
		}
	}

	return resource
}

func resourceCdnFrontDoorCustomDomainCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDCustomDomainsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := profiles.ParseProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	id := afdcustomdomains.NewCustomDomainID(subscriptionId, profileId.ResourceGroupName, profileId.ProfileName, d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_cdn_frontdoor_custom_domain", id.ID())
	}

	dnsZone := d.Get("dns_zone_id").(string)
	tls := d.Get("tls").([]interface{})

	props := afdcustomdomains.AFDDomain{
		Properties: &afdcustomdomains.AFDDomainProperties{
			HostName: d.Get("host_name").(string),
		},
	}

	if dnsZone != "" {
		props.Properties.AzureDnsZone = expandAfdResourceReference(dnsZone)
	}

	tlsSettings, err := expandAfdDomainTlsParameters(tls)
	if err != nil {
		return err
	}

	props.Properties.TlsSettings = tlsSettings

	if err := client.CreateThenPoll(ctx, id, props); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCdnFrontDoorCustomDomainRead(d, meta)
}

func resourceCdnFrontDoorCustomDomainRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDCustomDomainsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdcustomdomains.ParseCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.CustomDomainName)
	d.Set("cdn_frontdoor_profile_id", profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("host_name", props.HostName)

			dnsZoneId := flattenAfdDNSZoneResourceReference(props.AzureDnsZone)
			if err := d.Set("dns_zone_id", dnsZoneId); err != nil {
				return fmt.Errorf("setting `dns_zone_id`: %+v", err)
			}

			tls := flattenAfdDomainHttpsParameters(props.TlsSettings)
			if err := d.Set("tls", tls); err != nil {
				return fmt.Errorf("setting `tls`: %+v", err)
			}

			if validationProps := props.ValidationProperties; validationProps != nil {
				d.Set("expiration_date", validationProps.ExpirationDate)
				d.Set("validation_token", validationProps.ValidationToken)
			}
		}
	}

	return nil
}

func resourceCdnFrontDoorCustomDomainUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDCustomDomainsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdcustomdomains.ParseCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	props := afdcustomdomains.AFDDomainUpdateParameters{
		Properties: &afdcustomdomains.AFDDomainUpdatePropertiesParameters{},
	}

	if d.HasChange("dns_zone_id") {
		if dnsZone := d.Get("dns_zone_id").(string); dnsZone != "" {
			props.Properties.AzureDnsZone = expandAfdResourceReference(dnsZone)
		}
	}

	if d.HasChange("tls") {
		tlsSettings := d.Get("tls").([]interface{})
		tls, err := expandAfdDomainTlsParameters(tlsSettings)
		if err != nil {
			return err
		}
		props.Properties.TlsSettings = tls
	}

	if err := client.UpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceCdnFrontDoorCustomDomainRead(d, meta)
}

func resourceCdnFrontDoorCustomDomainDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDCustomDomainsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdcustomdomains.ParseCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func validateCipherSuiteConfiguration(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
	tlsRaw := diff.Get("tls").([]interface{})
	if len(tlsRaw) == 0 || tlsRaw[0] == nil {
		return nil
	}

	tls := tlsRaw[0].(map[string]interface{})
	cipherSuiteSetType := tls["cipher_suite_set_type"].(string)
	customizedCipherSuiteRaw := tls["customized_cipher_suite"].([]interface{})

	if cipherSuiteSetType == string(afdcustomdomains.AfdCipherSuiteSetTypeCustomized) {
		if len(customizedCipherSuiteRaw) == 0 {
			return fmt.Errorf("`customized_cipher_suite` is required when `cipher_suite_set_type` is `Customized`")
		}

		if customizedCipherSuiteRaw[0] == nil {
			return fmt.Errorf("at least one cipher suite must be selected in `customized_cipher_suite` when `cipher_suite_set_type` is set to `Customized`")
		}

		customizedCipherSuite := customizedCipherSuiteRaw[0].(map[string]interface{})
		tls12Suites := customizedCipherSuite["tls12_cipher_suites"].(*pluginsdk.Set)
		tls13Suites := customizedCipherSuite["tls13_cipher_suites"].(*pluginsdk.Set)

		if tls12Suites.Len() == 0 && tls13Suites.Len() == 0 {
			return fmt.Errorf("at least one cipher suite must be selected in `customized_cipher_suite` when `cipher_suite_set_type` is set to `Customized`")
		}

		if tls12Suites.Len() > 0 {
			supportedTls12 := []string{
				"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
				"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
				"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
				"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384",
				"TLS_DHE_RSA_WITH_AES_128_GCM_SHA256",
				"TLS_DHE_RSA_WITH_AES_256_GCM_SHA384",
			}
			for _, suite := range tls12Suites.List() {
				suiteStr := suite.(string)
				if !contains(supportedTls12, suiteStr) {
					return fmt.Errorf("unsupported TLS 1.2 cipher suite `%s` - Azure Front Door only supports RSA-based cipher suites: %v", suiteStr, supportedTls12)
				}
			}
		}

		if tls13Suites.Len() > 0 {
			supportedTls13 := []string{
				"TLS_AES_128_GCM_SHA256",
				"TLS_AES_256_GCM_SHA384",
			}
			for _, suite := range tls13Suites.List() {
				suiteStr := suite.(string)
				if !contains(supportedTls13, suiteStr) {
					return fmt.Errorf("unsupported TLS 1.3 cipher suite `%s` - supported cipher suites are: %v", suiteStr, supportedTls13)
				}
			}
		}
	} else if cipherSuiteSetType != "" && len(customizedCipherSuiteRaw) > 0 && customizedCipherSuiteRaw[0] != nil {
		return fmt.Errorf("`customized_cipher_suite` cannot be specified when `cipher_suite_set_type` is not `Customized`")
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func expandAfdDomainTlsParameters(input []interface{}) (*afdcustomdomains.AFDDomainHTTPSParameters, error) {
	// NOTE: With the Frontdoor service, they do not treat an empty object like an empty object
	// if it is not nil they assume it is fully defined and then end up throwing errors when they
	// attempt to get a value from one of the fields.
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	v := input[0].(map[string]interface{})

	certType := v["certificate_type"].(string)
	secretRaw := v["cdn_frontdoor_secret_id"].(string)
	minTlsVersion := v["minimum_tls_version"].(string)
	cipherSuiteSetType := v["cipher_suite_set_type"].(string)
	customizedCipherSuiteRaw := v["customized_cipher_suite"].([]interface{})

	// NOTE: Cert type has to always be passed in the else you will get a
	// "AfdDomain.TlsSettings.CertificateType' is required but it was not set" error
	tls := afdcustomdomains.AFDDomainHTTPSParameters{
		CertificateType: afdcustomdomains.AfdCertificateType(certType),
	}

	if certType == string(afdcustomdomains.AfdCertificateTypeCustomerCertificate) && secretRaw == "" {
		return nil, fmt.Errorf("the `cdn_frontdoor_secret_id` field must be set if the `certificate_type` is `CustomerCertificate`")
	} else if certType == string(afdcustomdomains.AfdCertificateTypeManagedCertificate) && secretRaw != "" {
		return nil, fmt.Errorf("the `cdn_frontdoor_secret_id` field is not supported if the `certificate_type` is `ManagedCertificate`")
	}

	// NOTE: Secret always needs to be passed if it is defined else you will
	// receive a 500 Internal Server Error
	if secretRaw != "" {
		secret, err := parse.FrontDoorSecretID(secretRaw)
		if err != nil {
			return nil, err
		}

		tls.Secret = expandAfdResourceReference(secret.ID())
	}

	if minTlsVersion != "" {
		tlsVer := afdcustomdomains.AfdMinimumTlsVersion(minTlsVersion)
		tls.MinimumTlsVersion = &tlsVer
	}

	if cipherSuiteSetType != "" {
		cipherType := afdcustomdomains.AfdCipherSuiteSetType(cipherSuiteSetType)
		tls.CipherSuiteSetType = &cipherType
	}

	if len(customizedCipherSuiteRaw) > 0 {
		tls.CustomizedCipherSuiteSet = expandAfdCustomizedCipherSuiteSet(customizedCipherSuiteRaw)
	}

	return &tls, nil
}

func expandAfdCustomizedCipherSuiteSet(input []interface{}) *afdcustomdomains.AFDDomainHTTPSCustomizedCipherSuiteSet {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	result := &afdcustomdomains.AFDDomainHTTPSCustomizedCipherSuiteSet{}

	if tls12Raw := v["tls12_cipher_suites"].(*pluginsdk.Set); tls12Raw.Len() > 0 {
		tls12Suites := make([]afdcustomdomains.AfdCustomizedCipherSuiteForTls12, 0)
		for _, suite := range tls12Raw.List() {
			tls12Suites = append(tls12Suites, afdcustomdomains.AfdCustomizedCipherSuiteForTls12(suite.(string)))
		}
		result.CipherSuiteSetForTls12 = &tls12Suites
	}

	if tls13Raw := v["tls13_cipher_suites"].(*pluginsdk.Set); tls13Raw.Len() > 0 {
		tls13Suites := make([]afdcustomdomains.AfdCustomizedCipherSuiteForTls13, 0)
		for _, suite := range tls13Raw.List() {
			tls13Suites = append(tls13Suites, afdcustomdomains.AfdCustomizedCipherSuiteForTls13(suite.(string)))
		}
		result.CipherSuiteSetForTls13 = &tls13Suites
	}

	return result
}

func expandAfdResourceReference(id string) *afdcustomdomains.ResourceReference {
	return &afdcustomdomains.ResourceReference{
		Id: pointer.To(id),
	}
}

func flattenAfdDomainHttpsParameters(input *afdcustomdomains.AFDDomainHTTPSParameters) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	secretId := ""
	if input.Secret != nil && input.Secret.Id != nil {
		secretId = *input.Secret.Id
	}

	minTlsVersion := ""
	if input.MinimumTlsVersion != nil {
		minTlsVersion = string(*input.MinimumTlsVersion)
	}

	cipherSuiteSetType := ""
	if input.CipherSuiteSetType != nil {
		cipherSuiteSetType = string(*input.CipherSuiteSetType)
	}

	return []interface{}{
		map[string]interface{}{
			"cdn_frontdoor_secret_id": secretId,
			"certificate_type":        string(input.CertificateType),
			"minimum_tls_version":     minTlsVersion,
			"cipher_suite_set_type":   cipherSuiteSetType,
			"customized_cipher_suite": flattenAfdCustomizedCipherSuiteSet(input.CustomizedCipherSuiteSet),
		},
	}
}

func flattenAfdCustomizedCipherSuiteSet(input *afdcustomdomains.AFDDomainHTTPSCustomizedCipherSuiteSet) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	tls12Suites := make([]interface{}, 0)
	if input.CipherSuiteSetForTls12 != nil {
		for _, suite := range *input.CipherSuiteSetForTls12 {
			tls12Suites = append(tls12Suites, string(suite))
		}
	}

	tls13Suites := make([]interface{}, 0)
	if input.CipherSuiteSetForTls13 != nil {
		for _, suite := range *input.CipherSuiteSetForTls13 {
			tls13Suites = append(tls13Suites, string(suite))
		}
	}

	if len(tls12Suites) == 0 && len(tls13Suites) == 0 {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"tls12_cipher_suites": tls12Suites,
			"tls13_cipher_suites": tls13Suites,
		},
	}
}

func flattenAfdDNSZoneResourceReference(input *afdcustomdomains.ResourceReference) string {
	if input == nil || input.Id == nil {
		return ""
	}
	return *input.Id
}
