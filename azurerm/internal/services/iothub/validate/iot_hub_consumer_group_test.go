package validate

import "testing"

func TestIoTHubConsumerGroupName(t *testing.T) {
	validNames := []string{
		"valid-name",
		"valid02-name",
		"validName1",
		"-validname1",
		"valid_name",
		"double-hyphen--valid",
		"hello.world",
	}
	for _, v := range validNames {
		_, errors := IoTHubConsumerGroupName(v, "example")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid IoT Hub Consumer Group Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"",
		"invalid!",
		"!@£",
	}
	for _, v := range invalidNames {
		_, errors := IoTHubConsumerGroupName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid IoT Hub Consumer Group Name", v)
		}
	}
}
