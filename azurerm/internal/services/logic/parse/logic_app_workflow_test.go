package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
	"testing"
)

var _ resourceid.Formatter = LogicAppWorkflowId{}

func TestLogicAppWorkflowIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	actual := NewLogicAppWorkflowID("group1", "workflow1").ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Logic/workflows/workflow1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestLogicAppWorkflowID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *LogicAppWorkflowId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Resource Groups Segment",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111",
			Error: true,
		},
		{
			Name:  "No Resource Groups Value",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111/resourceGroups/",
			Error: true,
		},
		{
			Name:  "No Workflow Name",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111/resourceGroups/resGroup1/providers/Microsoft.Logic/workflows",
			Error: true,
		},
		{
			Name:  "Correct case",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111/resourceGroups/resGroup1/providers/Microsoft.Logic/workflows/workflow1",
			Expect: &LogicAppWorkflowId{
				ResourceGroup: "resGroup1",
				Name:          "workflow1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := LogicAppWorkflowID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}
	}
}
