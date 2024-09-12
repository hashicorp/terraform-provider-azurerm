// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6serverlogging

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/internal/logging"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ConfigureProviderClientCapabilities generates a TRACE "Announced client capabilities" log.
func ConfigureProviderClientCapabilities(ctx context.Context, capabilities *tfprotov6.ConfigureProviderClientCapabilities) {
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
func ReadDataSourceClientCapabilities(ctx context.Context, capabilities *tfprotov6.ReadDataSourceClientCapabilities) {
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
func ReadResourceClientCapabilities(ctx context.Context, capabilities *tfprotov6.ReadResourceClientCapabilities) {
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
func PlanResourceChangeClientCapabilities(ctx context.Context, capabilities *tfprotov6.PlanResourceChangeClientCapabilities) {
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
func ImportResourceStateClientCapabilities(ctx context.Context, capabilities *tfprotov6.ImportResourceStateClientCapabilities) {
	if capabilities == nil {
		logging.ProtocolTrace(ctx, "No announced client capabilities", map[string]interface{}{})
		return
	}

	responseFields := map[string]interface{}{
		logging.KeyClientCapabilityDeferralAllowed: capabilities.DeferralAllowed,
	}

	logging.ProtocolTrace(ctx, "Announced client capabilities", responseFields)
}
