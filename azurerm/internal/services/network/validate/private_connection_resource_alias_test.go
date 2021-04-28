package validate

import "testing"

func TestPrivateConnectionResourceAlias(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{

		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// missing suffix
			Input: "example-privatelinkservice.d20286c8-4ea5-11eb-9584-8f53157226c6.centralus",
			Valid: false,
		},

		{
			// valid
			Input: "example-privatelinkservice.d20286c8-4ea5-11eb-9584-8f53157226c6.centralus.azure.privatelinkservice",
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := PrivateConnectionResourceAlias(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
