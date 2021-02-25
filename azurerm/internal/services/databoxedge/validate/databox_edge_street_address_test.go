package validate

import "testing"

func TestDataboxEdgeStreetAddress(t *testing.T) {
	testData := []struct {
		input    []interface{}
		expected bool
	}{
		{
			input:    make([]interface{}, 0),
			expected: false,
		},
		{
			input:    []interface{}{""},
			expected: false,
		},
		{
			input:    []interface{}{"12 Grimmauld Place"},
			expected: true,
		},
		{
			input:    []interface{}{"12 Grimmauld Place", "740 Evergreen Terrace"},
			expected: true,
		},
		{
			input:    []interface{}{"12 Grimmauld Place", "740 Evergreen Terrace", "129 West 81st Street, Apartment: 5A"},
			expected: true,
		},
		{
			input:    []interface{}{"129 West 81st Street , Apartment: 5A"},
			expected: false,
		},
		{
			input:    []interface{}{"740 Evergreen Terrace", "129 West 81st Street , Apartment: 5A"},
			expected: false,
		},
		{
			input:    []interface{}{"12 Grimmauld Place", "740 Evergreen Terrace", "129 West 81st Street , Apartment: 5A"},
			expected: false,
		},
		{
			input:    []interface{}{"12 Grimmauld Place", "740 Evergreen Terrace", "129 West 81st Street , Apartment: 5A", "What? Four lines... no!"},
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := DataboxEdgeStreetAddress(v.input, "address")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
