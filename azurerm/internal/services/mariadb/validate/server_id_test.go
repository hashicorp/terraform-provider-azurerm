package validate

import "testing"

func TestValidateServerID(t *testing.T) {
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
			// invalid
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg",
			expected: false,
		},
		{
			// valid
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DBforMariaDB/servers/test-mariadb",
			expected: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ServerID(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
