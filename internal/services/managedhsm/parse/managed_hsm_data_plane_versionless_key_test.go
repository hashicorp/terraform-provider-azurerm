// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

func TestParseManagedHSMDataPlaneVersionlessKeyID_NoDomainSuffix_InvalidValuesFail(t *testing.T) {
	values := []string{
		"",                                       // empty = invalid
		"https://example.com/keys/abc123/foobar", // Hostname is incomplete
		"https://example.keyvault.azure.net/keys/abc123",             // Key Vault Key
		"https://example.managedhsm.azure.net/",                      // no path
		"https://example.managedhsm.azure.net/keys/",                 // trailing slash - no key
		"https://example.managedhsm.azure.net/keys/abc123/bcd234",    // with version
		"https://example.managedhsm.azure.net/numbers/abc123/bcd234", // wrong type
		"http://example.managedhsm.azure.net:80/keys/abc123",         // HTTP rather than HTTPS
	}
	for _, input := range values {
		t.Logf("Validating %q", input)
		actual, err := ManagedHSMDataPlaneVersionlessKeyID(input, nil)
		if err != nil {
			continue
		}
		t.Fatalf("unexpected value for %q: %q", input, actual.ID())
	}
}

func TestParseManagedHSMDataPlaneVersionlessKeyID_NoDomainSuffix_ValidValues(t *testing.T) {
	// NOTE: this scenario tests the Validation use-case - where a Domain Suffix isn't known
	// since the Environment/Credentials aren't available until Provider initialization.
	values := map[string]ManagedHSMDataPlaneVersionlessKeyId{
		"https://example.managedhsm.azure.net/keys/abc123": {
			// Public
			ManagedHSMName: "example",
			DomainSuffix:   "managedhsm.azure.net",
			KeyName:        "abc123",
		},
		"https://EXAMPLE.managedhsm.azure.net/keys/abc123": {
			// Public but the uppercase name should be normalised
			ManagedHSMName: "example",
			DomainSuffix:   "managedhsm.azure.net",
			KeyName:        "abc123",
		},
		"https://example.managedhsm.azure.net:443/keys/abc123": {
			// in this instance the domain suffix includes a port which should be removed
			// this is a bug in the Azure API response data, so we should filter it out
			ManagedHSMName: "example",
			DomainSuffix:   "managedhsm.azure.net",
			KeyName:        "abc123",
		},
		"https://example.managedhsm.azure.cn/keys/abc123": {
			// China
			ManagedHSMName: "example",
			DomainSuffix:   "managedhsm.azure.cn",
			KeyName:        "abc123",
		},
		"https://example.managedhsm.usgovcloudapi.net/keys/abc123": {
			// US Gov
			ManagedHSMName: "example",
			DomainSuffix:   "managedhsm.usgovcloudapi.net",
			KeyName:        "abc123",
		},
	}
	for input, expected := range values {
		t.Logf("Validating %q", input)
		actual, err := ManagedHSMDataPlaneVersionlessKeyID(input, nil)
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
	}
}

func TestParseManagedHSMDataPlaneVersionlessKeyID_WithDomainSuffix_InvalidValuesFail(t *testing.T) {
	values := []string{
		"",                                       // empty = invalid
		"https://example.com/keys/abc123/foobar", // Hostname is incomplete
		"https://example.keyvault.azure.net/keys/abc123",             // Key Vault Key
		"https://example.managedhsm.azure.net/",                      // no path
		"https://example.managedhsm.azure.net/keys/",                 // trailing slash - no key
		"https://example.managedhsm.azure.net/keys/abc123/bcd234",    // with version
		"https://example.managedhsm.azure.net/numbers/abc123/bcd234", // wrong type
		"http://example.managedhsm.azure.net:80/keys/abc123",         // HTTP rather than HTTPS
		"https://managedhsm.azure.net/keys/foo",                      // hostname is only the domainSuffix
		"https://example.managedhsm.some.domain/keys/foo",            // hostname doesn't contain the domainSuffix
	}
	for _, input := range values {
		t.Logf("Validating %q", input)
		actual, err := ManagedHSMDataPlaneVersionlessKeyID(input, pointer.To("managedhsm.azure.net"))
		if err == nil {
			t.Fatalf("unexpected value for %q: %q", input, actual.ID())
		}
	}
}

func TestParseManagedHSMDataPlaneVersionlessKeyID_WithDomainSuffix_ValidValues(t *testing.T) {
	// NOTE: this scenario tests the Validation use-case - where a Domain Suffix isn't known
	// since the Environment/Credentials aren't available until Provider initialization.
	values := map[string]struct {
		expected     ManagedHSMDataPlaneVersionlessKeyId
		domainSuffix string
	}{
		"https://example.managedhsm.azure.net/keys/abc123": {
			// Public
			domainSuffix: "managedhsm.azure.net",
			expected: ManagedHSMDataPlaneVersionlessKeyId{
				ManagedHSMName: "example",
				DomainSuffix:   "managedhsm.azure.net",
				KeyName:        "abc123",
			},
		},
		"https://EXAMPLE.managedhsm.azure.net/keys/abc123": {
			// Public but the uppercase name should be normalised
			domainSuffix: "managedhsm.azure.net",
			expected: ManagedHSMDataPlaneVersionlessKeyId{
				ManagedHSMName: "example",
				DomainSuffix:   "managedhsm.azure.net",
				KeyName:        "abc123",
			},
		},
		"https://example.managedhsm.azure.net:443/keys/abc123": {
			// in this instance the domain suffix includes a port which should be removed
			// this is a bug in the Azure API response data, so we should filter it out
			domainSuffix: "managedhsm.azure.net",
			expected: ManagedHSMDataPlaneVersionlessKeyId{
				ManagedHSMName: "example",
				DomainSuffix:   "managedhsm.azure.net",
				KeyName:        "abc123",
			},
		},
		"https://example.managedhsm.azure.cn/keys/abc123": {
			// China
			domainSuffix: "managedhsm.azure.cn",
			expected: ManagedHSMDataPlaneVersionlessKeyId{
				ManagedHSMName: "example",
				DomainSuffix:   "managedhsm.azure.cn",
				KeyName:        "abc123",
			},
		},
		"https://example.managedhsm.usgovcloudapi.net/keys/abc123": {
			// US Gov
			domainSuffix: "managedhsm.usgovcloudapi.net",
			expected: ManagedHSMDataPlaneVersionlessKeyId{
				ManagedHSMName: "example",
				DomainSuffix:   "managedhsm.usgovcloudapi.net",
				KeyName:        "abc123",
			},
		},
	}
	for input, item := range values {
		t.Logf("Validating %q", input)
		actual, err := ManagedHSMDataPlaneVersionlessKeyID(input, pointer.To(item.domainSuffix))
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
	}
}
