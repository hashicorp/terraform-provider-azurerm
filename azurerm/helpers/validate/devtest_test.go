package validate

import "testing"

func TestValidateDevTestLabName(t *testing.T) {
	validNames := []string{
		"valid-name",
		"valid02-name",
		"validName1",
		"-validname1",
		"valid_name",
		"double-hyphen--valid",
	}
	for _, v := range validNames {
		_, errors := DevTestLabName()(v, "example")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Dev Test Lab Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"invalid!",
		"!@£",
	}
	for _, v := range invalidNames {
		_, errors := DevTestLabName()(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Dev Test Lab Name", v)
		}
	}
}
