// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logging

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tfsdklog"
)

const (
	// SubsystemMux is the tfsdklog subsystem name for mux logging.
	SubsystemMux = "mux"
)

// MuxTrace emits a mux subsystem log at TRACE level.
func MuxTrace(ctx context.Context, msg string, additionalFields ...map[string]interface{}) {
	tfsdklog.SubsystemTrace(ctx, SubsystemMux, msg, additionalFields...)
}
