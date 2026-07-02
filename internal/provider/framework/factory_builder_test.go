// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package framework

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// TestProtoV5ProviderServerFactoryMuxSchemasMatch guards the manually-maintained
// parity between the Plugin SDKv2 provider schema (internal/provider) and the
// Plugin Framework provider schema (this package). tf5muxserver requires the
// provider block schema returned by both muxed servers to be identical; if they
// drift (e.g. a new provider argument added to only one half), GetProviderSchema
// surfaces an error diagnostic.
func TestProtoV5ProviderServerFactoryMuxSchemasMatch(t *testing.T) {
	ctx := context.Background()

	factory, _, err := ProtoV5ProviderServerFactory(ctx)
	if err != nil {
		t.Fatalf("constructing mux provider server factory: %+v", err)
	}

	resp, err := factory().GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("calling GetProviderSchema on the mux server: %+v", err)
	}

	for _, d := range resp.Diagnostics {
		if d.Severity == tfprotov5.DiagnosticSeverityError {
			t.Fatalf("mux provider schema mismatch: %s: %s", d.Summary, d.Detail)
		}
	}
}
