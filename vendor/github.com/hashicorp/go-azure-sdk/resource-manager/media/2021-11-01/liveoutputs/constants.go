package liveoutputs

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiveOutputResourceState string

const (
	LiveOutputResourceStateCreating LiveOutputResourceState = "Creating"
	LiveOutputResourceStateDeleting LiveOutputResourceState = "Deleting"
	LiveOutputResourceStateRunning  LiveOutputResourceState = "Running"
)

func PossibleValuesForLiveOutputResourceState() []string {
	return []string{
		string(LiveOutputResourceStateCreating),
		string(LiveOutputResourceStateDeleting),
		string(LiveOutputResourceStateRunning),
	}
}

func parseLiveOutputResourceState(input string) (*LiveOutputResourceState, error) {
	vals := map[string]LiveOutputResourceState{
		"creating": LiveOutputResourceStateCreating,
		"deleting": LiveOutputResourceStateDeleting,
		"running":  LiveOutputResourceStateRunning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LiveOutputResourceState(input)
	return &out, nil
}
