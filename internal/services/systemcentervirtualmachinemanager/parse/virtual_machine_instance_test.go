package parse

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SystemCenterVirtualMachineManagerVirtualMachineInstanceId{}

func TestSystemCenterVirtualMachineManagerVirtualMachineInstanceIDFormatter(t *testing.T) {
	actual := NewSystemCenterVirtualMachineManagerVirtualMachineInstanceID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.HybridCompute/machines/machine1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.ScVmm/virtualMachineInstances/default"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSystemCenterVirtualMachineManagerVirtualMachineInstanceID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SystemCenterVirtualMachineManagerVirtualMachineInstanceId
	}{
		{
			Input: "",
			Error: true,
		},
		{
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.HybridCompute/machines/machine1",
			Error: true,
		},
		{
			Input: "/providers/Microsoft.ScVmm/virtualMachineInstances/default",
			Error: true,
		},
		{
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.ScVmm/virtualMachineInstances/default",
			Expected: &SystemCenterVirtualMachineManagerVirtualMachineInstanceId{
				Scope: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.HybridCompute/machines/machine1",
			},
		},
		{
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.HYBRIDCOMPUTE/MACHINES/MACHINE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SystemCenterVirtualMachineManagerVirtualMachineInstanceID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.Scope != v.Expected.Scope {
			t.Fatalf("Expected %q but got %q for ScopeId", v.Expected.Scope, actual.Scope)
		}
	}
}
