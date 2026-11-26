// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tf5serverlogging

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/internal/logging"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ServerCapabilities generates a TRACE "Announced server capabilities" log.
func ServerCapabilities(ctx context.Context, capabilities *tfprotov5.ServerCapabilities) {
	responseFields := map[string]interface{}{
		logging.KeyServerCapabilityGetProviderSchemaOptional: false,
		logging.KeyServerCapabilityMoveResourceState:         false,
		logging.KeyServerCapabilityPlanDestroy:               false,
		logging.KeyServerCapabilityGenerateResourceConfig:    false,
	}

	if capabilities != nil {
		responseFields[logging.KeyServerCapabilityGetProviderSchemaOptional] = capabilities.GetProviderSchemaOptional
		responseFields[logging.KeyServerCapabilityMoveResourceState] = capabilities.MoveResourceState
		responseFields[logging.KeyServerCapabilityPlanDestroy] = capabilities.PlanDestroy
		responseFields[logging.KeyServerCapabilityGenerateResourceConfig] = capabilities.GenerateResourceConfig
	}

	logging.ProtocolTrace(ctx, "Announced server capabilities", responseFields)
}
