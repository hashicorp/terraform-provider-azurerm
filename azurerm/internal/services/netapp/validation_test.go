package netapp

import "testing"

func TestValidateNetAppVolumeName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// empty
			input:    "",
			expected: false,
		},
		{
			// basic example
			input:    "hello",
			expected: true,
		},
		{
			// can't start with an underscore
			input:    "_hello",
			expected: false,
		},
		{
			// can't end with a dash
			input:    "hello-",
			expected: true,
		},
		{
			// can't contain an exclamation mark
			input:    "hello!",
			expected: false,
		},
		{
			// dash in the middle
			input:    "malcolm-in-the-middle",
			expected: true,
		},
		{
			// can't end with a period
			input:    "hello.",
			expected: false,
		},
		{
			// 63 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk",
			expected: true,
		},
		{
			// 64 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkj",
			expected: true,
		},
		{
			// 65 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkja",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ValidateNetAppVolumeName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestValidateNetAppVolumeVolumePath(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// empty
			input:    "",
			expected: false,
		},
		{
			// basic example
			input:    "hello",
			expected: true,
		},
		{
			// can't start with an underscore
			input:    "_hello",
			expected: false,
		},
		{
			// can't end with a dash
			input:    "hello-",
			expected: true,
		},
		{
			// can't contain an exclamation mark
			input:    "hello!",
			expected: false,
		},
		{
			// dash in the middle
			input:    "malcolm-in-the-middle",
			expected: true,
		},
		{
			// can't end with a period
			input:    "hello.",
			expected: false,
		},
		{
			// 79 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijabcdefgheysudciac",
			expected: true,
		},
		{
			// 80 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkasbdjdssardwyupac",
			expected: true,
		},
		{
			// 81 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkjspoiuytrewqasdfac",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ValidateNetAppVolumeVolumePath(v.input, "volume_path")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestValidateNetAppSnapshotName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// empty
			input:    "",
			expected: false,
		},
		{
			// basic example
			input:    "hello",
			expected: true,
		},
		{
			// can't start with an underscore
			input:    "_hello",
			expected: false,
		},
		{
			// can't end with a dash
			input:    "hello-",
			expected: true,
		},
		{
			// can't contain an exclamation mark
			input:    "hello!",
			expected: false,
		},
		{
			// dash in the middle
			input:    "malcolm-in-the-middle",
			expected: true,
		},
		{
			// can't end with a period
			input:    "hello.",
			expected: false,
		},
		{
			// 63 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk",
			expected: true,
		},
		{
			// 64 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkj",
			expected: true,
		},
		{
			// 65 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkja",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ValidateNetAppSnapshotName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

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
