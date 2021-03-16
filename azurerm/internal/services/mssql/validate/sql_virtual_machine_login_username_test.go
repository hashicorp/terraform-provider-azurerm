package validate

import "testing"

func TestSqlVirtualMachineLoginUserName(t *testing.T) {
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{"dfasdlk", false},
		{"sdfs@ ", false},
		{"dfsjsiajfiweangfvnjaksdflaklsdjdjskfamlkcsdflamkldfklafamsdklfmlaksjfdkadklsfmklamdklsfakldsflamkslfmlkeamkldmfkamfmdkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk", true},
		{"60", false},
		{"7.d", true},
		{"u i", true},
		{"a", true},
	}

	for _, test := range testCases {
		_, es := SqlVirtualMachineLoginUserName(test.input, "name")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating name %q to fail", test.input)
		}
	}
}
