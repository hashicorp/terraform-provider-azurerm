package clientgroups

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClientGroupProvisioningState string

const (
	ClientGroupProvisioningStateCanceled  ClientGroupProvisioningState = "Canceled"
	ClientGroupProvisioningStateCreating  ClientGroupProvisioningState = "Creating"
	ClientGroupProvisioningStateDeleted   ClientGroupProvisioningState = "Deleted"
	ClientGroupProvisioningStateDeleting  ClientGroupProvisioningState = "Deleting"
	ClientGroupProvisioningStateFailed    ClientGroupProvisioningState = "Failed"
	ClientGroupProvisioningStateSucceeded ClientGroupProvisioningState = "Succeeded"
	ClientGroupProvisioningStateUpdating  ClientGroupProvisioningState = "Updating"
)

func PossibleValuesForClientGroupProvisioningState() []string {
	return []string{
		string(ClientGroupProvisioningStateCanceled),
		string(ClientGroupProvisioningStateCreating),
		string(ClientGroupProvisioningStateDeleted),
		string(ClientGroupProvisioningStateDeleting),
		string(ClientGroupProvisioningStateFailed),
		string(ClientGroupProvisioningStateSucceeded),
		string(ClientGroupProvisioningStateUpdating),
	}
}

func (s *ClientGroupProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClientGroupProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClientGroupProvisioningState(input string) (*ClientGroupProvisioningState, error) {
	vals := map[string]ClientGroupProvisioningState{
		"canceled":  ClientGroupProvisioningStateCanceled,
		"creating":  ClientGroupProvisioningStateCreating,
		"deleted":   ClientGroupProvisioningStateDeleted,
		"deleting":  ClientGroupProvisioningStateDeleting,
		"failed":    ClientGroupProvisioningStateFailed,
		"succeeded": ClientGroupProvisioningStateSucceeded,
		"updating":  ClientGroupProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClientGroupProvisioningState(input)
	return &out, nil
}
