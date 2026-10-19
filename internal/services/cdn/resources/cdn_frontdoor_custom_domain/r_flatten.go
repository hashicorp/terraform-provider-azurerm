// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_custom_domain

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/secrets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
)

func flattenAfdDomainHttpsParameters(input *afddomains.AFDDomainHTTPSParameters, includeDefaultCipherSuite bool) []any {
	if input == nil {
		return []any{}
	}

	secretId := ""
	if input.Secret != nil && input.Secret.Id != nil {
		if id, err := secrets.ParseSecretIDInsensitively(pointer.From(input.Secret.Id)); err == nil {
			secretId = id.ID()
		}
	}

	minTlsVersion := ""
	if input.MinimumTlsVersion != nil {
		minTlsVersion = pointer.FromEnum(input.MinimumTlsVersion)
	}

	// Azure omits `minimumTlsVersion` when the value is `TLS12`, so we default the field to
	// that setting to avoid spurious diffs when the user explicitly configured nothing.
	if minTlsVersion == "" {
		minTlsVersion = string(afddomains.AfdMinimumTlsVersionTLSOneTwo)
	}

	customCiphers := flattenAfdCustomizedCipherSuiteSet(input.CustomizedCipherSuiteSet)
	cipherSuiteType := ""
	if input.CipherSuiteSetType != nil {
		cipherSuiteType = pointer.FromEnum(input.CipherSuiteSetType)
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
	case string(afddomains.AfdCipherSuiteSetTypeTLSOneTwoTwoZeroTwoTwo):
		if includeDefaultCipherSuite {
			includeCipherSuite = true
		}
	default:
		includeCipherSuite = true
	}

	cipherSuite := make([]any, 0)
	if includeCipherSuite {
		cipherSuite = []any{
			map[string]any{
				"type":           cipherSuiteType,
				"custom_ciphers": customCiphers,
			},
		}
	}

	result := map[string]any{
		"cdn_frontdoor_secret_id": secretId,
		"certificate_type":        string(input.CertificateType),
		"cipher_suite":            cipherSuite,
		"minimum_version":         minTlsVersion,
	}

	return []any{result}
}

func flattenAfdCustomizedCipherSuiteSet(input *afddomains.AFDDomainHTTPSCustomizedCipherSuiteSet) []any {
	if input == nil {
		return []any{}
	}

	tls12Suites := make([]any, 0)
	if input.CipherSuiteSetForTls12 != nil {
		for _, suite := range *input.CipherSuiteSetForTls12 {
			tls12Suites = append(tls12Suites, string(suite))
		}
	}

	tls13Suites := make([]any, 0)
	if input.CipherSuiteSetForTls13 != nil {
		for _, suite := range *input.CipherSuiteSetForTls13 {
			tls13Suites = append(tls13Suites, string(suite))
		}
	}

	if len(tls12Suites) == 0 && len(tls13Suites) == 0 {
		return []any{}
	}

	return []any{
		map[string]any{
			"tls12": tls12Suites,
			"tls13": tls13Suites,
		},
	}
}

func flattenAfdDNSZoneResourceReference(input *afddomains.ResourceReference) string {
	if input == nil || input.Id == nil {
		return ""
	}

	if id, err := zones.ParseDnsZoneIDInsensitively(pointer.From(input.Id)); err == nil {
		return id.ID()
	}

	return ""
}
