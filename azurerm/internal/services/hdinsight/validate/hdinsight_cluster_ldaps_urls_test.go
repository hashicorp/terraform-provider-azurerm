package validate

import "testing"

func TestHDInsightClusterLdapsUrls(t *testing.T) {
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
			input:    "ldaps://",
			expected: true,
		},
		{
			// can't start with an underscore
			input:    "http://",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := HDInsightClusterLdapsUrls(v.input, "ldaps_urls")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
