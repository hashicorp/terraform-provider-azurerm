package integrationaccounts

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventLevel string

const (
	EventLevelCritical      EventLevel = "Critical"
	EventLevelError         EventLevel = "Error"
	EventLevelInformational EventLevel = "Informational"
	EventLevelLogAlways     EventLevel = "LogAlways"
	EventLevelVerbose       EventLevel = "Verbose"
	EventLevelWarning       EventLevel = "Warning"
)

func PossibleValuesForEventLevel() []string {
	return []string{
		string(EventLevelCritical),
		string(EventLevelError),
		string(EventLevelInformational),
		string(EventLevelLogAlways),
		string(EventLevelVerbose),
		string(EventLevelWarning),
	}
}

func (s *EventLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEventLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEventLevel(input string) (*EventLevel, error) {
	vals := map[string]EventLevel{
		"critical":      EventLevelCritical,
		"error":         EventLevelError,
		"informational": EventLevelInformational,
		"logalways":     EventLevelLogAlways,
		"verbose":       EventLevelVerbose,
		"warning":       EventLevelWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EventLevel(input)
	return &out, nil
}

type IntegrationAccountSkuName string

const (
	IntegrationAccountSkuNameBasic        IntegrationAccountSkuName = "Basic"
	IntegrationAccountSkuNameFree         IntegrationAccountSkuName = "Free"
	IntegrationAccountSkuNameNotSpecified IntegrationAccountSkuName = "NotSpecified"
	IntegrationAccountSkuNameStandard     IntegrationAccountSkuName = "Standard"
)

func PossibleValuesForIntegrationAccountSkuName() []string {
	return []string{
		string(IntegrationAccountSkuNameBasic),
		string(IntegrationAccountSkuNameFree),
		string(IntegrationAccountSkuNameNotSpecified),
		string(IntegrationAccountSkuNameStandard),
	}
}

func (s *IntegrationAccountSkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIntegrationAccountSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIntegrationAccountSkuName(input string) (*IntegrationAccountSkuName, error) {
	vals := map[string]IntegrationAccountSkuName{
		"basic":        IntegrationAccountSkuNameBasic,
		"free":         IntegrationAccountSkuNameFree,
		"notspecified": IntegrationAccountSkuNameNotSpecified,
		"standard":     IntegrationAccountSkuNameStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IntegrationAccountSkuName(input)
	return &out, nil
}

type KeyType string

const (
	KeyTypeNotSpecified KeyType = "NotSpecified"
	KeyTypePrimary      KeyType = "Primary"
	KeyTypeSecondary    KeyType = "Secondary"
)

func PossibleValuesForKeyType() []string {
	return []string{
		string(KeyTypeNotSpecified),
		string(KeyTypePrimary),
		string(KeyTypeSecondary),
	}
}

func (s *KeyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyType(input string) (*KeyType, error) {
	vals := map[string]KeyType{
		"notspecified": KeyTypeNotSpecified,
		"primary":      KeyTypePrimary,
		"secondary":    KeyTypeSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyType(input)
	return &out, nil
}

type TrackEventsOperationOptions string

const (
	TrackEventsOperationOptionsDisableSourceInfoEnrich TrackEventsOperationOptions = "DisableSourceInfoEnrich"
	TrackEventsOperationOptionsNone                    TrackEventsOperationOptions = "None"
)

func PossibleValuesForTrackEventsOperationOptions() []string {
	return []string{
		string(TrackEventsOperationOptionsDisableSourceInfoEnrich),
		string(TrackEventsOperationOptionsNone),
	}
}

func (s *TrackEventsOperationOptions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTrackEventsOperationOptions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTrackEventsOperationOptions(input string) (*TrackEventsOperationOptions, error) {
	vals := map[string]TrackEventsOperationOptions{
		"disablesourceinfoenrich": TrackEventsOperationOptionsDisableSourceInfoEnrich,
		"none":                    TrackEventsOperationOptionsNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TrackEventsOperationOptions(input)
	return &out, nil
}

type TrackingRecordType string

const (
	TrackingRecordTypeASTwoMDN                             TrackingRecordType = "AS2MDN"
	TrackingRecordTypeASTwoMessage                         TrackingRecordType = "AS2Message"
	TrackingRecordTypeCustom                               TrackingRecordType = "Custom"
	TrackingRecordTypeEdifactFunctionalGroup               TrackingRecordType = "EdifactFunctionalGroup"
	TrackingRecordTypeEdifactFunctionalGroupAcknowledgment TrackingRecordType = "EdifactFunctionalGroupAcknowledgment"
	TrackingRecordTypeEdifactInterchange                   TrackingRecordType = "EdifactInterchange"
	TrackingRecordTypeEdifactInterchangeAcknowledgment     TrackingRecordType = "EdifactInterchangeAcknowledgment"
	TrackingRecordTypeEdifactTransactionSet                TrackingRecordType = "EdifactTransactionSet"
	TrackingRecordTypeEdifactTransactionSetAcknowledgment  TrackingRecordType = "EdifactTransactionSetAcknowledgment"
	TrackingRecordTypeNotSpecified                         TrackingRecordType = "NotSpecified"
	TrackingRecordTypeXOneTwoFunctionalGroup               TrackingRecordType = "X12FunctionalGroup"
	TrackingRecordTypeXOneTwoFunctionalGroupAcknowledgment TrackingRecordType = "X12FunctionalGroupAcknowledgment"
	TrackingRecordTypeXOneTwoInterchange                   TrackingRecordType = "X12Interchange"
	TrackingRecordTypeXOneTwoInterchangeAcknowledgment     TrackingRecordType = "X12InterchangeAcknowledgment"
	TrackingRecordTypeXOneTwoTransactionSet                TrackingRecordType = "X12TransactionSet"
	TrackingRecordTypeXOneTwoTransactionSetAcknowledgment  TrackingRecordType = "X12TransactionSetAcknowledgment"
)

func PossibleValuesForTrackingRecordType() []string {
	return []string{
		string(TrackingRecordTypeASTwoMDN),
		string(TrackingRecordTypeASTwoMessage),
		string(TrackingRecordTypeCustom),
		string(TrackingRecordTypeEdifactFunctionalGroup),
		string(TrackingRecordTypeEdifactFunctionalGroupAcknowledgment),
		string(TrackingRecordTypeEdifactInterchange),
		string(TrackingRecordTypeEdifactInterchangeAcknowledgment),
		string(TrackingRecordTypeEdifactTransactionSet),
		string(TrackingRecordTypeEdifactTransactionSetAcknowledgment),
		string(TrackingRecordTypeNotSpecified),
		string(TrackingRecordTypeXOneTwoFunctionalGroup),
		string(TrackingRecordTypeXOneTwoFunctionalGroupAcknowledgment),
		string(TrackingRecordTypeXOneTwoInterchange),
		string(TrackingRecordTypeXOneTwoInterchangeAcknowledgment),
		string(TrackingRecordTypeXOneTwoTransactionSet),
		string(TrackingRecordTypeXOneTwoTransactionSetAcknowledgment),
	}
}

func (s *TrackingRecordType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTrackingRecordType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTrackingRecordType(input string) (*TrackingRecordType, error) {
	vals := map[string]TrackingRecordType{
		"as2mdn":                               TrackingRecordTypeASTwoMDN,
		"as2message":                           TrackingRecordTypeASTwoMessage,
		"custom":                               TrackingRecordTypeCustom,
		"edifactfunctionalgroup":               TrackingRecordTypeEdifactFunctionalGroup,
		"edifactfunctionalgroupacknowledgment": TrackingRecordTypeEdifactFunctionalGroupAcknowledgment,
		"edifactinterchange":                   TrackingRecordTypeEdifactInterchange,
		"edifactinterchangeacknowledgment":     TrackingRecordTypeEdifactInterchangeAcknowledgment,
		"edifacttransactionset":                TrackingRecordTypeEdifactTransactionSet,
		"edifacttransactionsetacknowledgment":  TrackingRecordTypeEdifactTransactionSetAcknowledgment,
		"notspecified":                         TrackingRecordTypeNotSpecified,
		"x12functionalgroup":                   TrackingRecordTypeXOneTwoFunctionalGroup,
		"x12functionalgroupacknowledgment":     TrackingRecordTypeXOneTwoFunctionalGroupAcknowledgment,
		"x12interchange":                       TrackingRecordTypeXOneTwoInterchange,
		"x12interchangeacknowledgment":         TrackingRecordTypeXOneTwoInterchangeAcknowledgment,
		"x12transactionset":                    TrackingRecordTypeXOneTwoTransactionSet,
		"x12transactionsetacknowledgment":      TrackingRecordTypeXOneTwoTransactionSetAcknowledgment,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TrackingRecordType(input)
	return &out, nil
}

type WorkflowState string

const (
	WorkflowStateCompleted    WorkflowState = "Completed"
	WorkflowStateDeleted      WorkflowState = "Deleted"
	WorkflowStateDisabled     WorkflowState = "Disabled"
	WorkflowStateEnabled      WorkflowState = "Enabled"
	WorkflowStateNotSpecified WorkflowState = "NotSpecified"
	WorkflowStateSuspended    WorkflowState = "Suspended"
)

func PossibleValuesForWorkflowState() []string {
	return []string{
		string(WorkflowStateCompleted),
		string(WorkflowStateDeleted),
		string(WorkflowStateDisabled),
		string(WorkflowStateEnabled),
		string(WorkflowStateNotSpecified),
		string(WorkflowStateSuspended),
	}
}

func (s *WorkflowState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWorkflowState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWorkflowState(input string) (*WorkflowState, error) {
	vals := map[string]WorkflowState{
		"completed":    WorkflowStateCompleted,
		"deleted":      WorkflowStateDeleted,
		"disabled":     WorkflowStateDisabled,
		"enabled":      WorkflowStateEnabled,
		"notspecified": WorkflowStateNotSpecified,
		"suspended":    WorkflowStateSuspended,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkflowState(input)
	return &out, nil
}
