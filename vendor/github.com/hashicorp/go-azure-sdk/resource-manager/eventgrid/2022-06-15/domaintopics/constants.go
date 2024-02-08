package domaintopics

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DomainTopicProvisioningState string

const (
	DomainTopicProvisioningStateCanceled  DomainTopicProvisioningState = "Canceled"
	DomainTopicProvisioningStateCreating  DomainTopicProvisioningState = "Creating"
	DomainTopicProvisioningStateDeleting  DomainTopicProvisioningState = "Deleting"
	DomainTopicProvisioningStateFailed    DomainTopicProvisioningState = "Failed"
	DomainTopicProvisioningStateSucceeded DomainTopicProvisioningState = "Succeeded"
	DomainTopicProvisioningStateUpdating  DomainTopicProvisioningState = "Updating"
)

func PossibleValuesForDomainTopicProvisioningState() []string {
	return []string{
		string(DomainTopicProvisioningStateCanceled),
		string(DomainTopicProvisioningStateCreating),
		string(DomainTopicProvisioningStateDeleting),
		string(DomainTopicProvisioningStateFailed),
		string(DomainTopicProvisioningStateSucceeded),
		string(DomainTopicProvisioningStateUpdating),
	}
}

func (s *DomainTopicProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDomainTopicProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDomainTopicProvisioningState(input string) (*DomainTopicProvisioningState, error) {
	vals := map[string]DomainTopicProvisioningState{
		"canceled":  DomainTopicProvisioningStateCanceled,
		"creating":  DomainTopicProvisioningStateCreating,
		"deleting":  DomainTopicProvisioningStateDeleting,
		"failed":    DomainTopicProvisioningStateFailed,
		"succeeded": DomainTopicProvisioningStateSucceeded,
		"updating":  DomainTopicProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DomainTopicProvisioningState(input)
	return &out, nil
}
