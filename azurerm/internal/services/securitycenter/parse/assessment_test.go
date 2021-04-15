package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = AssessmentId{}

func TestAssessmentIDFormatter(t *testing.T) {
	actual := NewAssessmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleset1", "assessment1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleset1/providers/Microsoft.Security/assessments/assessment1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestAssessmentID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *AssessmentId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Security resource provider",
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleset1",
			Error: true,
		},
		{
			Name:  "No target resource Segment",
			Input: "/providers/Microsoft.Security/assessments/assessment1",
			Error: true,
		},
		{
			Name:  "No Security Center Assessment Segment",
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleset1/providers/Microsoft.Security/",
			Error: true,
		},
		{
			Name:  "No Security Center Assessment name",
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleset1/providers/Microsoft.Security/assessments/",
			Error: true,
		},
		{
			Name:  "ID of Security Center Assessment",
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleset1/providers/Microsoft.Security/assessments/assessment1",
			Error: false,
			Expect: &AssessmentId{
				TargetResourceID: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleset1",
				Name:             "assessment1",
			},
		},
		{
			Name:  "Wrong Casing",
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.COMPUTE/VIRTUALMACHINESCALESETS/SCALESET1/PROVIDERS/MICROSOFT.SECURITY/ASSESSMENTS/ASSESSMENT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := AssessmentID(v.Input)
		if err != nil {
			if v.Expect == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}
		if actual.TargetResourceID != v.Expect.TargetResourceID {
			t.Fatalf("Expected %q but got %q for TargetResourceID", v.Expect.TargetResourceID, actual.TargetResourceID)
		}
	}
}
