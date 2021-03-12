package validate

import "testing"

func TestIoTHubName(t *testing.T) {
	validNames := []string{
		"valid-name",
		"valid02-name",
		"validName1",
		"-validname1",
		"double-hyphen--valid",
	}
	for _, v := range validNames {
		_, errors := IoTHubName(v, "example")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid IoT Hub Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"",
		"invalid_name",
		"invalid!",
		"!@£",
		"hello.world",
	}
	for _, v := range invalidNames {
		_, errors := IoTHubName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid IoT Hub Name", v)
		}
	}
}
