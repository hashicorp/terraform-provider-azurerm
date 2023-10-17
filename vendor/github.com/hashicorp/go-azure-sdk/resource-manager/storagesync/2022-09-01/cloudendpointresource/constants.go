package cloudendpointresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ChangeDetectionMode string

const (
	ChangeDetectionModeDefault   ChangeDetectionMode = "Default"
	ChangeDetectionModeRecursive ChangeDetectionMode = "Recursive"
)

func PossibleValuesForChangeDetectionMode() []string {
	return []string{
		string(ChangeDetectionModeDefault),
		string(ChangeDetectionModeRecursive),
	}
}

func (s *ChangeDetectionMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseChangeDetectionMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseChangeDetectionMode(input string) (*ChangeDetectionMode, error) {
	vals := map[string]ChangeDetectionMode{
		"default":   ChangeDetectionModeDefault,
		"recursive": ChangeDetectionModeRecursive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ChangeDetectionMode(input)
	return &out, nil
}

type CloudEndpointChangeEnumerationActivityState string

const (
	CloudEndpointChangeEnumerationActivityStateEnumerationInProgress        CloudEndpointChangeEnumerationActivityState = "EnumerationInProgress"
	CloudEndpointChangeEnumerationActivityStateInitialEnumerationInProgress CloudEndpointChangeEnumerationActivityState = "InitialEnumerationInProgress"
)

func PossibleValuesForCloudEndpointChangeEnumerationActivityState() []string {
	return []string{
		string(CloudEndpointChangeEnumerationActivityStateEnumerationInProgress),
		string(CloudEndpointChangeEnumerationActivityStateInitialEnumerationInProgress),
	}
}

func (s *CloudEndpointChangeEnumerationActivityState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCloudEndpointChangeEnumerationActivityState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCloudEndpointChangeEnumerationActivityState(input string) (*CloudEndpointChangeEnumerationActivityState, error) {
	vals := map[string]CloudEndpointChangeEnumerationActivityState{
		"enumerationinprogress":        CloudEndpointChangeEnumerationActivityStateEnumerationInProgress,
		"initialenumerationinprogress": CloudEndpointChangeEnumerationActivityStateInitialEnumerationInProgress,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CloudEndpointChangeEnumerationActivityState(input)
	return &out, nil
}

type CloudEndpointChangeEnumerationTotalCountsState string

const (
	CloudEndpointChangeEnumerationTotalCountsStateCalculating CloudEndpointChangeEnumerationTotalCountsState = "Calculating"
	CloudEndpointChangeEnumerationTotalCountsStateFinal       CloudEndpointChangeEnumerationTotalCountsState = "Final"
)

func PossibleValuesForCloudEndpointChangeEnumerationTotalCountsState() []string {
	return []string{
		string(CloudEndpointChangeEnumerationTotalCountsStateCalculating),
		string(CloudEndpointChangeEnumerationTotalCountsStateFinal),
	}
}

func (s *CloudEndpointChangeEnumerationTotalCountsState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCloudEndpointChangeEnumerationTotalCountsState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCloudEndpointChangeEnumerationTotalCountsState(input string) (*CloudEndpointChangeEnumerationTotalCountsState, error) {
	vals := map[string]CloudEndpointChangeEnumerationTotalCountsState{
		"calculating": CloudEndpointChangeEnumerationTotalCountsStateCalculating,
		"final":       CloudEndpointChangeEnumerationTotalCountsStateFinal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CloudEndpointChangeEnumerationTotalCountsState(input)
	return &out, nil
}
