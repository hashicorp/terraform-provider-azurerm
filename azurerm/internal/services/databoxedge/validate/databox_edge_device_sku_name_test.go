package validate

import "testing"

func TestDataboxEdgeSkuName(t *testing.T) {
	testData := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Empty sku_name value",
			input:    "",
			expected: false,
		},
		{
			name:     "Ignore extra delimiter",
			input:    "TEA_4Node_UPS_Heater-Standard-IgnoredInput",
			expected: true,
		},
		{
			name:     "Valid sku with no delimiter or tier",
			input:    "TEA_4Node_UPS_Heater",
			expected: true,
		},
		{
			name:     "valid sku with empty tier and extra delimiter",
			input:    "TEA_4Node_UPS_Heater--",
			expected: true,
		},
		{
			name:     "Valid sku with delimiter and no tier defined",
			input:    "TEA_4Node_UPS_Heater-",
			expected: true,
		},
		{
			name:     "Valid sku and tier",
			input:    "TEA_4Node_UPS_Heater-Standard",
			expected: true,
		},
		{
			name:     "Valid sku and tier with wrong sku casing",
			input:    "tea_4Node_UPS_Heater-Standard",
			expected: false,
		},
		{
			name:     "Valid sku and tier with wrong tier casing",
			input:    "TEA_4Node_UPS_Heater-standard",
			expected: false,
		},
		{
			name:     "No sku with delimiter and invalid tier",
			input:    "-YEET!",
			expected: false,
		},
		{
			name:     "Valid sku with delimiter and invalid Tier",
			input:    "Gateway-MemoryOptimized",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.name)

		_, errors := DataboxEdgeDeviceSkuName(v.input, "sku_name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
