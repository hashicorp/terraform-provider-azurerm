// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5server

import (
	"context"
	"errors"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/tfplugin5"
	"google.golang.org/grpc"
)

// GRPCProviderPlugin is an implementation of the
// github.com/hashicorp/go-plugin#Plugin and
// github.com/hashicorp/go-plugin#GRPCPlugin interfaces, indicating how to
// serve tfprotov5.ProviderServers as gRPC plugins for go-plugin.
type GRPCProviderPlugin struct {
	GRPCProvider func() tfprotov5.ProviderServer
	Opts         []ServeOpt
	Name         string
}

// Server always returns an error; we're only implementing the GRPCPlugin
// interface, not the Plugin interface.
func (p *GRPCProviderPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return nil, errors.New("terraform-plugin-go only implements gRPC servers")
}

// Client always returns an error; we're only implementing the GRPCPlugin
// interface, not the Plugin interface.
func (p *GRPCProviderPlugin) Client(*plugin.MuxBroker, *rpc.Client) (interface{}, error) {
	return nil, errors.New("terraform-plugin-go only implements gRPC servers")
}

// GRPCClient always returns an error; we're only implementing the server half
// of the interface.
func (p *GRPCProviderPlugin) GRPCClient(context.Context, *plugin.GRPCBroker, *grpc.ClientConn) (interface{}, error) {
	return nil, errors.New("terraform-plugin-go only implements gRPC servers")
}

// GRPCServer registers the gRPC provider server with the gRPC server that
// go-plugin is standing up.
func (p *GRPCProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	tfplugin5.RegisterProviderServer(s, New(p.Name, p.GRPCProvider(), p.Opts...))
	return nil
}
