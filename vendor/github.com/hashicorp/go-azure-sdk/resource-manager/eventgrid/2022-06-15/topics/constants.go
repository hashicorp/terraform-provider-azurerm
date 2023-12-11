package topics

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataResidencyBoundary string

const (
	DataResidencyBoundaryWithinGeopair DataResidencyBoundary = "WithinGeopair"
	DataResidencyBoundaryWithinRegion  DataResidencyBoundary = "WithinRegion"
)

func PossibleValuesForDataResidencyBoundary() []string {
	return []string{
		string(DataResidencyBoundaryWithinGeopair),
		string(DataResidencyBoundaryWithinRegion),
	}
}

func (s *DataResidencyBoundary) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataResidencyBoundary(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataResidencyBoundary(input string) (*DataResidencyBoundary, error) {
	vals := map[string]DataResidencyBoundary{
		"withingeopair": DataResidencyBoundaryWithinGeopair,
		"withinregion":  DataResidencyBoundaryWithinRegion,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataResidencyBoundary(input)
	return &out, nil
}

type IPActionType string

const (
	IPActionTypeAllow IPActionType = "Allow"
)

func PossibleValuesForIPActionType() []string {
	return []string{
		string(IPActionTypeAllow),
	}
}

func (s *IPActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIPActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIPActionType(input string) (*IPActionType, error) {
	vals := map[string]IPActionType{
		"allow": IPActionTypeAllow,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IPActionType(input)
	return &out, nil
}

type InputSchema string

const (
	InputSchemaCloudEventSchemaVOneZero InputSchema = "CloudEventSchemaV1_0"
	InputSchemaCustomEventSchema        InputSchema = "CustomEventSchema"
	InputSchemaEventGridSchema          InputSchema = "EventGridSchema"
)

func PossibleValuesForInputSchema() []string {
	return []string{
		string(InputSchemaCloudEventSchemaVOneZero),
		string(InputSchemaCustomEventSchema),
		string(InputSchemaEventGridSchema),
	}
}

func (s *InputSchema) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseInputSchema(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseInputSchema(input string) (*InputSchema, error) {
	vals := map[string]InputSchema{
		"cloudeventschemav1_0": InputSchemaCloudEventSchemaVOneZero,
		"customeventschema":    InputSchemaCustomEventSchema,
		"eventgridschema":      InputSchemaEventGridSchema,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InputSchema(input)
	return &out, nil
}

type InputSchemaMappingType string

const (
	InputSchemaMappingTypeJson InputSchemaMappingType = "Json"
)

func PossibleValuesForInputSchemaMappingType() []string {
	return []string{
		string(InputSchemaMappingTypeJson),
	}
}

func (s *InputSchemaMappingType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseInputSchemaMappingType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseInputSchemaMappingType(input string) (*InputSchemaMappingType, error) {
	vals := map[string]InputSchemaMappingType{
		"json": InputSchemaMappingTypeJson,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InputSchemaMappingType(input)
	return &out, nil
}

type PersistedConnectionStatus string

const (
	PersistedConnectionStatusApproved     PersistedConnectionStatus = "Approved"
	PersistedConnectionStatusDisconnected PersistedConnectionStatus = "Disconnected"
	PersistedConnectionStatusPending      PersistedConnectionStatus = "Pending"
	PersistedConnectionStatusRejected     PersistedConnectionStatus = "Rejected"
)

func PossibleValuesForPersistedConnectionStatus() []string {
	return []string{
		string(PersistedConnectionStatusApproved),
		string(PersistedConnectionStatusDisconnected),
		string(PersistedConnectionStatusPending),
		string(PersistedConnectionStatusRejected),
	}
}

func (s *PersistedConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePersistedConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePersistedConnectionStatus(input string) (*PersistedConnectionStatus, error) {
	vals := map[string]PersistedConnectionStatus{
		"approved":     PersistedConnectionStatusApproved,
		"disconnected": PersistedConnectionStatusDisconnected,
		"pending":      PersistedConnectionStatusPending,
		"rejected":     PersistedConnectionStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PersistedConnectionStatus(input)
	return &out, nil
}

type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled  PublicNetworkAccess = "Enabled"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		string(PublicNetworkAccessDisabled),
		string(PublicNetworkAccessEnabled),
	}
}

func (s *PublicNetworkAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicNetworkAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicNetworkAccess(input string) (*PublicNetworkAccess, error) {
	vals := map[string]PublicNetworkAccess{
		"disabled": PublicNetworkAccessDisabled,
		"enabled":  PublicNetworkAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccess(input)
	return &out, nil
}

type ResourceProvisioningState string

const (
	ResourceProvisioningStateCanceled  ResourceProvisioningState = "Canceled"
	ResourceProvisioningStateCreating  ResourceProvisioningState = "Creating"
	ResourceProvisioningStateDeleting  ResourceProvisioningState = "Deleting"
	ResourceProvisioningStateFailed    ResourceProvisioningState = "Failed"
	ResourceProvisioningStateSucceeded ResourceProvisioningState = "Succeeded"
	ResourceProvisioningStateUpdating  ResourceProvisioningState = "Updating"
)

func PossibleValuesForResourceProvisioningState() []string {
	return []string{
		string(ResourceProvisioningStateCanceled),
		string(ResourceProvisioningStateCreating),
		string(ResourceProvisioningStateDeleting),
		string(ResourceProvisioningStateFailed),
		string(ResourceProvisioningStateSucceeded),
		string(ResourceProvisioningStateUpdating),
	}
}

func (s *ResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceProvisioningState(input string) (*ResourceProvisioningState, error) {
	vals := map[string]ResourceProvisioningState{
		"canceled":  ResourceProvisioningStateCanceled,
		"creating":  ResourceProvisioningStateCreating,
		"deleting":  ResourceProvisioningStateDeleting,
		"failed":    ResourceProvisioningStateFailed,
		"succeeded": ResourceProvisioningStateSucceeded,
		"updating":  ResourceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceProvisioningState(input)
	return &out, nil
}

type TopicProvisioningState string

const (
	TopicProvisioningStateCanceled  TopicProvisioningState = "Canceled"
	TopicProvisioningStateCreating  TopicProvisioningState = "Creating"
	TopicProvisioningStateDeleting  TopicProvisioningState = "Deleting"
	TopicProvisioningStateFailed    TopicProvisioningState = "Failed"
	TopicProvisioningStateSucceeded TopicProvisioningState = "Succeeded"
	TopicProvisioningStateUpdating  TopicProvisioningState = "Updating"
)

func PossibleValuesForTopicProvisioningState() []string {
	return []string{
		string(TopicProvisioningStateCanceled),
		string(TopicProvisioningStateCreating),
		string(TopicProvisioningStateDeleting),
		string(TopicProvisioningStateFailed),
		string(TopicProvisioningStateSucceeded),
		string(TopicProvisioningStateUpdating),
	}
}

func (s *TopicProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTopicProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTopicProvisioningState(input string) (*TopicProvisioningState, error) {
	vals := map[string]TopicProvisioningState{
		"canceled":  TopicProvisioningStateCanceled,
		"creating":  TopicProvisioningStateCreating,
		"deleting":  TopicProvisioningStateDeleting,
		"failed":    TopicProvisioningStateFailed,
		"succeeded": TopicProvisioningStateSucceeded,
		"updating":  TopicProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TopicProvisioningState(input)
	return &out, nil
}
