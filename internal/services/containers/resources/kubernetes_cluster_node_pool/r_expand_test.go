// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_node_pool

import (
	"context"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/agentpools"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-cty/cty/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/helpers"
)

// Note: Verifies the same zero-value vs omitted behavior for standalone node pools
// (`azurerm_kubernetes_cluster_node_pool`) via `expandAgentPoolLinuxOSConfig`.
func TestExpandAgentPoolSysctlConfig_ZeroValues(t *testing.T) {
	s := map[string]*schema.Schema{
		"linux_os_config": helpers.SchemaNodePoolLinuxOSConfig(),
	}

	testExpand := func(t *testing.T, jsonConfig []byte) *agentpools.LinuxOSConfig {
		var capturedResult *agentpools.LinuxOSConfig
		var capturedErr error

		res := &schema.Resource{
			Schema: s,
			CreateContext: func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
				if linuxOSConfig := d.Get("linux_os_config").([]any); len(linuxOSConfig) > 0 {
					capturedResult, capturedErr = expandAgentPoolLinuxOSConfig(d, "linux_os_config.0.", linuxOSConfig)
				}
				return nil
			},
		}

		cfgVal, err := json.Unmarshal(jsonConfig, res.CoreConfigSchema().ImpliedType())
		if err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		diff, err := schema.DiffFromValues(context.Background(), cty.NullVal(res.CoreConfigSchema().ImpliedType()), cfgVal, cfgVal, res)
		if err != nil {
			t.Fatalf("diff error: %v", err)
		}

		_, diags := res.Apply(context.Background(), nil, diff, nil)
		if diags.HasError() {
			t.Fatalf("apply diags error: %v", diags)
		}
		if capturedErr != nil {
			t.Fatalf("expand error: %v", capturedErr)
		}
		if capturedResult == nil {
			t.Fatalf("expected non-empty result")
		}
		return capturedResult
	}

	t.Run("explicit zero values", func(t *testing.T) {
		jsonConfig := []byte(`{
			"linux_os_config": [{
				"sysctl_config": [{
					"vm_swappiness": 0,
					"vm_vfs_cache_pressure": 0
				}]
			}]
		}`)

		cfg := testExpand(t, jsonConfig)
		if cfg.Sysctls == nil {
			t.Fatalf("expected Sysctls to be populated")
		}
		if cfg.Sysctls.VMSwappiness == nil {
			t.Errorf("expected VMSwappiness to be non-nil when set to 0")
		} else if *cfg.Sysctls.VMSwappiness != 0 {
			t.Errorf("expected VMSwappiness to be 0, got %d", *cfg.Sysctls.VMSwappiness)
		}
		if cfg.Sysctls.VMVfsCachePressure == nil {
			t.Errorf("expected VMVfsCachePressure to be non-nil when set to 0")
		} else if *cfg.Sysctls.VMVfsCachePressure != 0 {
			t.Errorf("expected VMVfsCachePressure to be 0, got %d", *cfg.Sysctls.VMVfsCachePressure)
		}
	})

	t.Run("omitted values", func(t *testing.T) {
		jsonConfig := []byte(`{
			"linux_os_config": [{
				"sysctl_config": [{
					"fs_file_max": 100000
				}]
			}]
		}`)

		cfg := testExpand(t, jsonConfig)
		if cfg.Sysctls == nil {
			t.Fatalf("expected Sysctls to be populated")
		}
		if cfg.Sysctls.VMSwappiness != nil {
			t.Errorf("expected VMSwappiness to be nil when omitted, got %d", *cfg.Sysctls.VMSwappiness)
		}
		if cfg.Sysctls.VMVfsCachePressure != nil {
			t.Errorf("expected VMVfsCachePressure to be nil when omitted, got %d", *cfg.Sysctls.VMVfsCachePressure)
		}
	})

	t.Run("explicit non-zero values", func(t *testing.T) {
		jsonConfig := []byte(`{
			"linux_os_config": [{
				"sysctl_config": [{
					"vm_swappiness": 60,
					"vm_vfs_cache_pressure": 80
				}]
			}]
		}`)

		cfg := testExpand(t, jsonConfig)
		if cfg.Sysctls == nil {
			t.Fatalf("expected Sysctls to be populated")
		}
		if cfg.Sysctls.VMSwappiness == nil || *cfg.Sysctls.VMSwappiness != 60 {
			t.Errorf("expected VMSwappiness to be 60, got %v", cfg.Sysctls.VMSwappiness)
		}
		if cfg.Sysctls.VMVfsCachePressure == nil || *cfg.Sysctls.VMVfsCachePressure != 80 {
			t.Errorf("expected VMVfsCachePressure to be 80, got %v", cfg.Sysctls.VMVfsCachePressure)
		}
	})
}
