// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

func TestContainerGroupDataSourceProbeFixture(t *testing.T) {
	data := acceptance.TestData{
		RandomInteger: 12345,
		Locations:     acceptance.Regions{Primary: "eastus"},
	}
	config := ContainerGroupDataSource{}.complete(data)
	file, diagnostics := hclsyntax.ParseConfig([]byte(config), "test.tf", hcl.InitialPos)
	if diagnostics.HasErrors() {
		t.Fatal(diagnostics)
	}
	groups, sources := 0, 0
	probes := map[string]int{}
	for _, block := range file.Body.(*hclsyntax.Body).Blocks {
		if len(block.Labels) != 2 || block.Labels[0] != "azurerm_container_group" {
			continue
		}
		if block.Type == "data" {
			sources++
			continue
		}
		if block.Type != "resource" {
			continue
		}
		groups++
		osType, diagnostics := block.Body.Attributes["os_type"].Expr.Value(nil)
		if diagnostics.HasErrors() || osType.AsString() != "Linux" {
			t.Fatalf("expected the Linux complete fixture: %v", diagnostics)
		}
		for _, container := range block.Body.Blocks {
			if container.Type != "container" {
				continue
			}
			commands, diagnostics := container.Body.Attributes["commands"].Expr.Value(nil)
			if diagnostics.HasErrors() {
				t.Fatal(diagnostics)
			}
			var actual []string
			for _, command := range commands.AsValueSlice() {
				actual = append(actual, command.AsString())
			}
			expected := []string{"/bin/sh", "-c", "(sleep 5; touch /tmp/healthy) & exec node /usr/src/app/index.js"}
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("commands = %#v, want %#v", actual, expected)
			}
			for _, probe := range container.Body.Blocks {
				if probe.Type != "readiness_probe" && probe.Type != "liveness_probe" {
					continue
				}
				probes[probe.Type]++
				for key, expected := range map[string]int64{
					"initial_delay_seconds": 10,
					"period_seconds":        5,
					"failure_threshold":     3,
					"timeout_seconds":       5,
				} {
					value, diagnostics := probe.Body.Attributes[key].Expr.Value(nil)
					if diagnostics.HasErrors() {
						t.Fatal(diagnostics)
					}
					number := value.AsBigFloat()
					if actual, accuracy := number.Int64(); accuracy != big.Exact || actual != expected {
						t.Errorf("%s.%s = %s, want %d", probe.Type, key, number.Text('g', -1), expected)
					}
				}
			}
		}
	}
	if groups != 1 || sources != 1 || probes["readiness_probe"] != 1 || probes["liveness_probe"] != 1 {
		t.Fatalf("expected one resource, one data source and one of each probe, got %d, %d, %v", groups, sources, probes)
	}
}
