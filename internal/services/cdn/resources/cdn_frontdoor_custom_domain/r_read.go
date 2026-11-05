// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_custom_domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorCustomDomainRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.AFDCustomDomainsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afddomains.ParseCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.AFDCustomDomainsGet(ctx, *id)
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

			if err := d.Set("dns_zone_id", flattenAfdDNSZoneResourceReference(props.AzureDnsZone)); err != nil {
				return fmt.Errorf("setting `dns_zone_id`: %+v", err)
			}

			includeDefaultCipherSuite := resourceCdnFrontDoorCustomDomainCipherSuiteConfigured(d)
			if err := d.Set("tls", flattenAfdDomainHttpsParameters(props.TlsSettings, includeDefaultCipherSuite)); err != nil {
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

// Front Door currently rejects the DHE TLS 1.2 suites exposed by the SDK helper.
// The service team confirmed this supported set officially, so this is not a
// temporary workaround. We keep an explicit allowlist of the accepted ECDHE
// values here.
var frontDoorCustomDomainCustomizedCipherSuitesForTls12 = []string{
	string(afddomains.AfdCustomizedCipherSuiteForTls12ECDHERSAAESOneTwoEightGCMSHATwoFiveSix),
	string(afddomains.AfdCustomizedCipherSuiteForTls12ECDHERSAAESTwoFiveSixGCMSHAThreeEightFour),
	string(afddomains.AfdCustomizedCipherSuiteForTls12ECDHERSAAESOneTwoEightSHATwoFiveSix),
	string(afddomains.AfdCustomizedCipherSuiteForTls12ECDHERSAAESTwoFiveSixSHAThreeEightFour),
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

func frontDoorCustomDomainHostNameCustomizeDiff(_ context.Context, diff *pluginsdk.ResourceDiff, _ any) error {
	tlsAny := diff.Get("tls")
	tlsRaw, ok := tlsAny.([]any)
	if !ok {
		return errors.New("unexpected value for `tls`: expected list")
	}
	if len(tlsRaw) == 0 || tlsRaw[0] == nil {
		return nil
	}

	tls, ok := tlsRaw[0].(map[string]any)
	if !ok {
		return errors.New("unexpected value for `tls`: expected object")
	}

	certificateType := string(afddomains.AfdCertificateTypeManagedCertificate)
	if raw, exists := tls["certificate_type"]; exists && raw != nil {
		v, ok := raw.(string)
		if !ok {
			return errors.New("unexpected value for `tls.certificate_type`: expected string")
		}
		if v != "" {
			certificateType = v
		}
	}

	if certificateType != string(afddomains.AfdCertificateTypeManagedCertificate) {
		return nil
	}

	hostName, ok := diff.Get("host_name").(string)
	if !ok {
		return errors.New("unexpected value for `host_name`: expected string")
	}

	if len(hostName) > 64 {
		return errors.New("`host_name` cannot be longer than 64 characters when `tls.certificate_type` is `ManagedCertificate`")
	}

	return nil
}

func frontDoorCustomDomainTlsCustomizeDiff(_ context.Context, diff *pluginsdk.ResourceDiff, _ any) error {
	tlsAny := diff.Get("tls")
	tlsRaw, ok := tlsAny.([]any)
	if !ok {
		return errors.New("unexpected value for `tls`: expected list")
	}
	if len(tlsRaw) == 0 || tlsRaw[0] == nil {
		return nil
	}

	tls, ok := tlsRaw[0].(map[string]any)
	if !ok {
		return errors.New("unexpected value for `tls`: expected object")
	}

	rawConfig := diff.GetRawConfig()
	tls13Configured := false

	if !rawConfig.IsNull() {
		tlsConfig := rawConfig.GetAttr("tls")
		if !tlsConfig.IsNull() && tlsConfig.LengthInt() > 0 {
			tlsBlock := tlsConfig.AsValueSlice()[0]
			if !tlsBlock.IsNull() {
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

	cipherSuiteAny, exists := tls["cipher_suite"]
	if !exists || cipherSuiteAny == nil {
		return nil
	}

	cipherSuiteRaw, ok := cipherSuiteAny.([]any)
	if !ok {
		return errors.New("unexpected value for `tls.cipher_suite`: expected list")
	}
	if cipherSuiteRaw == nil {
		return nil
	}
	if len(cipherSuiteRaw) == 0 || cipherSuiteRaw[0] == nil {
		return nil
	}

	cipherSuite, ok := cipherSuiteRaw[0].(map[string]any)
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
	customCiphersRaw := make([]any, 0)
	if raw, exists := cipherSuite["custom_ciphers"]; exists && raw != nil {
		v, ok := raw.([]any)
		if !ok {
			return errors.New("unexpected value for `tls.cipher_suite.custom_ciphers`: expected list")
		}
		customCiphersRaw = v
	}

	if cipherSuiteType == string(afddomains.AfdCipherSuiteSetTypeCustomized) {
		if len(customCiphersRaw) == 0 {
			return errors.New("`custom_ciphers` is required when `type` is `Customized`")
		}

		if customCiphersRaw[0] == nil {
			return errors.New("at least one cipher suite must be selected in `custom_ciphers` when `type` is set to `Customized`")
		}

		customCiphers, ok := customCiphersRaw[0].(map[string]any)
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

		if rawMin := tls["minimum_version"]; rawMin != nil {
			if minStr, ok := rawMin.(string); ok {
				minimumVersion = minStr
			} else {
				return errors.New("unexpected value for `tls.minimum_version`: expected string")
			}
		}

		if minimumVersion == string(afddomains.AfdMinimumTlsVersionTLSOneTwo) && setLen(tls12Suites) == 0 {
			return errors.New("at least one TLS 1.2 cipher suite must be specified in `custom_ciphers.tls12` when `minimum_version` is set to `TLS12`")
		}

		if minimumVersion == string(afddomains.AfdMinimumTlsVersionTLSOneThree) && tls13Configured && setLen(tls13Suites) == 0 {
			return errors.New("at least one TLS 1.3 cipher suite must be specified in `custom_ciphers.tls13` when `minimum_version` is set to `TLS13`")
		}
	} else if len(customCiphersRaw) > 0 && customCiphersRaw[0] != nil {
		return errors.New("`custom_ciphers` cannot be specified when `type` is not `Customized`")
	}

	return nil
}
