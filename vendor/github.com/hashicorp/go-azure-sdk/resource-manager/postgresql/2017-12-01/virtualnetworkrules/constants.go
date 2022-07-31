package virtualnetworkrules

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkRuleState string

const (
	VirtualNetworkRuleStateDeleting     VirtualNetworkRuleState = "Deleting"
	VirtualNetworkRuleStateInProgress   VirtualNetworkRuleState = "InProgress"
	VirtualNetworkRuleStateInitializing VirtualNetworkRuleState = "Initializing"
	VirtualNetworkRuleStateReady        VirtualNetworkRuleState = "Ready"
	VirtualNetworkRuleStateUnknown      VirtualNetworkRuleState = "Unknown"
)

func PossibleValuesForVirtualNetworkRuleState() []string {
	return []string{
		string(VirtualNetworkRuleStateDeleting),
		string(VirtualNetworkRuleStateInProgress),
		string(VirtualNetworkRuleStateInitializing),
		string(VirtualNetworkRuleStateReady),
		string(VirtualNetworkRuleStateUnknown),
	}
}

func parseVirtualNetworkRuleState(input string) (*VirtualNetworkRuleState, error) {
	vals := map[string]VirtualNetworkRuleState{
		"deleting":     VirtualNetworkRuleStateDeleting,
		"inprogress":   VirtualNetworkRuleStateInProgress,
		"initializing": VirtualNetworkRuleStateInitializing,
		"ready":        VirtualNetworkRuleStateReady,
		"unknown":      VirtualNetworkRuleStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkRuleState(input)
	return &out, nil
}
