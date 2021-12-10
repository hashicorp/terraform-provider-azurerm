package validate

import (
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
			DNSPrefix: "a",
			Errors:    1,
		},
		{
			DNSPrefix: "aBc-123abc",
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
