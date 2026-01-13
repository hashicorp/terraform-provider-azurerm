package azbp003

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
)

func validCase1() *virtualmachines.VirtualMachinePriorityTypes {
	priority := "Spot"
	return pointer.ToEnum(virtualmachines.VirtualMachinePriorityTypes(priority))
}

func validCase2() *virtualmachines.OperatingSystemTypes {
	return pointer.ToEnum(virtualmachines.OperatingSystemTypes("Linux"))
}

func validCase3(config map[string]interface{}) *virtualmachines.VirtualMachinePriorityTypes {
	return pointer.ToEnum(virtualmachines.VirtualMachinePriorityTypes(config["priority"].(string)))
}

func validCase4() *string {
	return pointer.To("regular string")
}

func validCase5() *int {
	return pointer.To(42)
}

func invalidCase1() *virtualmachines.VirtualMachinePriorityTypes {
	priority := "Spot"
	return pointer.To(virtualmachines.VirtualMachinePriorityTypes(priority)) // want `AZBP003`
}

func invalidCase2() *virtualmachines.OperatingSystemTypes {
	return pointer.To(virtualmachines.OperatingSystemTypes("Linux")) // want `AZBP003`
}

func invalidCase3(config map[string]interface{}) *virtualmachines.VirtualMachinePriorityTypes {
	return pointer.To(virtualmachines.VirtualMachinePriorityTypes(config["priority"].(string))) // want `AZBP003`
}

func invalidCase4() *virtualmachines.OperatingSystemTypes {
	osType := "Windows"
	return pointer.To(virtualmachines.OperatingSystemTypes(osType)) // want `AZBP003`
}
