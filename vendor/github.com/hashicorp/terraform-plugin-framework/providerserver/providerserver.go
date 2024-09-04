// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package providerserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/proto5server"
	"github.com/hashicorp/terraform-plugin-framework/internal/proto6server"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
)

// NewProtocol5 returns a protocol version 5 ProviderServer implementation
// based on the given Provider and suitable for usage with the
// github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server.Serve()
// function and various terraform-plugin-mux functions.
func NewProtocol5(p provider.Provider) func() tfprotov5.ProviderServer {
	return func() tfprotov5.ProviderServer {
		return &proto5server.Server{
			FrameworkServer: fwserver.Server{
				Provider: p,
			},
		}
	}
}

// NewProtocol5WithError returns a protocol version 5 ProviderServer
// implementation based on the given Provider and suitable for usage with
// github.com/hashicorp/terraform-plugin-testing/helper/resource.TestCase.ProtoV5ProviderFactories.
//
// The error return is not currently used, but it may be in the future.
func NewProtocol5WithError(p provider.Provider) func() (tfprotov5.ProviderServer, error) {
	return func() (tfprotov5.ProviderServer, error) {
		return &proto5server.Server{
			FrameworkServer: fwserver.Server{
				Provider: p,
			},
		}, nil
	}
}

// NewProtocol6 returns a protocol version 6 ProviderServer implementation
// based on the given Provider and suitable for usage with the
// github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server.Serve()
// function and various terraform-plugin-mux functions.
func NewProtocol6(p provider.Provider) func() tfprotov6.ProviderServer {
	return func() tfprotov6.ProviderServer {
		return &proto6server.Server{
			FrameworkServer: fwserver.Server{
				Provider: p,
			},
		}
	}
}

// NewProtocol6WithError returns a protocol version 6 ProviderServer
// implementation based on the given Provider and suitable for usage with
// github.com/hashicorp/terraform-plugin-testing/helper/resource.TestCase.ProtoV6ProviderFactories.
//
// The error return is not currently used, but it may be in the future.
func NewProtocol6WithError(p provider.Provider) func() (tfprotov6.ProviderServer, error) {
	return func() (tfprotov6.ProviderServer, error) {
		return &proto6server.Server{
			FrameworkServer: fwserver.Server{
				Provider: p,
			},
		}, nil
	}
}

// Serve serves a provider, blocking until the context is canceled.
func Serve(ctx context.Context, providerFunc func() provider.Provider, opts ServeOpts) error {
	err := opts.validate(ctx)

	if err != nil {
		return fmt.Errorf("unable to validate ServeOpts: %w", err)
	}

	switch opts.ProtocolVersion {
	case 5:
		var tf5serverOpts []tf5server.ServeOpt

		if opts.Debug {
			tf5serverOpts = append(tf5serverOpts, tf5server.WithManagedDebug())
		}

		return tf5server.Serve(
			opts.Address,
			func() tfprotov5.ProviderServer {
				provider := providerFunc()

				return &proto5server.Server{
					FrameworkServer: fwserver.Server{
						Provider: provider,
					},
				}
			},
			tf5serverOpts...,
		)
	default:
		var tf6serverOpts []tf6server.ServeOpt

		if opts.Debug {
			tf6serverOpts = append(tf6serverOpts, tf6server.WithManagedDebug())
		}

		return tf6server.Serve(
			opts.Address,
			func() tfprotov6.ProviderServer {
				provider := providerFunc()

				return &proto6server.Server{
					FrameworkServer: fwserver.Server{
						Provider: provider,
					},
				}
			},
			tf6serverOpts...,
		)
	}
}
