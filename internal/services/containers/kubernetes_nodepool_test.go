// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-cty/cty/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Note: This unit test verifies that optional integer sysctl settings whose valid range includes 0
// (specifically `vm_swappiness` and `vm_vfs_cache_pressure`) correctly differentiate between:
//  1. Explicitly configured as `0` in HCL: sent as `pointer.To(0)` to the Azure API.
//  2. Omitted in HCL: sent as `nil` to preserve Azure / OS defaults, rather than sending `0`.
//  3. Configured with a non-zero value: sent as `pointer.To(v)`.
//
// Why this test setup is unusual:
// In `terraform-plugin-sdk`, standard test helpers like `schema.TestResourceDataRaw` do not populate
// `RawConfig` on `ResourceData` (it evaluates to `cty.NullVal`). Because the expander relies on
// `d.GetRawConfigAt(...)` to inspect whether the attribute was explicitly defined in configuration
// vs omitted (and `TypeInt` within a block defaults to `0` in Go), we must simulate the Terraform's apply
// lifecycle. By unmarshaling the config into a `cty.Value`, generating an `InstanceDiff`
// via `schema.DiffFromValues`, and invoking `res.Apply`, `d.GetRawConfigAt(...)` is populated to simulate
// a real `terraform apply`.
func TestExpandSysctlConfig_ZeroValues(t *testing.T) {
	s := map[string]*schema.Schema{
		"default_node_pool": SchemaDefaultNodePool(),
	}

	testExpand := func(t *testing.T, jsonConfig []byte) *managedclusters.ManagedClusterAgentPoolProfile {
		var capturedResult *[]managedclusters.ManagedClusterAgentPoolProfile
		var capturedErr error

		res := &schema.Resource{
			Schema: s,
			CreateContext: func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
				capturedResult, capturedErr = ExpandDefaultNodePool(d)
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
		if capturedResult == nil || len(*capturedResult) == 0 {
			t.Fatalf("expected non-empty result")
		}
		return &(*capturedResult)[0]
	}

	t.Run("explicit zero values", func(t *testing.T) {
		jsonConfig := []byte(`{
			"default_node_pool": [{
				"name": "default",
				"node_count": 1,
				"vm_size": "Standard_DS2_v2",
				"linux_os_config": [{
					"sysctl_config": [{
						"vm_swappiness": 0,
						"vm_vfs_cache_pressure": 0
					}]
				}]
			}]
		}`)

		pool := testExpand(t, jsonConfig)
		if pool.LinuxOSConfig == nil || pool.LinuxOSConfig.Sysctls == nil {
			t.Fatalf("expected LinuxOSConfig and Sysctls to be populated")
		}
		sysctls := pool.LinuxOSConfig.Sysctls
		if sysctls.VMSwappiness == nil {
			t.Errorf("expected VMSwappiness to be non-nil when set to 0")
		} else if *sysctls.VMSwappiness != 0 {
			t.Errorf("expected VMSwappiness to be 0, got %d", *sysctls.VMSwappiness)
		}
		if sysctls.VMVfsCachePressure == nil {
			t.Errorf("expected VMVfsCachePressure to be non-nil when set to 0")
		} else if *sysctls.VMVfsCachePressure != 0 {
			t.Errorf("expected VMVfsCachePressure to be 0, got %d", *sysctls.VMVfsCachePressure)
		}
	})

	t.Run("omitted values", func(t *testing.T) {
		jsonConfig := []byte(`{
			"default_node_pool": [{
				"name": "default",
				"node_count": 1,
				"vm_size": "Standard_DS2_v2",
				"linux_os_config": [{
					"sysctl_config": [{
						"fs_file_max": 100000
					}]
				}]
			}]
		}`)

		pool := testExpand(t, jsonConfig)
		if pool.LinuxOSConfig == nil || pool.LinuxOSConfig.Sysctls == nil {
			t.Fatalf("expected LinuxOSConfig and Sysctls to be populated")
		}
		sysctls := pool.LinuxOSConfig.Sysctls
		if sysctls.VMSwappiness != nil {
			t.Errorf("expected VMSwappiness to be nil when omitted, got %d", *sysctls.VMSwappiness)
		}
		if sysctls.VMVfsCachePressure != nil {
			t.Errorf("expected VMVfsCachePressure to be nil when omitted, got %d", *sysctls.VMVfsCachePressure)
		}
	})

	t.Run("explicit non-zero values", func(t *testing.T) {
		jsonConfig := []byte(`{
			"default_node_pool": [{
				"name": "default",
				"node_count": 1,
				"vm_size": "Standard_DS2_v2",
				"linux_os_config": [{
					"sysctl_config": [{
						"vm_swappiness": 60,
						"vm_vfs_cache_pressure": 80
					}]
				}]
			}]
		}`)

		pool := testExpand(t, jsonConfig)
		if pool.LinuxOSConfig == nil || pool.LinuxOSConfig.Sysctls == nil {
			t.Fatalf("expected LinuxOSConfig and Sysctls to be populated")
		}
		sysctls := pool.LinuxOSConfig.Sysctls
		if sysctls.VMSwappiness == nil || *sysctls.VMSwappiness != 60 {
			t.Errorf("expected VMSwappiness to be 60, got %v", sysctls.VMSwappiness)
		}
		if sysctls.VMVfsCachePressure == nil || *sysctls.VMVfsCachePressure != 80 {
			t.Errorf("expected VMVfsCachePressure to be 80, got %v", sysctls.VMVfsCachePressure)
		}
	})
}

// Note: Verifies the same zero-value vs omitted behavior for standalone node pools
// (`azurerm_kubernetes_cluster_node_pool`) via `expandAgentPoolLinuxOSConfig`.
func TestExpandAgentPoolSysctlConfig_ZeroValues(t *testing.T) {
	s := map[string]*schema.Schema{
		"linux_os_config": schemaNodePoolLinuxOSConfig(),
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
