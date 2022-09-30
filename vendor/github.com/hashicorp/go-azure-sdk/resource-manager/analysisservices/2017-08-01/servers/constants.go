package servers

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionMode string

const (
	ConnectionModeAll      ConnectionMode = "All"
	ConnectionModeReadOnly ConnectionMode = "ReadOnly"
)

func PossibleValuesForConnectionMode() []string {
	return []string{
		string(ConnectionModeAll),
		string(ConnectionModeReadOnly),
	}
}

func parseConnectionMode(input string) (*ConnectionMode, error) {
	vals := map[string]ConnectionMode{
		"all":      ConnectionModeAll,
		"readonly": ConnectionModeReadOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionMode(input)
	return &out, nil
}

type ManagedMode int64

const (
	ManagedModeOne  ManagedMode = 1
	ManagedModeZero ManagedMode = 0
)

func PossibleValuesForManagedMode() []int64 {
	return []int64{
		int64(ManagedModeOne),
		int64(ManagedModeZero),
	}
}

func parseManagedMode(input int64) (*ManagedMode, error) {
	vals := map[int64]ManagedMode{
		1: ManagedModeOne,
		0: ManagedModeZero,
	}
	if v, ok := vals[input]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedMode(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStatePaused       ProvisioningState = "Paused"
	ProvisioningStatePausing      ProvisioningState = "Pausing"
	ProvisioningStatePreparing    ProvisioningState = "Preparing"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateResuming     ProvisioningState = "Resuming"
	ProvisioningStateScaling      ProvisioningState = "Scaling"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateSuspended    ProvisioningState = "Suspended"
	ProvisioningStateSuspending   ProvisioningState = "Suspending"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStatePaused),
		string(ProvisioningStatePausing),
		string(ProvisioningStatePreparing),
		string(ProvisioningStateProvisioning),
		string(ProvisioningStateResuming),
		string(ProvisioningStateScaling),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateSuspended),
		string(ProvisioningStateSuspending),
		string(ProvisioningStateUpdating),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"paused":       ProvisioningStatePaused,
		"pausing":      ProvisioningStatePausing,
		"preparing":    ProvisioningStatePreparing,
		"provisioning": ProvisioningStateProvisioning,
		"resuming":     ProvisioningStateResuming,
		"scaling":      ProvisioningStateScaling,
		"succeeded":    ProvisioningStateSucceeded,
		"suspended":    ProvisioningStateSuspended,
		"suspending":   ProvisioningStateSuspending,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type ServerMonitorMode int64

const (
	ServerMonitorModeOne  ServerMonitorMode = 1
	ServerMonitorModeZero ServerMonitorMode = 0
)

func PossibleValuesForServerMonitorMode() []int64 {
	return []int64{
		int64(ServerMonitorModeOne),
		int64(ServerMonitorModeZero),
	}
}

func parseServerMonitorMode(input int64) (*ServerMonitorMode, error) {
	vals := map[int64]ServerMonitorMode{
		1: ServerMonitorModeOne,
		0: ServerMonitorModeZero,
	}
	if v, ok := vals[input]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerMonitorMode(input)
	return &out, nil
}

type SkuTier string

const (
	SkuTierBasic       SkuTier = "Basic"
	SkuTierDevelopment SkuTier = "Development"
	SkuTierStandard    SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierBasic),
		string(SkuTierDevelopment),
		string(SkuTierStandard),
	}
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"basic":       SkuTierBasic,
		"development": SkuTierDevelopment,
		"standard":    SkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuTier(input)
	return &out, nil
}

type State string

const (
	StateDeleting     State = "Deleting"
	StateFailed       State = "Failed"
	StatePaused       State = "Paused"
	StatePausing      State = "Pausing"
	StatePreparing    State = "Preparing"
	StateProvisioning State = "Provisioning"
	StateResuming     State = "Resuming"
	StateScaling      State = "Scaling"
	StateSucceeded    State = "Succeeded"
	StateSuspended    State = "Suspended"
	StateSuspending   State = "Suspending"
	StateUpdating     State = "Updating"
)

func PossibleValuesForState() []string {
	return []string{
		string(StateDeleting),
		string(StateFailed),
		string(StatePaused),
		string(StatePausing),
		string(StatePreparing),
		string(StateProvisioning),
		string(StateResuming),
		string(StateScaling),
		string(StateSucceeded),
		string(StateSuspended),
		string(StateSuspending),
		string(StateUpdating),
	}
}

func parseState(input string) (*State, error) {
	vals := map[string]State{
		"deleting":     StateDeleting,
		"failed":       StateFailed,
		"paused":       StatePaused,
		"pausing":      StatePausing,
		"preparing":    StatePreparing,
		"provisioning": StateProvisioning,
		"resuming":     StateResuming,
		"scaling":      StateScaling,
		"succeeded":    StateSucceeded,
		"suspended":    StateSuspended,
		"suspending":   StateSuspending,
		"updating":     StateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := State(input)
	return &out, nil
}

type Status int64

const (
	StatusZero Status = 0
)

func PossibleValuesForStatus() []int64 {
	return []int64{
		int64(StatusZero),
	}
}

func parseStatus(input int64) (*Status, error) {
	vals := map[int64]Status{
		0: StatusZero,
	}
	if v, ok := vals[input]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Status(input)
	return &out, nil
}
