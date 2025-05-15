package partnertopics

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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

type PartnerTopicActivationState string

const (
	PartnerTopicActivationStateActivated      PartnerTopicActivationState = "Activated"
	PartnerTopicActivationStateDeactivated    PartnerTopicActivationState = "Deactivated"
	PartnerTopicActivationStateNeverActivated PartnerTopicActivationState = "NeverActivated"
)

func PossibleValuesForPartnerTopicActivationState() []string {
	return []string{
		string(PartnerTopicActivationStateActivated),
		string(PartnerTopicActivationStateDeactivated),
		string(PartnerTopicActivationStateNeverActivated),
	}
}

func (s *PartnerTopicActivationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePartnerTopicActivationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePartnerTopicActivationState(input string) (*PartnerTopicActivationState, error) {
	vals := map[string]PartnerTopicActivationState{
		"activated":      PartnerTopicActivationStateActivated,
		"deactivated":    PartnerTopicActivationStateDeactivated,
		"neveractivated": PartnerTopicActivationStateNeverActivated,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PartnerTopicActivationState(input)
	return &out, nil
}

type PartnerTopicProvisioningState string

const (
	PartnerTopicProvisioningStateCanceled                                 PartnerTopicProvisioningState = "Canceled"
	PartnerTopicProvisioningStateCreating                                 PartnerTopicProvisioningState = "Creating"
	PartnerTopicProvisioningStateDeleting                                 PartnerTopicProvisioningState = "Deleting"
	PartnerTopicProvisioningStateFailed                                   PartnerTopicProvisioningState = "Failed"
	PartnerTopicProvisioningStateIdleDueToMirroredChannelResourceDeletion PartnerTopicProvisioningState = "IdleDueToMirroredChannelResourceDeletion"
	PartnerTopicProvisioningStateSucceeded                                PartnerTopicProvisioningState = "Succeeded"
	PartnerTopicProvisioningStateUpdating                                 PartnerTopicProvisioningState = "Updating"
)

func PossibleValuesForPartnerTopicProvisioningState() []string {
	return []string{
		string(PartnerTopicProvisioningStateCanceled),
		string(PartnerTopicProvisioningStateCreating),
		string(PartnerTopicProvisioningStateDeleting),
		string(PartnerTopicProvisioningStateFailed),
		string(PartnerTopicProvisioningStateIdleDueToMirroredChannelResourceDeletion),
		string(PartnerTopicProvisioningStateSucceeded),
		string(PartnerTopicProvisioningStateUpdating),
	}
}

func (s *PartnerTopicProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePartnerTopicProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePartnerTopicProvisioningState(input string) (*PartnerTopicProvisioningState, error) {
	vals := map[string]PartnerTopicProvisioningState{
		"canceled": PartnerTopicProvisioningStateCanceled,
		"creating": PartnerTopicProvisioningStateCreating,
		"deleting": PartnerTopicProvisioningStateDeleting,
		"failed":   PartnerTopicProvisioningStateFailed,
		"idleduetomirroredchannelresourcedeletion": PartnerTopicProvisioningStateIdleDueToMirroredChannelResourceDeletion,
		"succeeded": PartnerTopicProvisioningStateSucceeded,
		"updating":  PartnerTopicProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PartnerTopicProvisioningState(input)
	return &out, nil
}
