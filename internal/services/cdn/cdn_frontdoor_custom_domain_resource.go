// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-04-15/afdcustomdomains"
	dnsValidate "github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

// Front Door currently rejects the DHE TLS 1.2 suites exposed by the SDK helper.
// The service team confirmed this supported set officially, so this is not a
// temporary workaround. We keep an explicit allowlist of the accepted ECDHE
// values here.
var frontDoorCustomDomainCustomizedCipherSuitesForTls12 = []string{
	string(afdcustomdomains.AfdCustomizedCipherSuiteForTls12ECDHERSAAESOneTwoEightGCMSHATwoFiveSix),
	string(afdcustomdomains.AfdCustomizedCipherSuiteForTls12ECDHERSAAESTwoFiveSixGCMSHAThreeEightFour),
	string(afdcustomdomains.AfdCustomizedCipherSuiteForTls12ECDHERSAAESOneTwoEightSHATwoFiveSix),
	string(afdcustomdomains.AfdCustomizedCipherSuiteForTls12ECDHERSAAESTwoFiveSixSHAThreeEightFour),
}

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
			pluginsdk.CustomizeDiffShim(frontDoorCustomDomainHostNameCustomizeDiff),
			pluginsdk.CustomizeDiffShim(frontDoorCustomDomainTlsCustomizeDiff),
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
				Type:         pluginsdk.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validate.FrontDoorCustomDomainHostName,
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
							Default:  string(afdcustomdomains.AfdCertificateTypeManagedCertificate),
							ValidateFunc: validation.StringInSlice([]string{
								string(afdcustomdomains.AfdCertificateTypeCustomerCertificate),
								string(afdcustomdomains.AfdCertificateTypeManagedCertificate),
							}, false),
						},

						"minimum_version": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(afdcustomdomains.AfdMinimumTlsVersionTLSOneTwo),
							ValidateFunc: validation.StringInSlice([]string{
								string(afdcustomdomains.AfdMinimumTlsVersionTLSOneTwo),
							}, false),
						},

						"cdn_frontdoor_secret_id": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							// O+C because if the secret is managed by FrontDoor this will cause a perpetual diff
							Computed:     true,
							ValidateFunc: validate.FrontDoorSecretID,
						},

						"cipher_suite": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"type": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(afdcustomdomains.AfdCipherSuiteSetTypeCustomized),
											string(afdcustomdomains.AfdCipherSuiteSetTypeTLSOneTwoTwoZeroTwoThree),
											string(afdcustomdomains.AfdCipherSuiteSetTypeTLSOneTwoTwoZeroTwoTwo),
											// Explicitly exclude TLS10_2019 - TLS 1.0/1.1 support retired March 1, 2025
										}, false),
									},

									"custom_ciphers": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"tls12": {
													Type:     pluginsdk.TypeSet,
													Optional: true,
													AtLeastOneOf: []string{
														"tls.0.cipher_suite.0.custom_ciphers.0.tls12",
														"tls.0.cipher_suite.0.custom_ciphers.0.tls13",
													},
													Elem: &pluginsdk.Schema{
														Type: pluginsdk.TypeString,
														ValidateFunc: validation.StringInSlice(
															frontDoorCustomDomainCustomizedCipherSuitesForTls12,
															false,
														),
													},
												},

												"tls13": {
													Type: pluginsdk.TypeSet,
													AtLeastOneOf: []string{
														"tls.0.cipher_suite.0.custom_ciphers.0.tls12",
														"tls.0.cipher_suite.0.custom_ciphers.0.tls13",
													},
													Optional: true,
													// NOTE: O+C Azure Front Door returns TLS 1.3 cipher suites even when `tls13` is not specified (see https://learn.microsoft.com/azure/frontdoor/standard-premium/tls-policy#custom-tls-policy)
													Computed: true,
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
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"tls.0.minimum_version"},
			// NOTE: O+C so both `minimum_tls_version` and `minimum_version` appear in state during v4.x for backward compatibility
			Computed:   true,
			Deprecated: "`minimum_tls_version` has been deprecated in favour of `minimum_version` and will be removed in v5.0 of the AzureRM provider",
			ValidateFunc: validation.StringInSlice([]string{
				string(afdcustomdomains.AfdMinimumTlsVersionTLSOneTwo),
				string(afdcustomdomains.AfdMinimumTlsVersionTLSOneZero),
			}, false),
		}

		resource.Schema["tls"].Elem.(*pluginsdk.Resource).Schema["minimum_version"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"tls.0.minimum_tls_version"},
			// NOTE: O+C so both `minimum_tls_version` and `minimum_version` appear in state during v4.x for backward compatibility
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(afdcustomdomains.AfdMinimumTlsVersionTLSOneTwo),
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

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_custom_domain", id.ID())
		}
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

	tlsSettings, err := expandAfdDomainTlsParameters(d, tls)
	if err != nil {
		return err
	}

	props.Properties.TlsSettings = tlsSettings

	if err := client.CreateCallbackThenPoll(ctx, id, props, sdk.SetIDCallback(meta, &id, d)); err != nil {
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

			includeDefaultCipherSuite := resourceCdnFrontDoorCustomDomainCipherSuiteConfigured(d)
			tls := flattenAfdDomainHttpsParameters(props.TlsSettings, includeDefaultCipherSuite)
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
		tls, err := expandAfdDomainTlsParameters(d, tlsSettings)
		if err != nil {
			return err
		}
		props.Properties.TlsSettings = tls
	}

	// The AFDx service currently rejects custom-domain updates until its internal
	// validation and backend synchronization flow completes. The service team
	// confirmed clients must wait for approval here, otherwise Update can return
	// 409 Conflict while the sync job is still in progress.
	approvalPollerType := custompollers.NewFrontDoorCustomDomainWaitForApprovedPoller(client, *id)
	approvalPoller := pollers.NewPoller(approvalPollerType, 30*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := approvalPoller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be approved: %+v", *id, err)
	}

	resp, err := client.Update(ctx, *id, props)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err := resp.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", *id, err)
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

	result, err := client.Delete(ctx, *id)
	if err != nil {
		if response.WasNotFound(result.HttpResponse) {
			return nil
		}
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	pollerType := custompollers.NewFrontDoorCustomDomainDeletePoller(client, *id)
	poller := pollers.NewPoller(pollerType, 30*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)

	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandAfdDomainTlsParameters(d *pluginsdk.ResourceData, input []interface{}) (*afdcustomdomains.AFDDomainHTTPSParameters, error) {
	// NOTE: With the Frontdoor service, they do not treat an empty object like an empty object
	// if it is not nil they assume it is fully defined and then end up throwing errors when they
	// attempt to get a value from one of the fields.
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	v := input[0].(map[string]interface{})

	certType := v["certificate_type"].(string)
	secretRaw := v["cdn_frontdoor_secret_id"].(string)
	secretWasConfigured := false
	minimumVersionConfigured := false
	minimumTlsVersionConfigured := false

	var minTlsVersion string

	if rawConfig := d.GetRawConfig(); !rawConfig.IsNull() {
		tlsConfig := rawConfig.GetAttr("tls")
		if !tlsConfig.IsNull() && tlsConfig.LengthInt() > 0 {
			tlsBlock := tlsConfig.AsValueSlice()[0]
			if !tlsBlock.IsNull() {
				secretWasConfigured = !tlsBlock.GetAttr("cdn_frontdoor_secret_id").IsNull()
				minimumVersionConfigured = !tlsBlock.GetAttr("minimum_version").IsNull()
				if !features.FivePointOh() {
					minimumTlsVersionConfigured = !tlsBlock.GetAttr("minimum_tls_version").IsNull()
				}
			}
		}
	}

	if !features.FivePointOh() {
		switch {
		case minimumVersionConfigured:
			minTlsVersion = v["minimum_version"].(string)
		case minimumTlsVersionConfigured:
			minTlsVersion = v["minimum_tls_version"].(string)
		default:
			minTlsVersion = string(afdcustomdomains.AfdMinimumTlsVersionTLSOneTwo)
		}
	} else {
		minTlsVersion = v["minimum_version"].(string)
	}

	cipherSuiteRaw := v["cipher_suite"].([]interface{})

	tls := afdcustomdomains.AFDDomainHTTPSParameters{
		CertificateType: afdcustomdomains.AfdCertificateType(certType),
	}

	if certType == string(afdcustomdomains.AfdCertificateTypeCustomerCertificate) && secretRaw == "" {
		return nil, errors.New("the `cdn_frontdoor_secret_id` field must be set if the `certificate_type` is `CustomerCertificate`")
	} else if certType == string(afdcustomdomains.AfdCertificateTypeManagedCertificate) {
		// Ignore computed `cdn_frontdoor_secret_id` for managed certs unless the
		// user explicitly configured it.
		if secretRaw != "" && secretWasConfigured {
			return nil, errors.New("the `cdn_frontdoor_secret_id` field is not supported if the `certificate_type` is `ManagedCertificate`")
		}
		if !secretWasConfigured {
			secretRaw = ""
		}
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

	if len(cipherSuiteRaw) > 0 && cipherSuiteRaw[0] != nil {
		cipherSuite := cipherSuiteRaw[0].(map[string]interface{})

		if cipherSuiteType := cipherSuite["type"].(string); cipherSuiteType != "" {
			cipherType := afdcustomdomains.AfdCipherSuiteSetType(cipherSuiteType)
			tls.CipherSuiteSetType = &cipherType
		}

		if customCiphersRaw := cipherSuite["custom_ciphers"].([]interface{}); len(customCiphersRaw) > 0 {
			tls.CustomizedCipherSuiteSet = expandAfdCustomizedCipherSuiteSet(customCiphersRaw)
		}
	}

	return &tls, nil
}

func expandAfdCustomizedCipherSuiteSet(input []interface{}) *afdcustomdomains.AFDDomainHTTPSCustomizedCipherSuiteSet {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	result := &afdcustomdomains.AFDDomainHTTPSCustomizedCipherSuiteSet{}

	if tls12Raw := v["tls12"].(*pluginsdk.Set); tls12Raw.Len() > 0 {
		tls12Suites := make([]afdcustomdomains.AfdCustomizedCipherSuiteForTls12, 0)
		for _, suite := range tls12Raw.List() {
			tls12Suites = append(tls12Suites, afdcustomdomains.AfdCustomizedCipherSuiteForTls12(suite.(string)))
		}
		result.CipherSuiteSetForTls12 = &tls12Suites
	}

	if tls13Raw := v["tls13"].(*pluginsdk.Set); tls13Raw.Len() > 0 {
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

func resourceCdnFrontDoorCustomDomainCipherSuiteConfigured(d *pluginsdk.ResourceData) bool {
	raw := d.GetRawConfig()
	if raw.IsNull() {
		return false
	}

	tlsConfig := raw.GetAttr("tls")
	if tlsConfig.IsNull() || tlsConfig.LengthInt() == 0 {
		return false
	}

	tlsBlock := tlsConfig.AsValueSlice()[0]
	if tlsBlock.IsNull() {
		return false
	}

	cipherConfig := tlsBlock.GetAttr("cipher_suite")
	if cipherConfig.IsNull() || cipherConfig.LengthInt() == 0 {
		return false
	}

	cipherBlock := cipherConfig.AsValueSlice()[0]
	return !cipherBlock.IsNull()
}

func flattenAfdDomainHttpsParameters(input *afdcustomdomains.AFDDomainHTTPSParameters, includeDefaultCipherSuite bool) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	secretId := ""
	if input.Secret != nil && input.Secret.Id != nil {
		if id, err := parse.FrontDoorSecretIDInsensitively(pointer.From(input.Secret.Id)); err == nil {
			secretId = id.ID()
		}
	}

	minTlsVersion := ""
	if input.MinimumTlsVersion != nil {
		minTlsVersion = string(pointer.From(input.MinimumTlsVersion))
	}

	// Azure omits `minimumTlsVersion` when the value is `TLS12`, so we default the field to
	// that setting to avoid spurious diffs when the user explicitly configured nothing.
	if minTlsVersion == "" {
		minTlsVersion = string(afdcustomdomains.AfdMinimumTlsVersionTLSOneTwo)
	}

	customCiphers := flattenAfdCustomizedCipherSuiteSet(input.CustomizedCipherSuiteSet)
	cipherSuiteType := ""
	if input.CipherSuiteSetType != nil {
		cipherSuiteType = string(pointer.From(input.CipherSuiteSetType))
	}

	// Azure always returns the default `TLS12_2022` cipher suite even when users never
	// configured a block, so we ignore that value to avoid perpetual diffs. Only emit the
	// block when the service reports a non-default type, custom entries exist, or the user
	// explicitly configured the block (tracked via includeDefaultCipherSuite).
	includeCipherSuite := len(customCiphers) > 0

	switch cipherSuiteType {
	case "":
		if includeDefaultCipherSuite {
			includeCipherSuite = true
		}
	case string(afdcustomdomains.AfdCipherSuiteSetTypeTLSOneTwoTwoZeroTwoTwo):
		if includeDefaultCipherSuite {
			includeCipherSuite = true
		}
	default:
		includeCipherSuite = true
	}

	cipherSuite := make([]interface{}, 0)
	if includeCipherSuite {
		cipherSuite = []interface{}{
			map[string]interface{}{
				"type":           cipherSuiteType,
				"custom_ciphers": customCiphers,
			},
		}
	}

	result := map[string]interface{}{
		"cdn_frontdoor_secret_id": secretId,
		"certificate_type":        string(input.CertificateType),
		"cipher_suite":            cipherSuite,
		"minimum_version":         minTlsVersion,
	}

	if !features.FivePointOh() {
		result["minimum_tls_version"] = minTlsVersion
	}

	return []interface{}{result}
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
			"tls12": tls12Suites,
			"tls13": tls13Suites,
		},
	}
}

func flattenAfdDNSZoneResourceReference(input *afdcustomdomains.ResourceReference) string {
	if input == nil || input.Id == nil {
		return ""
	}

	if id, err := dnsValidate.ParseDnsZoneIDInsensitively(pointer.From(input.Id)); err == nil {
		return id.ID()
	}

	return ""
}

func frontDoorCustomDomainHostNameCustomizeDiff(_ context.Context, diff *pluginsdk.ResourceDiff, _ interface{}) error {
	tlsAny := diff.Get("tls")
	tlsRaw, ok := tlsAny.([]interface{})
	if !ok {
		return errors.New("unexpected value for `tls`: expected list")
	}
	if len(tlsRaw) == 0 || tlsRaw[0] == nil {
		return nil
	}

	tls, ok := tlsRaw[0].(map[string]interface{})
	if !ok {
		return errors.New("unexpected value for `tls`: expected object")
	}

	certificateType := string(afdcustomdomains.AfdCertificateTypeManagedCertificate)
	if raw, exists := tls["certificate_type"]; exists && raw != nil {
		v, ok := raw.(string)
		if !ok {
			return errors.New("unexpected value for `tls.certificate_type`: expected string")
		}
		if v != "" {
			certificateType = v
		}
	}

	if certificateType != string(afdcustomdomains.AfdCertificateTypeManagedCertificate) {
		return nil
	}

	hostName, ok := diff.Get("host_name").(string)
	if !ok {
		return errors.New("unexpected value for `host_name`: expected string")
	}

	if len(hostName) > 64 {
		return errors.New("`host_name` cannot be longer than 64 characters when `tls.certificate_type` is `ManagedCertificate`")
	}

	if strings.HasPrefix(hostName, "*.") {
		return errors.New("`host_name` cannot be a wildcard domain when `tls.certificate_type` is `ManagedCertificate`")
	}

	return nil
}

func frontDoorCustomDomainTlsCustomizeDiff(_ context.Context, diff *pluginsdk.ResourceDiff, _ interface{}) error {
	tlsAny := diff.Get("tls")
	tlsRaw, ok := tlsAny.([]interface{})
	if !ok {
		return errors.New("unexpected value for `tls`: expected list")
	}
	if len(tlsRaw) == 0 || tlsRaw[0] == nil {
		return nil
	}

	tls, ok := tlsRaw[0].(map[string]interface{})
	if !ok {
		return errors.New("unexpected value for `tls`: expected object")
	}

	rawConfig := diff.GetRawConfig()
	minimumVersionConfigured := false
	minimumTlsVersionConfigured := false
	tls13Configured := false

	if !rawConfig.IsNull() {
		tlsConfig := rawConfig.GetAttr("tls")
		if !tlsConfig.IsNull() && tlsConfig.LengthInt() > 0 {
			tlsBlock := tlsConfig.AsValueSlice()[0]
			if !tlsBlock.IsNull() {
				minimumVersionConfigured = !tlsBlock.GetAttr("minimum_version").IsNull()
				if !features.FivePointOh() {
					minimumTlsVersionConfigured = !tlsBlock.GetAttr("minimum_tls_version").IsNull()
				}

				cipherConfig := tlsBlock.GetAttr("cipher_suite")
				if !cipherConfig.IsNull() && cipherConfig.LengthInt() > 0 {
					cipherBlock := cipherConfig.AsValueSlice()[0]
					if !cipherBlock.IsNull() {
						customCiphersConfig := cipherBlock.GetAttr("custom_ciphers")
						if !customCiphersConfig.IsNull() && customCiphersConfig.LengthInt() > 0 {
							customCiphersBlock := customCiphersConfig.AsValueSlice()[0]
							if !customCiphersBlock.IsNull() {
								tls13Config := customCiphersBlock.GetAttr("tls13")
								tls13Configured = !tls13Config.IsNull()
							}
						}
					}
				}
			}
		}
	}

	if !features.FivePointOh() {
		if minimumTlsVersionConfigured {
			if minTlsVersion := tls["minimum_tls_version"].(string); strings.EqualFold(minTlsVersion, string(afdcustomdomains.AfdMinimumTlsVersionTLSOneZero)) {
				return errors.New("support for TLS 1.0 and 1.1 was retired on March 1, 2025. Please use `minimum_version = \"TLS12\"` instead")
			}
		}
	}

	cipherSuiteAny, exists := tls["cipher_suite"]
	if !exists || cipherSuiteAny == nil {
		return nil
	}

	cipherSuiteRaw, ok := cipherSuiteAny.([]interface{})
	if !ok {
		return errors.New("unexpected value for `tls.cipher_suite`: expected list")
	}
	if cipherSuiteRaw == nil {
		return nil
	}
	if len(cipherSuiteRaw) == 0 || cipherSuiteRaw[0] == nil {
		return nil
	}

	cipherSuite, ok := cipherSuiteRaw[0].(map[string]interface{})
	if !ok {
		return errors.New("unexpected value for `tls.cipher_suite`: expected object")
	}
	cipherSuiteTypeAny, exists := cipherSuite["type"]
	if !exists || cipherSuiteTypeAny == nil {
		return errors.New("unexpected value for `tls.cipher_suite.type`: expected string")
	}
	cipherSuiteType, ok := cipherSuiteTypeAny.(string)
	if !ok {
		return errors.New("unexpected value for `tls.cipher_suite.type`: expected string")
	}
	customCiphersRaw := make([]interface{}, 0)
	if raw, exists := cipherSuite["custom_ciphers"]; exists && raw != nil {
		v, ok := raw.([]interface{})
		if !ok {
			return errors.New("unexpected value for `tls.cipher_suite.custom_ciphers`: expected list")
		}
		customCiphersRaw = v
	}

	if cipherSuiteType == string(afdcustomdomains.AfdCipherSuiteSetTypeCustomized) {
		if len(customCiphersRaw) == 0 {
			return errors.New("`custom_ciphers` is required when `type` is `Customized`")
		}

		if customCiphersRaw[0] == nil {
			return errors.New("at least one cipher suite must be selected in `custom_ciphers` when `type` is set to `Customized`")
		}

		customCiphers, ok := customCiphersRaw[0].(map[string]interface{})
		if !ok {
			return errors.New("unexpected value for `tls.cipher_suite.custom_ciphers`: expected object")
		}

		setLen := func(s *pluginsdk.Set) int {
			if s == nil {
				return 0
			}
			return s.Len()
		}

		var tls12Suites *pluginsdk.Set
		if raw, exists := customCiphers["tls12"]; exists && raw != nil {
			v, ok := raw.(*pluginsdk.Set)
			if !ok {
				return errors.New("unexpected value for `custom_ciphers.tls12`: expected set")
			}
			tls12Suites = v
		}

		var tls13Suites *pluginsdk.Set
		if raw, exists := customCiphers["tls13"]; exists && raw != nil {
			v, ok := raw.(*pluginsdk.Set)
			if !ok {
				return errors.New("unexpected value for `custom_ciphers.tls13`: expected set")
			}
			tls13Suites = v
		}

		if setLen(tls12Suites) == 0 && setLen(tls13Suites) == 0 {
			return errors.New("at least one cipher suite must be selected in `custom_ciphers` when `type` is set to `Customized`")
		}

		if tls13Configured {
			has128 := false
			has256 := false
			if tls13Suites != nil {
				for _, raw := range tls13Suites.List() {
					v, ok := raw.(string)
					if !ok {
						continue
					}
					switch v {
					case "TLS_AES_128_GCM_SHA256":
						has128 = true
					case "TLS_AES_256_GCM_SHA384":
						has256 = true
					}
				}
			}

			if !has128 || !has256 {
				return errors.New("`custom_ciphers.tls13` must contain both `TLS_AES_128_GCM_SHA256` and `TLS_AES_256_GCM_SHA384` when specified")
			}
		}

		minimumVersion := ""

		if !features.FivePointOh() {
			switch {
			case minimumVersionConfigured:
				minimumVersion = tls["minimum_version"].(string)
			case minimumTlsVersionConfigured:
				minimumVersion = tls["minimum_tls_version"].(string)
			default:
				minimumVersion = string(afdcustomdomains.AfdMinimumTlsVersionTLSOneTwo)
			}
		} else {
			if rawMin := tls["minimum_version"]; rawMin != nil {
				if minStr, ok := rawMin.(string); ok {
					minimumVersion = minStr
				} else {
					return errors.New("unexpected value for `tls.minimum_version`: expected string")
				}
			}
		}

		if minimumVersion == string(afdcustomdomains.AfdMinimumTlsVersionTLSOneTwo) && setLen(tls12Suites) == 0 {
			return errors.New("at least one TLS 1.2 cipher suite must be specified in `custom_ciphers.tls12` when `minimum_version` is set to `TLS12`")
		}

		if minimumVersion == string(afdcustomdomains.AfdMinimumTlsVersionTLSOneThree) && tls13Configured && setLen(tls13Suites) == 0 {
			return errors.New("at least one TLS 1.3 cipher suite must be specified in `custom_ciphers.tls13` when `minimum_version` is set to `TLS13`")
		}
	} else if len(customCiphersRaw) > 0 && customCiphersRaw[0] != nil {
		return errors.New("`custom_ciphers` cannot be specified when `type` is not `Customized`")
	}

	return nil
}
