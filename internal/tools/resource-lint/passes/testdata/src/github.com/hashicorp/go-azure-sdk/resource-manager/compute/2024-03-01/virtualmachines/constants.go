package virtualmachines

// Mock Azure SDK enum types for testing

type VirtualMachinePriorityTypes string

const (
	VirtualMachinePriorityTypesLow     VirtualMachinePriorityTypes = "Low"
	VirtualMachinePriorityTypesRegular VirtualMachinePriorityTypes = "Regular"
	VirtualMachinePriorityTypesSpot    VirtualMachinePriorityTypes = "Spot"
)

type OperatingSystemTypes string

const (
	OperatingSystemTypesLinux   OperatingSystemTypes = "Linux"
	OperatingSystemTypesWindows OperatingSystemTypes = "Windows"
)

func PossibleValuesForVirtualMachinePriorityTypes() []string {
    return []string{
        string(VirtualMachinePriorityTypesLow),
        string(VirtualMachinePriorityTypesRegular),
        string(VirtualMachinePriorityTypesSpot),
    }
}

func PossibleValuesForOperatingSystemTypes() []string {
    return []string{
        string(OperatingSystemTypesLinux),
        string(OperatingSystemTypesWindows),
    }
}
