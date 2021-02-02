package validate

import "testing"

func TestValidateDataFactoryPipelineAndTriggerName(t *testing.T) {
	validNames := []string{
		"validname",
		"valid02name",
		"validName1",
	}
	for _, v := range validNames {
		_, errors := DataFactoryPipelineAndTriggerName()(v, "valid")
		if len(errors) != 0 {
			t.Fatalf("%q should be an invalid DataFactory Pipeline or Trigger Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"invalid.",
		":@£",
		">invalid",
		"invalid&name",
	}
	for _, v := range invalidNames {
		_, errors := DataFactoryPipelineAndTriggerName()(v, "invalid")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid DataFactory Pipeline or Trigger Name", v)
		}
	}
}

func TestValidateDataFactoryName(t *testing.T) {
	validNames := []string{
		"valid-name",
		"valid02-name",
		"validName1",
	}
	for _, v := range validNames {
		_, errors := DataFactoryName()(v, "valid")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid DataFactory Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"invalid.",
		":@£",
		">invalid",
		"invalid&name",
	}
	for _, v := range invalidNames {
		_, errors := DataFactoryName()(v, "invalid")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid DataFactory Name", v)
		}
	}
}
