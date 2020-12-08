package validate

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/validate"
)

func TestAzureRMRedisFirewallRuleName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "ab",
			ErrCount: 0,
		},
		{
			Value:    "abc",
			ErrCount: 0,
		},
		{
			Value:    "webapp1",
			ErrCount: 0,
		},
		{
			Value:    "hello-world",
			ErrCount: 1,
		},
		{
			Value:    "hello_world",
			ErrCount: 0,
		},
		{
			Value:    "helloworld21!",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validate.FirewallRuleName(tc.Value, "azurerm_redis_firewall_rule")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Redis Firewall Rule Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}
