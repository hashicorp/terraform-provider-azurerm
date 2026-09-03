// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestKubernetesAdminUserName(t *testing.T) {
	cases := []struct {
		AdminUserName string
		Errors        int
	}{
		{
			AdminUserName: "",
			Errors:        1,
		},
		{
			AdminUserName: "Abc-123_abc",
			Errors:        0,
		},
		{
			AdminUserName: "123abc",
			Errors:        1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.AdminUserName, func(t *testing.T) {
			_, errors := KubernetesAdminUserName(tc.AdminUserName, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected AdminUserName to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}

func TestKubernetesClusterName(t *testing.T) {
	cases := []struct {
		ClusterName string
		Errors      int
	}{
		{
			ClusterName: "",
			Errors:      1,
		},
		{
			ClusterName: "A",
			Errors:      0,
		},
		{
			ClusterName: "abc123",
			Errors:      0,
		},
		{
			ClusterName: "hello-world",
			Errors:      0,
		},
		{
			ClusterName: "hello_world",
			Errors:      0,
		},
		{
			ClusterName: "Hello-World",
			Errors:      0,
		},
		{
			ClusterName: "thisnamendsinahyphen-",
			Errors:      1,
		},
		{
			ClusterName: "thisnameis64charactersitisverylongandlongandlong1234567890123456",
			Errors:      1,
		},
		{
			ClusterName: "thisnameisonly63charactersbutitisstillverylongandlongandlong123",
			Errors:      0,
		},
		{
			ClusterName: "ABC!123",
			Errors:      1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.ClusterName, func(t *testing.T) {
			_, errors := KubernetesClusterName(tc.ClusterName, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected ClusterName to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}

func TestKubernetesAgentPoolName(t *testing.T) {
	cases := []struct {
		AgentPoolName string
		Errors        int
	}{
		{
			AgentPoolName: "",
			Errors:        1,
		},
		{
			AgentPoolName: "ABC123",
			Errors:        1,
		},
		{
			AgentPoolName: "abc123",
			Errors:        0,
		},
		{
			AgentPoolName: "123abc",
			Errors:        1,
		},
		{
			AgentPoolName: "hi",
			Errors:        0,
		},
		{
			AgentPoolName: "hello",
			Errors:        0,
		},
		{
			AgentPoolName: "hello-world",
			Errors:        1,
		},
		{
			AgentPoolName: "helloworld123",
			Errors:        1,
		},
		{
			AgentPoolName: "hello_world",
			Errors:        1,
		},
		{
			AgentPoolName: "Hello-World",
			Errors:        1,
		},
		{
			AgentPoolName: "20202020",
			Errors:        1,
		},
		{
			AgentPoolName: "h20202020",
			Errors:        0,
		},
		{
			AgentPoolName: "ABC123!@£",
			Errors:        1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.AgentPoolName, func(t *testing.T) {
			_, errors := KubernetesAgentPoolName(tc.AgentPoolName, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected AgentPoolName to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}

func TestKubernetesDNSPrefix(t *testing.T) {
	cases := []struct {
		DNSPrefix string
		Errors    int
	}{
		{
			DNSPrefix: "",
			Errors:    1,
		},
		{
			DNSPrefix: "aBc-123ab-",
			Errors:    1,
		},
		{
			DNSPrefix: "-aBc-123abc",
			Errors:    1,
		},
		{
			DNSPrefix: "a",
			Errors:    0,
		},
		{
			DNSPrefix: "aBc-123abc",
			Errors:    0,
		},
		{
			DNSPrefix: "ThisIsAKubernetesDNSPrefixThatIsExactlyFiftyFourCharac",
			Errors:    0,
		},
		{
			DNSPrefix: "ThisIsAKubernetesDNSPrefixThatIsNotExactlyFiftyFourChar",
			Errors:    1,
		},
		{
			DNSPrefix: "2",
			Errors:    0,
		},
		{
			DNSPrefix: "2ndCluster",
			Errors:    0,
		},
		{
			DNSPrefix: "aBc-123abc2",
			Errors:    0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.DNSPrefix, func(t *testing.T) {
			_, errors := KubernetesDNSPrefix(tc.DNSPrefix, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected DNSPrefix to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}

func TestKubernetesGitRepositoryUrl(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// start with https://
			Input: "https://github.com/Azure/arc-k8s-demo",
			Valid: true,
		},

		{
			// start with http://
			Input: "http://github.com/Azure/arc-k8s-demo",
			Valid: true,
		},

		{
			// start with git@
			Input: "git@github.com:Azure/arc-k8s-demo.git",
			Valid: true,
		},

		{
			// start with ssh://
			Input: "ssh://git@github.com:Azure/arc-k8s-demo.git",
			Valid: true,
		},

		{
			// random string
			Input: "randomstring",
			Valid: false,
		},
	}
	validationFunction := KubernetesGitRepositoryUrl()
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := validationFunction(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

func TestKubernetesNodeTaint(t *testing.T) {
	cases := []struct {
		NodeTaint string
		Errors    int
	}{
		{
			NodeTaint: "",
			Errors:    1,
		},
		{
			NodeTaint: "key=value:NoSchedule",
			Errors:    0,
		},
		{
			NodeTaint: "key=value:PreferNoSchedule",
			Errors:    0,
		},
		{
			NodeTaint: "key=value:NoExecute",
			Errors:    0,
		},
		{
			// the value is optional
			NodeTaint: "key=:NoSchedule",
			Errors:    0,
		},
		{
			// keys may be prefixed with a DNS subdomain
			NodeTaint: "example.com/key=value:NoSchedule",
			Errors:    0,
		},
		{
			NodeTaint: "Key_1-2.3=Value_1-2.3:NoSchedule",
			Errors:    0,
		},
		{
			// the key is not optional
			NodeTaint: "=value:NoSchedule",
			Errors:    1,
		},
		{
			// missing the effect
			NodeTaint: "key=value",
			Errors:    1,
		},
		{
			// missing the value separator
			NodeTaint: "key:NoSchedule",
			Errors:    1,
		},
		{
			// effects are case sensitive
			NodeTaint: "key=value:noschedule",
			Errors:    1,
		},
		{
			NodeTaint: "key=value:NotAnEffect",
			Errors:    1,
		},
		{
			NodeTaint: "key=value:NoSchedule:NoExecute",
			Errors:    1,
		},
		{
			NodeTaint: "key with spaces=value:NoSchedule",
			Errors:    1,
		},
		{
			NodeTaint: "key=value with spaces:NoSchedule",
			Errors:    1,
		},
		{
			NodeTaint: "key=value:NoSchedule ",
			Errors:    1,
		},
		{
			// a 253 character key is the maximum
			NodeTaint: strings.Repeat("k", 253) + "=value:NoSchedule",
			Errors:    0,
		},
		{
			NodeTaint: strings.Repeat("k", 254) + "=value:NoSchedule",
			Errors:    1,
		},
		{
			// a 63 character value is the maximum
			NodeTaint: "key=" + strings.Repeat("v", 63) + ":NoSchedule",
			Errors:    0,
		},
		{
			NodeTaint: "key=" + strings.Repeat("v", 64) + ":NoSchedule",
			Errors:    1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.NodeTaint, func(t *testing.T) {
			_, errors := KubernetesNodeTaint(tc.NodeTaint, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected NodeTaint to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}
