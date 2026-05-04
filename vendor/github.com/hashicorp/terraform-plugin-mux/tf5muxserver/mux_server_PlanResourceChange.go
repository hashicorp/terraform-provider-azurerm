// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// PlanResourceChange calls the PlanResourceChange method, passing `req`, on
// the provider that returned the resource specified by req.TypeName in its
// schema.
func (s *muxServer) PlanResourceChange(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	rpc := "PlanResourceChange"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getResourceServer(ctx, req.TypeName)

	if err != nil {
		return nil, err
	}

	if diagnosticsHasError(diags) {
		return &tfprotov5.PlanResourceChangeResponse{
			Diagnostics: diags,
		}, nil
	}

	ctx = logging.Tfprotov5ProviderServerContext(ctx, server)

	// Prevent ServerCapabilities.PlanDestroy from sending destroy plans to
	// servers which do not enable the capability.
	if !serverSupportsPlanDestroy(s.resourceCapabilities[req.TypeName]) {
		if req.ProposedNewState == nil {
			logging.MuxTrace(ctx, "server does not enable destroy plans, returning without calling downstream server")

			resp := &tfprotov5.PlanResourceChangeResponse{
				// Presumably, we must preserve any prior private state so it
				// is still available during ApplyResourceChange.
				PlannedPrivate: req.PriorPrivate,
			}

			return resp, nil
		}

		isDestroyPlan, err := req.ProposedNewState.IsNull()

		if err != nil {
			return nil, fmt.Errorf("unable to determine if request is destroy plan: %w", err)
		}

		if isDestroyPlan {
			logging.MuxTrace(ctx, "server does not enable destroy plans, returning without calling downstream server")

			resp := &tfprotov5.PlanResourceChangeResponse{
				// Presumably, we must preserve any prior private state so it
				// is still available during ApplyResourceChange.
				PlannedPrivate: req.PriorPrivate,
				PlannedState:   req.ProposedNewState,
			}

			return resp, nil
		}
	}

	logging.MuxTrace(ctx, "calling downstream server")

	return server.PlanResourceChange(ctx, req)
}
