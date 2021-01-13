package validate

import "testing"

func TestCacheBackupFrequency_validation(t *testing.T) {
	cases := []struct {
		Value    int
		ErrCount int
	}{
		{Value: 1, ErrCount: 1},
		{Value: 15, ErrCount: 0},
		{Value: 30, ErrCount: 0},
		{Value: 45, ErrCount: 1},
		{Value: 60, ErrCount: 0},
		{Value: 120, ErrCount: 1},
		{Value: 240, ErrCount: 1},
		{Value: 360, ErrCount: 0},
		{Value: 720, ErrCount: 0},
		{Value: 1440, ErrCount: 0},
	}

	for _, tc := range cases {
		_, errors := CacheBackupFrequency(tc.Value, "azurerm_redis_cache")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the AzureRM Redis Cache Backup Frequency to trigger a validation error for '%d'", tc.Value)
		}
	}
}
