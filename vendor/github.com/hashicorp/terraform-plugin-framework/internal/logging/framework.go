// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logging

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tfsdklog"
)

const (
	// SubsystemFramework is the tfsdklog subsystem name for framework.
	SubsystemFramework = "framework"
)

// FrameworkDebug emits a framework subsystem log at DEBUG level.
func FrameworkDebug(ctx context.Context, msg string, additionalFields ...map[string]interface{}) {
	tfsdklog.SubsystemDebug(ctx, SubsystemFramework, msg, additionalFields...)
}

// FrameworkError emits a framework subsystem log at ERROR level.
func FrameworkError(ctx context.Context, msg string, additionalFields ...map[string]interface{}) {
	tfsdklog.SubsystemError(ctx, SubsystemFramework, msg, additionalFields...)
}

// FrameworkTrace emits a framework subsystem log at TRACE level.
func FrameworkTrace(ctx context.Context, msg string, additionalFields ...map[string]interface{}) {
	tfsdklog.SubsystemTrace(ctx, SubsystemFramework, msg, additionalFields...)
}

// FrameworkWarn emits a framework subsystem log at WARN level.
func FrameworkWarn(ctx context.Context, msg string, additionalFields ...map[string]interface{}) {
	tfsdklog.SubsystemWarn(ctx, SubsystemFramework, msg, additionalFields...)
}

// FrameworkWithAttributePath returns a new Context with KeyAttributePath set.
// The attribute path is expected to be string, so the logging package does not
// need to import path handling code.
func FrameworkWithAttributePath(ctx context.Context, attributePath string) context.Context {
	ctx = tfsdklog.SubsystemSetField(ctx, SubsystemFramework, KeyAttributePath, attributePath)
	return ctx
}
