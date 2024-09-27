// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

func TestParseManagedHSMDataPlaneVersionedKeyID_NoDomainSuffix_InvalidValuesFail(t *testing.T) {
	values := []string{
		"",                                       // empty = invalid
		"https://example.com/keys/abc123/foobar", // Hostname is incomplete
		"https://example.keyvault.azure.net/keys/abc123/foobar",                     // Key Vault Key
		"https://example.managedhsm.azure.net/",                                     // no path
		"https://example.managedhsm.azure.net/keys/",                                // trailing slash - no key/no version
		"https://example.managedhsm.azure.net/keys/abc123",                          // no version
		"https://example.managedhsm.azure.net/keys/abc123/",                         // trailing slash but no version
		"https://example.managedhsm.azure.net/keys/abc123/bcd234/foobar",            // too many segments
		"https://example.managedhsm.azure.net/numbers/abc123/bcd234",                // wrong type
		"http://example.managedhsm.azure.net:80/keys/abc123/123456789oijhgfdertyhj", // HTTP rather than HTTPS
	}
	for _, input := range values {
		t.Logf("Validating %q", input)
		actual, err := ManagedHSMDataPlaneVersionedKeyID(input, nil)
		if err != nil {
			continue
		}
		t.Fatalf("unexpected value for %q: %q", input, actual.ID())
	}
}

func TestParseManagedHSMDataPlaneVersionedKeyID_NoDomainSuffix_ValidValues(t *testing.T) {
	// NOTE: this scenario tests the Validation use-case - where a Domain Suffix isn't known
	// since the Environment/Credentials aren't available until Provider initialization.
	values := map[string]ManagedHSMDataPlaneVersionedKeyId{
		"https://example.managedhsm.azure.net/keys/abc123/123456789oijhgfdertyhj": {
			// Public
			ManagedHSMName: "example",
			DomainSuffix:   "managedhsm.azure.net",
			KeyName:        "abc123",
			KeyVersion:     "123456789oijhgfdertyhj",
		},
		"https://EXAMPLE.managedhsm.azure.net/keys/abc123/123456789oijhgfdertyhj": {
			// Public but the uppercase name should be normalised
			ManagedHSMName: "example",
			DomainSuffix:   "managedhsm.azure.net",
			KeyName:        "abc123",
			KeyVersion:     "123456789oijhgfdertyhj",
		},
		"https://example.managedhsm.azure.net:443/keys/abc123/123456789oijhgfdertyhj": {
			// in this instance the domain suffix includes a port which should be removed
			// this is a bug in the Azure API response data, so we should filter it out
			ManagedHSMName: "example",
			DomainSuffix:   "managedhsm.azure.net",
			KeyName:        "abc123",
			KeyVersion:     "123456789oijhgfdertyhj",
		},
		"https://example.managedhsm.azure.cn/keys/abc123/123456789oijhgfdertyhj": {
			// China
			ManagedHSMName: "example",
			DomainSuffix:   "managedhsm.azure.cn",
			KeyName:        "abc123",
			KeyVersion:     "123456789oijhgfdertyhj",
		},
		"https://example.managedhsm.usgovcloudapi.net/keys/abc123/123456789oijhgfdertyhj": {
			// US Gov
			ManagedHSMName: "example",
			DomainSuffix:   "managedhsm.usgovcloudapi.net",
			KeyName:        "abc123",
			KeyVersion:     "123456789oijhgfdertyhj",
		},
	}
	for input, expected := range values {
		t.Logf("Validating %q", input)
		actual, err := ManagedHSMDataPlaneVersionedKeyID(input, nil)
		if err != nil {
			t.Fatalf("unexpected error: %+v", err.Error())
		}

		if actual.ManagedHSMName != expected.ManagedHSMName {
			t.Fatalf("expected `ManagedHSMName` to be %q but got %q", expected.ManagedHSMName, actual.ManagedHSMName)
		}
		if actual.DomainSuffix != expected.DomainSuffix {
			t.Fatalf("expected `DomainSuffix` to be %q but got %q", expected.DomainSuffix, actual.DomainSuffix)
		}
		if actual.KeyName != expected.KeyName {
			t.Fatalf("expected `KeyName` to be %q but got %q", expected.KeyName, actual.KeyName)
		}
		if actual.KeyVersion != expected.KeyVersion {
			t.Fatalf("expected `KeyVersion` to be %q but got %q", expected.KeyVersion, actual.KeyVersion)
		}
	}
}

func TestParseManagedHSMDataPlaneVersionedKeyID_WithDomainSuffix_InvalidValuesFail(t *testing.T) {
	values := []string{
		"",                                       // empty = invalid
		"https://example.com/keys/abc123/foobar", // Hostname is incomplete
		"https://example.keyvault.azure.net/keys/abc123/foobar",                     // Key Vault Key
		"https://example.managedhsm.azure.net/",                                     // no path
		"https://example.managedhsm.azure.net/keys/",                                // trailing slash - no key/no version
		"https://example.managedhsm.azure.net/keys/abc123",                          // no version
		"https://example.managedhsm.azure.net/keys/abc123/",                         // trailing slash but no version
		"https://example.managedhsm.azure.net/keys/abc123/bcd234/foobar",            // too many segments
		"https://example.managedhsm.azure.net/numbers/abc123/bcd234",                // wrong type
		"http://example.managedhsm.azure.net:80/keys/abc123/123456789oijhgfdertyhj", // HTTP rather than HTTPS
		"https://managedhsm.azure.net/keys/foo/bar",                                 // hostname is only the domainSuffix
		"https://example.managedhsm.some.domain/keys/foo/bar",                       // hostname doesn't contain the domainSuffix
	}
	for _, input := range values {
		t.Logf("Validating %q", input)
		actual, err := ManagedHSMDataPlaneVersionedKeyID(input, pointer.To("managedhsm.azure.net"))
		if err == nil {
			t.Fatalf("unexpected value for %q: %q", input, actual.ID())
		}
	}
}

func TestParseManagedHSMDataPlaneVersionedKeyID_WithDomainSuffix_ValidValues(t *testing.T) {
	// NOTE: this scenario tests the Validation use-case - where a Domain Suffix isn't known
	// since the Environment/Credentials aren't available until Provider initialization.
	values := map[string]struct {
		expected     ManagedHSMDataPlaneVersionedKeyId
		domainSuffix string
	}{
		"https://example.managedhsm.azure.net/keys/abc123/123456789oijhgfdertyhj": {
			// Public
			domainSuffix: "managedhsm.azure.net",
			expected: ManagedHSMDataPlaneVersionedKeyId{
				ManagedHSMName: "example",
				DomainSuffix:   "managedhsm.azure.net",
				KeyName:        "abc123",
				KeyVersion:     "123456789oijhgfdertyhj",
			},
		},
		"https://EXAMPLE.managedhsm.azure.net/keys/abc123/123456789oijhgfdertyhj": {
			// Public but the uppercase name should be normalised
			domainSuffix: "managedhsm.azure.net",
			expected: ManagedHSMDataPlaneVersionedKeyId{
				ManagedHSMName: "example",
				DomainSuffix:   "managedhsm.azure.net",
				KeyName:        "abc123",
				KeyVersion:     "123456789oijhgfdertyhj",
			},
		},
		"https://example.managedhsm.azure.net:443/keys/abc123/123456789oijhgfdertyhj": {
			// in this instance the domain suffix includes a port which should be removed
			// this is a bug in the Azure API response data, so we should filter it out
			domainSuffix: "managedhsm.azure.net",
			expected: ManagedHSMDataPlaneVersionedKeyId{
				ManagedHSMName: "example",
				DomainSuffix:   "managedhsm.azure.net",
				KeyName:        "abc123",
				KeyVersion:     "123456789oijhgfdertyhj",
			},
		},
		"https://example.managedhsm.azure.cn/keys/abc123/123456789oijhgfdertyhj": {
			// China
			domainSuffix: "managedhsm.azure.cn",
			expected: ManagedHSMDataPlaneVersionedKeyId{
				ManagedHSMName: "example",
				DomainSuffix:   "managedhsm.azure.cn",
				KeyName:        "abc123",
				KeyVersion:     "123456789oijhgfdertyhj",
			},
		},
		"https://example.managedhsm.usgovcloudapi.net/keys/abc123/123456789oijhgfdertyhj": {
			// US Gov
			domainSuffix: "managedhsm.usgovcloudapi.net",
			expected: ManagedHSMDataPlaneVersionedKeyId{
				ManagedHSMName: "example",
				DomainSuffix:   "managedhsm.usgovcloudapi.net",
				KeyName:        "abc123",
				KeyVersion:     "123456789oijhgfdertyhj",
			},
		},
	}
	for input, item := range values {
		t.Logf("Validating %q", input)
		actual, err := ManagedHSMDataPlaneVersionedKeyID(input, pointer.To(item.domainSuffix))
		if err != nil {
			t.Fatalf("unexpected error: %+v", err.Error())
		}

		if actual.ManagedHSMName != item.expected.ManagedHSMName {
			t.Fatalf("expected `ManagedHSMName` to be %q but got %q", item.expected.ManagedHSMName, actual.ManagedHSMName)
		}
		if actual.DomainSuffix != item.expected.DomainSuffix {
			t.Fatalf("expected `DomainSuffix` to be %q but got %q", item.expected.DomainSuffix, actual.DomainSuffix)
		}
		if actual.KeyName != item.expected.KeyName {
			t.Fatalf("expected `KeyName` to be %q but got %q", item.expected.KeyName, actual.KeyName)
		}
		if actual.KeyVersion != item.expected.KeyVersion {
			t.Fatalf("expected `KeyVersion` to be %q but got %q", item.expected.KeyVersion, actual.KeyVersion)
		}
	}
}
