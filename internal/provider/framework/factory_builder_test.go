// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package framework_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

// TestProtoV5ProviderServerFactory_SchemaEquivalence checks the lazy validation in MUX for schema equivalence and
// calls out any mis-matched items.
func TestProtoV5ProviderServerFactory_SchemaEquivalence(t *testing.T) {
	t.Parallel()

	factory, _, err := framework.ProtoV5ProviderServerFactory(context.Background())
	if err != nil {
		t.Fatalf("Failed to initialize Mux Server: %s", err)
	}

	server := factory()
	req := &tfprotov5.GetProviderSchemaRequest{}
	resp, err := server.GetProviderSchema(context.Background(), req)
	if err != nil {
		t.Fatalf("GetProviderSchema returned an error: %s", err)
	}

	for _, diag := range resp.Diagnostics {
		if diag.Severity == tfprotov5.DiagnosticSeverityError {
			t.Errorf("MUX'd schema validation error: %s\nDetail: %s", diag.Summary, diag.Detail)
		}
	}
}
