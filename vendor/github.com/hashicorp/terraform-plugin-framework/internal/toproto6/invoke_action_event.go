// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
)

func ProgressInvokeActionEventType(ctx context.Context, event fwserver.InvokeProgressEvent) tfprotov6.InvokeActionEvent {
	return tfprotov6.InvokeActionEvent{
		Type: tfprotov6.ProgressInvokeActionEventType{
			Message: event.Message,
		},
	}
}

func CompletedInvokeActionEventType(ctx context.Context, event *fwserver.InvokeActionResponse) tfprotov6.InvokeActionEvent {
	return tfprotov6.InvokeActionEvent{
		Type: tfprotov6.CompletedInvokeActionEventType{
			Diagnostics: Diagnostics(ctx, event.Diagnostics),
		},
	}
}
