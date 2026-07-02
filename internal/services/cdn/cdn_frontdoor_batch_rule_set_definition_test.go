// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
)

func TestFlattenRouteConfigurationOverrideAction_CacheConfiguration_Disabled(t *testing.T) {
	input := rules.RouteConfigurationOverrideActionParameters{}

	output, err := flattenRouteConfigurationOverrideAction(input)
	if err != nil {
		t.Fatalf("expected no error but got %q", err)
	}

	if output.CacheBehavior != string(rules.RuleIsCompressionEnabledDisabled) {
		t.Fatalf("expected cache_behavior %q but got %q", rules.RuleIsCompressionEnabledDisabled, output.CacheBehavior)
	}
}

func TestFlattenRouteConfigurationOverrideAction_OriginGroupOverride_WithoutOriginGroup(t *testing.T) {
	input := rules.RouteConfigurationOverrideActionParameters{
		OriginGroupOverride: &rules.OriginGroupOverride{
			ForwardingProtocol: pointer.To(rules.ForwardingProtocolHTTPSOnly),
		},
	}

	output, err := flattenRouteConfigurationOverrideAction(input)
	if err != nil {
		t.Fatalf("expected no error but got %q", err)
	}

	if output.CdnFrontDoorOriginGroupID != "" {
		t.Fatalf("expected empty cdn_frontdoor_origin_group_id but got %q", output.CdnFrontDoorOriginGroupID)
	}

	if output.ForwardingProtocol != string(rules.ForwardingProtocolHTTPSOnly) {
		t.Fatalf("expected forwarding_protocol %q but got %q", rules.ForwardingProtocolHTTPSOnly, output.ForwardingProtocol)
	}
}
