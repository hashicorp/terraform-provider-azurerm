// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package proto5server

import (
	"context"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

var _ tfprotov5.ProviderServer = &Server{}

// Provider server implementation.
type Server struct {
	FrameworkServer fwserver.Server

	contextCancels   []context.CancelFunc
	contextCancelsMu sync.Mutex
}

func (s *Server) registerContext(in context.Context) context.Context {
	ctx, cancel := context.WithCancel(in)
	s.contextCancelsMu.Lock()
	defer s.contextCancelsMu.Unlock()
	s.contextCancels = append(s.contextCancels, cancel)
	return ctx
}

func (s *Server) cancelRegisteredContexts(_ context.Context) {
	s.contextCancelsMu.Lock()
	defer s.contextCancelsMu.Unlock()
	for _, cancel := range s.contextCancels {
		cancel()
	}
	s.contextCancels = nil
}

// StopProvider satisfies the tfprotov5.ProviderServer interface.
func (s *Server) StopProvider(ctx context.Context, _ *tfprotov5.StopProviderRequest) (*tfprotov5.StopProviderResponse, error) {
	s.cancelRegisteredContexts(ctx)

	return &tfprotov5.StopProviderResponse{}, nil
}
