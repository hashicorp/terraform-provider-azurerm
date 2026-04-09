package datastores

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatastoreProvisioningState string

const (
	DatastoreProvisioningStateCanceled  DatastoreProvisioningState = "Canceled"
	DatastoreProvisioningStateCancelled DatastoreProvisioningState = "Cancelled"
	DatastoreProvisioningStateCreating  DatastoreProvisioningState = "Creating"
	DatastoreProvisioningStateDeleting  DatastoreProvisioningState = "Deleting"
	DatastoreProvisioningStateFailed    DatastoreProvisioningState = "Failed"
	DatastoreProvisioningStatePending   DatastoreProvisioningState = "Pending"
	DatastoreProvisioningStateSucceeded DatastoreProvisioningState = "Succeeded"
	DatastoreProvisioningStateUpdating  DatastoreProvisioningState = "Updating"
)

func PossibleValuesForDatastoreProvisioningState() []string {
	return []string{
		string(DatastoreProvisioningStateCanceled),
		string(DatastoreProvisioningStateCancelled),
		string(DatastoreProvisioningStateCreating),
		string(DatastoreProvisioningStateDeleting),
		string(DatastoreProvisioningStateFailed),
		string(DatastoreProvisioningStatePending),
		string(DatastoreProvisioningStateSucceeded),
		string(DatastoreProvisioningStateUpdating),
	}
}

func (s *DatastoreProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatastoreProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatastoreProvisioningState(input string) (*DatastoreProvisioningState, error) {
	vals := map[string]DatastoreProvisioningState{
		"canceled":  DatastoreProvisioningStateCanceled,
		"cancelled": DatastoreProvisioningStateCancelled,
		"creating":  DatastoreProvisioningStateCreating,
		"deleting":  DatastoreProvisioningStateDeleting,
		"failed":    DatastoreProvisioningStateFailed,
		"pending":   DatastoreProvisioningStatePending,
		"succeeded": DatastoreProvisioningStateSucceeded,
		"updating":  DatastoreProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatastoreProvisioningState(input)
	return &out, nil
}

type DatastoreStatus string

const (
	DatastoreStatusAccessible        DatastoreStatus = "Accessible"
	DatastoreStatusAttached          DatastoreStatus = "Attached"
	DatastoreStatusDeadOrError       DatastoreStatus = "DeadOrError"
	DatastoreStatusDetached          DatastoreStatus = "Detached"
	DatastoreStatusInaccessible      DatastoreStatus = "Inaccessible"
	DatastoreStatusLostCommunication DatastoreStatus = "LostCommunication"
	DatastoreStatusUnknown           DatastoreStatus = "Unknown"
)

func PossibleValuesForDatastoreStatus() []string {
	return []string{
		string(DatastoreStatusAccessible),
		string(DatastoreStatusAttached),
		string(DatastoreStatusDeadOrError),
		string(DatastoreStatusDetached),
		string(DatastoreStatusInaccessible),
		string(DatastoreStatusLostCommunication),
		string(DatastoreStatusUnknown),
	}
}

func (s *DatastoreStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatastoreStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatastoreStatus(input string) (*DatastoreStatus, error) {
	vals := map[string]DatastoreStatus{
		"accessible":        DatastoreStatusAccessible,
		"attached":          DatastoreStatusAttached,
		"deadorerror":       DatastoreStatusDeadOrError,
		"detached":          DatastoreStatusDetached,
		"inaccessible":      DatastoreStatusInaccessible,
		"lostcommunication": DatastoreStatusLostCommunication,
		"unknown":           DatastoreStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatastoreStatus(input)
	return &out, nil
}

type MountOptionEnum string

const (
	MountOptionEnumATTACH MountOptionEnum = "ATTACH"
	MountOptionEnumMOUNT  MountOptionEnum = "MOUNT"
)

func PossibleValuesForMountOptionEnum() []string {
	return []string{
		string(MountOptionEnumATTACH),
		string(MountOptionEnumMOUNT),
	}
}

func (s *MountOptionEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMountOptionEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMountOptionEnum(input string) (*MountOptionEnum, error) {
	vals := map[string]MountOptionEnum{
		"attach": MountOptionEnumATTACH,
		"mount":  MountOptionEnumMOUNT,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MountOptionEnum(input)
	return &out, nil
}
