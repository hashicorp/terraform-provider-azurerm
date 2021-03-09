package validate

import "testing"

func TestRedisEnterpriseClusterSkuName(t *testing.T) {
	testData := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Invalid Empty sku_name value",
			input:    "",
			expected: false,
		},
		{
			name:     "Invalid sku name with no delimiter or capacity",
			input:    "EnterpriseFlash_F1500",
			expected: false,
		},
		{
			name:     "Invalid sku name with invalid capacity above max int32",
			input:    "EnterpriseFlash_F1500-2147483649",
			expected: false,
		},
		{
			name:     "Invalid sku name with empty capacity and extra delimiter",
			input:    "TEA_4Node_UPS_Heater--",
			expected: false,
		},
		{
			name:     "Invalid sku name with delimiter and no capacity defined",
			input:    "EnterpriseFlash_F1500-",
			expected: false,
		},
		{
			name:     "Invalid wrong sku name casing",
			input:    "Enterpriseflash_F1500-15",
			expected: false,
		},
		{
			name:     "Invalid No sku with delimiter and invalid capacity",
			input:    "-YEET!",
			expected: false,
		},
		{
			name:     "Invalid sku name with delimiter and invalid capacity",
			input:    "EnterpriseFlash_F1500-2",
			expected: false,
		},
		{
			name:     "Valid Ignore extra delimiter",
			input:    "Enterprise_E100-2-IgnoredInput",
			expected: true,
		},
		{
			name:     "Valid sku and capacity",
			input:    "EnterpriseFlash_F1500-15",
			expected: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.name)

		_, errors := RedisEnterpriseClusterSkuName(v.input, "sku_name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
