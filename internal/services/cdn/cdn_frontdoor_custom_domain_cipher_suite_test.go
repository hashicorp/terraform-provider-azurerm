// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceCdnFrontDoorCustomDomainCipherSuiteConfigured(t *testing.T) {
	resource := resourceCdnFrontDoorCustomDomain()

	// A configured `cipher_suite` block must be detected even when the raw
	// configuration is unavailable (the refresh path, where Terraform only sends
	// prior state). `schema.TestResourceDataRaw` builds a ResourceData whose
	// `GetRawConfig()` is null while `d.Get("tls")` returns the supplied value,
	// faithfully reproducing the refresh scenario that previously caused the
	// default `TLS12_2022` cipher suite to be dropped from state.
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"tls": []interface{}{
			map[string]interface{}{
				"cipher_suite": []interface{}{
					map[string]interface{}{
						"type": "TLS12_2022",
					},
				},
			},
		},
	})
	if !resourceCdnFrontDoorCustomDomainCipherSuiteConfigured(d) {
		t.Fatal("expected a configured cipher_suite block to be detected")
	}

	// No cipher_suite block configured: the helper must return false so the
	// service-returned default `TLS12_2022` value is dropped and does not cause a
	// perpetual diff.
	dWithoutCipherSuite := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"tls": []interface{}{
			map[string]interface{}{},
		},
	})
	if resourceCdnFrontDoorCustomDomainCipherSuiteConfigured(dWithoutCipherSuite) {
		t.Fatal("expected the absence of a cipher_suite block to be detected")
	}
}
