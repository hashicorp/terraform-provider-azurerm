// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
)

func ProgressInvokeActionEventType(ctx context.Context, event fwserver.InvokeProgressEvent) tfprotov5.InvokeActionEvent {
	return tfprotov5.InvokeActionEvent{
		Type: tfprotov5.ProgressInvokeActionEventType{
			Message: event.Message,
		},
	}
}

func CompletedInvokeActionEventType(ctx context.Context, event *fwserver.InvokeActionResponse) tfprotov5.InvokeActionEvent {
	return tfprotov5.InvokeActionEvent{
		Type: tfprotov5.CompletedInvokeActionEventType{
			Diagnostics: Diagnostics(ctx, event.Diagnostics),
		},
	}
}
