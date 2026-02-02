// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package proto6server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fromproto6"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/toproto6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// invokeActionErrorDiagnostics returns a value suitable for
// [InvokeActionServerStream.Events]. It yields a single result that contains
// the given error diagnostics.
func invokeActionErrorDiagnostics(ctx context.Context, diags diag.Diagnostics) (*tfprotov6.InvokeActionServerStream, error) {
	return &tfprotov6.InvokeActionServerStream{
		Events: func(push func(tfprotov6.InvokeActionEvent) bool) {
			push(tfprotov6.InvokeActionEvent{
				Type: tfprotov6.CompletedInvokeActionEventType{
					Diagnostics: toproto6.Diagnostics(ctx, diags),
				},
			})
		},
	}, nil
}

// InvokeAction satisfies the tfprotov6.ProviderServer interface.
func (s *Server) InvokeAction(ctx context.Context, proto6Req *tfprotov6.InvokeActionRequest) (*tfprotov6.InvokeActionServerStream, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.InvokeActionResponse{}

	action, diags := s.FrameworkServer.Action(ctx, proto6Req.ActionType)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return invokeActionErrorDiagnostics(ctx, fwResp.Diagnostics)
	}

	actionSchema, diags := s.FrameworkServer.ActionSchema(ctx, proto6Req.ActionType)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return invokeActionErrorDiagnostics(ctx, fwResp.Diagnostics)
	}

	fwReq, diags := fromproto6.InvokeActionRequest(ctx, proto6Req, action, actionSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return invokeActionErrorDiagnostics(ctx, fwResp.Diagnostics)
	}

	protoStream := &tfprotov6.InvokeActionServerStream{
		Events: func(push func(tfprotov6.InvokeActionEvent) bool) {
			// Create a channel for framework to receive progress events
			progressChan := make(chan fwserver.InvokeProgressEvent)
			fwResp.ProgressEvents = progressChan

			// Create a channel to be triggered when the invoke action method has finished
			completedChan := make(chan any)
			go func() {
				s.FrameworkServer.InvokeAction(ctx, fwReq, fwResp)
				close(completedChan)
			}()

			for {
				select {
				// Actions can only push one completed event and it's automatically handled by the framework
				// by closing the completed channel above.
				case <-completedChan:
					push(toproto6.CompletedInvokeActionEventType(ctx, fwResp))
					return

				// Actions can push multiple progress events
				case progressEvent := <-fwResp.ProgressEvents:
					if !push(toproto6.ProgressInvokeActionEventType(ctx, progressEvent)) {
						return
					}
				}
			}
		},
	}

	return protoStream, nil
}
