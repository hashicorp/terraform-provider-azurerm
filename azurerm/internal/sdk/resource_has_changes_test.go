package sdk

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

type SimpleType struct {
	Name string `tfschema:"name"`
}

type hasChangeTestData struct {
	Schema       map[string]*schema.Schema
	ResourceData *schema.ResourceData
	State        *terraform.InstanceState
	Diff         *terraform.InstanceDiff
	Input        SimpleType
	Expected     bool
}

func TestHasChange_Changed(t *testing.T) {
	hasChangeTestData{
		State: nil,
		Diff: &terraform.InstanceDiff{
			Attributes: map[string]*terraform.ResourceAttrDiff{
				"name": {
					Old: "",
					New: "foo",
				},
			},
		},
		Input:    SimpleType{},
		Expected: true,
	}.test(t)
}

func (testData hasChangeTestData) test(t *testing.T) {
	/* This needs further thought on how to test properly
	debugLogger := ConsoleLogger{}
	schema, err := schema.InternalMap(testData.Schema).Data(testData.State, testData.Diff)
	if err != nil {
		t.Fatalf("err: %s", err)
	}


	changed := hasChangesReflectedType(testData.Input.Name, schema, debugLogger)
	if testData.Expected != changed {
		t.Fatalf("\nExpected: %+v\n\n Received %+v\n\n", testData.Expected, changed)
	}
	*/
}
