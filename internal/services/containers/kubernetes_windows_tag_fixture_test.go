// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/agentpools"
)

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
