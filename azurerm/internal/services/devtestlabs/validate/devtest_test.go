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

func TestValidateDevTestVirtualMachineName(t *testing.T) {
	validNames := []string{
		"valid-name",
		"valid02-ne",
		"validName1",
		"1hello",
		"hello1",
		"1hello1",
		"dbl--valid",
	}
	for _, v := range validNames {
		_, errors := DevTestVirtualMachineName(10)(v, "example")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Dev Test Virtual Machine Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"",
		"-invalidname1",
		"thisnameiswaytoolong",
		"12345",
		"in_valid",
		"-hello",
		"hello-",
		"1hello-",
		"invalid!",
		"!@£",
	}
	for _, v := range invalidNames {
		_, errors := DevTestVirtualMachineName(10)(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Dev Test Virtual Machine Name", v)
		}
	}
}
