package networksecurityperimeterlinks

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NspLinkProvisioningState string

const (
	NspLinkProvisioningStateAccepted                NspLinkProvisioningState = "Accepted"
	NspLinkProvisioningStateCreating                NspLinkProvisioningState = "Creating"
	NspLinkProvisioningStateDeleting                NspLinkProvisioningState = "Deleting"
	NspLinkProvisioningStateFailed                  NspLinkProvisioningState = "Failed"
	NspLinkProvisioningStateSucceeded               NspLinkProvisioningState = "Succeeded"
	NspLinkProvisioningStateUpdating                NspLinkProvisioningState = "Updating"
	NspLinkProvisioningStateWaitForRemoteCompletion NspLinkProvisioningState = "WaitForRemoteCompletion"
)

func PossibleValuesForNspLinkProvisioningState() []string {
	return []string{
		string(NspLinkProvisioningStateAccepted),
		string(NspLinkProvisioningStateCreating),
		string(NspLinkProvisioningStateDeleting),
		string(NspLinkProvisioningStateFailed),
		string(NspLinkProvisioningStateSucceeded),
		string(NspLinkProvisioningStateUpdating),
		string(NspLinkProvisioningStateWaitForRemoteCompletion),
	}
}

func (s *NspLinkProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNspLinkProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNspLinkProvisioningState(input string) (*NspLinkProvisioningState, error) {
	vals := map[string]NspLinkProvisioningState{
		"accepted":                NspLinkProvisioningStateAccepted,
		"creating":                NspLinkProvisioningStateCreating,
		"deleting":                NspLinkProvisioningStateDeleting,
		"failed":                  NspLinkProvisioningStateFailed,
		"succeeded":               NspLinkProvisioningStateSucceeded,
		"updating":                NspLinkProvisioningStateUpdating,
		"waitforremotecompletion": NspLinkProvisioningStateWaitForRemoteCompletion,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NspLinkProvisioningState(input)
	return &out, nil
}

type NspLinkStatus string

const (
	NspLinkStatusApproved     NspLinkStatus = "Approved"
	NspLinkStatusDisconnected NspLinkStatus = "Disconnected"
	NspLinkStatusPending      NspLinkStatus = "Pending"
	NspLinkStatusRejected     NspLinkStatus = "Rejected"
)

func PossibleValuesForNspLinkStatus() []string {
	return []string{
		string(NspLinkStatusApproved),
		string(NspLinkStatusDisconnected),
		string(NspLinkStatusPending),
		string(NspLinkStatusRejected),
	}
}

func (s *NspLinkStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNspLinkStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNspLinkStatus(input string) (*NspLinkStatus, error) {
	vals := map[string]NspLinkStatus{
		"approved":     NspLinkStatusApproved,
		"disconnected": NspLinkStatusDisconnected,
		"pending":      NspLinkStatusPending,
		"rejected":     NspLinkStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NspLinkStatus(input)
	return &out, nil
}
