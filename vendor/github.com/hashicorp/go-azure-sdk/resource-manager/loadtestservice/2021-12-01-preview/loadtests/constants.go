package loadtests

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceState string

const (
	ResourceStateCanceled  ResourceState = "Canceled"
	ResourceStateDeleted   ResourceState = "Deleted"
	ResourceStateFailed    ResourceState = "Failed"
	ResourceStateSucceeded ResourceState = "Succeeded"
)

func PossibleValuesForResourceState() []string {
	return []string{
		string(ResourceStateCanceled),
		string(ResourceStateDeleted),
		string(ResourceStateFailed),
		string(ResourceStateSucceeded),
	}
}

func parseResourceState(input string) (*ResourceState, error) {
	vals := map[string]ResourceState{
		"canceled":  ResourceStateCanceled,
		"deleted":   ResourceStateDeleted,
		"failed":    ResourceStateFailed,
		"succeeded": ResourceStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceState(input)
	return &out, nil
}
