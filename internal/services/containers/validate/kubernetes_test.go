// Copyright IBM Corp. 2014, 2026
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
			AgentPoolName: "a",
			Errors:        0,
		},
		{
			AgentPoolName: strings.Repeat("a", 12),
			Errors:        0,
		},
		{
			AgentPoolName: "linuxworkerpool01",
			Errors:        0,
		},
		{
			AgentPoolName: strings.Repeat("a", 26),
			Errors:        0,
		},
		{
			AgentPoolName: strings.Repeat("a", 27),
			Errors:        0,
		},
		{
			AgentPoolName: "a" + strings.Repeat("1", 38),
			Errors:        0,
		},
		{
			AgentPoolName: strings.Repeat("a", 40),
			Errors:        1,
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
			Errors:        0,
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

func TestKubernetesAgentPoolNameWithFeatures(t *testing.T) {
	tests := []struct {
		name      string
		poolName  string
		maxLength int
		wantError bool
	}{
		{name: "one character", poolName: "a", maxLength: 12},
		{name: "default maximum", poolName: strings.Repeat("a", 12), maxLength: 12},
		{name: "over default maximum", poolName: strings.Repeat("a", 13), maxLength: 12, wantError: true},
		{name: "long name without opt in", poolName: "linuxworkerpool01", maxLength: 12, wantError: true},
		{name: "custom maximum", poolName: "linuxworkerpool01", maxLength: 17},
		{name: "over custom maximum", poolName: strings.Repeat("a", 18), maxLength: 17, wantError: true},
		{name: "long name with opt in", poolName: "linuxworkerpool01", maxLength: 39},
		{name: "provider maximum", poolName: strings.Repeat("a", 39), maxLength: 39},
		{name: "over provider maximum", poolName: strings.Repeat("a", 40), maxLength: 39, wantError: true},
		{name: "empty name", maxLength: 39, wantError: true},
		{name: "initial digit", poolName: "1pool", maxLength: 39, wantError: true},
		{name: "uppercase", poolName: "Pool", maxLength: 39, wantError: true},
		{name: "hyphen", poolName: "pool-name", maxLength: 39, wantError: true},
		{name: "maximum below range", poolName: "pool", maxLength: 11, wantError: true},
		{name: "maximum above range", poolName: "pool", maxLength: 40, wantError: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := KubernetesAgentPoolNameWithFeatures(test.poolName, "name", test.maxLength)
			if (err != nil) != test.wantError {
				t.Fatalf("expected error=%t, got %v", test.wantError, err)
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
