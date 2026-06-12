// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestSchemaLocalDNSOverrideHashIncludesSettings(t *testing.T) {
	localDNSOverride := schemaLocalDNSOverride()

	localDNSOverrides := localDNSOverride.ZeroValue().(*schema.Set)
	localDNSOverrides.Add(map[string]interface{}{
		"domain":                    ".",
		"cache_duration_in_seconds": 3600,
		"forward_destination":       "VnetDNS",
		"forward_policy":            "Sequential",
		"max_concurrent":            1000,
		"protocol":                  "PreferUDP",
		"query_logging":             "Error",
		"serve_stale":               "Verify",
		"serve_stale_duration":      3600,
	})
	localDNSOverrides.Add(map[string]interface{}{
		"domain":                    ".",
		"cache_duration_in_seconds": 7200,
		"forward_destination":       "VnetDNS",
		"forward_policy":            "RoundRobin",
		"max_concurrent":            2000,
		"protocol":                  "PreferUDP",
		"query_logging":             "Error",
		"serve_stale":               "Verify",
		"serve_stale_duration":      3600,
	})

	if localDNSOverrides.Len() != 2 {
		t.Fatalf("expected two local DNS override hashes, got %d", localDNSOverrides.Len())
	}
}

func TestValidateLocalDNSOverrideRawConfigSkipsUnknownDomain(t *testing.T) {
	input := cty.ListVal([]cty.Value{
		cty.ObjectVal(map[string]cty.Value{
			"domain": cty.UnknownVal(cty.String),
		}),
		cty.ObjectVal(map[string]cty.Value{
			"domain": cty.UnknownVal(cty.String),
		}),
	})

	if err := validateLocalDNSOverrideRawConfig(input, "local_dns_profile.kube_dns_override"); err != nil {
		t.Fatalf("expected unknown domains to be skipped, got %q", err)
	}
}

func TestValidateLocalDNSOverrideRawConfigRejectsDuplicateDomain(t *testing.T) {
	input := cty.ListVal([]cty.Value{
		cty.ObjectVal(map[string]cty.Value{
			"domain": cty.StringVal("example.com"),
		}),
		cty.ObjectVal(map[string]cty.Value{
			"domain": cty.StringVal("example.com"),
		}),
	})

	if err := validateLocalDNSOverrideRawConfig(input, "local_dns_profile.kube_dns_override"); err == nil {
		t.Fatal("expected duplicate domain validation error")
	}
}
