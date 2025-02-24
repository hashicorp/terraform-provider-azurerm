// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5serverlogging

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/internal/logging"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ValidateResourceTypeConfigClientCapabilities generates a TRACE "Announced client capabilities" log.
func ValidateResourceTypeConfigClientCapabilities(ctx context.Context, capabilities *tfprotov5.ValidateResourceTypeConfigClientCapabilities) {
	if capabilities == nil {
		logging.ProtocolTrace(ctx, "No announced client capabilities", map[string]interface{}{})
		return
	}

	responseFields := map[string]interface{}{
		logging.KeyClientCapabilityWriteOnlyAttributesAllowed: capabilities.WriteOnlyAttributesAllowed,
	}

	logging.ProtocolTrace(ctx, "Announced client capabilities", responseFields)
}

// ConfigureProviderClientCapabilities generates a TRACE "Announced client capabilities" log.
func ConfigureProviderClientCapabilities(ctx context.Context, capabilities *tfprotov5.ConfigureProviderClientCapabilities) {
	if capabilities == nil {
		logging.ProtocolTrace(ctx, "No announced client capabilities", map[string]interface{}{})
		return
	}

	responseFields := map[string]interface{}{
		logging.KeyClientCapabilityDeferralAllowed: capabilities.DeferralAllowed,
	}

	logging.ProtocolTrace(ctx, "Announced client capabilities", responseFields)
}

// ReadDataSourceClientCapabilities generates a TRACE "Announced client capabilities" log.
func ReadDataSourceClientCapabilities(ctx context.Context, capabilities *tfprotov5.ReadDataSourceClientCapabilities) {
	if capabilities == nil {
		logging.ProtocolTrace(ctx, "No announced client capabilities", map[string]interface{}{})
		return
	}

	responseFields := map[string]interface{}{
		logging.KeyClientCapabilityDeferralAllowed: capabilities.DeferralAllowed,
	}

	logging.ProtocolTrace(ctx, "Announced client capabilities", responseFields)
}

// ReadResourceClientCapabilities generates a TRACE "Announced client capabilities" log.
func ReadResourceClientCapabilities(ctx context.Context, capabilities *tfprotov5.ReadResourceClientCapabilities) {
	if capabilities == nil {
		logging.ProtocolTrace(ctx, "No announced client capabilities", map[string]interface{}{})
		return
	}

	responseFields := map[string]interface{}{
		logging.KeyClientCapabilityDeferralAllowed: capabilities.DeferralAllowed,
	}

	logging.ProtocolTrace(ctx, "Announced client capabilities", responseFields)
}

// PlanResourceChangeClientCapabilities generates a TRACE "Announced client capabilities" log.
func PlanResourceChangeClientCapabilities(ctx context.Context, capabilities *tfprotov5.PlanResourceChangeClientCapabilities) {
	if capabilities == nil {
		logging.ProtocolTrace(ctx, "No announced client capabilities", map[string]interface{}{})
		return
	}

	responseFields := map[string]interface{}{
		logging.KeyClientCapabilityDeferralAllowed: capabilities.DeferralAllowed,
	}

	logging.ProtocolTrace(ctx, "Announced client capabilities", responseFields)
}

// ImportResourceStateClientCapabilities generates a TRACE "Announced client capabilities" log.
func ImportResourceStateClientCapabilities(ctx context.Context, capabilities *tfprotov5.ImportResourceStateClientCapabilities) {
	if capabilities == nil {
		logging.ProtocolTrace(ctx, "No announced client capabilities", map[string]interface{}{})
		return
	}

	responseFields := map[string]interface{}{
		logging.KeyClientCapabilityDeferralAllowed: capabilities.DeferralAllowed,
	}

	logging.ProtocolTrace(ctx, "Announced client capabilities", responseFields)
}

// OpenEphemeralResourceClientCapabilities generates a TRACE "Announced client capabilities" log.
func OpenEphemeralResourceClientCapabilities(ctx context.Context, capabilities *tfprotov5.OpenEphemeralResourceClientCapabilities) {
	if capabilities == nil {
		logging.ProtocolTrace(ctx, "No announced client capabilities", map[string]interface{}{})
		return
	}

	responseFields := map[string]interface{}{
		logging.KeyClientCapabilityDeferralAllowed: capabilities.DeferralAllowed,
	}

	logging.ProtocolTrace(ctx, "Announced client capabilities", responseFields)
}
