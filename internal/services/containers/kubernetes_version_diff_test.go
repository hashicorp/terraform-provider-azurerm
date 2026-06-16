package containers

import "testing"

func TestSuppressKubernetesVersionDiff(t *testing.T) {
	testCases := []struct {
		name     string
		old      string // value currently in state
		new      string // configured value
		expected bool   // true == diff suppressed
	}{
		{
			name:     "minor alias matches running patch version is suppressed",
			old:      "1.27.14",
			new:      "1.27",
			expected: true,
		},
		{
			name:     "minor alias with different minor is a real change",
			old:      "1.27.14",
			new:      "1.28",
			expected: false,
		},
		{
			name:     "minor alias with different major is a real change",
			old:      "1.27.14",
			new:      "2.27",
			expected: false,
		},
		{
			name:     "exact patch pin differing from running patch is a real change",
			old:      "1.27.14",
			new:      "1.27.15",
			expected: false,
		},
		{
			name:     "exact patch pin matching running patch is not suppressed by this func",
			old:      "1.27.14",
			new:      "1.27.14",
			expected: false,
		},
		{
			name:     "empty state (create) is not suppressed",
			old:      "",
			new:      "1.27",
			expected: false,
		},
		{
			name:     "empty config is not suppressed",
			old:      "1.27.14",
			new:      "",
			expected: false,
		},
		{
			name:     "alias matching an alias state with same minor is suppressed",
			old:      "1.27",
			new:      "1.27",
			expected: true,
		},
		{
			name:     "alias matching an alias state with different minor is a real change",
			old:      "1.27",
			new:      "1.28",
			expected: false,
		},
		{
			name:     "double-digit minor alias matches running patch version is suppressed",
			old:      "1.34.3",
			new:      "1.34",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := suppressKubernetesVersionDiff("", tc.old, tc.new, nil)
			if actual != tc.expected {
				t.Fatalf("suppressKubernetesVersionDiff(old=%q, new=%q) = %t, expected %t", tc.old, tc.new, actual, tc.expected)
			}
		})
	}
}
