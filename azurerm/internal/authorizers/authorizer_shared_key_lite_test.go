package authorizers

import (
	"testing"
)

func TestBuildCanonicalizedStringForSharedKeyLite(t *testing.T) {
	testData := []struct {
		name                  string
		headers               map[string][]string
		canonicalizedHeaders  string
		canonicalizedResource string
		verb                  string
		expected              string
	}{
		{
			name: "completed",
			verb: "NOM",
			headers: map[string][]string{
				"Content-MD5":  {"abc123"},
				"Content-Type": {"vnd/panda-pops+v1"},
			},
			canonicalizedHeaders:  "all-the-headers",
			canonicalizedResource: "all-the-resources",
			expected:              "NOM\n\nvnd/panda-pops+v1\n\nall-the-headers\nall-the-resources",
		},
	}

	for _, test := range testData {
		t.Logf("Test: %q", test.name)
		actual := buildCanonicalizedStringForSharedKeyLite(test.verb, test.headers, test.canonicalizedHeaders, test.canonicalizedResource)
		if actual != test.expected {
			t.Fatalf("Expected %q but got %q", test.expected, actual)
		}
	}
}
