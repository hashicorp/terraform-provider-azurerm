package networksecurityperimeteraccessrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessRuleDirection string

const (
	AccessRuleDirectionInbound  AccessRuleDirection = "Inbound"
	AccessRuleDirectionOutbound AccessRuleDirection = "Outbound"
)

func PossibleValuesForAccessRuleDirection() []string {
	return []string{
		string(AccessRuleDirectionInbound),
		string(AccessRuleDirectionOutbound),
	}
}

func (s *AccessRuleDirection) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccessRuleDirection(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccessRuleDirection(input string) (*AccessRuleDirection, error) {
	vals := map[string]AccessRuleDirection{
		"inbound":  AccessRuleDirectionInbound,
		"outbound": AccessRuleDirectionOutbound,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessRuleDirection(input)
	return &out, nil
}

type NspProvisioningState string

const (
	NspProvisioningStateAccepted  NspProvisioningState = "Accepted"
	NspProvisioningStateCreating  NspProvisioningState = "Creating"
	NspProvisioningStateDeleting  NspProvisioningState = "Deleting"
	NspProvisioningStateFailed    NspProvisioningState = "Failed"
	NspProvisioningStateSucceeded NspProvisioningState = "Succeeded"
	NspProvisioningStateUpdating  NspProvisioningState = "Updating"
)

func PossibleValuesForNspProvisioningState() []string {
	return []string{
		string(NspProvisioningStateAccepted),
		string(NspProvisioningStateCreating),
		string(NspProvisioningStateDeleting),
		string(NspProvisioningStateFailed),
		string(NspProvisioningStateSucceeded),
		string(NspProvisioningStateUpdating),
	}
}

func (s *NspProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNspProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNspProvisioningState(input string) (*NspProvisioningState, error) {
	vals := map[string]NspProvisioningState{
		"accepted":  NspProvisioningStateAccepted,
		"creating":  NspProvisioningStateCreating,
		"deleting":  NspProvisioningStateDeleting,
		"failed":    NspProvisioningStateFailed,
		"succeeded": NspProvisioningStateSucceeded,
		"updating":  NspProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NspProvisioningState(input)
	return &out, nil
}
