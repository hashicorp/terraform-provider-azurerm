// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6serverlogging

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/internal/logging"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// Deferred generates a TRACE "Received downstream deferred response" log if populated.
func Deferred(ctx context.Context, deferred *tfprotov6.Deferred) {
	if deferred == nil {
		return
	}

	responseFields := map[string]interface{}{
		logging.KeyDeferredReason: deferred.Reason.String(),
	}

	logging.ProtocolTrace(ctx, "Received downstream deferred response", responseFields)
}
