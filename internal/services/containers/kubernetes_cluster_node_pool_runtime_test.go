// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestKubernetesClusterNodePoolWorkloadRuntimePlan(t *testing.T) {
	schemas := map[string]*pluginsdk.Schema{
		"additional": resourceKubernetesClusterNodePoolSchema()["workload_runtime"],
		"default":    SchemaDefaultNodePool().Elem.(*pluginsdk.Resource).Schema["workload_runtime"],
	}

	tests := []struct {
		name       string
		state      string
		config     map[string]any
		unknown    bool
		wantChange bool
	}{
		{name: "omitted OCI", state: "OCIContainer", config: map[string]any{}},
		{name: "omitted Kata", state: "KataVmIsolation", config: map[string]any{}},
		{name: "null OCI", state: "OCIContainer", config: map[string]any{"workload_runtime": nil}},
		{name: "null Kata", state: "KataVmIsolation", config: map[string]any{"workload_runtime": nil}},
		{name: "explicit OCI", state: "OCIContainer", config: map[string]any{"workload_runtime": "OCIContainer"}},
		{name: "explicit Kata", state: "KataVmIsolation", config: map[string]any{"workload_runtime": "KataVmIsolation"}},
		{name: "explicit change", state: "OCIContainer", config: map[string]any{"workload_runtime": "KataVmIsolation"}, wantChange: true},
		{name: "unknown", state: "OCIContainer", unknown: true, wantChange: true},
	}

	for pool, field := range schemas {
		t.Run(pool, func(t *testing.T) {
			resource := &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{"workload_runtime": field}}
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					state := &terraform.InstanceState{
						ID:         "existing-pool",
						Attributes: map[string]string{"workload_runtime": test.state},
					}
					config := terraform.NewResourceConfigRaw(test.config)
					if test.unknown {
						config = terraform.NewResourceConfigShimmed(cty.ObjectVal(map[string]cty.Value{
							"workload_runtime": cty.UnknownVal(cty.String),
						}), resource.CoreConfigSchema())
					}
					diff, err := resource.SimpleDiff(context.Background(), state, config, nil)
					if err != nil {
						t.Fatal(err)
					}
					if got := !diff.Empty(); got != test.wantChange {
						t.Fatalf("change = %t, want %t: %#v", got, test.wantChange, diff.Attributes)
					}
					if diff.RequiresNew() {
						t.Fatal("workload runtime planning unexpectedly requires replacement")
					}
					if test.wantChange {
						runtimeDiff := diff.Attributes["workload_runtime"]
						if runtimeDiff == nil {
							t.Fatal("expected a workload_runtime diff")
						}
						if runtimeDiff.NewComputed != test.unknown {
							t.Fatalf("workload runtime unknown = %t, want %t", runtimeDiff.NewComputed, test.unknown)
						}
						if runtimeDiff.NewRemoved {
							t.Fatal("workload runtime planning unexpectedly removes the value")
						}
						if !test.unknown && runtimeDiff.New != test.config["workload_runtime"] {
							t.Fatalf("workload runtime = %q, want %q", runtimeDiff.New, test.config["workload_runtime"])
						}
					}
				})
			}
			for _, invalid := range []string{"", "invalid", "ocicontainer"} {
				if diagnostics := resource.Validate(terraform.NewResourceConfigRaw(map[string]any{"workload_runtime": invalid})); !diagnostics.HasError() {
					t.Errorf("expected validation to reject %q", invalid)
				}
			}
		})
	}
}
