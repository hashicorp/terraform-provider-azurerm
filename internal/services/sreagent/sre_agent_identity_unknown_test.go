// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent

import (
	"context"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestSreAgentRootIdentityUnknownReferences(t *testing.T) {
	wrapped := sdk.WrappedResource(SreAgentResource{})
	customize := wrapped.CustomizeDiff
	for _, tc := range []struct {
		name     string
		kind     cty.Value
		ids      cty.Value
		wantFail bool
	}{
		{"single-unknown-id", cty.StringVal("SystemAssigned, UserAssigned"), cty.SetVal([]cty.Value{cty.UnknownVal(cty.String)}), false},
		{"mixed-known-unknown-ids", cty.StringVal("UserAssigned"), cty.SetVal([]cty.Value{cty.StringVal(testIdentityID), cty.UnknownVal(cty.String)}), false},
		{"unknown-set", cty.StringVal("UserAssigned"), cty.UnknownVal(cty.Set(cty.String)), false},
		{"unknown-type", cty.UnknownVal(cty.String), cty.SetVal([]cty.Value{cty.StringVal(testIdentityID)}), false},
		{"known-empty-invalid", cty.StringVal("UserAssigned"), cty.SetValEmpty(cty.String), true},
		{"system-only-valid", cty.StringVal("SystemAssigned"), cty.SetValEmpty(cty.String), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			wrapped.CustomizeDiff = func(ctx context.Context, diff *pluginsdk.ResourceDiff, meta interface{}) error {
				t.Logf("known block=%t type=%t ids=%t rawNull=%t", diff.NewValueKnown("identity"),
					diff.NewValueKnown("identity.0.type"), diff.NewValueKnown("identity.0.identity_ids"), diff.GetRawConfig().IsNull())
				return customize(ctx, diff, meta)
			}
			values := protocolValues(wrapped)
			values["identity"] = cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"type": tc.kind, "identity_ids": tc.ids,
				"principal_id": cty.NullVal(cty.String), "tenant_id": cty.NullVal(cty.String),
			})})
			_, err := wrapped.Diff(context.Background(), nil, terraform.NewResourceConfigShimmed(cty.ObjectVal(values), wrapped.CoreConfigSchema()), &clients.Client{})
			if (err != nil) != tc.wantFail {
				t.Fatalf("wantFail=%t, plan error=%v", tc.wantFail, err)
			}
		})
	}
}
