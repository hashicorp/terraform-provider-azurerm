package validate

import "testing"

func TestMaxMemoryPolicy_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{Value: "noeviction", ErrCount: 0},
		{Value: "allkeys-lru", ErrCount: 0},
		{Value: "volatile-lru", ErrCount: 0},
		{Value: "allkeys-random", ErrCount: 0},
		{Value: "volatile-random", ErrCount: 0},
		{Value: "volatile-ttl", ErrCount: 0},
		{Value: "allkeys-lfu", ErrCount: 0},
		{Value: "volatile-lfu", ErrCount: 0},
		{Value: "something-else", ErrCount: 1},
	}

	for _, tc := range cases {
		_, errors := MaxMemoryPolicy(tc.Value, "azurerm_redis_cache")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM Redis Cache Max Memory Policy to trigger a validation error")
		}
	}
}
