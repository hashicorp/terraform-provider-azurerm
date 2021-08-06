package netapp

import (
	"testing"
)

func TestValidateSlicesEquality(t *testing.T) {
	testData := []struct {
		input1   []string
		input2   []string
		input3   bool
		expected bool
	}{
		{
			// Same order, case sensitive
			input1:   []string{"CIFS", "NFSv3"},
			input2:   []string{"CIFS", "NFSv3"},
			input3:   true,
			expected: true,
		},
		{
			// Same order, case insensitive
			input1:   []string{"CIFS", "NFSv3"},
			input2:   []string{"cifs", "nfsv3"},
			input3:   false,
			expected: true,
		},
		{
			// Reversed order, case sensitive
			input1:   []string{"CIFS", "NFSv3"},
			input2:   []string{"NFSv3", "CIFS"},
			input3:   true,
			expected: true,
		},
		{
			// Reversed order, case insensitive
			input1:   []string{"cifs", "nfsv3"},
			input2:   []string{"NFSv3", "CIFS"},
			input3:   false,
			expected: true,
		},

		{
			// Different, case sensitive
			input1:   []string{"CIFS", "NFSv3"},
			input2:   []string{"NFSv3"},
			input3:   true,
			expected: false,
		},
		{
			// Different, case insensitive
			input1:   []string{"CIFS", "NFSv3"},
			input2:   []string{"nfsv3"},
			input3:   false,
			expected: false,
		},
		{
			// Different, single slices, case sensitive
			input1:   []string{"CIFS"},
			input2:   []string{"NFSv3"},
			input3:   true,
			expected: false,
		},
		{
			// Different, single slices, case insensitive
			input1:   []string{"cifs"},
			input2:   []string{"NFSv3"},
			input3:   false,
			expected: false,
		},
		{
			// Same, single slices, case sensitive
			input1:   []string{"CIFS"},
			input2:   []string{"CIFS"},
			input3:   true,
			expected: true,
		},
		{
			// Different, single slices, case insensitive
			input1:   []string{"cifs"},
			input2:   []string{"CIFS"},
			input3:   false,
			expected: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %+v and %+v for %v where 'caseSensitive' = %v result..", v.input1, v.input2, v.expected, v.input3)

		actual := ValidateSlicesEquality(v.input1, v.input2, v.input3)
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
