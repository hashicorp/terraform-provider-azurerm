// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/agentpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

func TestKubernetesClusterNodePoolWindowsTagFixture(t *testing.T) {
	data := acceptance.TestData{
		RandomInteger: 12345,
		Locations:     acceptance.Regions{Primary: "eastus"},
	}
	resource := KubernetesClusterNodePoolResource{}
	for _, tagValue := range []string{"dev", "prod"} {
		t.Run(tagValue, func(t *testing.T) {
			config := resource.windowsNodePoolWithTags(data, tagValue)
			if !strings.Contains(config, resource.templateWindowsConfig(data)) {
				t.Fatal("Windows tag fixture must retain the Windows cluster prerequisites")
			}
			_, nodePool, found := strings.Cut(config, `resource "azurerm_kubernetes_cluster_node_pool" "test" {`)
			if !found {
				t.Fatal("Windows tag fixture must contain the target node pool")
			}
			if strings.Contains(nodePool, "windows_profile") {
				t.Fatal("Windows tag fixture must omit the node-pool Windows profile from creation")
			}
			if !regexp.MustCompile(`os_type\s*=\s*"Windows"`).MatchString(nodePool) || !regexp.MustCompile(`node_count\s*=\s*1\b`).MatchString(nodePool) {
				t.Fatal("Windows tag fixture must run one Windows node")
			}
			if !regexp.MustCompile(`os_sku\s*=\s*"Windows2022"`).MatchString(nodePool) {
				t.Fatal("Windows tag fixture must select the tested Windows2022 image")
			}
			if !strings.Contains(nodePool, fmt.Sprintf("Environment = %q", tagValue)) {
				t.Fatal("Windows tag fixture must render the requested tag value")
			}
		})
	}
}

func TestKubernetesClusterNodePoolWindowsProfileState(t *testing.T) {
	for _, test := range []struct {
		name           string
		profile        *agentpools.AgentPoolWindowsProfile
		count          string
		enabled        string
		requireProfile bool
		expectError    bool
	}{
		{name: "initial_nil_profile", count: "0"},
		{name: "initial_nil_flag", profile: &agentpools.AgentPoolWindowsProfile{}, count: "0"},
		{name: "updated_nil_profile", count: "0", requireProfile: true, expectError: true},
		{name: "updated_nil_flag", profile: &agentpools.AgentPoolWindowsProfile{}, count: "0", requireProfile: true, expectError: true},
		{name: "enabled", profile: &agentpools.AgentPoolWindowsProfile{DisableOutboundNat: pointer.To(false)}, count: "1", enabled: "true", requireProfile: true},
		{name: "disabled", profile: &agentpools.AgentPoolWindowsProfile{DisableOutboundNat: pointer.To(true)}, count: "1", enabled: "false", requireProfile: true},
		{name: "missing_state", profile: &agentpools.AgentPoolWindowsProfile{DisableOutboundNat: pointer.To(false)}, requireProfile: true, expectError: true},
		{name: "mismatched_state", profile: &agentpools.AgentPoolWindowsProfile{DisableOutboundNat: pointer.To(false)}, count: "1", enabled: "false", requireProfile: true, expectError: true},
		{name: "stale_state", count: "1", enabled: "true", expectError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := checkWindowsProfileState(test.profile, map[string]string{
				"windows_profile.#":                      test.count,
				"windows_profile.0.outbound_nat_enabled": test.enabled,
			}, test.requireProfile)
			if (err != nil) != test.expectError {
				t.Fatalf("expected error=%t, got %v", test.expectError, err)
			}
		})
	}
}
