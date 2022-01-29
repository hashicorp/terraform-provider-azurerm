package validate

import "testing"

func TestLocationID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "/",
			Valid: false,
		},
		{
			Input: "/subscriptions/",
			Valid: false,
		},
		{
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Valid: false,
		},
		{
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/locations",
			Valid: true,
		},
		{
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/LOCATIONS",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := LocationID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
