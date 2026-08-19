// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func TestEffectiveReadTimeout(t *testing.T) {
	declared := 5 * time.Minute
	cases := []struct {
		name      string
		rawConfig cty.Value
		expected  time.Duration
	}{
		{
			name:      "null config",
			rawConfig: cty.NullVal(cty.Object(map[string]cty.Type{})),
			expected:  declared,
		},
		{
			name:      "no timeouts attribute",
			rawConfig: cty.ObjectVal(map[string]cty.Value{"name": cty.StringVal("example")}),
			expected:  declared,
		},
		{
			name: "null timeouts block",
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"timeouts": cty.NullVal(cty.Object(map[string]cty.Type{"read": cty.String})),
			}),
			expected: declared,
		},
		{
			name: "timeouts block without read",
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"timeouts": cty.ObjectVal(map[string]cty.Value{"create": cty.StringVal("1m")}),
			}),
			expected: declared,
		},
		{
			name: "null read value",
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"timeouts": cty.ObjectVal(map[string]cty.Value{"read": cty.NullVal(cty.String)}),
			}),
			expected: declared,
		},
		{
			name: "invalid read value",
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"timeouts": cty.ObjectVal(map[string]cty.Value{"read": cty.StringVal("not-a-duration")}),
			}),
			expected: declared,
		},
		{
			name: "read value set",
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"timeouts": cty.ObjectVal(map[string]cty.Value{"read": cty.StringVal("1m")}),
			}),
			expected: time.Minute,
		},
		{
			name: "read value above the sdk default",
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"timeouts": cty.ObjectVal(map[string]cty.Value{"read": cty.StringVal("35m")}),
			}),
			expected: 35 * time.Minute,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if actual := effectiveReadTimeout(tc.rawConfig, declared); actual != tc.expected {
				t.Fatalf("expected %s but got %s", tc.expected, actual)
			}
		})
	}
}

func TestWrapDataSourceReadTimeoutsLegacyRead(t *testing.T) {
	declared := 5 * time.Minute
	var observed time.Duration

	r := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{Read: &declared},
		Read: func(d *schema.ResourceData, meta interface{}) error {
			// legacy read functions derive their deadline via timeouts.ForRead - this should
			// now honour the declared read timeout rather than the SDK default of 20 minutes
			ctx, cancel := timeouts.ForRead(context.Background(), d)
			defer cancel()
			deadline, ok := ctx.Deadline()
			if !ok {
				t.Fatal("expected a deadline to be set on the context but none was")
			}
			observed = time.Until(deadline)
			return nil
		},
	}
	wrapDataSourceReadTimeouts(r)

	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
	if err := r.Read(d, nil); err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if observed > declared || observed < declared-time.Minute {
		t.Fatalf("expected a deadline of ~%s from now but got %s", declared, observed)
	}
}

func TestWrapDataSourceReadTimeoutsReadContext(t *testing.T) {
	declared := 5 * time.Minute
	var observed time.Duration

	r := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{Read: &declared},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			deadline, ok := ctx.Deadline()
			if !ok {
				t.Fatal("expected a deadline to be set on the context but none was")
			}
			observed = time.Until(deadline)
			return nil
		},
	}
	wrapDataSourceReadTimeouts(r)

	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
	if diags := r.ReadContext(context.Background(), d, nil); diags.HasError() {
		t.Fatalf("unexpected error: %+v", diags)
	}
	if observed > declared || observed < declared-time.Minute {
		t.Fatalf("expected a deadline of ~%s from now but got %s", declared, observed)
	}
}

func TestWrapDataSourceReadTimeoutsNoDeclaredTimeout(t *testing.T) {
	r := &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error { return nil },
	}
	wrapDataSourceReadTimeouts(r)

	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})

	// no declared timeout means the resource is left untouched - ForRead falls back to the SDK default
	ctx, cancel := timeouts.ForRead(context.Background(), d)
	defer cancel()
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected a deadline to be set on the context but none was")
	}
	if remaining := time.Until(deadline); remaining > 20*time.Minute || remaining < 19*time.Minute {
		t.Fatalf("expected a deadline of ~20m from now but got %s", remaining)
	}
}
