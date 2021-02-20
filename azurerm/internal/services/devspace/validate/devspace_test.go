package validate

import "testing"

func TestValidateDevSpaceName(t *testing.T) {
	validNames := []string{
		"valid-name",
		"valid02-name",
		"validName1",
	}
	for _, v := range validNames {
		_, errors := DevSpaceName()(v, "valid")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid DevSpace Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"invalid!",
		"!@£",
		"-invalid",
		"double-hyphen--invalid",
		"invalid_name",
	}
	for _, v := range invalidNames {
		_, errors := DevSpaceName()(v, "invalid")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid DevSpace Name", v)
		}
	}
}
