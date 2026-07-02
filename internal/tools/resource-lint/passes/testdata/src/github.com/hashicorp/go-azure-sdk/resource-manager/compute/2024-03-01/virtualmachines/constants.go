// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

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

type ShutdownOnIdleMode string

const (
	ShutdownOnIdleModeNone        ShutdownOnIdleMode = "None"
	ShutdownOnIdleModeUserAbsence ShutdownOnIdleMode = "UserAbsence"
	ShutdownOnIdleModeLowUsage    ShutdownOnIdleMode = "LowUsage"
)

func PossibleValuesForShutdownOnIdleMode() []string {
	return []string{
		string(ShutdownOnIdleModeNone),
		string(ShutdownOnIdleModeUserAbsence),
		string(ShutdownOnIdleModeLowUsage),
	}
}

type GetOperationOptions struct {
	Expand *string
}

func DefaultGetOperationOptions() GetOperationOptions {
	return GetOperationOptions{}
}

type DeleteOperationOptions struct {
	Force *bool
}

func DefaultDeleteOperationOptions() DeleteOperationOptions {
	return DeleteOperationOptions{}
}

type ListOperationOptions struct {
	Filter *string
}
