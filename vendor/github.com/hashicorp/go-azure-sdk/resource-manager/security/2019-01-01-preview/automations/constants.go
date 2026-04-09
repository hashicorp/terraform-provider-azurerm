package automations

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionType string

const (
	ActionTypeEventHub  ActionType = "EventHub"
	ActionTypeLogicApp  ActionType = "LogicApp"
	ActionTypeWorkspace ActionType = "Workspace"
)

func PossibleValuesForActionType() []string {
	return []string{
		string(ActionTypeEventHub),
		string(ActionTypeLogicApp),
		string(ActionTypeWorkspace),
	}
}

func (s *ActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseActionType(input string) (*ActionType, error) {
	vals := map[string]ActionType{
		"eventhub":  ActionTypeEventHub,
		"logicapp":  ActionTypeLogicApp,
		"workspace": ActionTypeWorkspace,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActionType(input)
	return &out, nil
}

type EventSource string

const (
	EventSourceAlerts                                 EventSource = "Alerts"
	EventSourceAssessments                            EventSource = "Assessments"
	EventSourceAssessmentsSnapshot                    EventSource = "AssessmentsSnapshot"
	EventSourceRegulatoryComplianceAssessment         EventSource = "RegulatoryComplianceAssessment"
	EventSourceRegulatoryComplianceAssessmentSnapshot EventSource = "RegulatoryComplianceAssessmentSnapshot"
	EventSourceSecureScoreControls                    EventSource = "SecureScoreControls"
	EventSourceSecureScoreControlsSnapshot            EventSource = "SecureScoreControlsSnapshot"
	EventSourceSecureScores                           EventSource = "SecureScores"
	EventSourceSecureScoresSnapshot                   EventSource = "SecureScoresSnapshot"
	EventSourceSubAssessments                         EventSource = "SubAssessments"
	EventSourceSubAssessmentsSnapshot                 EventSource = "SubAssessmentsSnapshot"
)

func PossibleValuesForEventSource() []string {
	return []string{
		string(EventSourceAlerts),
		string(EventSourceAssessments),
		string(EventSourceAssessmentsSnapshot),
		string(EventSourceRegulatoryComplianceAssessment),
		string(EventSourceRegulatoryComplianceAssessmentSnapshot),
		string(EventSourceSecureScoreControls),
		string(EventSourceSecureScoreControlsSnapshot),
		string(EventSourceSecureScores),
		string(EventSourceSecureScoresSnapshot),
		string(EventSourceSubAssessments),
		string(EventSourceSubAssessmentsSnapshot),
	}
}

func (s *EventSource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEventSource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEventSource(input string) (*EventSource, error) {
	vals := map[string]EventSource{
		"alerts":                                 EventSourceAlerts,
		"assessments":                            EventSourceAssessments,
		"assessmentssnapshot":                    EventSourceAssessmentsSnapshot,
		"regulatorycomplianceassessment":         EventSourceRegulatoryComplianceAssessment,
		"regulatorycomplianceassessmentsnapshot": EventSourceRegulatoryComplianceAssessmentSnapshot,
		"securescorecontrols":                    EventSourceSecureScoreControls,
		"securescorecontrolssnapshot":            EventSourceSecureScoreControlsSnapshot,
		"securescores":                           EventSourceSecureScores,
		"securescoressnapshot":                   EventSourceSecureScoresSnapshot,
		"subassessments":                         EventSourceSubAssessments,
		"subassessmentssnapshot":                 EventSourceSubAssessmentsSnapshot,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EventSource(input)
	return &out, nil
}

type Operator string

const (
	OperatorContains             Operator = "Contains"
	OperatorEndsWith             Operator = "EndsWith"
	OperatorEquals               Operator = "Equals"
	OperatorGreaterThan          Operator = "GreaterThan"
	OperatorGreaterThanOrEqualTo Operator = "GreaterThanOrEqualTo"
	OperatorLesserThan           Operator = "LesserThan"
	OperatorLesserThanOrEqualTo  Operator = "LesserThanOrEqualTo"
	OperatorNotEquals            Operator = "NotEquals"
	OperatorStartsWith           Operator = "StartsWith"
)

func PossibleValuesForOperator() []string {
	return []string{
		string(OperatorContains),
		string(OperatorEndsWith),
		string(OperatorEquals),
		string(OperatorGreaterThan),
		string(OperatorGreaterThanOrEqualTo),
		string(OperatorLesserThan),
		string(OperatorLesserThanOrEqualTo),
		string(OperatorNotEquals),
		string(OperatorStartsWith),
	}
}

func (s *Operator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperator(input string) (*Operator, error) {
	vals := map[string]Operator{
		"contains":             OperatorContains,
		"endswith":             OperatorEndsWith,
		"equals":               OperatorEquals,
		"greaterthan":          OperatorGreaterThan,
		"greaterthanorequalto": OperatorGreaterThanOrEqualTo,
		"lesserthan":           OperatorLesserThan,
		"lesserthanorequalto":  OperatorLesserThanOrEqualTo,
		"notequals":            OperatorNotEquals,
		"startswith":           OperatorStartsWith,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Operator(input)
	return &out, nil
}

type PropertyType string

const (
	PropertyTypeBoolean PropertyType = "Boolean"
	PropertyTypeInteger PropertyType = "Integer"
	PropertyTypeNumber  PropertyType = "Number"
	PropertyTypeString  PropertyType = "String"
)

func PossibleValuesForPropertyType() []string {
	return []string{
		string(PropertyTypeBoolean),
		string(PropertyTypeInteger),
		string(PropertyTypeNumber),
		string(PropertyTypeString),
	}
}

func (s *PropertyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePropertyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePropertyType(input string) (*PropertyType, error) {
	vals := map[string]PropertyType{
		"boolean": PropertyTypeBoolean,
		"integer": PropertyTypeInteger,
		"number":  PropertyTypeNumber,
		"string":  PropertyTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PropertyType(input)
	return &out, nil
}
