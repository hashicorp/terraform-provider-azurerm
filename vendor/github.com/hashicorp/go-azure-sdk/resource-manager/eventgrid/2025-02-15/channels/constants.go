package channels

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ChannelProvisioningState string

const (
	ChannelProvisioningStateCanceled                              ChannelProvisioningState = "Canceled"
	ChannelProvisioningStateCreating                              ChannelProvisioningState = "Creating"
	ChannelProvisioningStateDeleting                              ChannelProvisioningState = "Deleting"
	ChannelProvisioningStateFailed                                ChannelProvisioningState = "Failed"
	ChannelProvisioningStateIdleDueToMirroredPartnerTopicDeletion ChannelProvisioningState = "IdleDueToMirroredPartnerTopicDeletion"
	ChannelProvisioningStateSucceeded                             ChannelProvisioningState = "Succeeded"
	ChannelProvisioningStateUpdating                              ChannelProvisioningState = "Updating"
)

func PossibleValuesForChannelProvisioningState() []string {
	return []string{
		string(ChannelProvisioningStateCanceled),
		string(ChannelProvisioningStateCreating),
		string(ChannelProvisioningStateDeleting),
		string(ChannelProvisioningStateFailed),
		string(ChannelProvisioningStateIdleDueToMirroredPartnerTopicDeletion),
		string(ChannelProvisioningStateSucceeded),
		string(ChannelProvisioningStateUpdating),
	}
}

func (s *ChannelProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseChannelProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseChannelProvisioningState(input string) (*ChannelProvisioningState, error) {
	vals := map[string]ChannelProvisioningState{
		"canceled":                              ChannelProvisioningStateCanceled,
		"creating":                              ChannelProvisioningStateCreating,
		"deleting":                              ChannelProvisioningStateDeleting,
		"failed":                                ChannelProvisioningStateFailed,
		"idleduetomirroredpartnertopicdeletion": ChannelProvisioningStateIdleDueToMirroredPartnerTopicDeletion,
		"succeeded":                             ChannelProvisioningStateSucceeded,
		"updating":                              ChannelProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ChannelProvisioningState(input)
	return &out, nil
}

type ChannelType string

const (
	ChannelTypePartnerTopic ChannelType = "PartnerTopic"
)

func PossibleValuesForChannelType() []string {
	return []string{
		string(ChannelTypePartnerTopic),
	}
}

func (s *ChannelType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseChannelType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseChannelType(input string) (*ChannelType, error) {
	vals := map[string]ChannelType{
		"partnertopic": ChannelTypePartnerTopic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ChannelType(input)
	return &out, nil
}

type EventDefinitionKind string

const (
	EventDefinitionKindInline EventDefinitionKind = "Inline"
)

func PossibleValuesForEventDefinitionKind() []string {
	return []string{
		string(EventDefinitionKindInline),
	}
}

func (s *EventDefinitionKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEventDefinitionKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEventDefinitionKind(input string) (*EventDefinitionKind, error) {
	vals := map[string]EventDefinitionKind{
		"inline": EventDefinitionKindInline,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EventDefinitionKind(input)
	return &out, nil
}

type ReadinessState string

const (
	ReadinessStateActivated      ReadinessState = "Activated"
	ReadinessStateNeverActivated ReadinessState = "NeverActivated"
)

func PossibleValuesForReadinessState() []string {
	return []string{
		string(ReadinessStateActivated),
		string(ReadinessStateNeverActivated),
	}
}

func (s *ReadinessState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReadinessState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReadinessState(input string) (*ReadinessState, error) {
	vals := map[string]ReadinessState{
		"activated":      ReadinessStateActivated,
		"neveractivated": ReadinessStateNeverActivated,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReadinessState(input)
	return &out, nil
}
