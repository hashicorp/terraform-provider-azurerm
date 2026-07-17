package namespacetopics

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventInputSchema string

const (
	EventInputSchemaCloudEventSchemaVOneZero EventInputSchema = "CloudEventSchemaV1_0"
)

func PossibleValuesForEventInputSchema() []string {
	return []string{
		string(EventInputSchemaCloudEventSchemaVOneZero),
	}
}

func (s *EventInputSchema) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEventInputSchema(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEventInputSchema(input string) (*EventInputSchema, error) {
	vals := map[string]EventInputSchema{
		"cloudeventschemav1_0": EventInputSchemaCloudEventSchemaVOneZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EventInputSchema(input)
	return &out, nil
}

type NamespaceTopicProvisioningState string

const (
	NamespaceTopicProvisioningStateCanceled      NamespaceTopicProvisioningState = "Canceled"
	NamespaceTopicProvisioningStateCreateFailed  NamespaceTopicProvisioningState = "CreateFailed"
	NamespaceTopicProvisioningStateCreating      NamespaceTopicProvisioningState = "Creating"
	NamespaceTopicProvisioningStateDeleteFailed  NamespaceTopicProvisioningState = "DeleteFailed"
	NamespaceTopicProvisioningStateDeleted       NamespaceTopicProvisioningState = "Deleted"
	NamespaceTopicProvisioningStateDeleting      NamespaceTopicProvisioningState = "Deleting"
	NamespaceTopicProvisioningStateFailed        NamespaceTopicProvisioningState = "Failed"
	NamespaceTopicProvisioningStateSucceeded     NamespaceTopicProvisioningState = "Succeeded"
	NamespaceTopicProvisioningStateUpdatedFailed NamespaceTopicProvisioningState = "UpdatedFailed"
	NamespaceTopicProvisioningStateUpdating      NamespaceTopicProvisioningState = "Updating"
)

func PossibleValuesForNamespaceTopicProvisioningState() []string {
	return []string{
		string(NamespaceTopicProvisioningStateCanceled),
		string(NamespaceTopicProvisioningStateCreateFailed),
		string(NamespaceTopicProvisioningStateCreating),
		string(NamespaceTopicProvisioningStateDeleteFailed),
		string(NamespaceTopicProvisioningStateDeleted),
		string(NamespaceTopicProvisioningStateDeleting),
		string(NamespaceTopicProvisioningStateFailed),
		string(NamespaceTopicProvisioningStateSucceeded),
		string(NamespaceTopicProvisioningStateUpdatedFailed),
		string(NamespaceTopicProvisioningStateUpdating),
	}
}

func (s *NamespaceTopicProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNamespaceTopicProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNamespaceTopicProvisioningState(input string) (*NamespaceTopicProvisioningState, error) {
	vals := map[string]NamespaceTopicProvisioningState{
		"canceled":      NamespaceTopicProvisioningStateCanceled,
		"createfailed":  NamespaceTopicProvisioningStateCreateFailed,
		"creating":      NamespaceTopicProvisioningStateCreating,
		"deletefailed":  NamespaceTopicProvisioningStateDeleteFailed,
		"deleted":       NamespaceTopicProvisioningStateDeleted,
		"deleting":      NamespaceTopicProvisioningStateDeleting,
		"failed":        NamespaceTopicProvisioningStateFailed,
		"succeeded":     NamespaceTopicProvisioningStateSucceeded,
		"updatedfailed": NamespaceTopicProvisioningStateUpdatedFailed,
		"updating":      NamespaceTopicProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NamespaceTopicProvisioningState(input)
	return &out, nil
}

type PublisherType string

const (
	PublisherTypeCustom PublisherType = "Custom"
)

func PossibleValuesForPublisherType() []string {
	return []string{
		string(PublisherTypeCustom),
	}
}

func (s *PublisherType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublisherType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublisherType(input string) (*PublisherType, error) {
	vals := map[string]PublisherType{
		"custom": PublisherTypeCustom,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublisherType(input)
	return &out, nil
}
