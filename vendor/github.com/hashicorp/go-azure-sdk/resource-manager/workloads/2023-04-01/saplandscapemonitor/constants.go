package saplandscapemonitor

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SapLandscapeMonitorProvisioningState string

const (
	SapLandscapeMonitorProvisioningStateAccepted  SapLandscapeMonitorProvisioningState = "Accepted"
	SapLandscapeMonitorProvisioningStateCanceled  SapLandscapeMonitorProvisioningState = "Canceled"
	SapLandscapeMonitorProvisioningStateCreated   SapLandscapeMonitorProvisioningState = "Created"
	SapLandscapeMonitorProvisioningStateFailed    SapLandscapeMonitorProvisioningState = "Failed"
	SapLandscapeMonitorProvisioningStateSucceeded SapLandscapeMonitorProvisioningState = "Succeeded"
)

func PossibleValuesForSapLandscapeMonitorProvisioningState() []string {
	return []string{
		string(SapLandscapeMonitorProvisioningStateAccepted),
		string(SapLandscapeMonitorProvisioningStateCanceled),
		string(SapLandscapeMonitorProvisioningStateCreated),
		string(SapLandscapeMonitorProvisioningStateFailed),
		string(SapLandscapeMonitorProvisioningStateSucceeded),
	}
}

func parseSapLandscapeMonitorProvisioningState(input string) (*SapLandscapeMonitorProvisioningState, error) {
	vals := map[string]SapLandscapeMonitorProvisioningState{
		"accepted":  SapLandscapeMonitorProvisioningStateAccepted,
		"canceled":  SapLandscapeMonitorProvisioningStateCanceled,
		"created":   SapLandscapeMonitorProvisioningStateCreated,
		"failed":    SapLandscapeMonitorProvisioningStateFailed,
		"succeeded": SapLandscapeMonitorProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SapLandscapeMonitorProvisioningState(input)
	return &out, nil
}
