package liveoutputs

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AsyncOperationStatus string

const (
	AsyncOperationStatusFailed     AsyncOperationStatus = "Failed"
	AsyncOperationStatusInProgress AsyncOperationStatus = "InProgress"
	AsyncOperationStatusSucceeded  AsyncOperationStatus = "Succeeded"
)

func PossibleValuesForAsyncOperationStatus() []string {
	return []string{
		string(AsyncOperationStatusFailed),
		string(AsyncOperationStatusInProgress),
		string(AsyncOperationStatusSucceeded),
	}
}

func parseAsyncOperationStatus(input string) (*AsyncOperationStatus, error) {
	vals := map[string]AsyncOperationStatus{
		"failed":     AsyncOperationStatusFailed,
		"inprogress": AsyncOperationStatusInProgress,
		"succeeded":  AsyncOperationStatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AsyncOperationStatus(input)
	return &out, nil
}

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
